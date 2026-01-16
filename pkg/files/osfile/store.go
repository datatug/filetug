package osfile

import (
	"os"
)

var osReadDir = os.ReadDir

type Store struct {
	root string
}

func (o Store) ReadDir(name string) ([]os.DirEntry, error) {
	return osReadDir(name)
}

func NewStore(root string) *Store {
	if root == "" {
		panic("root is empty")
	}
	return &Store{root: root}
}
