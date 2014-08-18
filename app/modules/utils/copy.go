package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type CopyTarget struct {
	TargetPath   string
	SourcePath   string
	SourceLength int
}

// Main function to copy directory recursively
func (ct *CopyTarget) CopyDirectory(source string, target string) error {
	ct.TargetPath = target
	ct.SourcePath = source
	ct.SourceLength = len(ct.SourcePath)

	err := filepath.Walk(source, ct.copyFilesWalkFunc)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Copies single file.
func CopyFile(source string, destination string) error {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(destination, si.Mode())
		}

	}

	return err
}

// Walk function to recursively copy files and directories in given path.
func (ct *CopyTarget) copyFilesWalkFunc(source string, info os.FileInfo, err error) error {
	copyPath := source[ct.SourceLength:]
	destination := ct.TargetPath + copyPath

	if err != nil {
		return err
	}

	if info.IsDir() {
		// Create directory
		os.Mkdir(destination, info.Mode())
		if err != nil {
			return err
		}
		return nil
	}

	err = CopyFile(source, destination)

	return err
}
