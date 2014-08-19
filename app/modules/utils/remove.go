package utils

import (
	"log"
	"os"
)

func RemoveDirectory(path string) error {
	log.Printf("Removing path: %v.\n", path)
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}

func RemoveFile(fileName string) error {
	log.Printf("Removing file or symlink: %v.\n", fileName)
	err := os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}
