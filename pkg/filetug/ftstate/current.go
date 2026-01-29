package ftstate

import (
	"os"

	"github.com/filetug/filetug/pkg/files"
)

type Current struct {
	dir *files.DirContext
}

func (c *Current) Store() files.Store {
	if c == nil || c.dir == nil {
		return nil
	}
	return c.dir.Store()
}

func (c *Current) ChangeDir(path string) {
	// Do not check c for nil
	//if c == nil {
	//	return
	//}
	dir := c.NewDirContext(path, nil)
	c.SetDir(dir)
}

func (c *Current) SetDir(dir *files.DirContext) {
	// Do not check c for nil
	//if c == nil {
	//	return
	//}
	c.dir = dir
}

func (c *Current) Dir() *files.DirContext {
	// Do not check c for nil
	//if c == nil {
	//	return nil
	//}
	return c.dir
}

func (c *Current) NewDirContext(path string, children []os.DirEntry) *files.DirContext {
	if c == nil || c.dir == nil {
		return files.NewDirContext(nil, path, children)
	}
	return files.NewDirContext(c.dir.Store(), path, children)
}
