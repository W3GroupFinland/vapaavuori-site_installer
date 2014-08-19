package models

import (
	"log"
)

type TemplateRollBackFunction func(*InstallTemplate) error

type SiteRollBack struct {
	Template      *InstallTemplate
	DBFunctions   []TemplateRollBackFunction
	FileFunctions []FileRollBack
}

type FileRollBack struct {
	FileName string
	Function func(fileName string) error
}

func NewSiteRollBack(tmpl *InstallTemplate) *SiteRollBack {
	return &SiteRollBack{Template: tmpl}
}

func (sb *SiteRollBack) AddDBFunction(rollBack TemplateRollBackFunction) {
	sb.DBFunctions = append(sb.DBFunctions, rollBack)
}

func (sb *SiteRollBack) AddFileFunction(fn func(fileName string) error, fileName string) {
	fileRB := FileRollBack{FileName: fileName, Function: fn}
	sb.FileFunctions = append(sb.FileFunctions, fileRB)
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
}
