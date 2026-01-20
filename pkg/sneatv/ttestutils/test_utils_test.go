package ttestutils

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewSimScreen(t *testing.T) {
	width, height := 80, 24
	s := NewSimScreen(t, "UTF-8", width, height)
	if s == nil {
		t.Fatal("NewSimScreen returned nil")
	}
	w, h := s.Size()
	if w != width || h != height {
		t.Errorf("expected size %dx%d, got %dx%d", width, height, w, h)
	}
}

func TestReadLine(t *testing.T) {
	width, height := 10, 2
	s := NewSimScreen(t, "UTF-8", width, height)

	// Test empty line (filled with spaces by ReadLine)
	line := ReadLine(s, 0, width)
	expected := "          " // 10 spaces
	if line != expected {
		t.Errorf("expected empty line to be %q, got %q", expected, line)
	}

	// Test line with some characters
	s.SetContent(0, 1, 'H', nil, tcell.StyleDefault)
	s.SetContent(1, 1, 'e', nil, tcell.StyleDefault)
	s.SetContent(2, 1, 'l', nil, tcell.StyleDefault)
	s.SetContent(3, 1, 'l', nil, tcell.StyleDefault)
	s.SetContent(4, 1, 'o', nil, tcell.StyleDefault)

	line = ReadLine(s, 1, width)
	expected = "Hello     "
	if line != expected {
		t.Errorf("expected line to be %q, got %q", expected, line)
	}

	// Test with mixed content (though simulation screen might not easily return "")
	// but we've covered the common cases.
}
