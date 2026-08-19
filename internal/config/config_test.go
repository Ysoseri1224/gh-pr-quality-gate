package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultShape(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := filepath.Join(root, "config.yml")
	data := `version: 1
local_gate:
  windows: ["go test ./..."]
  posix: ["go test ./..."]
required_checks: ["CI"]
required_files: ["AGENTS.md"]
issue_reference:
  required: true
  allow_partial_refs: true
agent_authority:
  merge: explicit-authorisation
  draft_transition: explicit-authorisation
  force_push: prohibited
branch_protection:
  branch: main
  required_approvals: 1
  dismiss_stale_reviews: true
  require_last_push_approval: true
  require_conversation_resolution: true
  enforce_admins: true
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.BranchProtection.Branch != "main" || len(cfg.LocalGate.Windows) != 1 {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}

func TestLoadRejectsUnknownFields(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte("version: 1\nunknown: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("Load() accepted an unknown field")
	}
}

func TestLoadRejectsMultipleDocuments(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte("version: 1\n---\nversion: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("Load() accepted multiple YAML documents")
	}
}

func TestValidateRejectsPolicyWeakeningAndTraversal(t *testing.T) {
	t.Parallel()
	cfg := Default()
	cfg.AgentAuthority.ForcePush = "allowed"
	cfg.RequiredFiles = []string{"../secret"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() accepted unsafe configuration")
	}
}
