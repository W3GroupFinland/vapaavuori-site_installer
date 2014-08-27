package models

import (
	"log"
	"os"
)

type (
	TemplateRollBackFunction   func(*InstallTemplate) error
	DatabaseIdRollBackFunction func(id int64) error
	FileRollBackFunction       func(fileName string) error
	FileRecoverFunction        func(fb *FileBackup) error
)

type SiteRollBack struct {
	Template      *InstallTemplate
	DBFunctions   []TemplateRollBackFunction
	DBIDFunctions []DatabaseIdRollBack
	FileFunctions []FileRollBack
	FileRecovers  []FileRecover
}

type DatabaseIdRollBack struct {
	Id       int64
	Function DatabaseIdRollBackFunction
}

type FileRollBack struct {
	FileName string
	Function FileRollBackFunction
}

type FileRecover struct {
	FileBackup *FileBackup
	Function   FileRecoverFunction
}

type FileBackup struct {
	NewFile string
	Backup  string
}

func NewSiteRollBack(tmpl *InstallTemplate) *SiteRollBack {
	return &SiteRollBack{Template: tmpl}
}

func (sb *SiteRollBack) AddDBFunction(rollBack TemplateRollBackFunction) {
	sb.DBFunctions = append(sb.DBFunctions, rollBack)
}

func (sb *SiteRollBack) AddDBIdFunction(fn DatabaseIdRollBackFunction, id int64) {
	dbRB := DatabaseIdRollBack{Id: id, Function: fn}
	sb.DBIDFunctions = append(sb.DBIDFunctions, dbRB)
}

func (sb *SiteRollBack) AddFileFunction(fn FileRollBackFunction, fileName string) {
	fileRB := FileRollBack{FileName: fileName, Function: fn}
	sb.FileFunctions = append(sb.FileFunctions, fileRB)
}

func (sb *SiteRollBack) AddFileRecoverFunction(fn FileRecoverFunction, fb *FileBackup) {
	recoverRB := FileRecover{FileBackup: fb, Function: fn}
	sb.FileRecovers = append(sb.FileRecovers, recoverRB)
}

func (fb *FileBackup) Delete() error {
	err := os.Remove(fb.Backup)
	if err != nil {
		return err
	}

	return nil
}

func (sb *SiteRollBack) DeleteBackupFiles() {
	for _, backup := range sb.FileRecovers {
		err := backup.FileBackup.Delete()
		if err != nil {
			log.Println(err)
		}
	}

	sb.FileRecovers = sb.FileRecovers[:0]
}

func (sb *SiteRollBack) Execute() {
	log.Println("Rolling back..")

	// Process roll back functions in reverse order.
	for i := len(sb.DBFunctions) - 1; i >= 0; i-- {
		v := sb.DBFunctions[i]
		err := v(sb.Template)
		if err != nil {
			log.Printf("Rollback error: %v\n", err.Error())
		}
	}

	for i := len(sb.DBIDFunctions) - 1; i >= 0; i-- {
		v := sb.DBIDFunctions[i]
		err := v.Function(v.Id)
		if err != nil {
			log.Printf("Rollback error: %v\n", err.Error())
		}
	}

	for i := len(sb.FileFunctions) - 1; i >= 0; i-- {
		v := sb.FileFunctions[i]
		err := v.Function(v.FileName)
		if err != nil {
			log.Printf("Rollback error: %v\n", err.Error())
		}
	}

	for i := len(sb.FileRecovers) - 1; i >= 0; i-- {
		v := sb.FileRecovers[i]
		err := v.Function(v.FileBackup)
		if err != nil {
			log.Println("Rollback error: %v\n", err.Error())
		}
	}
}
