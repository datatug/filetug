package filetug

import "testing"

func TestNewTree(t *testing.T) {
	tree := NewTree()
	if tree == nil {
		t.Fatal("tree is nil")
	}
	if tree.GetRoot() == nil {
		t.Fatal("root is nil")
	}
	if tree.GetBox() == nil {
		t.Fatal("box is nil")
	}
}
