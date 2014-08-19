package controllers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
	"github.com/tuomasvapaavuori/site_installer/app/modules/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
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

	// Rollback: Remove install directory.
	templ.RollBack.AddFileFunction(
		utils.RemoveDirectory,
		filepath.Join(info.DrupalRoot, "sites", info.SubDirectory))

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

	err = s.Drush.VariableSet(templ, "file_private_path", filepath.Join("sites", info.SubDirectory, "private", "files"))
	err = s.Drush.VariableSet(templ, "file_public_path", filepath.Join("sites", info.SubDirectory, "files"))
	err = s.Drush.VariableSet(templ, "file_temporary_path", filepath.Join("sites", info.SubDirectory, "private", "temp"))

	if err != nil {
		log.Printf("Error: Failed to set variables with drush. Error message: %v\n", err.Error())
		return db, err
	}

	if templ.SSLServer.DomainInfo.DomainName != "" {
		s.Drush.VariableSet(templ, "securepages_basepath_ssl", "//"+templ.SSLServer.DomainInfo.DomainName)
	}

	log.Println("Database succesfully imported.")

	return db, nil
}

func (s *Site) GetSiteTemplateDomains(templ *models.InstallTemplate) *models.SiteDomains {
	domains := models.NewSiteDomains()

	// Get domains
	domains.SetDomain(templ.InstallInfo.DomainInfo)
	domains.SetDomain(templ.HttpServer.DomainInfo)
	domains.SetDomain(templ.SSLServer.DomainInfo)

	for _, domain := range templ.HttpServer.DomainAliases {
		domains.SetDomain(domain)
	}
	for _, domain := range templ.SSLServer.DomainAliases {
		domains.SetDomain(domain)
	}

	domains.SubDirectory = templ.InstallInfo.SubDirectory
	domains.SiteName = templ.InstallInfo.SiteName

	return domains
}

func (s *Site) CreateDomainSymlinks(templ *models.InstallTemplate, domains *models.SiteDomains) {
	for _, domain := range domains.Domains {
		pathToSubDir := filepath.Join(templ.InstallInfo.DrupalRoot, "sites", templ.InstallInfo.SubDirectory)
		pathToDomain := filepath.Join(templ.InstallInfo.DrupalRoot, "sites", domain.DomainName)

		if _, err := os.Stat(pathToDomain); err != nil {
			if os.IsNotExist(err) {
				err := os.Symlink(pathToSubDir, pathToDomain)
				if err != nil {
					log.Printf("Error creating symlink: %v\n", err.Error())
				}
				templ.RollBack.AddFileFunction(utils.RemoveFile, pathToDomain)
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

func (s *Site) ReadHostsFile(hostsFile string) (models.HostsMap, error) {
	fi, err := os.Open(hostsFile)

	if err != nil {
		log.Println(err)
		return models.HostsMap{}, err
	}

	defer fi.Close()

	return s.ParseHostsContent(fi)
}

func (s *Site) ParseHostsContent(content io.Reader) (models.HostsMap, error) {
	hostsMap := models.NewHostsMap()
	var (
		// Make a read buffer
		r         = bufio.NewReader(content)
		readState = models.READ_NOT_STARTED
	)

	for {
		// Read strings from file.
		ln, err := r.ReadString(models.NEW_LINE_BYTE)

		// Trim spaces from string.
		ln = strings.TrimSpace(ln)

		// Get application hosts read start.
		if ln == models.HOSTS_START_READ_STR {
			log.Println("Application hosts read started.")
			readState = models.READ_STARTED
			continue
		}

		// Get application hosts read end and brake.
		if ln == models.HOSTS_END_READ_STR {
			log.Println("Application hosts read ended.")
			readState = models.READ_ENDED
			break
		}

		if readState == models.READ_STARTED {
			hd := models.NewHostDomains()
			err := hd.Parse(ln)
			if err != nil {
				log.Println(err)
			}
			hostsMap.AddHostDomains(hd)
		}

		// If error wasn't nil return error.
		if err != nil && err != io.EOF {
			return hostsMap, err
		}

		// If no bytes where read return from loop.
		if err == io.EOF {
			break
		}
	}

	if readState != models.READ_ENDED {
		msg := fmt.Sprintf("No application hosts found on hosts file.\n")
		return hostsMap, errors.New(msg)
	}

	return hostsMap, nil
}

func (s *Site) WriteNewHosts(hostsFile string, hostsMap *models.HostsMap) error {
	fi, err := os.Open(hostsFile)
	if err != nil {
		log.Println(err)
		return err
	}
	r := bufio.NewReader(fi)

	var (
		// Create temporary slice of bytes.
		temp      []byte
		readState = models.READ_NOT_STARTED
	)
	for {
		// Prefixed indicates if the current line was only partially read.
		// If is prefixed we won't later add new line byte to slice.
		ln, prefixed, err := r.ReadLine()
		str := string(ln)

		// If error wasn't nil return error.
		if err != nil && err != io.EOF {
			log.Println(err)
			return err
		}

		// If no bytes where read return from loop.
		if err == io.EOF {
			break
		}

		if str == models.HOSTS_START_READ_STR {
			// Set read state started.
			readState = models.READ_STARTED
			// Write hosts map as bytes.
			temp = append(temp, hostsMap.Bytes(models.SPACE_BYTE, models.NEW_LINE_BYTE)...)
		}

		if str == models.HOSTS_END_READ_STR {
			readState = models.READ_ENDED
			// Empty the string so we don't add snippet end line twice.
			ln = ln[:0]
			str = ""
		}

		// If read started is indicated, don't write contents.
		// Hosts in snippet area are rebuild.
		if readState != models.READ_STARTED {
			if len(ln) > 0 {
				temp = append(temp, ln...)
				// If prefixed we don't add new line byte to slice.
				if prefixed {
					continue
				}

				temp = append(temp, models.NEW_LINE_BYTE)
			}
		}
	}
	ok, err := s.HostsContentAssertsTrue(bytes.NewReader(temp))
	if !ok {
		return err
	}

	_, err = utils.CreateBackupFile(hostsFile)
	err = ioutil.WriteFile(hostsFile, temp, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Site) HostsContentAssertsTrue(content io.Reader) (bool, error) {
	_, err := s.ParseHostsContent(content)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Site) AddToHosts(templ *models.InstallTemplate, domains *models.SiteDomains) error {
	hostsFile := "/etc/hosts"

	_, err := s.ReadHostsFile(hostsFile)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = os.OpenFile(hostsFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	//str := fmt.Sprintf("%v %v\n", "127.0.0.1", domainStr)
	/*_, err = fi.WriteString(str)
	if err != nil {
		log.Println(err)
		return err
	}*/

	return nil
}
