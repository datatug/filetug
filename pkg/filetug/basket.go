package filetug

import "os"

type Basket struct {
	entries []os.DirEntry
}

func (b *Basket) AddToBasket(entry os.DirEntry) {
	b.entries = append(b.entries, entry)
}

func (b *Basket) Clear() {
	b.entries = []os.DirEntry{}
}
