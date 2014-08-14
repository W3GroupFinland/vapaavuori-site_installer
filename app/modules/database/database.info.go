package database

import (
	"errors"
	"fmt"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"log"
)

const (
	UsernameLength    = 16
	PasswordLength    = UsernameLength
	DatabaseNameLenth = UsernameLength
)

type DatabaseInfo struct {
	User        *models.RandomValue
	Password    *models.RandomValue
	DbName      *models.RandomValue
	Privileges  []string
	GrantOption bool
	Hosts       []string
	DataStore   *DataStore
}

func NewDatabase(ds *DataStore) *DatabaseInfo {
	return &DatabaseInfo{
		User:        &models.RandomValue{},
		Password:    &models.RandomValue{},
		DbName:      &models.RandomValue{},
		GrantOption: false,
		DataStore:   ds,
	}
}

// This is pretty random.
func (di *DatabaseInfo) Randomize() *DatabaseInfo {
	di.User.Random = true
	di.Password.Random = true
	di.DbName.Random = true

	// When unused user name is found break.
	for {
		di.User.Randomize(UsernameLength)
		if !di.DataStore.CheckUserExists(di.User.Value) {
			break
		}
	}
	// When unused database name is found break.
	for {
		di.DbName.Randomize(DatabaseNameLenth)
		if !di.DataStore.CheckDatabaseExists(di.DbName.Value) {
			break
		}
	}
	di.Password.Randomize(PasswordLength)

	return di
}

func (di *DatabaseInfo) SetUser(user *models.RandomValue, pass *models.RandomValue, hosts []string) *DatabaseInfo {
	// Set given user info to Database info.
	di.User = user
	di.Password = pass

	if di.User.Random {
		di.User.Randomize(UsernameLength)
	}
	if di.Password.Random {
		di.Password.Randomize(PasswordLength)
	}

	di.Hosts = hosts

	return di
}

func (di *DatabaseInfo) SetUserPrivileges(privileges []string, grantOption bool) *DatabaseInfo {
	di.Privileges = privileges
	di.GrantOption = grantOption

	return di
}

func (di *DatabaseInfo) SetHosts(hosts []string) *DatabaseInfo {
	di.Hosts = hosts

	return di
}

func (di *DatabaseInfo) SetDBName(dbName *models.RandomValue) *DatabaseInfo {
	// Set given Database name to Database info.
	di.DbName = dbName

	if di.DbName.Random {
		di.DbName.Randomize(DatabaseNameLenth)
	}

	return di
}

func (di *DatabaseInfo) CreateDatabase() (*DatabaseInfo, error) {
	// Check all values have been given.
	if di.User.Value == "" || di.Password.Value == "" || di.DbName.Value == "" || len(di.Privileges) == 0 || len(di.Hosts) == 0 {

		msg := "Can't create database with following info:\n"
		msg += "Username: %v\n"
		msg += "Password: %v\n"
		msg += "Hosts: %v\n"
		msg += "Database name: %v\n"
		msg += "Privileges: %v\n"
		msg += "With grant option: %v\n\n"

		priv := fmt.Sprint(di.Privileges)
		hosts := fmt.Sprint(di.Hosts)

		log.Printf(msg, di.User.Value, di.Password.Value, hosts, di.DbName.Value, priv, di.GrantOption)
		return di, errors.New("Database creating failed.")
	}

	// Create database.
	err := di.DataStore.CreateDatabase(di.DbName.Value)
	if err != nil {
		log.Println(err)
		return di, err
	}

	// Create User on hosts.
	u := models.User{Username: di.User.Value, Password: di.Password.Value}
	err = di.DataStore.CreateUserOnHosts(&u, di.Hosts)
	if err != nil {
		log.Println(err)
		return di, err
	}

	// Grant user privileges on hosts.
	err = di.DataStore.GrantUserPrivilegesOnHosts(&u, di.DbName.Value, di.Hosts, di.Privileges, di.GrantOption)
	if err != nil {
		log.Println(err)
		return di, err
	}

	msg := "Created database with following info:\n"
	msg += "Username: %v\n"
	msg += "Password: %v\n"
	msg += "Hosts: %v\n"
	msg += "Database name: %v\n"
	msg += "Privileges: %v\n"
	msg += "With grant option: %v\n\n"

	priv := fmt.Sprint(di.Privileges)
	hosts := fmt.Sprint(di.Hosts)

	log.Printf(msg, di.User.Value, di.Password.Value, hosts, di.DbName.Value, priv, di.GrantOption)
	return di, nil
}
