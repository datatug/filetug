package filetug

import "github.com/datatug/filetug/pkg/gitutils"

type DirInfo struct {
	Git *DirGitInfo
}

type DirGitInfo struct {
	Repo *gitutils.RepoStatus
}
