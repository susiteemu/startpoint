package paths

import (
	"os"
	"path/filepath"
	"strings"
)

func ShortenPath(path string) string {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		if strings.HasPrefix(path, homeDir) {
			path = "~" + path[len(homeDir):]
		}
	}

	// remove separator char from the suffix
	path, _ = strings.CutSuffix(path, string(filepath.Separator))
	parts := strings.Split(path, string(filepath.Separator))

	// Abbreviate all directories except for the last one
	for i := 0; i < len(parts)-1; i++ {
		if len(parts[i]) > 1 {
			parts[i] = string(parts[i][0]) // Take the first character of each directory
		}
	}
	// Join the shortened parts back together
	return strings.Join(parts, string(filepath.Separator))
}
