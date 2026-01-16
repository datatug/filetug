package files

import "os"

type Store interface {
	RootTitle() string
	ReadDir(name string) ([]os.DirEntry, error)
}
