package filetug

import (
	"github.com/rivo/tview"
)

type Tree struct {
	*tview.TreeView
	currDirRoot *tview.TreeNode

	selectedDirNode *tview.TreeNode
}

func (t *Tree) GetBox() *tview.Box {
	return t.Box
}

func NewTree(nav *Navigator) *Tree {
	t := &Tree{
		TreeView: tview.NewTreeView(),
	}

	t.currDirRoot = tview.NewTreeNode("~")
	t.SetRoot(t.currDirRoot)
	t.SetChangedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if dir, ok := ref.(string); ok {
			nav.showDir(dir, node)
		}
	})

	return t
}
