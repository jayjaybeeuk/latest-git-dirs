package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLatestModified_BasicDir(t *testing.T) {
	// Create a temp directory structure
	tmp := t.TempDir()

	// Create a file with known mod time
	f := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	latest, err := LatestModified(tmp, 10, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}

	if latest.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}

	// Should be recent (within last 5 seconds)
	if time.Since(latest) > 5*time.Second {
		t.Fatalf("expected recent timestamp, got %v", latest)
	}
}

func TestLatestModified_ExcludesDir(t *testing.T) {
	tmp := t.TempDir()

	// Create excluded dir with a newer file
	excluded := filepath.Join(tmp, "node_modules")
	os.MkdirAll(excluded, 0755)
	os.WriteFile(filepath.Join(excluded, "pkg.json"), []byte("{}"), 0644)

	// Create a normal file that's older
	normal := filepath.Join(tmp, "old.txt")
	os.WriteFile(normal, []byte("old"), 0644)
	oldTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes(normal, oldTime, oldTime)
	os.Chtimes(tmp, oldTime, oldTime)

	excludes := map[string]bool{"node_modules": true}
	latest, err := LatestModified(tmp, 10, excludes)
	if err != nil {
		t.Fatal(err)
	}

	// Should match the old file time, not the excluded dir
	if time.Since(latest) < 30*time.Minute {
		t.Fatal("expected old timestamp since node_modules should be excluded")
	}
}

func TestLatestModified_RespectsMaxDepth(t *testing.T) {
	tmp := t.TempDir()

	// Create deeply nested file
	deep := filepath.Join(tmp, "a", "b", "c", "d")
	os.MkdirAll(deep, 0755)
	deepFile := filepath.Join(deep, "deep.txt")
	os.WriteFile(deepFile, []byte("deep"), 0644)

	// Create shallow old file
	shallow := filepath.Join(tmp, "shallow.txt")
	os.WriteFile(shallow, []byte("shallow"), 0644)
	oldTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes(shallow, oldTime, oldTime)
	os.Chtimes(tmp, oldTime, oldTime)
	os.Chtimes(filepath.Join(tmp, "a"), oldTime, oldTime)

	// maxDepth=1 should not reach depth 4
	latest, err := LatestModified(tmp, 1, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}

	if time.Since(latest) < 30*time.Minute {
		t.Fatal("expected old timestamp since deep files should be beyond max depth")
	}
}

func TestLatestModified_EmptyDir(t *testing.T) {
	tmp := t.TempDir()

	latest, err := LatestModified(tmp, 10, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}

	// Empty dir still has its own mod time
	if latest.IsZero() {
		t.Fatal("expected non-zero timestamp for directory itself")
	}
}

func TestDepth(t *testing.T) {
	tests := []struct {
		path     string
		expected int
	}{
		{"a", 1},
		{"a/b", 2},
		{"a/b/c", 3},
		{"/a/b", 3}, // leading slash = empty + a + b
	}

	for _, tt := range tests {
		got := depth(tt.path)
		if got != tt.expected {
			t.Errorf("depth(%q) = %d, want %d", tt.path, got, tt.expected)
		}
	}
}
