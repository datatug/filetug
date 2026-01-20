package gitutils

import "testing"

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
