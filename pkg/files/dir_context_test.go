package files

import (
	"context"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	root url.URL
}

func (m mockStore) RootTitle() string { return "Mock" }
func (m mockStore) RootURL() url.URL  { return m.root }
func (m mockStore) ReadDir(ctx context.Context, name string) ([]os.DirEntry, error) {
	_, _ = ctx, name
	return nil, nil
}
func (m mockStore) CreateDir(ctx context.Context, path string) error {
	_, _ = ctx, path
	return nil
}
func (m mockStore) CreateFile(ctx context.Context, path string) error {
	_, _ = ctx, path
	return nil
}
func (m mockStore) Delete(ctx context.Context, path string) error {
	_, _ = ctx, path
	return nil
}

func TestDirContextMethods(t *testing.T) {
	tempDir := filepath.ToSlash(t.TempDir())
	ctx := NewDirContext(mockStore{root: url.URL{Scheme: "file"}}, tempDir, nil)

	ctx.SetChildren([]os.DirEntry{NewDirEntry("a.txt", false)})
	assert.Len(t, ctx.Children(), 1)

	entries := ctx.Entries()
	if assert.Len(t, entries, 1) {
		assert.Equal(t, "a.txt", entries[0].Name())
		assert.Equal(t, tempDir, entries[0].DirPath())
	}

	assert.Equal(t, path.Dir(tempDir), ctx.DirPath())
	assert.Equal(t, tempDir, ctx.FullName())
	assert.Equal(t, tempDir, ctx.String())
	assert.Equal(t, path.Base(tempDir), ctx.Name())
	assert.True(t, ctx.IsDir())
	assert.Equal(t, os.ModeDir, ctx.Type())
	info, err := ctx.Info()
	assert.NoError(t, err)
	assert.NotNil(t, info)

	root := NewDirContext(nil, "/", nil)
	assert.Equal(t, "/", root.Name())

	empty := NewDirContext(nil, "", nil)
	assert.Equal(t, "", empty.DirPath())
	assert.Equal(t, "", empty.Name())
	info, err = empty.Info()
	assert.NoError(t, err)
	assert.Nil(t, info)

	nonFileStore := mockStore{root: url.URL{Scheme: "ftp"}}
	nonFileCtx := NewDirContext(nonFileStore, tempDir, nil)
	info, err = nonFileCtx.Info()
	assert.NoError(t, err)
	assert.Nil(t, info)
}
