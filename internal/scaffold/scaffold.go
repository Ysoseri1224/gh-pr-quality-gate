package scaffold

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed all:assets
var assets embed.FS

type Status string

const (
	Add       Status = "add"
	Unchanged Status = "unchanged"
	Conflict  Status = "conflict"
)

type Change struct {
	Path    string `json:"path"`
	Status  Status `json:"status"`
	Content []byte `json:"-"`
}

func Plan(root string) ([]Change, error) {
	var changes []Change
	err := fs.WalkDir(assets, "assets", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		content, readErr := assets.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		relative := strings.TrimSuffix(strings.TrimPrefix(filepath.ToSlash(path), "assets/"), ".tmpl")
		target := filepath.Join(root, filepath.FromSlash(relative))
		status := Add
		current, readCurrentErr := os.ReadFile(target)
		switch {
		case readCurrentErr == nil && string(current) == string(content):
			status = Unchanged
		case readCurrentErr == nil:
			status = Conflict
		case !errors.Is(readCurrentErr, os.ErrNotExist):
			return readCurrentErr
		}
		changes = append(changes, Change{Path: relative, Status: status, Content: content})
		return nil
	})
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	return changes, err
}

func Apply(root string, changes []Change) error {
	var conflicts []string
	for _, change := range changes {
		if change.Status == Conflict {
			conflicts = append(conflicts, change.Path)
		}
	}
	if len(conflicts) > 0 {
		return fmt.Errorf("refusing to overwrite existing files: %s", strings.Join(conflicts, ", "))
	}
	for _, change := range changes {
		if change.Status != Add {
			continue
		}
		target := filepath.Join(root, filepath.FromSlash(change.Path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		mode := os.FileMode(0o644)
		if strings.HasSuffix(target, ".sh") {
			mode = 0o755
		}
		if err := os.WriteFile(target, change.Content, mode); err != nil {
			return err
		}
	}
	return nil
}
