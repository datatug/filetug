package httpfile

import (
	"os"
)

type httpDirEntry struct {
	name  string
	isDir bool
}

func (d httpDirEntry) Name() string { return d.name }
func (d httpDirEntry) IsDir() bool  { return d.isDir }
func (d httpDirEntry) Type() os.FileMode {
	if d.isDir {
		return os.ModeDir
	}
	return 0
}
func (d httpDirEntry) Info() (os.FileInfo, error) { return httpFileInfo{d}, nil }
