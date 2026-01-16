package files

import (
	"context"
	"net/url"
	"os"
)

type Store interface {
	RootTitle() string
	RootURL() url.URL
	ReadDir(ctx context.Context, name string) ([]os.DirEntry, error)
}
