package viewers

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/filetug/filetug/pkg/files"
	"github.com/strongo/dsstore"
)

type errorLexer struct {
	config *chroma.Config
}

func (l *errorLexer) Config() *chroma.Config {
	return l.config
}

func (l *errorLexer) Tokenise(options *chroma.TokeniseOptions, text string) (chroma.Iterator, error) {
	_, _ = options, text
	return nil, errors.New("tokenise failure")
}

func (l *errorLexer) SetRegistry(registry *chroma.LexerRegistry) chroma.Lexer {
	_ = registry
	return l
}

func (l *errorLexer) SetAnalyser(analyser func(text string) float32) chroma.Lexer {
	_ = analyser
	return l
}

func (l *errorLexer) AnalyseText(text string) float32 {
	_ = text
	return 1
}

func waitForUpdate(t *testing.T, done <-chan struct{}) {
	select {
	case <-done:
		return
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for preview update")
	}
}

func waitForText(t *testing.T, previewer *TextPreviewer, needle string) {
	deadline := time.After(2 * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			t.Fatalf("timeout waiting for text %q", needle)
		case <-ticker.C:
			text := previewer.GetText(false)
			if strings.Contains(text, needle) {
				return
			}
		}
	}
}

func TestTextPreviewerPreviewPlainText(t *testing.T) {
	previewer := NewTextPreviewer()
	data := []byte("plain text")
	dir := filepath.Dir("note.unknownext")
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "note.unknownext"},
		Dir:      dir,
	}

	done := make(chan struct{})
	queueUpdateDraw := func(fn func()) {
		fn()
		close(done)
	}

	previewer.Preview(entry, data, queueUpdateDraw)
	waitForUpdate(t, done)

	text := previewer.GetText(false)
	expected := string(data)
	assert.Equal(t, expected, text)
}

func TestTextPreviewerPreviewWithLexer(t *testing.T) {
	previewer := NewTextPreviewer()
	data := []byte("package main\n")
	dir := filepath.Dir("main.go")
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "main.go"},
		Dir:      dir,
	}

	done := make(chan struct{})
	queueUpdateDraw := func(fn func()) {
		fn()
		close(done)
	}

	previewer.Preview(entry, data, queueUpdateDraw)
	waitForUpdate(t, done)

	text := previewer.GetText(false)
	assert.Contains(t, text, "package")
}

func TestTextPreviewerPreviewWithLexerError(t *testing.T) {
	lexers.Register(&errorLexer{
		config: &chroma.Config{
			Name:      "ErrLexer",
			Filenames: []string{"*.errlex"},
		},
	})

	previewer := NewTextPreviewer()
	data := []byte("content")
	dir := filepath.Dir("file.errlex")
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "file.errlex"},
		Dir:      dir,
	}

	done := make(chan struct{})
	queueUpdateDraw := func(fn func()) {
		fn()
		close(done)
	}

	previewer.Preview(entry, data, queueUpdateDraw)
	waitForUpdate(t, done)

	text := previewer.GetText(false)
	assert.Contains(t, text, "Failed to format file")
}

func TestTextPreviewerPreviewReadsFile(t *testing.T) {
	previewer := NewTextPreviewer()

	content := []byte("file content")
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "note.unknownext")
	err := os.WriteFile(path, content, 0644)
	assert.NoError(t, err)

	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "note.unknownext"},
		Dir:      tmpDir,
	}

	done := make(chan struct{})
	queueUpdateDraw := func(fn func()) {
		fn()
		close(done)
	}

	previewer.Preview(entry, nil, queueUpdateDraw)
	waitForUpdate(t, done)

	text := previewer.GetText(false)
	expected := string(content)
	assert.Equal(t, expected, text)
}

func TestTextPreviewerPreviewReadFileError(t *testing.T) {
	previewer := NewTextPreviewer()
	tmpDir := t.TempDir()
	name := filepath.Base(tmpDir)
	dir := filepath.Dir(tmpDir)
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: name, isDir: true},
		Dir:      dir,
	}

	queueUpdateDraw := func(func()) {}
	previewer.Preview(entry, nil, queueUpdateDraw)
	waitForText(t, previewer, "Failed to read file")
}

func TestTextPreviewerPreviewQueueUpdateNil(t *testing.T) {
	previewer := NewTextPreviewer()
	data := []byte("queue nil")
	dir := filepath.Dir("note.unknownext")
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "note.unknownext"},
		Dir:      dir,
	}

	previewer.Preview(entry, data, nil)
	waitForText(t, previewer, "queue nil")
}

func TestTextPreviewerPreviewStalePlain(t *testing.T) {
	previewer := NewTextPreviewer()
	dir := filepath.Dir("note.unknownext")
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "note.unknownext"},
		Dir:      dir,
	}

	allowFirst := make(chan struct{})
	doneFirst := make(chan struct{})
	doneSecond := make(chan struct{})

	queueUpdateFirst := func(fn func()) {
		<-allowFirst
		fn()
		close(doneFirst)
	}
	queueUpdateSecond := func(fn func()) {
		fn()
		close(doneSecond)
	}

	previewer.Preview(entry, []byte("first"), queueUpdateFirst)
	previewer.Preview(entry, []byte("second"), queueUpdateSecond)
	waitForUpdate(t, doneSecond)

	close(allowFirst)
	waitForUpdate(t, doneFirst)

	waitForText(t, previewer, "second")
}

func TestTextPreviewerPreviewStaleLexer(t *testing.T) {
	previewer := NewTextPreviewer()
	lexerDir := filepath.Dir("main.go")
	plainDir := filepath.Dir("note.unknownext")
	lexerEntry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "main.go"},
		Dir:      lexerDir,
	}
	plainEntry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "note.unknownext"},
		Dir:      plainDir,
	}

	allowFirst := make(chan struct{})
	doneFirst := make(chan struct{})
	doneSecond := make(chan struct{})

	queueUpdateFirst := func(fn func()) {
		<-allowFirst
		fn()
		close(doneFirst)
	}
	queueUpdateSecond := func(fn func()) {
		fn()
		close(doneSecond)
	}

	previewer.Preview(lexerEntry, []byte("package main\n"), queueUpdateFirst)
	previewer.Preview(plainEntry, []byte("second"), queueUpdateSecond)
	waitForUpdate(t, doneSecond)

	close(allowFirst)
	waitForUpdate(t, doneFirst)

	waitForText(t, previewer, "second")
}
func TestTextPreviewerMetaAndMain(t *testing.T) {
	previewer := NewTextPreviewer()
	meta := previewer.Meta()
	main := previewer.Main()
	if meta != nil {
		t.Errorf("expected nil meta, got %v", meta)
	}
	if main != previewer.TextView {
		t.Errorf("expected main to be text view")
	}
}

func TestPrettyJSONSuccess(t *testing.T) {
	input := "{\"a\":1}"
	output, err := prettyJSON(input)
	assert.NoError(t, err)
	assert.Contains(t, output, "\n  \"a\": 1\n")
}

func TestJsonPreviewerPreviewReadsFile(t *testing.T) {
	previewer := NewJsonPreviewer()

	content := []byte("{\"a\":1}")
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "data.unknownext")
	err := os.WriteFile(path, content, 0644)
	assert.NoError(t, err)

	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "data.unknownext"},
		Dir:      tmpDir,
	}

	done := make(chan struct{})
	queueUpdateDraw := func(fn func()) {
		fn()
		close(done)
	}

	previewer.Preview(entry, nil, queueUpdateDraw)
	waitForUpdate(t, done)

	text := previewer.GetText(false)
	assert.Contains(t, text, "\n  \"a\": 1\n")
}

func TestJsonPreviewerPreviewReadFileError(t *testing.T) {
	previewer := NewJsonPreviewer()
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "missing.json"},
		Dir:      t.TempDir(),
	}
	queueUpdateDraw := func(func()) {}
	previewer.Preview(entry, nil, queueUpdateDraw)
}
func TestJsonPreviewerPreviewWithData(t *testing.T) {
	previewer := NewJsonPreviewer()
	data := []byte("{\"a\":1}")
	dir := filepath.Dir("data.unknownext")
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "data.unknownext"},
		Dir:      dir,
	}

	done := make(chan struct{})
	queueUpdateDraw := func(fn func()) {
		fn()
		close(done)
	}

	previewer.Preview(entry, data, queueUpdateDraw)
	waitForUpdate(t, done)

	text := previewer.GetText(false)
	assert.Contains(t, text, "\n  \"a\": 1\n")
}

func TestDsstorePreviewerPreviewSuccess(t *testing.T) {
	previewer := NewDsstorePreviewer()

	store := dsstore.Store{
		Records: []dsstore.Record{
			{
				FileName: "example",
				Type:     "bool",
				Data:     []byte{1},
				DataLen:  0,
			},
		},
	}
	var buffer bytes.Buffer
	err := store.Write(&buffer)
	assert.NoError(t, err)

	dir := filepath.Dir("test.DS_Store")
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "test.DS_Store"},
		Dir:      dir,
	}

	done := make(chan struct{})
	queueUpdateDraw := func(fn func()) {
		fn()
		close(done)
	}

	data := buffer.Bytes()
	previewer.Preview(entry, data, queueUpdateDraw)
	waitForUpdate(t, done)

	text := previewer.GetText(false)
	assert.Contains(t, text, "example: bool")
}

func TestDsstorePreviewerPreviewReadsFile(t *testing.T) {
	previewer := NewDsstorePreviewer()

	store := dsstore.Store{
		Records: []dsstore.Record{
			{
				FileName: "example",
				Type:     "bool",
				Data:     []byte{1},
				DataLen:  0,
			},
		},
	}
	var buffer bytes.Buffer
	err := store.Write(&buffer)
	assert.NoError(t, err)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "good.DS_Store")
	data := buffer.Bytes()
	err = os.WriteFile(path, data, 0644)
	assert.NoError(t, err)

	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "good.DS_Store"},
		Dir:      tmpDir,
	}

	done := make(chan struct{})
	queueUpdateDraw := func(fn func()) {
		fn()
		close(done)
	}

	previewer.Preview(entry, nil, queueUpdateDraw)
	waitForUpdate(t, done)

	text := previewer.GetText(false)
	assert.Contains(t, text, "example: bool")
}

func TestDsstorePreviewerPreviewReadFileError(t *testing.T) {
	previewer := NewDsstorePreviewer()
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "missing.DS_Store"},
		Dir:      t.TempDir(),
	}
	queueUpdateDraw := func(func()) {}
	previewer.Preview(entry, nil, queueUpdateDraw)
}

func TestDsstorePreviewerPreviewError(t *testing.T) {
	previewer := NewDsstorePreviewer()
	dir := filepath.Dir("bad.DS_Store")
	entry := files.EntryWithDirPath{
		DirEntry: mockDirEntry{name: "bad.DS_Store"},
		Dir:      dir,
	}

	data := []byte("not a dsstore")
	queueUpdateDraw := func(func()) {}
	previewer.Preview(entry, data, queueUpdateDraw)

	text := previewer.GetText(false)
	assert.Contains(t, text, "Failed to read")
}
