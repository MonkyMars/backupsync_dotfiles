package main

import (
	"backupsync/compare"
	"backupsync/config"
	"backupsync/copy"
	"fmt"
	"log"
	"os"
)

func main() {
	config, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	entries, err := os.ReadDir("/mnt/")
	if err != nil {
		log.Fatal("Failed to read directory with name /mnt", err)
	}

	found := false
	for _, e := range entries {
		if e.Name() == config.DriveName {
			found = true
			dir, err := os.ReadDir(config.DestinationDirectory)
			if err != nil {
				fmt.Println("Failed to read destination directory")
			}
			if dir == nil {
				os.Mkdir(config.DestinationDirectory, 0750)
			}
		}
	}

	if !found {
		err := os.Mkdir(config.DestinationDirectory, 0750)
		if err != nil {
			log.Fatal("Failed to create destination directory", err)
		}
	}

	fmt.Println("Comparing directories...")

	isIdentical, err := compare.CompareDirs(config.SourceDirectory, config.DestinationDirectory)
	if err != nil {
		fmt.Println("Non fatal error occured: failed to compare directories; copying directories anyway", err)
		isIdentical = false
	}

	fmt.Println("Completed directory comparison")

	if !isIdentical {
		fmt.Println("Copying files from source to destination directory")
		err = copy.CopyFolder(config.SourceDirectory, config.DestinationDirectory)
		if err != nil {
			log.Fatal("Failed to copy dir", err)
		}
	} else {
		fmt.Println("Directories are identical; skipped backing up to target dir")
	}
}
