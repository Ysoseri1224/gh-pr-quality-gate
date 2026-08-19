package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/config"
)

func TestAuditIssueReferences(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("rules"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Default()
	cfg.RequiredFiles = []string{"AGENTS.md"}
	cfg.RequiredChecks = nil

	tests := []struct {
		name     string
		body     string
		wantFail bool
	}{
		{name: "closing", body: "Closes #42"},
		{name: "partial", body: "Refs owner/repo#42"},
		{name: "missing", body: "No issue here", wantFail: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Audit(root, cfg, test.body)
			if got := result.HasFailures(); got != test.wantFail {
				t.Fatalf("HasFailures() = %t, want %t; findings: %#v", got, test.wantFail, result.Findings)
			}
		})
	}
}

func TestAuditMissingRequiredFile(t *testing.T) {
	t.Parallel()
	cfg := config.Default()
	cfg.RequiredFiles = []string{"missing.md"}
	if !Audit(t.TempDir(), cfg, "Closes #1").HasFailures() {
		t.Fatal("Audit() did not fail for a missing required file")
	}
}

func TestAuditMatchesRequiredCheckToWorkflowJobName(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	workflowDir := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	workflow := "name: CI\non: [push]\njobs:\n  test:\n    name: Required CI\n    runs-on: ubuntu-latest\n    steps: []\n"
	if err := os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte(workflow), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Default()
	cfg.RequiredFiles = nil
	cfg.RequiredChecks = []string{"Required CI"}
	if Audit(root, cfg, "Closes #1").HasFailures() {
		t.Fatal("Audit() did not recognize the workflow job name")
	}
	cfg.RequiredChecks = []string{"Missing CI"}
	if !Audit(root, cfg, "Closes #1").HasFailures() {
		t.Fatal("Audit() accepted a missing workflow job name")
	}
}
