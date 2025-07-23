// Package copy contains 'CopyFolder' and the excluded array.
// The CopyFolder func is used for copying a folder into another folder.
// The excluded array contains all folders that should NOT be copied into the destination directory.
package copy

import (
	"backupsync/include"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyFolder(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err) // use %w instead of %e
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		if relPath != "." && !include.Include(relPath, d) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath)
	})
}

func copyFile(srcFile, dstFile string) error {
	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Optionally preserve permissions
	info, err := os.Stat(srcFile)
	if err != nil {
		return err
	}
	return os.Chmod(dstFile, info.Mode())
}
