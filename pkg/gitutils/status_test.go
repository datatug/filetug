package gitutils

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
)

func TestFileGitStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status FileGitStatus
		want   string
	}{
		{"zero", FileGitStatus{0, 0}, "[lightgray]±0[-]"},
		{"insertions only", FileGitStatus{5, 0}, "[green]+5[-]"},
		{"deletions only", FileGitStatus{0, 3}, "[red]-3[-]"},
		{"both", FileGitStatus{5, 3}, "[green]+5[-][red]-3[-]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("FileGitStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileStatus_UntrackedFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitutils-file-status-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	_, err = git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	filePath := filepath.Join(tempDir, "file.txt")
	err = os.WriteFile(filePath, []byte("line1\nline2\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repo, err := git.PlainOpen(tempDir)
	if err != nil {
		t.Fatalf("Failed to open git repo: %v", err)
	}

	ctx := context.Background()
	status := GetFileStatus(ctx, repo, filePath)
	if status == nil {
		t.Fatal("Expected non-nil status")
	}
	if status.FilesChanged != 1 {
		t.Fatalf("Expected FilesChanged=1, got %d", status.FilesChanged)
	}
	if status.Branch == "" {
		t.Fatal("Expected non-empty branch")
	}
}
