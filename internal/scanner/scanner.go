package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func depth(path string) int {
	clean := filepath.Clean(path)
	return len(strings.Split(clean, string(filepath.Separator)))
}

func LatestModified(path string, maxDepth int, excludes map[string]bool) (time.Time, error) {
	var latest time.Time

	rootDepth := depth(path)

	err := filepath.WalkDir(path, func(current string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() && excludes[d.Name()] {
			return filepath.SkipDir
		}

		currentDepth := depth(current) - rootDepth
		if currentDepth > maxDepth {
			return filepath.SkipDir
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().After(latest) {
			latest = info.ModTime()
		}

		return nil
	})

	return latest, err
}
