package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMainRoot(t *testing.T) {
	runCalled := false

	oldRun := run
	defer func() {
		run = oldRun
	}()
	run = func(app application) {
		runCalled = true
	}

	main()

	if !runCalled {
		t.Fatal("expected main function to call run")
	}
}

func Test_newApp(t *testing.T) {
	app := newApp()
	if app == nil {
		t.Errorf("newApp returned nil")
	}
}

type fakeApp struct {
	err error
}

func (f fakeApp) Run() error {
	return fmt.Errorf("app failed: %w", f.err)
}

func Test_run(t *testing.T) {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	defer func() {
		os.Stderr = oldStderr
	}()

	var expectedErr = errors.New("test error")
	run(fakeApp{err: expectedErr})

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, expectedErr.Error()) {
		t.Errorf("expected stderr to contain %q, got %q", expectedErr.Error(), output)
	}
}
