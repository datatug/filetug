package filetug

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestNewFavorites(t *testing.T) {
	f := newFavorites()
	if f == nil {
		t.Fatal("f is nil")
	}
	if f.GetRoot() == nil {
		t.Fatal("root is nil")
	}
	assert.Equal(t, 2, len(f.GetRoot().GetChildren()))
}
