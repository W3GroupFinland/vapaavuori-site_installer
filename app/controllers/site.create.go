package controllers

import (
	"fmt"
	a "github.com/tuomasvapaavuori/site_installer/app/app_base"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
	"log"
)

type Site struct {
	Drush *Drush
	Base  *a.AppBase
}

func (s *Site) Create(info *models.SiteInstallInfo) *database.DatabaseInfo {
	newDB := database.NewDatabase(&s.Base.DataStore)
	db, err := newDB.Randomize().
		SetHosts([]string{"localhost", "127.0.0.1"}).
		SetUserPrivileges([]string{"ALL"}, true).
		CreateDatabase()

	if err != nil {
		log.Println(err, db)
	}

	mysqlStr := fmt.Sprintf("--db-url=mysql://%v:%v@%v:%v/%v",
		db.User.Value,
		db.Password.Value,
		s.Base.Config.Mysql.Host,
		s.Base.Config.Mysql.Port,
		db.DbName.Value)

	siteNameStr := fmt.Sprintf("--site-name=%v", info.SiteName)
	subDirStr := fmt.Sprintf("--sites-subdir=%v", info.SubDirectory)

	out, err := s.Drush.Run("-y", "-r", info.InstallRoot, "site-install", info.InstallType, mysqlStr, siteNameStr, subDirStr)
	if err != nil {
		log.Println(err)
	}
	log.Println(out)

	return db
}
