package ftpfile

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/stretchr/testify/assert"
)

type mockFtpClient struct {
	LoginFunc func(user, password string) error
	ListFunc  func(path string) ([]*ftp.Entry, error)
	QuitFunc  func() error
}

func (m *mockFtpClient) Login(user, password string) error {
	if m.LoginFunc != nil {
		return m.LoginFunc(user, password)
	}
	return nil
}

func (m *mockFtpClient) List(path string) ([]*ftp.Entry, error) {
	if m.ListFunc != nil {
		return m.ListFunc(path)
	}
	return nil, nil
}

func (m *mockFtpClient) Quit() error {
	if m.QuitFunc != nil {
		return m.QuitFunc()
	}
	return nil
}

func TestStore_ReadDir_Mock(t *testing.T) {
	root, _ := url.Parse("ftp://demo:password@example.com/")
	mockClient := &mockFtpClient{
		ListFunc: func(path string) ([]*ftp.Entry, error) {
			return []*ftp.Entry{
				{Name: "file1.txt", Type: ftp.EntryTypeFile, Size: 100},
				{Name: "dir1", Type: ftp.EntryTypeFolder},
			}, nil
		},
	}

	factory := func(addr string, options ...ftp.DialOption) (FtpClient, error) {
		return mockClient, nil
	}

	s := NewStore(*root, WithFtpClientFactory(factory))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	entries, err := s.ReadDir(ctx, ".")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entries))
	assert.Equal(t, "file1.txt", entries[0].Name())
	assert.False(t, entries[0].IsDir())
	assert.Equal(t, "dir1", entries[1].Name())
	assert.True(t, entries[1].IsDir())
}

func TestStore_ReadDir(t *testing.T) {
	if os.Getenv("RUN_FTP_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_FTP_INTEGRATION_TESTS=true to run")
	}
	const host = "test.rebex.net"
	const port = 21
	root := url.URL{
		Scheme: "ftp",
		Host:   fmt.Sprintf("%s:%d", host, port),
		User:   url.UserPassword("demo", "password"),
	}
	t.Run("host_with_port", func(t *testing.T) {
		root := root
		root.Host = fmt.Sprintf("%s:%d", host, port)
		s := NewStore(root)
		testReadDir(t, s)
	})

	t.Run("plain_default_port", func(t *testing.T) {
		root := root
		root.Host = host
		s := NewStore(root)
		testReadDir(t, s)
	})

	t.Run("explicit_TLS", func(t *testing.T) {
		t.Skip("test.rebex.net requires TLS session resumption which github.com/jlaffaye/ftp might not support or needs more config")
		s := NewStore(root)
		s.SetTLS(true, false)
		testReadDir(t, s)
	})

	t.Run("implicit_TLS", func(t *testing.T) {
		t.Skip("test.rebex.net requires TLS session resumption which github.com/jlaffaye/ftp might not support or needs more config")
		s := NewStore(root)
		s.SetTLS(false, true)
		testReadDir(t, s)
	})
}

func testReadDir(t *testing.T, s *Store) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	entries, err := s.ReadDir(ctx, ".")
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}

	if len(entries) == 0 {
		t.Error("expected at least one entry, got 0")
	}

	for _, entry := range entries {
		t.Logf("Entry: %s, IsDir: %v", entry.Name(), entry.IsDir())
	}
}
