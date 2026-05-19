package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetRoot_InGitRepo(t *testing.T) {
	repoDir := setupTestRepo(t)

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("failed to change to repo dir: %v", err)
	}
	defer os.Chdir(origDir)

	root, err := GetRoot()
	if err != nil {
		t.Fatalf("GetRoot failed: %v", err)
	}

	// Resolve symlinks for comparison (macOS /var -> /private/var)
	expectedRoot, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("failed to resolve symlinks: %v", err)
	}
	actualRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("failed to resolve symlinks: %v", err)
	}

	if actualRoot != expectedRoot {
		t.Errorf("expected root %s, got %s", expectedRoot, actualRoot)
	}
}

func TestGetRoot_NotInGitRepo(t *testing.T) {
	tmpDir := t.TempDir() // Not a git repo

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change dir: %v", err)
	}
	defer os.Chdir(origDir)

	_, err = GetRoot()
	if err == nil {
		t.Error("expected error when not in git repo")
	}
}
