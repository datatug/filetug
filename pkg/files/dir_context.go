package files

import (
	"os"
	"path"
	"strings"
	"time"
)

var _ EntryWithDirPath = (*DirContext)(nil)

type DirContext struct {
	store    Store
	path     string
	children []os.DirEntry
	time     time.Time
}

func NewDirContext(store Store, fullPath string, children []os.DirEntry) *DirContext {
	return &DirContext{
		store:    store,
		path:     fullPath,
		children: children,
	}
}

func (d *DirContext) Timestamp() time.Time {
	return d.time
}

func (d *DirContext) SetChildren(entries []os.DirEntry) {
	d.time = time.Now()
	d.children = entries
}

func (d *DirContext) Entries() []EntryWithDirPath {
	entries := make([]EntryWithDirPath, len(d.children))
	for i, child := range d.children {
		entries[i] = NewEntryWithDirPath(child, d.path)
	}
	return entries
}

func (d *DirContext) Children() []os.DirEntry {
	if d.children == nil {
		return nil
	}
	children := make([]os.DirEntry, len(d.children))
	copy(children, d.children)
	return children
}

func (d *DirContext) DirPath() string {
	if d.path == "" {
		return ""
	}
	return path.Dir(d.path)
}

func (d *DirContext) FullName() string {
	return d.path
}

func (d *DirContext) String() string {
	return d.path
}

func (d *DirContext) Name() string {
	if d.path == "" {
		return ""
	}
	if d.path == "/" {
		return "/"
	}
	trimmed := strings.TrimSuffix(d.path, "/")
	return path.Base(trimmed)
}

func (d *DirContext) IsDir() bool {
	return true
}

func (d *DirContext) Type() os.FileMode {
	return os.ModeDir
}

func (d *DirContext) Info() (os.FileInfo, error) {
	if d.path == "" {
		return nil, nil
	}
	if d.store != nil && d.store.RootURL().Scheme == "file" {
		return os.Stat(d.path)
	}
	return nil, nil
}

func (d *DirContext) Store() Store {
	return d.store
}

func (d *DirContext) Path() string {
	return d.path
}
