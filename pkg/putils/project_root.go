package putils

import (
	"os"
	"path/filepath"
)

// Returns project's root by looking for specific file or directory - marker.
// If no marker provided it will use "go.mod" by default.
// If provided several markers, it will use first one.
func GetProjectRoot(startDir string, marker ...string) (string, error) {
	stopFlag := "go.mod"
	if len(marker) > 0 {
		stopFlag = marker[0]
	}
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, stopFlag)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}
