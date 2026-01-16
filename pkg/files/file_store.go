package files

import "os"

type Store interface {
	ReadDir(name string) ([]os.DirEntry, error)
}
