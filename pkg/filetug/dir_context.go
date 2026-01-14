package filetug

import "os"

type DirItem struct {
	Name  string
	IsDir bool
	Size  int64
}

type DirContext struct {
	Path     string
	children []os.DirEntry
}

func newDirContext(path string, children []os.DirEntry) *DirContext {
	return &DirContext{
		Path:     path,
		children: children,
	}
}
