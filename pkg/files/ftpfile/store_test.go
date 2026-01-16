package ftpfile

import (
	"testing"
)

func TestStore_ReadDir(t *testing.T) {
	t.Run("plain", func(t *testing.T) {
		s := NewStore("test.rebex.net:21", "demo", "password")
		testReadDir(t, s)
	})

	t.Run("plain_default_port", func(t *testing.T) {
		s := NewStore("test.rebex.net", "demo", "password")
		testReadDir(t, s)
	})

	t.Run("explicit_TLS", func(t *testing.T) {
		t.Skip("test.rebex.net requires TLS session resumption which github.com/jlaffaye/ftp might not support or needs more config")
		s := NewStore("test.rebex.net:21", "demo", "password")
		s.SetTLS(true, false)
		testReadDir(t, s)
	})

	t.Run("implicit_TLS", func(t *testing.T) {
		t.Skip("test.rebex.net requires TLS session resumption which github.com/jlaffaye/ftp might not support or needs more config")
		s := NewStore("test.rebex.net:990", "demo", "password")
		s.SetTLS(false, true)
		testReadDir(t, s)
	})
}

func testReadDir(t *testing.T, s *Store) {
	entries, err := s.ReadDir(".")
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
