package gitutils

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestStaging(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitutils-staging-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// 1. Test CanBeStaged for non-git directory
	fileInNonGit := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(fileInNonGit, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	can, err := CanBeStaged(fileInNonGit)
	if err != nil {
		t.Errorf("CanBeStaged failed for non-git: %v", err)
	}
	if can {
		t.Error("CanBeStaged returned true for non-git directory")
	}

	// 2. Initialize git repo
	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// 3. Test CanBeStaged for untracked file
	can, err = CanBeStaged(fileInNonGit)
	if err != nil {
		t.Errorf("CanBeStaged failed for untracked file: %v", err)
	}
	if !can {
		t.Error("Expected CanBeStaged to be true for untracked file")
	}

	// 4. Test StageFile
	err = StageFile(fileInNonGit)
	if err != nil {
		t.Errorf("StageFile failed: %v", err)
	}

	// 5. Verify it's staged
	worktree, _ := repo.Worktree()
	status, _ := worktree.Status()
	fileStatus := status.File("file.txt")
	if fileStatus.Staging == git.Unmodified {
		t.Error("File should be staged")
	}

	// 6. Test CanBeStaged for staged but unmodified in worktree
	// It should still return true because it has changes relative to HEAD (though here HEAD is empty)
	// Actually, go-git status for a new staged file shows Staging=Added, Worktree=Unmodified.
	can, err = CanBeStaged(fileInNonGit)
	if err != nil {
		t.Errorf("CanBeStaged failed for staged file: %v", err)
	}
	if !can {
		t.Error("Expected CanBeStaged to be true for staged file (Added)")
	}

	// 7. Commit and check
	_, err = worktree.Commit("initial", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com", When: time.Now()},
	})
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	can, err = CanBeStaged(fileInNonGit)
	if err != nil {
		t.Errorf("CanBeStaged failed for clean file: %v", err)
	}
	if can {
		t.Error("Expected CanBeStaged to be false for clean file")
	}

	// 8. Modify and check
	if err := os.WriteFile(fileInNonGit, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	can, err = CanBeStaged(fileInNonGit)
	if err != nil {
		t.Errorf("CanBeStaged failed for modified file: %v", err)
	}
	if !can {
		t.Error("Expected CanBeStaged to be true for modified file")
	}

	// 9. Test UnstageFile
	// First stage the modified file
	if err := StageFile(fileInNonGit); err != nil {
		t.Fatalf("Failed to stage file for unstage test: %v", err)
	}

	// Verify it is staged
	status, _ = worktree.Status()
	fileStatus = status.File("file.txt")
	if fileStatus.Staging != git.Modified {
		t.Errorf("Expected status Staging=Modified, got %v", fileStatus.Staging)
	}

	// Now unstage it
	if err := UnstageFile(fileInNonGit); err != nil {
		t.Fatalf("UnstageFile failed: %v", err)
	}

	// Verify it is unstaged
	status, _ = worktree.Status()
	fileStatus = status.File("file.txt")
	if fileStatus.Staging != git.Unmodified {
		t.Errorf("Expected status Staging=Unmodified after unstage, got %v", fileStatus.Staging)
	}
	if fileStatus.Worktree != git.Modified {
		t.Errorf("Expected status Worktree=Modified after unstage, got %v", fileStatus.Worktree)
	}
}

func TestStageDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitutils-stagedir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	worktree, _ := repo.Worktree()

	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	nestedDir := filepath.Join(subDir, "nested")
	if err := os.Mkdir(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	file1 := filepath.Join(subDir, "file1.txt")
	if err := os.WriteFile(file1, []byte("file1"), 0644); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}

	file2 := filepath.Join(nestedDir, "file2.txt")
	if err := os.WriteFile(file2, []byte("file2"), 0644); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	t.Run("non-recursive", func(t *testing.T) {
		if err := StageDir(subDir, false); err != nil {
			t.Fatalf("StageDir non-recursive failed: %v", err)
		}

		status, _ := worktree.Status()
		if status.File("subdir/file1.txt").Staging == git.Unmodified {
			t.Error("subdir/file1.txt should be staged")
		}
		if status.File("subdir/nested/file2.txt").Staging != git.Untracked {
			t.Errorf("subdir/nested/file2.txt should NOT be staged, got %v", status.File("subdir/nested/file2.txt").Staging)
		}

		// Reset for next test
		_ = worktree.Reset(&git.ResetOptions{Mode: git.HardReset})
		// Clean untracked files if any (Reset Hard doesn't remove untracked)
		_ = os.WriteFile(file1, []byte("file1"), 0644)
		_ = os.WriteFile(file2, []byte("file2"), 0644)
	})

	t.Run("recursive", func(t *testing.T) {
		if err := StageDir(subDir, true); err != nil {
			t.Fatalf("StageDir recursive failed: %v", err)
		}

		status, _ := worktree.Status()
		if status.File("subdir/file1.txt").Staging == git.Unmodified {
			t.Error("subdir/file1.txt should be staged")
		}
		if status.File("subdir/nested/file2.txt").Staging == git.Unmodified {
			t.Error("subdir/nested/file2.txt should be staged")
		}
	})
}
