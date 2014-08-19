package controllers

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"github.com/tuomasvapaavuori/site_installer/app/modules/database"
	"log"
)

func (s *Site) CreateDatabase(templ *models.InstallTemplate) (*database.DatabaseInfo, error) {
	newDB := database.NewDatabase(&s.Base.DataStore)
	db := newDB.SetUser(&templ.MysqlUser, &templ.MysqlPassword, templ.MysqlUserHosts.Hosts).
		SetUserPrivileges(templ.MysqlUserPrivileges.Privileges, templ.MysqlGrantOption.Value).SetDBName(&templ.DatabaseName)

	err := db.CreateDatabase()

	if err != nil {
		log.Println(err, db)
		return db, err
	}

	// Rollback function.
	templ.RollBack.AddDBFunction(s.RemoveDatabase)

	err = db.CreateUser()

	if err != nil {
		log.Println(err, db)
		return db, err
	}

	// Rollback user remove function.
	templ.RollBack.AddDBFunction(s.RemoveUser)

	return db, nil
}

func (s *Site) RemoveDatabase(templ *models.InstallTemplate) error {
	return s.Base.DataStore.RemoveDatabase(templ.DatabaseName.Value)
}

func (s *Site) RemoveUser(templ *models.InstallTemplate) error {
	return s.Base.DataStore.RemoveUserOnHosts(templ.MysqlUser.Value, templ.MysqlUserHosts.Hosts)
}
