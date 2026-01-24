package files

import (
	"context"
	"errors"
	"net/url"
	"os"
)

// type DirEntry = os.DirEntry
type Store interface {
	RootTitle() string
	RootURL() url.URL
	ReadDir(ctx context.Context, path string) ([]os.DirEntry, error)
	Delete(ctx context.Context, path string) error // TODO(unsure): should it be Remove to match os.Remove?
	CreateDir(ctx context.Context, path string) error
	CreateFile(ctx context.Context, path string) error
}

var ErrNotSupportedOperation = errors.New("no supported operation")
var ErrNotImplemented = errors.New("no implemented")
