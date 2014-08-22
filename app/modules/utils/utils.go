package utils

import (
	"bitbucket.org/kardianos/osext"
	"bytes"
	"code.google.com/p/gcfg"
	"crypto/rand"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func RandomString(length int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}

	return string(bytes)
}

// Get file path in relation to application executable.
func GetFileFullPath(file string) (string, error) {
	// Check if first charachter of string is /.
	// In that case given file path is absolute.
	matched, err := regexp.MatchString(`^[/].{0,}$`, file)
	if err != nil {
		return "", err
	}
	if matched {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return "", err
		}

		return file, nil
	}

	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	path := filepath.Join(folderPath, file)
	absPath, err := GetAbsDirectory(path)
	if err != nil {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}

		absPath, err = GetAbsDirectory(dir)
		if err != nil {
			return "", err
		}

		path = filepath.Join(absPath, file)
	}

	return path, nil
}

func ReadConfigData(data []byte, c interface{}) error {
	r := bytes.NewReader(data)

	err := gcfg.ReadInto(c, r)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfigFile(file string, c interface{}) error {
	path, err := GetFileFullPath(file)
	if err != nil {
		return err
	}

	data, err := ReadFileContents(path)
	if err != nil {
		return err
	}

	err = gcfg.ReadInto(c, bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}

func GetAbsDirectory(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", err
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func FileExists(fp string) bool {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		log.Printf("No such file or directory: %v", fp)
		return false
	}

	return true
}

func StripPathWhiteSpace(in []byte) []byte {
	var (
		bytes    []byte
		inLength = len(in)
	)
	for idx, b := range in {
		if (idx == 0 || idx == (inLength-1)) && b == 10 {
			continue
		}
		bytes = append(bytes, b)
	}

	return bytes
}
