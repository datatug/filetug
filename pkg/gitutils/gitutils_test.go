package gitutils

import (
	"context"
	"strings"
	"testing"
)

func TestDirGitStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status *RepoStatus
		want   string
	}{
		{
			name:   "nil",
			status: nil,
			want:   "",
		},
		{
			name:   "clean",
			status: &RepoStatus{Branch: "main"},
			want:   "[gray]🌿main±0[-]",
		},
		{
			name: "dirty",
			status: &RepoStatus{Branch: "feature", DirGitChangesStats: DirGitChangesStats{
				FilesChanged:  2,
				FileGitStatus: FileGitStatus{Insertions: 10, Deletions: 5},
			}},
			want: "[gray]🌿feature📄2[-][green]+10[-][red]-5[-]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("RepoStatus.String() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestGetGitStatus(t *testing.T) {
	status := GetRepositoryStatus(context.Background(), ".")
	if status != nil {
		s := status.String()
		if !strings.HasPrefix(s, "[gray]🌿") {
			t.Errorf("Expected status string starting with '[gray]🌿', got '%s'", s)
		}
	}
}
