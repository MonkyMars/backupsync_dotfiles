// Package include contains an array and a function of directories that need to be included from copying.
package include

import (
	"io/fs"
	"os"
	"path/filepath"
)

var Included = []string{"hypr", "nvim", "kitty", "rofi", "waybar"}

func Include(path string, d fs.DirEntry) bool {
	for _, inc := range Included {
		if path == inc || filepath.HasPrefix(path, inc+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}
