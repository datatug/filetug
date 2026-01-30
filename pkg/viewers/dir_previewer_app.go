package viewers

import "github.com/filetug/filetug/pkg/sneatv"

// DirPreviewerApp defines what methods of a tview.Application are used by a previewer.
// User navigator.NewApp with options or navigator.MockApp for tests
type DirPreviewerApp interface {
	sneatv.TabsApp
}
