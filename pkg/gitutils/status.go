package gitutils

import (
	"fmt"
	"strings"
)

type FileGitStatus struct {
	Insertions int
	Deletions  int
}

func (s *FileGitStatus) String() string {
	var sb strings.Builder
	if s.Insertions > 0 {
		_, _ = fmt.Fprintf(&sb, "[green]+%d[-]", s.Insertions)
	}
	if s.Deletions > 0 {
		_, _ = fmt.Fprintf(&sb, "[red]-%d[-]", s.Deletions)
	}
	if sb.Len() == 0 {
		return "[lightgray]±0[-]"
	}
	return sb.String()
}

type DirGitChangesStats struct {
	FilesChanged int
	FileGitStatus
}

type RepoStatus struct {
	Branch string
	DirGitChangesStats
}

func (s *RepoStatus) String() string {
	const separator = "[gray]┆[-]"
	if s == nil {
		return ""
	}
	var noChanges DirGitChangesStats
	if s.DirGitChangesStats == noChanges {
		return separator + fmt.Sprintf("[darkgray]%s[-]%s", s.Branch, s.FileGitStatus.String())
	}
	return separator + fmt.Sprintf("[darkgray]%s[-]%s[darkgray]ƒ%d[-]%s", s.Branch, separator, s.FilesChanged, s.FileGitStatus.String())
}
