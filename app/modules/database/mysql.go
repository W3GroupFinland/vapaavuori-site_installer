package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"log"
)

type DataStore struct {
	DB *sql.DB
}

type DatabaseInfo struct {
	User      string
	Password  string
	DbName    string
	DataStore *DataStore
}

func NewDatabase(ds *DataStore) *DatabaseInfo {
	return &DatabaseInfo{DataStore: ds}
}

func (di *DatabaseInfo) RandomNames() *DatabaseInfo {
	//var userName string
	for {
		//userName = a.RandomString(16)

	}
}

func (d *DataStore) OpenConn(user string, passWord string, protocol string, host string, port string, dbName string) {
	dbString := fmt.Sprintf("%v:%v@%v(%v:%v)/%v",
		user,
		passWord,
		protocol,
		host,
		port,
		dbName)

	var (
		err error
		db  *sql.DB
	)
	// Open db connection
	db, err = sql.Open("mysql", dbString)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	d.DB = db

	log.Println("Successfully opened database connection.")
}

func (d *DataStore) CreateDatabase(name string) error {
	if d.DB == nil {
		log.Fatalln("No database connection active.")
	}

	// Check if database exists.
	if d.CheckDatabaseExists(name) {
		msg := fmt.Sprintf("Database %v already exists.\n", name)
		log.Println(msg)
		return errors.New(msg)
	}

	tx, err := d.DB.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	q := fmt.Sprintf("CREATE DATABASE %v", name)
	res, err := tx.Exec(q)
	if err != nil {
		log.Println(err)
		// Rollback
		err := tx.Rollback()
		if err != nil {
			log.Fatalln(err)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("%v rows affected.", affected)
	log.Printf("Created database %v.", name)
	return nil
}

func (d *DataStore) CreateUserOnHosts(u *models.User, hosts []string) error {
	tx, err := d.DB.Begin()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	for _, hostName := range hosts {
		// If user already exists for hostname, continue..
		if d.CheckIfUserExists(u.Username, hostName) {
			log.Printf("User %v already exists on hostname %v.", u.Username, hostName)
			continue
		}
		err := d.CreateUser(tx, u, hostName)
		if err != nil {
			log.Println(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	return nil
}

func (d *DataStore) CreateUser(tx *sql.Tx, u *models.User, host string) error {
	q := fmt.Sprintf("CREATE USER '%v'@'%v' IDENTIFIED BY '%v'", u.Username, host, u.Password)

	res, err := tx.Exec(q)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("Created user %v on hostname %v.", u.Username, host)
	return nil
}

func (d *DataStore) GrantUserPrivilegesOnHosts(
	u *models.User,
	dbName string,
	hosts []string,
	privileges []string,
	grantOpt bool) error {

	// Check database exists.
	if !d.CheckDatabaseExists(dbName) {
		msg := fmt.Sprintf("Database with name %v doesn't exist.", dbName)
		log.Println(msg)
		return errors.New(msg)
	}

	var privStr string
	privLen := len(privileges)
	for idx, priv := range privileges {
		if idx == (privLen - 1) {
			privStr += priv
			break
		}

		privStr += priv + ","
	}

	tx, err := d.DB.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, hostName := range hosts {
		// Continue if user doesn't exist.
		if !d.CheckIfUserExists(u.Username, hostName) {
			continue
		}

		q := fmt.Sprintf("GRANT %v PRIVILEGES ON %v.* TO '%v'@'%v'", privStr, dbName, u.Username, hostName)
		if grantOpt {
			q += " WITH GRANT OPTION"
		}

		_, err := tx.Exec(q)
		if err != nil {
			log.Println(err)
		}

		log.Println(q)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}

	return nil
}

func (d *DataStore) CheckIfUserExists(username string, host string) bool {
	q := "SELECT u.Host, u.User FROM mysql.user u WHERE u.User = ? AND u.Host = ?"
	rows, err := d.DB.Query(q, username, host)
	if err != nil {
		log.Println(err)
	}

	var values []bool
	for rows.Next() {
		values = append(values, true)
	}

	if len(values) > 0 {
		return true
	}

	return false
}

func (d *DataStore) CheckDatabaseExists(name string) bool {
	q := "SELECT schema_name FROM information_schema.schemata WHERE schema_name = ?"
	rows, err := d.DB.Query(q, name)
	if err != nil {
		log.Fatalln(err)
	}

	var values []bool
	for rows.Next() {
		values = append(values, true)
	}

	if len(values) > 0 {
		return true
	}

	return false
}
