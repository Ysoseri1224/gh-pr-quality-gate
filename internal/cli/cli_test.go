package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitIsDryRunByDefault(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"init", "--repo", root}, &stdout, &stderr, "test"); err != nil {
		t.Fatalf("Run(init) error = %v; stderr=%s", err, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatal("dry-run init created files")
	}
	if !strings.Contains(stdout.String(), "Dry run only") {
		t.Fatalf("dry-run message missing: %s", stdout.String())
	}
}

func TestInitJSONContainsNoPlainTextSuffix(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"init", "--repo", root, "--json"}, &stdout, &stderr, "test"); err != nil {
		t.Fatalf("Run(init --json) error = %v; stderr=%s", err, stderr.String())
	}
	var changes []map[string]interface{}
	decoder := json.NewDecoder(&stdout)
	if err := decoder.Decode(&changes); err != nil {
		t.Fatalf("invalid JSON output: %v; output=%s", err, stdout.String())
	}
	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		t.Fatalf("unexpected JSON suffix: %v; output=%s", err, stdout.String())
	}
}

func TestInitApplyCreatesFiles(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"init", "--repo", root, "--apply"}, &stdout, &stderr, "test"); err != nil {
		t.Fatalf("Run(init --apply) error = %v; stderr=%s", err, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); err != nil {
		t.Fatalf("AGENTS.md not created: %v", err)
	}
}
