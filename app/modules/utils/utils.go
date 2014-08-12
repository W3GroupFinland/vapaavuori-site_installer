package utils

import (
	"bitbucket.org/kardianos/osext"
	"code.google.com/p/gcfg"
	"crypto/rand"
	"log"
	"os"
	"path/filepath"
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

func ReadConfigFile(file string, c interface{}) {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(folderPath, file)
	absPath, err := GetAbsDirectory(path)
	if err != nil {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}

		absPath, err = GetAbsDirectory(dir)
		if err != nil {
			log.Fatalln(err)
		}

		path = filepath.Join(absPath, file)
	}

	err = gcfg.ReadFileInto(c, path)
	if err != nil {
		log.Fatal(err)
	}
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
