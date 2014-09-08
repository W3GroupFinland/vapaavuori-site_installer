package controllers

import (
	"bufio"
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	a "github.com/tuomasvapaavuori/site_installer/app/modules/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Site struct {
	Drush *Drush
	Base  *a.AppBase
}

func (s *Site) Create(tmpl *models.InstallTemplate, sp *models.SubProcess) (*database.DatabaseInfo, error) {
	info := &tmpl.InstallInfo
	var err error

	exists, _, err := s.InstallRootStatus(info.DrupalRoot)
	if err != nil {
		log.Println(err)
		return &database.DatabaseInfo{}, err
	}

	if !exists {
		return &database.DatabaseInfo{}, errors.New("No platform in given path.")
	}

	var db *database.DatabaseInfo
	switch info.InstallType {
	case "standard":
		db, err = s.StandardInstallation(tmpl, info.InstallType, sp)
	case "template":
		db, err = s.SiteTemplateInstallation(tmpl, sp)
	}

	if err != nil {
		return db, err
	}

	return db, nil
}

func (s *Site) StandardInstallation(tmpl *models.InstallTemplate, installType string, sp *models.SubProcess) (*database.DatabaseInfo, error) {
	db, err := s.CreateDatabase(tmpl)
	if err != nil {
		return db, err
	}

	info := tmpl.InstallInfo

	var (
		mysqlStr    = s.Drush.FormatDatabaseStr(db)
		siteNameStr = s.Drush.FormatSiteNameStr(&info)
		subDirStr   = s.Drush.FormatSiteSubDirStr(&info)
	)

	_, err = s.Drush.Run("-y", "-r", info.DrupalRoot, "site-install", installType, mysqlStr, siteNameStr, subDirStr)
	if err != nil {
		return db, err
	}

	// Rollback: Remove install directory.
	tmpl.RollBack.AddFileFunction(
		utils.RemoveDirectory,
		filepath.Join(info.DrupalRoot, "sites", info.SubDirectory))

	// Chown install subdirectory to deploy user.
	err = utils.ChownRecursive(filepath.Join(info.DrupalRoot, "sites", info.SubDirectory),
		s.Base.Config.DeployUser.User, s.Base.Config.DeployUser.Group)
	if err != nil {
		log.Println(err)
		return db, err
	}

	return db, nil
}

func (s *Site) SiteTemplateInstallation(tmpl *models.InstallTemplate, sp *models.SubProcess) (*database.DatabaseInfo, error) {
	sp.Start()

	info := tmpl.InstallInfo

	// Install new site.
	db, err := s.StandardInstallation(tmpl, "standard", sp)
	if err != nil {
		return db, err
	}

	var (
		siteSubDirectory = filepath.Join(info.DrupalRoot, "sites", info.SubDirectory)
	)

	sp.Update("Remove files from new created site.")
	// Remove files folder under new created site.
	err = os.RemoveAll(filepath.Join(info.DrupalRoot, "sites", info.SubDirectory, "files"))
	if err != nil {
		log.Println(err)
		return db, err
	}

	sp.Update("Copy template files to new site.")
	// Copy files from given template to new site.
	ct := utils.CopyTarget{}
	err = ct.CopyDirectory(filepath.Join(info.TemplatePath, "site-files"), siteSubDirectory)
	if err != nil {
		log.Println(err)
		return db, err
	}

	sp.Update("Empty database from new site.")
	// Empty database from new site.
	_, err = s.Drush.Run("-y", "-r", info.DrupalRoot, info.SubDirectory, "sql-drop")
	if err != nil {
		return db, err
	}

	ds, err := database.NewDataStore().OpenConn(
		tmpl.MysqlUser.Value,
		tmpl.MysqlPassword.Value,
		s.Base.Config.Mysql.Protocol,
		s.Base.Config.Mysql.Host,
		s.Base.Config.Mysql.Port,
		tmpl.DatabaseName.Value,
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

	sp.Update("Importing template database.")
	err = ds.SqlImport(r)
	if err != nil {
		log.Println(err)
		return db, err
	}

	sp.Update("Chown web user folders.")

	// Chown install subdirectory to deploy user.
	err = utils.ChownRecursive(filepath.Join(info.DrupalRoot, "sites", info.SubDirectory),
		s.Base.Config.DeployUser.User, s.Base.Config.DeployUser.Group)
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

	sp.Update("Setting site directory variables with drush.")
	err = s.Drush.VariableSet(tmpl, "file_private_path", filepath.Join("sites", info.SubDirectory, "private", "files"))
	err = s.Drush.VariableSet(tmpl, "file_public_path", filepath.Join("sites", info.SubDirectory, "files"))
	err = s.Drush.VariableSet(tmpl, "file_temporary_path", filepath.Join("sites", info.SubDirectory, "private", "temp"))

	if err != nil {
		log.Printf("Error: Failed to set variables with drush. Error message: %v\n", err.Error())
		return db, err
	}

	if tmpl.SSLServer.DomainInfo.DomainName != "" {
		s.Drush.VariableSet(tmpl, "securepages_basepath_ssl", "//"+tmpl.SSLServer.DomainInfo.DomainName)
	}

	log.Println("Database succesfully imported.")

	sp.Finish()
	return db, nil
}

func (s *Site) GetSiteTemplateDomains(tmpl *models.InstallTemplate) *models.SiteDomains {
	domains := models.NewSiteDomains()

	if tmpl.HttpServer.Include {
		// Set types to domains so they can be later identified.
		tmpl.HttpServer.DomainInfo.Type = models.DomainTypeServerName

		// Get domains
		domains.SetDomain(tmpl.HttpServer.DomainInfo)

		for _, domain := range tmpl.HttpServer.DomainAliases {
			domain.Type = models.DomainTypeServerAlias
			domains.SetDomain(domain)
		}
	}

	if tmpl.SSLServer.Include {
		// Set types to domains so they can be later identified.
		tmpl.SSLServer.DomainInfo.Type = models.DomainTypeServerName

		// Get domains
		domains.SetDomain(tmpl.SSLServer.DomainInfo)

		for _, domain := range tmpl.SSLServer.DomainAliases {
			domain.Type = models.DomainTypeServerAlias
			domains.SetDomain(domain)
		}
	}

	domains.SubDirectory = tmpl.InstallInfo.SubDirectory
	domains.SiteName = tmpl.InstallInfo.SiteName

	return domains
}

func (s *Site) CreateDomainSymlinks(tmpl *models.InstallTemplate, domains *models.SiteDomains, sp *models.SubProcess) {
	sp.Start()

	for _, domain := range domains.Domains {
		pathToSubDir := filepath.Join(tmpl.InstallInfo.DrupalRoot, "sites", tmpl.InstallInfo.SubDirectory)
		pathToDomain := filepath.Join(tmpl.InstallInfo.DrupalRoot, "sites", domain.DomainName)

		if _, err := os.Stat(pathToDomain); err != nil {
			if os.IsNotExist(err) {
				err := os.Symlink(pathToSubDir, pathToDomain)
				if err != nil {
					log.Printf("Error creating symlink: %v\n", err.Error())
				}

				err = utils.ChownRecursive(pathToDomain, s.Base.Config.DeployUser.User, s.Base.Config.DeployUser.Group)
				if err != nil {
					log.Println(err)
				}
				tmpl.RollBack.AddFileFunction(utils.RemoveFile, pathToDomain)
			}
		}
	}

	sp.Finish()
}

func (s *Site) InstallRootStatus(path string) (bool, *models.SiteRootInfo, error) {
	out, err := s.Drush.Run("-r", path, "status")
	if err != nil {
		log.Println(err)
		return false, &models.SiteRootInfo{}, err
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

	if rootInfo.DrupalRoot != path && rootInfo.DrupalVersion == "" {
		return false, &models.SiteRootInfo{}, nil
	}

	return true, &rootInfo, nil
}
