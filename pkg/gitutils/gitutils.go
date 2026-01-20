package gitutils

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	// gitStatusSemaphore limits concurrent git status calls to avoid system hang
	gitStatusSemaphore = make(chan struct{}, 2)
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

// GetDirStatus returns a brief git status for the given directory.
// It uses a context to allow cancellation and a semaphore to limit concurrency.
func GetDirStatus(ctx context.Context, dir string) *RepoStatus {
	repoRoot := GetRepositoryRoot(dir)
	if repoRoot == "" {
		return nil
	}

	select {
	case <-ctx.Done():
		return nil
	case gitStatusSemaphore <- struct{}{}:
		defer func() { <-gitStatusSemaphore }()
	}

	repo, err := git.PlainOpen(repoRoot)
	if err != nil {
		return nil
	}

	res := &RepoStatus{}

	head, err := repo.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) || err.Error() == "reference not found" {
			res.Branch = "master"
		} else {
			// This covers some other error during repo.Head()
			res.Branch = "unknown"
		}
	} else if head == nil || head.Hash().IsZero() {
		res.Branch = "unknown"
	} else {
		if head.Name().IsBranch() {
			res.Branch = head.Name().Short()
		} else {
			hashStr := head.Hash().String()
			if len(hashStr) >= 7 {
				res.Branch = hashStr[:7]
			} else {
				res.Branch = hashStr
			}
		}
	}

	select {
	case <-ctx.Done():
		return res
	default:
	}

	var headHash plumbing.Hash
	if head != nil {
		headHash = head.Hash()
	}
	headCommit, _ := repo.CommitObject(headHash)

	worktree, err := repo.Worktree()
	if err != nil {
		return res
	}

	status, err := worktree.Status()
	if err != nil {
		return res
	}

	if status.IsClean() {
		return res
	}

	relPath, err := filepath.Rel(repoRoot, dir)
	if err != nil {
		relPath = ""
	}
	if relPath == "." {
		relPath = ""
	}
	if relPath != "" && !strings.HasSuffix(relPath, string(filepath.Separator)) {
		relPath += string(filepath.Separator)
	}

	res.FilesChanged = 0
	for fileName, s := range status {
		if relPath != "" && !strings.HasPrefix(fileName, relPath) {
			continue
		}
		if s.Worktree != git.Unmodified || s.Staging != git.Unmodified {
			res.FilesChanged++
		}
	}

	if res.FilesChanged == 0 {
		return res
	}

	// To get insertions/deletions, we need to diff
	if headCommit != nil {
		headTree, err := headCommit.Tree()
		if err == nil {
			for fileName, fileStatus := range status {
				select {
				case <-ctx.Done():
					return res
				default:
				}

				if relPath != "" && !strings.HasPrefix(fileName, relPath) {
					continue
				}

				if fileStatus.Worktree == git.Unmodified && fileStatus.Staging == git.Unmodified {
					continue
				}

				// If file is untracked, we can count its lines as insertions
				if fileStatus.Worktree == git.Untracked {
					if f, err := worktree.Filesystem.Open(fileName); err == nil {
						// Limit reading to avoid performance issues with large files
						const maxRead = 1 * 1024 * 1024 // 1MB
						b := make([]byte, maxRead)
						n, _ := f.Read(b)
						res.Insertions += strings.Count(string(b[:n]), "\n")
						_ = f.Close()
					}
					continue
				}

				// If file is deleted, we can count its lines in head as deletions
				if fileStatus.Worktree == git.Deleted || fileStatus.Staging == git.Deleted {
					if f, err := headTree.File(fileName); err == nil {
						if content, err := f.Contents(); err == nil {
							res.Deletions += strings.Count(content, "\n")
						}
					}
					continue
				}
			}
		}
	}

	return res
}
