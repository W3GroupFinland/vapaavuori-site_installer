package controllers

import (
	"errors"
	"fmt"
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
	"log"
	"os"
	"strings"
)

type Site struct {
	Drush *Drush
	Base  *a.AppBase
}

func (s *Site) Create(templ *models.InstallTemplate) (*database.DatabaseInfo, error) {
	info := &templ.InstallInfo

	_, err := s.InstallRootStatus(info)
	if err != nil {
		log.Println(err)
		return &database.DatabaseInfo{}, err
	}

	var db *database.DatabaseInfo
	switch info.InstallType {
	case "standard":
		db, err = s.StandardInstallation(templ)
	case "archive-dump":
		//db, err := s.ArchiveDumpInstallation(templ)
	}

	return db, nil
}

func (s *Site) CreateDatabase(templ *models.InstallTemplate) (*database.DatabaseInfo, error) {
	newDB := database.NewDatabase(&s.Base.DataStore)
	db, err := newDB.SetUser(&templ.MysqlUser, &templ.MysqlPassword, templ.MysqlUserHosts.Hosts).
		SetUserPrivileges(templ.MysqlUserPrivileges.Privileges, templ.MysqlGrantOption.Value).SetDBName(&templ.DatabaseName).
		CreateDatabase()

	if err != nil {
		log.Println(err, db)
		return db, err
	}

	return db, err
}

func (s *Site) StandardInstallation(templ *models.InstallTemplate) (*database.DatabaseInfo, error) {
	db, err := s.CreateDatabase(templ)
	if err != nil {
		return db, err
	}

	info := templ.InstallInfo

	var (
		mysqlStr    = s.Drush.FormatDatabaseStr(db)
		siteNameStr = s.Drush.FormatSiteNameStr(&info)
		subDirStr   = s.Drush.FormatSiteSubDirStr(&info)
	)

	_, err = s.Drush.Run("-y", "-r", info.DrupalRoot, "site-install", info.InstallType, mysqlStr, siteNameStr, subDirStr)
	if err != nil {
		log.Println(err)
		return db, err
	}

	return db, nil
}

func (s *Site) ArchiveDumpInstallation(templ *models.InstallTemplate) (*database.DatabaseInfo, error) {
	db, err := s.CreateDatabase(templ)
	if err != nil {
		return db, err
	}

	info := templ.InstallInfo

	var (
		mysqlStr    = s.Drush.FormatDatabaseStr(db)
		siteNameStr = s.Drush.FormatSiteNameStr(&info)
		subDirStr   = s.Drush.FormatSiteSubDirStr(&info)
	)

	// drush archive-restore [new-backup-archive] [site] --destination=./[new-folder-name] --db-url=mysql://[msql-user]:[mysql-password]@[target-server]/[db-name]
	_, err = s.Drush.Run("-y", "-r", info.DrupalRoot, "archive-restore", info.InstallType, mysqlStr, siteNameStr, subDirStr)
	if err != nil {
		log.Println(err)
		return db, err
	}

	return db, nil
}

func (s *Site) InstallRootStatus(info *models.SiteInstallConfig) (*models.SiteRootInfo, error) {
	out, err := s.Drush.Run("-r", info.DrupalRoot, "status")
	if err != nil {
		log.Println(err)
		return &models.SiteRootInfo{}, err
	}

	var rows []string
	var rowRunes []rune
	for _, r := range out {
		if r == 10 {
			rows = append(rows, string(rowRunes))
			rowRunes = rowRunes[:0]
			continue
		}

		rowRunes = append(rowRunes, r)
	}

	statusMap := make(map[string]string)
	for _, row := range rows {
		parts := strings.Split(row, ":")
		// Trim whitespace from parts.
		for idx, part := range parts {
			parts[idx] = strings.TrimSpace(part)
		}
		if len(parts) == 2 {
			statusMap[parts[0]] = parts[1]
		}
	}

	rootInfo := models.SiteRootInfo{}
	if val, ok := statusMap["Drupal version"]; ok {
		rootInfo.DrupalVersion = val
	}
	if val, ok := statusMap["Default theme"]; ok {
		rootInfo.DefaultTheme = val
	}
	if val, ok := statusMap["Administration theme"]; ok {
		rootInfo.AdministrationTheme = val
	}
	if val, ok := statusMap["PHP configuration"]; ok {
		rootInfo.PHPConfig = val
	}
	if val, ok := statusMap["PHP OS"]; ok {
		rootInfo.PHPOs = val
	}
	if val, ok := statusMap["Drush configuration"]; ok {
		rootInfo.DrushConfiguration = val
	}
	if val, ok := statusMap["Drush version"]; ok {
		rootInfo.DrushVersion = val
	}
	if val, ok := statusMap["Drush alias files"]; ok {
		rootInfo.DrushAliasFiles = val
	}
	if val, ok := statusMap["Drupal root"]; ok {
		rootInfo.DrupalRoot = val
	}

	if rootInfo.DrupalRoot != info.DrupalRoot && rootInfo.DrupalVersion == "" {
		msg := fmt.Sprintf("No drupal installation found from path %v.\n", info.DrupalRoot)
		log.Print(msg)
		return &models.SiteRootInfo{}, errors.New(msg)
	}

	return &rootInfo, nil
}

func (s *Site) AddToHosts(templ *models.InstallTemplate) error {
	fi, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	str := fmt.Sprintf("%v %v\n", "127.0.0.1", templ.InstallInfo.SiteName)
	_, err = fi.WriteString(str)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
