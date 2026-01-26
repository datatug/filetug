package filetug

func createLeft(nav *Navigator) {
	nav.left = NewContainer(0, nav)
	nav.left.SetContent(nav.dirsTree)

	nav.favorites.flex.SetFocusFunc(func() {
		nav.activeCol = 0
	})
	nav.favoritesFocusFunc = func() {
		nav.activeCol = 0
	}
}
