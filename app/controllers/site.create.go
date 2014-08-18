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

	// Install new site.
	db, err := s.StandardInstallation(templ, "standard")
	if err != nil {
		return db, err
	}

	var (
		siteSubDirectory = filepath.Join(info.DrupalRoot, "sites", info.SubDirectory)
	)

	// Remove files folder under new created site.
	err = os.RemoveAll(filepath.Join(info.DrupalRoot, "sites", info.SubDirectory, "files"))
	if err != nil {
		log.Println(err)
		return db, err
	}
	// Copy files from given template to new site.
	ct := utils.CopyTarget{}
	err = ct.CopyDirectory(filepath.Join(info.TemplatePath, "site-files"), siteSubDirectory)
	if err != nil {
		log.Println(err)
		return db, err
	}

	// Empty database from new site.
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

	if templ.SSLServer.ServerName != "" {
		s.Drush.VariableSet(templ, "securepages_basepath_ssl", "//"+templ.SSLServer.ServerName)
	}

	log.Println("Database succesfully imported.")

	return db, nil
}

func (s *Site) GetSiteTemplateDomains(templ *models.InstallTemplate) *models.SiteDomains {
	domains := models.NewSiteDomains()

	// Get domains
	domains.SetDomain(templ.InstallInfo.ServerName)
	domains.SetDomain(templ.HttpServer.ServerName)
	domains.SetDomain(templ.SSLServer.ServerName)

	for _, domain := range templ.HttpServer.ServerAliases {
		domains.SetDomain(domain)
	}
	for _, domain := range templ.SSLServer.ServerAliases {
		domains.SetDomain(domain)
	}

	domains.SubDirectory = templ.InstallInfo.SubDirectory
	domains.SiteName = templ.InstallInfo.SiteName

	return domains
}

func (s *Site) CreateDomainSymlinks(templ *models.InstallTemplate, domains *models.SiteDomains) {
	for _, domain := range domains.Domains {
		pathToSubDir := filepath.Join(templ.InstallInfo.DrupalRoot, "sites", templ.InstallInfo.SubDirectory)
		pathToDomain := filepath.Join(templ.InstallInfo.DrupalRoot, "sites", domain)

		if _, err := os.Stat(pathToDomain); err != nil {
			if os.IsNotExist(err) {
				err := os.Symlink(pathToSubDir, pathToDomain)
				if err != nil {
					log.Printf("Error creating symlink: %v\n", err.Error())
				}
			}
		}
	}
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

func (s *Site) AddToHosts(templ *models.InstallTemplate, domains *models.SiteDomains) error {
	fi, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	var (
		domainStr    string
		i            = 0
		totalDomains = len(domains.Domains)
	)

	for _, domain := range domains.Domains {
		if i == (totalDomains - 1) {
			domainStr += domain
			break
		}

		domainStr += domain + " "
		i++
	}

	str := fmt.Sprintf("%v %v\n", "127.0.0.1", domainStr)
	_, err = fi.WriteString(str)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
