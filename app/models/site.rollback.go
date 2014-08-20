package models

import (
	"log"
	"os"
)

type (
	TemplateRollBackFunction func(*InstallTemplate) error
	FileRollBackFunction     func(fileName string) error
	FileRecoverFunction      func(fb *FileBackup) error
)

type SiteRollBack struct {
	Template      *InstallTemplate
	DBFunctions   []TemplateRollBackFunction
	FileFunctions []FileRollBack
	FileRecovers  []FileRecover
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

	for _, fn := range sb.DBFunctions {
		err := fn(sb.Template)
		if err != nil {
			log.Printf("Rollback error: %v\n", err.Error())
		}
	}

	for _, fileFunc := range sb.FileFunctions {
		err := fileFunc.Function(fileFunc.FileName)
		if err != nil {
			log.Printf("Rollback error: %v\n", err.Error())
		}
	}

	for _, fileRecover := range sb.FileRecovers {
		err := fileRecover.Function(fileRecover.FileBackup)
		if err != nil {
			log.Println("Rollback error: %v\n", err.Error())
		}
	}
}
