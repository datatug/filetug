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

	// Try to trigger Init() error by using an invalid charset if possible,
	// but tcell's simulation screen often ignores it or defaults to UTF-8.
	// However, we can at least test with an empty charset.
	s2 := NewSimScreen(t, "", width, height)
	if s2 == nil {
		t.Fatal("NewSimScreen with empty charset returned nil")
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

	// Test with multi-byte characters
	s.SetContent(0, 1, '世', nil, tcell.StyleDefault)
	s.SetContent(1, 1, ' ', nil, tcell.StyleDefault)
	line = ReadLine(s, 1, 2)
	expected = "世 "
	if line != expected {
		t.Errorf("expected line with multi-byte char to be %q, got %q", expected, line)
	}

	s.SetContent(0, 0, 'A', nil, tcell.StyleDefault)
	s.SetContent(1, 0, '\x00', nil, tcell.StyleDefault)
	line = ReadLine(s, 0, 2)
	expected = "A "
	if line != expected {
		t.Errorf("expected line with null char to be %q, got %q", expected, line)
	}

	// Test with "" (empty string)
	// Simulation screen might return "" if we don't set anything?
	s2 := NewSimScreen(t, "UTF-8", 1, 1)
	line = ReadLine(s2, 0, 1)
	expected = " "
	if line != expected {
		t.Errorf("expected line with empty content to be %q, got %q", expected, line)
	}
}
