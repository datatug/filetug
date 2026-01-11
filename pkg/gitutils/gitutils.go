package gitutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	// gitStatusSemaphore limits concurrent git status calls to avoid system hang
	gitStatusSemaphore = make(chan struct{}, 2)
)

type DirGitStatus struct {
	Branch       string
	FilesChanged int
	Insertions   int
	Deletions    int
}

func (s *DirGitStatus) String() string {
	if s == nil {
		return ""
	}
	if s.FilesChanged == 0 && s.Insertions == 0 && s.Deletions == 0 {
		return fmt.Sprintf("[gray]🌿%s±0[-]", s.Branch)
	}
	return fmt.Sprintf("[gray]🌿%s📄%d[-][green]+%d[-][red]-%d[-]", s.Branch, s.FilesChanged, s.Insertions, s.Deletions)
}

// GetGitStatus returns a brief git status for the given directory.
// It uses a context to allow cancellation and a semaphore to limit concurrency.
func GetGitStatus(ctx context.Context, dir string) *DirGitStatus {
	// Quick check if .git exists to avoid expensive go-git calls for non-git dirs
	dotGit := filepath.Join(dir, ".git")
	if _, err := os.Stat(dotGit); os.IsNotExist(err) {
		return nil
	}

	select {
	case <-ctx.Done():
		return nil
	case gitStatusSemaphore <- struct{}{}:
		defer func() { <-gitStatusSemaphore }()
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil
	}

	res := &DirGitStatus{}

	head, err := repo.Head()
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			res.Branch = "master"
		} else {
			return nil
		}
	} else {
		if head.Name().IsBranch() {
			res.Branch = head.Name().Short()
		} else {
			res.Branch = head.Hash().String()[:7]
		}
	}

	select {
	case <-ctx.Done():
		return res
	default:
	}

	headCommit, _ := repo.CommitObject(head.Hash())

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

	res.FilesChanged = len(status)

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
						go func() {
							_ = f.Close()
						}()
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
