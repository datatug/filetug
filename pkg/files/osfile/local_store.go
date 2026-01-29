package osfile

import "github.com/filetug/filetug/pkg/files"

var localFileStore = NewStore("/")

func NewLocalDir(fullPath string) *files.DirContext {
	return files.NewDirContext(localFileStore, fullPath, nil)
}
