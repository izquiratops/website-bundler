package directory

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
)

func MoveTo(entryPath, destPath string) error {
	inputFile, err := os.Open(entryPath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %v", err)
	}

	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %v", err)
	}

	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("couldn't copy to dest from source: %v", err)
	}

	return nil
}

func Clean(distLocalPath string) {
	whitelist := []string{"index.html", "favicon.ico"}

	if _, err := os.Stat(distLocalPath); os.IsNotExist(err) {
		message := fmt.Sprintf("directory %s does not exist.\n", distLocalPath)
		log.Fatal(message)
	}

	err := filepath.Walk(distLocalPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the dist directory itself and whitelisted files like index.html
		if path == distLocalPath || slices.Contains(whitelist, info.Name()) {
			return nil
		}

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %v", path, err)
		}

		fmt.Printf("removed: %s\n", path)
		return nil
	})

	if err != nil {
		log.Fatal("error cleaning directory:", err)
	}
}
