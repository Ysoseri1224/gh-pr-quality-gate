package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlanApplyAndDetectConflict(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	changes, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("Plan() returned no templates")
	}
	if err := Apply(root, changes); err != nil {
		t.Fatal(err)
	}

	second, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range second {
		if change.Status != Unchanged {
			t.Fatalf("second plan contains %s for %s", change.Status, change.Path)
		}
	}

	target := filepath.Join(root, filepath.FromSlash(second[0].Path))
	if err := os.WriteFile(target, []byte("different"), 0o644); err != nil {
		t.Fatal(err)
	}
	conflicting, err := Plan(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := Apply(root, conflicting); err == nil {
		t.Fatal("Apply() overwrote a conflict")
	}
}
