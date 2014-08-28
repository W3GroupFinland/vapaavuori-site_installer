package utils

import (
	"fmt"
	"log"
	"os"
)

type RequireFiles struct {
	List []*FileRequired
}

type FileRequired struct {
	Role     string
	FileName string
}

func NewRequireFiles() *RequireFiles {
	return &RequireFiles{}
}

func (rf *RequireFiles) Add(file string, role string) *RequireFiles {
	rf.List = append(rf.List, &FileRequired{
		Role:     role,
		FileName: file,
	})

	return rf
}

func (rf *RequireFiles) Require() {
	var errList []string
	for _, f := range rf.List {
		if !FileExists(f.FileName) {
			msg := fmt.Sprintf("Fatal error: %v name \"%v\" doesn't exist.", f.Role, f.FileName)
			errList = append(errList, msg)
		}
	}

	var str string
	if len(errList) > 0 {
		for _, msg := range errList {
			str += msg + "\n"
		}

		log.Fatalln(str)
	}
}

func FileExists(fp string) bool {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		log.Printf("No such file or directory: %v", fp)
		return false
	}

	return true
}
