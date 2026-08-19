package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const FileName = ".github/pr-quality-gate.yml"

type Config struct {
	Version          int              `yaml:"version"`
	LocalGate        LocalGate        `yaml:"local_gate"`
	RequiredChecks   []string         `yaml:"required_checks"`
	RequiredFiles    []string         `yaml:"required_files"`
	IssueReference   IssueReference   `yaml:"issue_reference"`
	AgentAuthority   AgentAuthority   `yaml:"agent_authority"`
	BranchProtection BranchProtection `yaml:"branch_protection"`
}

type LocalGate struct {
	Windows []string `yaml:"windows"`
	Posix   []string `yaml:"posix"`
}

type IssueReference struct {
	Required         bool `yaml:"required"`
	AllowPartialRefs bool `yaml:"allow_partial_refs"`
}

type AgentAuthority struct {
	Merge           string `yaml:"merge"`
	DraftTransition string `yaml:"draft_transition"`
	ForcePush       string `yaml:"force_push"`
}

type BranchProtection struct {
	Branch                        string `yaml:"branch"`
	RequiredApprovals             int    `yaml:"required_approvals"`
	DismissStaleReviews           bool   `yaml:"dismiss_stale_reviews"`
	RequireLastPushApproval       bool   `yaml:"require_last_push_approval"`
	RequireConversationResolution bool   `yaml:"require_conversation_resolution"`
	EnforceAdmins                 bool   `yaml:"enforce_admins"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", path, err)
	}
	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return Config{}, fmt.Errorf("parse %s: multiple YAML documents are not allowed", path)
		}
		return Config{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate %s: %w", path, err)
	}
	return cfg, nil
}

func Default() Config {
	return Config{
		Version:        1,
		RequiredChecks: []string{"Repository quality gate", "PR policy"},
		RequiredFiles: []string{
			"AGENTS.md",
			"docs/repo_rule.md",
			".github/pull_request_template.md",
			".github/workflows/pr-policy.yml",
			".github/workflows/repository-quality-gate.yml",
		},
		IssueReference: IssueReference{Required: true, AllowPartialRefs: true},
		AgentAuthority: AgentAuthority{
			Merge:           "explicit-authorisation",
			DraftTransition: "explicit-authorisation",
			ForcePush:       "prohibited",
		},
		BranchProtection: BranchProtection{
			Branch:                        "main",
			RequiredApprovals:             1,
			DismissStaleReviews:           true,
			RequireLastPushApproval:       true,
			RequireConversationResolution: true,
			EnforceAdmins:                 true,
		},
	}
}

func (c Config) Validate() error {
	var errs []error
	if c.Version != 1 {
		errs = append(errs, fmt.Errorf("version must be 1, got %d", c.Version))
	}
	if c.AgentAuthority.Merge != "explicit-authorisation" {
		errs = append(errs, errors.New("agent_authority.merge must be explicit-authorisation"))
	}
	if c.AgentAuthority.DraftTransition != "explicit-authorisation" {
		errs = append(errs, errors.New("agent_authority.draft_transition must be explicit-authorisation"))
	}
	if c.AgentAuthority.ForcePush != "prohibited" {
		errs = append(errs, errors.New("agent_authority.force_push must be prohibited"))
	}
	if c.BranchProtection.Branch == "" {
		errs = append(errs, errors.New("branch_protection.branch is required"))
	}
	if c.BranchProtection.RequiredApprovals < 0 || c.BranchProtection.RequiredApprovals > 6 {
		errs = append(errs, errors.New("branch_protection.required_approvals must be between 0 and 6"))
	}
	seenChecks := make(map[string]bool)
	for _, check := range c.RequiredChecks {
		if strings.TrimSpace(check) == "" {
			errs = append(errs, errors.New("required_checks cannot contain an empty name"))
		}
		if seenChecks[check] {
			errs = append(errs, fmt.Errorf("required_checks contains duplicate %q", check))
		}
		seenChecks[check] = true
	}
	for _, name := range c.RequiredFiles {
		clean := filepath.Clean(filepath.FromSlash(name))
		if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			errs = append(errs, fmt.Errorf("required_files must stay inside the repository: %q", name))
		}
	}
	return errors.Join(errs...)
}
