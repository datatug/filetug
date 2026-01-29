package filetug

import "testing"

func Test_createLeft(t *testing.T) {
	nav := &Navigator{}
	nav.favorites = newFavoritesPanel(nav)
	createLeft(nav)
}
