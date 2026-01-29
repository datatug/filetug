package ftfav

import (
	"net/url"
	"path/filepath"
)

type Favorite struct {
	Store       url.URL
	Path        string
	Shortcut    rune
	Description string
}

func (f Favorite) Key() string {
	key := f.Store
	key.Path = filepath.Join(key.Path, f.Path)
	return key.String()
}
