package httpfile

import (
	"os"
	"time"
)

type httpFileInfo struct {
	httpDirEntry
}

func (f httpFileInfo) Name() string       { return f.name }
func (f httpFileInfo) Size() int64        { return 0 }
func (f httpFileInfo) Mode() os.FileMode  { return f.Type() }
func (f httpFileInfo) ModTime() time.Time { return time.Time{} }
func (f httpFileInfo) IsDir() bool        { return f.isDir }
func (f httpFileInfo) Sys() any           { return nil }
