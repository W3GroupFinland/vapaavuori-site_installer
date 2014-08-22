package templates

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type Templates struct {
	Templates *template.Template
	CDelims   CDelims
	BasePath  string
}

const (
	DELIM_LEFT  = "{{%"
	DELIM_RIGHT = "%}}"
)

type CDelims struct {
	Use   bool
	Left  string
	Right string
}

func (t *Templates) CustomDelims(args ...string) *Templates {
	if len(args) < 2 {
		t.CDelims = CDelims{Use: true, Left: DELIM_LEFT, Right: DELIM_RIGHT}
		return t
	}

	t.CDelims = CDelims{Use: true, Left: args[0], Right: args[1]}

	return t
}

func (t *Templates) ReadDir(basePath string) error {
	t.BasePath = basePath
	err := filepath.Walk(basePath, t.PathWalkFunc)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (t *Templates) PathWalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	templateName := path[len(t.BasePath):]

	if t.Templates == nil {
		t.Templates = template.New(templateName)
		if t.CDelims.Use {
			t.Templates.Delims(t.CDelims.Left, t.CDelims.Right)
		}
		_, err = t.Templates.ParseFiles(path)
	} else {
		_, err = t.Templates.New(templateName).ParseFiles(path)
	}
	log.Printf("Processed template %s\n", templateName)

	return err
}
