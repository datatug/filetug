package filetug

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/datatug/filetug/pkg/files/osfile"
	"github.com/go-git/go-git/v5"
	"github.com/rivo/tview"
)

func TestTree_SetDirContext_GitOptimization(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tree-git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Initialize git repo
	_, err = git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Create subdirectories
	subDir1 := filepath.Join(tempDir, "subdir1")
	subDir2 := filepath.Join(tempDir, "subdir2")
	_ = os.Mkdir(subDir1, 0755)
	_ = os.Mkdir(subDir2, 0755)

	app := tview.NewApplication()
	nav := NewNavigator(app)
	nav.store = osfile.NewStore(tempDir)

	// Mock queueUpdateDraw to avoid hanging
	nav.queueUpdateDraw = func(f func()) {
		f()
	}

	tree := NewTree(nav)
	node := tview.NewTreeNode("root")

	dirEntries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read dir: %v", err)
	}

	dirContext := &DirContext{
		Path:     tempDir,
		children: dirEntries,
		Store:    nav.store,
	}

	ctx := context.Background()
	tree.setDirContext(ctx, node, dirContext)

	// Give some time for goroutines to start and call updateGitStatus
	time.Sleep(100 * time.Millisecond)

	// Check if children were added
	children := node.GetChildren()
	if len(children) < 2 {
		t.Errorf("Expected at least 2 children (subdir1, subdir2), got %d", len(children))
	}
}

func TestNavigator_ShowDir_GitStatusText(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nav-git-text-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Initialize git repo
	_, err = git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Create a subdirectory
	subDirName := "subdir1"
	subDirPath := filepath.Join(tempDir, subDirName)
	_ = os.Mkdir(subDirPath, 0755)

	app := tview.NewApplication()
	nav := NewNavigator(app)
	nav.store = osfile.NewStore(tempDir)

	// Mock queueUpdateDraw to execute immediately
	nav.queueUpdateDraw = func(f func()) {
		f()
	}

	// Create a tree node for the subdirectory as it would be in the tree
	// In the tree, it would have a prefix like "📁subdir1"
	node := tview.NewTreeNode("📁" + subDirName).SetReference(subDirPath)

	ctx := context.Background()

	// When showDir is called (e.g., when a node is selected)
	nav.showDir(ctx, node, subDirPath, false)

	// Wait for goroutines
	time.Sleep(200 * time.Millisecond)

	text := node.GetText()
	if strings.Contains(text, tempDir) {
		t.Errorf("Node text contains full path, but it should only contain dir name and git status. Got: %q", text)
	}

	if !strings.HasPrefix(text, "📁"+subDirName) {
		t.Errorf("Node text should start with original name %q, got: %q", "📁"+subDirName, text)
	}
}
