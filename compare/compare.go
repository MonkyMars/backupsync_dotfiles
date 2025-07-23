// Package compare contains functions for comparing directories.
// This will be used to check if copying files over from the source to the target is neccesary.
// If not, it saves time and resources.
package compare

import (
	"backupsync/include"
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func readDirFiles(root string) (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		fmt.Println("Checking path:", relPath)

		if relPath != "." && !include.Include(relPath, d) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			files[relPath] = path
		}
		return nil
	})
	return files, err
}

func filesAreEqual(file1, file2 string) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	s1 := bufio.NewScanner(f1)
	s2 := bufio.NewScanner(f2)

	for s1.Scan() {
		if !s2.Scan() || s1.Text() != s2.Text() {
			return false, nil
		}
	}
	// If f2 still has more lines
	if s2.Scan() {
		return false, nil
	}
	return true, nil
}

func CompareDirs(dir1, dir2 string) (bool, error) {
	files1, err := readDirFiles(dir1)
	if err != nil {
		return false, err
	}
	files2, err := readDirFiles(dir2)
	if err != nil {
		return false, err
	}

	if len(files1) != len(files2) {
		fmt.Println("Different number of files.")
		return false, nil
	}

	for relPath, path1 := range files1 {
		path2, exists := files2[relPath]
		if !exists {
			fmt.Printf("File %s is missing in %s\n", relPath, dir2)
			return false, nil
		}

		equal, err := filesAreEqual(path1, path2)
		if err != nil {
			return false, err
		}
		if !equal {
			fmt.Printf("File %s differs.\n", relPath)
			return false, nil
		}
	}

	return true, nil
}
