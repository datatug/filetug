package osfile

import (
	"os"

	"github.com/filetug/filetug/pkg/files"
)

var _ files.DirReader = (*dirReader)(nil)

type dirReader struct {
	file *os.File
}

func newDirReader(dirPath string) (files.DirReader, error) {
	file, err := osOpen(dirPath)
	if err != nil {
		return nil, err
	}
	return &dirReader{file: file}, nil
}

func (d dirReader) Close() error {
	return d.file.Close()
}

func (d dirReader) Readdir() ([]os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}
