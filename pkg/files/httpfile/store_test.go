package httpfile

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_httpFileStore_ReadDir(t *testing.T) {
	root, _ := url.Parse("https://cdn.kernel.org/")
	store := NewStore(*root)
	t.Run("Root", func(t *testing.T) {
		entries, err := store.ReadDir("/pub/")
		assert.NoError(t, err)
		assert.Greater(t, len(entries), 0)
		expectedNames := []string{"linux", "scm", "tools"}
		for _, name := range expectedNames {
			found := false
			for _, entry := range entries {
				if entry.Name() == name {
					found = true
					break
				}
			}
			assert.True(t, found, "expected to find %s in /pub/", name)
		}
	})
	t.Run("linux", func(t *testing.T) {
		entries, err := store.ReadDir("/pub/linux/")
		assert.NoError(t, err)
		assert.Greater(t, len(entries), 0)
		expectedNames := []string{"kernel", "utils"}
		for _, name := range expectedNames {
			found := false
			for _, entry := range entries {
				if entry.Name() == name {
					found = true
					break
				}
			}
			assert.True(t, found, "expected to find %s in /pub/linux/", name)
		}
	})
}
