package filetug

import "github.com/filetug/filetug/pkg/viewers"

// PreviewerApp defines what methods of a tview.Application are used by a previewer.
// User navigator.NewApp with options or navigator.MockApp for tests
type PreviewerApp interface {
	viewers.DirPreviewerApp
}
