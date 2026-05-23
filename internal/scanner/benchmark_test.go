package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkLatestModified_Shallow(b *testing.B) {
	tmp := b.TempDir()

	// Create 100 files
	for i := 0; i < 100; i++ {
		os.WriteFile(filepath.Join(tmp, fmt.Sprintf("file%d.txt", i)), []byte("x"), 0644)
	}

	excludes := map[string]bool{".git": true, "node_modules": true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LatestModified(tmp, 10, excludes)
	}
}

func BenchmarkLatestModified_Deep(b *testing.B) {
	tmp := b.TempDir()

	// Create nested structure: 10 dirs × 10 files each
	for d := 0; d < 10; d++ {
		dir := filepath.Join(tmp, fmt.Sprintf("dir%d", d))
		os.MkdirAll(dir, 0755)
		for f := 0; f < 10; f++ {
			os.WriteFile(filepath.Join(dir, fmt.Sprintf("f%d.txt", f)), []byte("x"), 0644)
		}
	}

	excludes := map[string]bool{".git": true, "node_modules": true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LatestModified(tmp, 10, excludes)
	}
}
