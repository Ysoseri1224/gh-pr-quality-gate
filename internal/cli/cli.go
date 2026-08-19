package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/config"
	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/gate"
	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/protection"
	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/report"
	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/scaffold"
)

func Run(args []string, stdout, stderr io.Writer, version string) error {
	if len(args) == 0 {
		writeUsage(stdout)
		return nil
	}
	switch args[0] {
	case "help", "--help", "-h":
		writeUsage(stdout)
		return nil
	case "version", "--version", "-v":
		fmt.Fprintln(stdout, version)
		return nil
	case "audit":
		return runAudit(args[1:], stdout, stderr, false)
	case "validate":
		return runAudit(args[1:], stdout, stderr, true)
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "protect":
		return runProtect(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func writeUsage(w io.Writer) {
	fmt.Fprintln(w, `gh pr-quality-gate audits and installs repository quality gates.

Usage:
  gh pr-quality-gate audit [flags]
  gh pr-quality-gate init [--apply] [flags]
  gh pr-quality-gate validate [--run-local] [flags]
  gh pr-quality-gate protect [--show-current | --apply --confirm owner/repo:branch] [flags]
  gh pr-quality-gate version

Mutating commands are dry-run by default. This extension never merges pull requests,
changes draft state, force-pushes, or closes issues.`)
}

type commonFlags struct {
	repo       string
	configPath string
	json       bool
}

func addCommon(flags *flag.FlagSet, values *commonFlags) {
	flags.StringVar(&values.repo, "repo", ".", "repository working tree")
	flags.StringVar(&values.configPath, "config", config.FileName, "configuration path, relative to repository")
	flags.BoolVar(&values.json, "json", false, "write machine-readable JSON")
}

func runAudit(args []string, stdout, stderr io.Writer, validate bool) error {
	name := "audit"
	if validate {
		name = "validate"
	}
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(stderr)
	var common commonFlags
	var bodyFile string
	var runLocal bool
	addCommon(flags, &common)
	flags.StringVar(&bodyFile, "pr-body-file", "", "file containing the pull request body")
	if validate {
		flags.BoolVar(&runLocal, "run-local", false, "run configured local quality commands")
	}
	if err := flags.Parse(args); err != nil {
		return err
	}
	root, err := repositoryRoot(common.repo)
	if err != nil {
		return err
	}
	cfg, err := loadConfig(root, common.configPath)
	if err != nil {
		if !validate && errors.Is(err, os.ErrNotExist) {
			cfg = config.Default()
			result := gate.Audit(root, cfg, "")
			result.Findings = append([]report.Finding{{Level: report.Fail, Check: "configuration", Message: common.configPath + " is missing"}}, result.Findings...)
			if writeErr := result.Write(stdout, common.json); writeErr != nil {
				return writeErr
			}
			return errors.New("quality gate audit found missing configuration")
		}
		return err
	}
	body, err := readPRBody(bodyFile)
	if err != nil {
		return err
	}
	result := gate.Audit(root, cfg, body)
	if err := result.Write(stdout, common.json); err != nil {
		return err
	}
	if result.HasFailures() {
		return errors.New("quality gate validation failed")
	}
	if runLocal {
		return gate.RunLocal(root, cfg, stdout, stderr)
	}
	return nil
}

func runInit(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var repo string
	var apply bool
	var asJSON bool
	flags.StringVar(&repo, "repo", ".", "repository working tree")
	flags.BoolVar(&apply, "apply", false, "create files after a conflict-free dry run")
	flags.BoolVar(&asJSON, "json", false, "write machine-readable JSON")
	if err := flags.Parse(args); err != nil {
		return err
	}
	root, err := repositoryRoot(repo)
	if err != nil {
		return err
	}
	changes, err := scaffold.Plan(root)
	if err != nil {
		return err
	}
	if asJSON {
		if err := json.NewEncoder(stdout).Encode(changes); err != nil {
			return err
		}
	} else {
		for _, change := range changes {
			fmt.Fprintf(stdout, "%-9s %s\n", change.Status, change.Path)
		}
	}
	if !apply {
		if !asJSON {
			fmt.Fprintln(stdout, "Dry run only. Re-run with --apply to create non-conflicting files.")
		}
		return nil
	}
	if err := scaffold.Apply(root, changes); err != nil {
		return err
	}
	if !asJSON {
		fmt.Fprintln(stdout, "Quality gate files created. Configure local commands before enabling required checks.")
	}
	return nil
}

func runProtect(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("protect", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var common commonFlags
	var repository string
	var apply bool
	var confirm string
	var showCurrent bool
	addCommon(flags, &common)
	flags.StringVar(&repository, "repository", "", "GitHub repository as owner/name; defaults to gh repo view")
	flags.BoolVar(&apply, "apply", false, "apply the displayed protection plan")
	flags.StringVar(&confirm, "confirm", "", "exact owner/name:branch confirmation")
	flags.BoolVar(&showCurrent, "show-current", false, "read current branch protection")
	if err := flags.Parse(args); err != nil {
		return err
	}
	root, err := repositoryRoot(common.repo)
	if err != nil {
		return err
	}
	cfg, err := loadConfig(root, common.configPath)
	if err != nil {
		return err
	}
	if repository == "" {
		repository, err = protection.ResolveRepository(root)
		if err != nil {
			return err
		}
	}
	if err := protection.ValidateRepository(repository); err != nil {
		return err
	}
	plan := protection.BuildPlan(repository, cfg)
	if showCurrent {
		return protection.Current(stdout, root, plan)
	}
	if common.json {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(plan); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(stdout, "Repository: %s\nBranch: %s\nRequired checks: %s\nApprovals: %d\nDismiss stale reviews: %t\nRequire last-push approval: %t\nResolve conversations: %t\nEnforce admins: %t\nForce pushes: false\nBranch deletion: false\n",
			plan.Repository, plan.Branch, strings.Join(cfg.RequiredChecks, ", "),
			cfg.BranchProtection.RequiredApprovals, cfg.BranchProtection.DismissStaleReviews,
			cfg.BranchProtection.RequireLastPushApproval,
			cfg.BranchProtection.RequireConversationResolution,
			cfg.BranchProtection.EnforceAdmins)
	}
	if !apply {
		if !common.json {
			fmt.Fprintln(stdout, "Dry run only. Review the plan before applying it.")
		}
		return nil
	}
	expected := plan.Repository + ":" + plan.Branch
	if confirm != expected {
		return fmt.Errorf("refusing to apply: --confirm must exactly equal %q", expected)
	}
	if len(cfg.RequiredChecks) == 0 {
		return errors.New("refusing to protect branch without configured required checks")
	}
	if err := protection.EnsureUnprotected(root, plan); err != nil {
		return err
	}
	return protection.Apply(root, plan)
}

func repositoryRoot(value string) (string, error) {
	root, err := filepath.Abs(value)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(root)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("repository path is not a directory: %s", root)
	}
	return root, nil
}

func loadConfig(root, value string) (config.Config, error) {
	path := value
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, filepath.FromSlash(value))
	}
	return config.Load(path)
}

func readPRBody(path string) (string, error) {
	if path == "" {
		return os.Getenv("PR_BODY"), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
