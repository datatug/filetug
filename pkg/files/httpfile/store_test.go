package httpfile

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func Test_NewStore(t *testing.T) {
	t.Run("https://example.com/pub/", func(t *testing.T) {
		root, _ := url.Parse("https://example.com/pub/")
		store := NewStore(*root)
		assert.NotNil(t, store)
	})

	t.Run("https://example.com/pub/", func(t *testing.T) {
		root, _ := url.Parse("https://example.com/pub/")
		store := NewStore(*root, WithHttpClient(&http.Client{}))
		assert.NotNil(t, store)
	})
}

func Test_httpFileStore_ReadDir(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mockClient := &http.Client{
		Transport: &mockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				var body string
				switch req.URL.Path {
				case "/pub/":
					body = `<a href="linux/">linux/</a><a href="scm/">scm/</a><a href="tools/">tools/</a>`
				case "/pub/linux/":
					body = `<a href="kernel/">kernel/</a><a href="utils/">utils/</a>`
				default:
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewBufferString("Not Found")),
					}, nil
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			},
		},
	}

	root, _ := url.Parse("https://cdn.kernel.org/")
	store := NewStore(*root, WithHttpClient(mockClient))
	t.Run("Root", func(t *testing.T) {
		entries, err := store.ReadDir(ctx, "/pub/")
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
		entries, err := store.ReadDir(ctx, "/pub/linux/")
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
