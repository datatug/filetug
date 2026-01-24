package filetug

import "os"

type basket struct {
	entries []os.DirEntry
}

func (b *basket) AddToBasket(entry os.DirEntry) {
	b.entries = append(b.entries, entry)
}

func (b *basket) Clear() {
	b.entries = []os.DirEntry{}
}
