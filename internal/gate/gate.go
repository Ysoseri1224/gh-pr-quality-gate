package gate

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/config"
	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/report"
	"gopkg.in/yaml.v3"
)

var (
	closingReference = regexp.MustCompile(`(?i)\b(close[sd]?|fix(e[sd])?|resolve[sd]?)\s+(?:[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+)?#\d+\b`)
	partialReference = regexp.MustCompile(`(?i)(\b(refs?|references?|relates\s+to)\s+(?:[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+)?#\d+\b|(?:^|\s)#[0-9]+\b)`)
)

func Audit(root string, cfg config.Config, prBody string) report.Report {
	result := report.Report{Repository: root}
	for _, name := range cfg.RequiredFiles {
		path := filepath.Join(root, filepath.FromSlash(name))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			result.Findings = append(result.Findings, report.Finding{Level: report.Pass, Check: "required-file", Message: name})
		} else {
			result.Findings = append(result.Findings, report.Finding{Level: report.Fail, Check: "required-file", Message: name + " is missing"})
		}
	}

	if cfg.IssueReference.Required {
		switch {
		case prBody == "":
			result.Findings = append(result.Findings, report.Finding{Level: report.Warn, Check: "issue-reference", Message: "PR body was not supplied; pass --pr-body-file or set PR_BODY"})
		case closingReference.MatchString(prBody):
			result.Findings = append(result.Findings, report.Finding{Level: report.Pass, Check: "issue-reference", Message: "closing issue reference found"})
		case cfg.IssueReference.AllowPartialRefs && partialReference.MatchString(prBody):
			result.Findings = append(result.Findings, report.Finding{Level: report.Pass, Check: "issue-reference", Message: "non-closing issue reference found"})
		default:
			result.Findings = append(result.Findings, report.Finding{Level: report.Fail, Check: "issue-reference", Message: "PR body must reference an issue"})
		}
	}

	if len(cfg.RequiredChecks) == 0 {
		result.Findings = append(result.Findings, report.Finding{Level: report.Warn, Check: "required-checks", Message: "no required check names are configured"})
	} else {
		jobs, findings := workflowJobNames(root)
		result.Findings = append(result.Findings, findings...)
		for _, required := range cfg.RequiredChecks {
			if jobs[required] {
				result.Findings = append(result.Findings, report.Finding{Level: report.Pass, Check: "required-check", Message: required})
			} else {
				result.Findings = append(result.Findings, report.Finding{Level: report.Fail, Check: "required-check", Message: required + " is not a workflow job name"})
			}
		}
	}

	result.Findings = append(result.Findings,
		report.Finding{Level: report.Pass, Check: "merge-authority", Message: "explicit authorisation required"},
		report.Finding{Level: report.Pass, Check: "draft-authority", Message: "explicit authorisation required"},
		report.Finding{Level: report.Pass, Check: "force-push", Message: "prohibited"},
	)
	return result
}

func workflowJobNames(root string) (map[string]bool, []report.Finding) {
	jobs := make(map[string]bool)
	directory := filepath.Join(root, ".github", "workflows")
	entries, err := os.ReadDir(directory)
	if err != nil {
		return jobs, []report.Finding{{Level: report.Fail, Check: "workflows", Message: ".github/workflows cannot be read"}}
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		extension := strings.ToLower(filepath.Ext(entry.Name()))
		if extension != ".yml" && extension != ".yaml" {
			continue
		}
		path := filepath.Join(directory, entry.Name())
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return jobs, []report.Finding{{Level: report.Fail, Check: "workflow-parse", Message: entry.Name() + " cannot be read"}}
		}
		var workflow struct {
			Jobs map[string]struct {
				Name string `yaml:"name"`
			} `yaml:"jobs"`
		}
		if parseErr := yaml.Unmarshal(data, &workflow); parseErr != nil {
			return jobs, []report.Finding{{Level: report.Fail, Check: "workflow-parse", Message: entry.Name() + ": " + parseErr.Error()}}
		}
		for key, job := range workflow.Jobs {
			name := strings.TrimSpace(job.Name)
			if name == "" {
				name = key
			}
			jobs[name] = true
		}
	}
	return jobs, nil
}

func RunLocal(root string, cfg config.Config, stdout, stderr io.Writer) error {
	commands := cfg.LocalGate.Posix
	if runtime.GOOS == "windows" {
		commands = cfg.LocalGate.Windows
	}
	if len(commands) == 0 {
		return fmt.Errorf("no local gate commands configured for %s", runtime.GOOS)
	}
	for _, command := range commands {
		fmt.Fprintf(stdout, "> %s\n", command)
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
		cmd.Dir = root
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("local gate failed: %s: %w", command, err)
		}
	}
	return nil
}
