package filetug

import (
	"context"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type recordingStore struct {
	mu    sync.Mutex
	paths []string
}

func (r *recordingStore) RootTitle() string { return "recording" }
func (r *recordingStore) RootURL() url.URL  { return url.URL{Path: "/"} }
func (r *recordingStore) ReadDir(ctx context.Context, path string) ([]os.DirEntry, error) {
	_, _ = ctx, path
	return nil, nil
}
func (r *recordingStore) Delete(ctx context.Context, path string) error {
	_, _ = ctx, path
	return nil
}
func (r *recordingStore) CreateDir(ctx context.Context, path string) error {
	_, _ = ctx, path
	r.mu.Lock()
	r.paths = append(r.paths, path)
	r.mu.Unlock()
	return nil
}
func (r *recordingStore) CreateFile(ctx context.Context, path string) error {
	_, _ = ctx, path
	return nil
}

func TestGeneratedNestedDirs_DefaultFormat(t *testing.T) {
	store := &recordingStore{}
	err := GeneratedNestedDirs(context.Background(), store, "/root", "", 2, 2)
	assert.NoError(t, err)

	expected := map[string]struct{}{
		"/root":                       {},
		"/root/Directory0":            {},
		"/root/Directory1":            {},
		"/root/Directory0/Directory0": {},
		"/root/Directory0/Directory1": {},
		"/root/Directory1/Directory0": {},
		"/root/Directory1/Directory1": {},
	}

	store.mu.Lock()
	paths := append([]string(nil), store.paths...)
	store.mu.Unlock()

	assert.Len(t, paths, len(expected))

	got := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		got[p] = struct{}{}
	}
	assert.Len(t, got, len(expected))

	for p := range expected {
		if _, ok := got[p]; !ok {
			t.Errorf("expected path %q to be created", p)
		}
	}
}

func TestGeneratedNestedDirs_DepthZero(t *testing.T) {
	store := &recordingStore{}
	err := GeneratedNestedDirs(context.Background(), store, "/root", "Dir%d", 0, 3)
	assert.NoError(t, err)

	store.mu.Lock()
	paths := append([]string(nil), store.paths...)
	store.mu.Unlock()

	assert.Equal(t, []string{"/root"}, paths)
}
