package filetug

func createLeft(nav *Navigator) {
	nav.left = newContainer(0, nav)
	nav.left.SetContent(nav.dirsTree)

	nav.favorites.SetFocusFunc(func() {
		nav.activeCol = 0
	})
	nav.favoritesFocusFunc = func() {
		nav.activeCol = 0
	}
}
