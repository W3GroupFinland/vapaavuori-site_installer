package controllers

import (
	"bufio"
	"errors"
	"fmt"
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	//"io"
	"log"
	"os"
	//"os/exec"
	"path/filepath"
	"strings"
)

type Site struct {
	Drush *Drush
	Base  *a.AppBase
}

func (s *Site) Create(templ *models.InstallTemplate) (*database.DatabaseInfo, error) {
	info := &templ.InstallInfo
	var err error

	_, err = s.InstallRootStatus(info)
	if err != nil {
		log.Println(err)
		return &database.DatabaseInfo{}, err
	}

	var db *database.DatabaseInfo
	switch info.InstallType {
	case "standard":
		db, err = s.StandardInstallation(templ, info.InstallType)
	case "template":
		db, err = s.SiteTemplateInstallation(templ)
	}

	if err != nil {
		return db, err
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

func (s *Site) StandardInstallation(templ *models.InstallTemplate, installType string) (*database.DatabaseInfo, error) {
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

	_, err = s.Drush.Run("-y", "-r", info.DrupalRoot, "site-install", installType, mysqlStr, siteNameStr, subDirStr)
	if err != nil {
		return db, err
	}

	return db, nil
}

func (s *Site) SiteTemplateInstallation(templ *models.InstallTemplate) (*database.DatabaseInfo, error) {
	info := templ.InstallInfo

	// TODO: Install new site.
	db, err := s.StandardInstallation(templ, "standard")
	if err != nil {
		return db, err
	}

	var (
		siteSubDirectory = filepath.Join(info.DrupalRoot, "sites", info.SubDirectory)
	)

	// TODO: Remove files folder under new created site.
	err = os.RemoveAll(filepath.Join(info.DrupalRoot, "sites", info.SubDirectory, "files"))
	if err != nil {
		log.Println(err)
		return db, err
	}
	// TODO: Copy files from given template to new site.
	ct := utils.CopyTarget{}
	err = ct.CopyDirectory(filepath.Join(info.TemplatePath, "site-files"), siteSubDirectory)
	if err != nil {
		log.Println(err)
		return db, err
	}

	// TODO: Empty database from new site.
	_, err = s.Drush.Run("-y", "-r", info.DrupalRoot, info.SubDirectory, "sql-drop")
	if err != nil {
		return db, err
	}

	ds, err := database.NewDataStore().OpenConn(
		templ.MysqlUser.Value,
		templ.MysqlPassword.Value,
		s.Base.Config.Mysql.Protocol,
		s.Base.Config.Mysql.Host,
		s.Base.Config.Mysql.Port,
		templ.DatabaseName.Value,
	)

	defer ds.DB.Close()

	if err != nil {
		log.Println(err)
		return db, err
	}

	// Open input file.
	fi, err := os.Open(filepath.Join(info.TemplatePath, "db.sql"))
	if err != nil {
		log.Println(err)
		return db, err
	}

	r := bufio.NewReader(fi)

	err = ds.SqlImport(r)
	if err != nil {
		log.Println(err)
		return db, err
	}

	/*// TODO: Import template database to site database.
	cmd := exec.Command("drush", "-y", "-r", info.DrupalRoot, info.SubDirectory, "sqlc")
	writeCloser, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
		return db, err
	}

	// Close fi on exit and check for its returned error.
	defer func() (*database.DatabaseInfo, error) {
		if err := fi.Close(); err != nil {
			log.Println(err)
			return db, err
		}

		if err := writeCloser.Close(); err != nil {
			log.Println(err)
			return db, err
		}

		return db, nil
	}()

	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return db, err
	}

	fiInfo, err := fi.Stat()
	if err != nil {
		log.Println(err)
		return db, err
	}
	fileLength := fiInfo.Size()

	// Make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	var bytesRead int64
	var i int
	for {
		// read a chunk
		n, err := r.Read(buf)
		bytesRead = bytesRead + int64(len(buf))
		if i == 1000 {
			log.Printf("%d bytes read of %d", bytesRead, fileLength)
			i = 0
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// Write a chunk
		if _, err := writeCloser.Write(buf[:n]); err != nil {
			panic(err)
		}

		i++
	}*/

	err = utils.ChownRecursive(filepath.Join(info.DrupalRoot, "sites", info.SubDirectory, "private"), info.HttpUser, info.HttpGroup)
	if err != nil {
		log.Println(err)
		return db, err
	}

	err = utils.ChownRecursive(filepath.Join(info.DrupalRoot, "sites", info.SubDirectory, "files"), info.HttpUser, info.HttpGroup)
	if err != nil {
		log.Println(err)
		return db, err
	}

	s.Drush.VariableSet(templ, "file_private_path", filepath.Join("sites", info.SubDirectory, "private", "files"))
	s.Drush.VariableSet(templ, "file_public_path", filepath.Join("sites", info.SubDirectory, "files"))
	s.Drush.VariableSet(templ, "file_temporary_path", filepath.Join("sites", info.SubDirectory, "private", "temp"))

	log.Println("Database succesfully imported.")

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
