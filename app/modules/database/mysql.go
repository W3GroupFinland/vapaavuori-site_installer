package database

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"io"
	"log"
	"regexp"
	"strings"
)

type DataStore struct {
	DB *sql.DB
}

func NewDataStore() *DataStore {
	return &DataStore{}
}

func (d *DataStore) OpenConn(user string, passWord string, protocol string, host string, port string, dbName string) (*DataStore, error) {
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
		log.Println(err)
		return d, err
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return d, err
	}

	d.DB = db

	log.Println("Successfully opened database connection.")
	return d, nil
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

	log.Printf("%v rows affected.\n", affected)
	log.Printf("Created database %v.\n", name)
	return nil
}

func (d *DataStore) RemoveDatabase(name string) error {
	tx, err := d.DB.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	q := fmt.Sprintf("DROP DATABASE %v", name)

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
	log.Printf("Succesfully removed database %v.\n", name)
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
		if d.CheckUserExistsWithHost(u.Username, hostName) {
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

	log.Printf("Created user %v on hostname %v.\n", u.Username, host)
	return nil
}

func (d *DataStore) RemoveUser(tx *sql.Tx, userName string, host string) error {
	q := fmt.Sprintf("DROP USER '%v'@'%v'", userName, host)

	res, err := tx.Exec(q)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("Removed user %v on hostname %v.\n", userName, host)
	return nil
}

func (d *DataStore) RemoveUserOnHosts(userName string, hosts []string) error {
	tx, err := d.DB.Begin()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	for _, hostName := range hosts {
		err := d.RemoveUser(tx, userName, hostName)
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
		if !d.CheckUserExistsWithHost(u.Username, hostName) {
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

func (d *DataStore) CheckUserExistsWithHost(username string, host string) bool {
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

func (d *DataStore) CheckUserExists(username string) bool {
	q := "SELECT u.User FROM mysql.user u WHERE u.User = ?"
	rows, err := d.DB.Query(q, username)
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

func (di *DataStore) SqlImport(r *bufio.Reader) error {
	var (
		bytesRead int64
		query     string
	)

	// Start database transaction.
	tx, err := di.DB.Begin()
	if err != nil {
		return err
	}

	for {
		// Read string / line.
		str, err := r.ReadString(10)
		bytesRead = bytesRead + int64(len(str))

		// If error wasn't nil return error.
		if err != nil && err != io.EOF {
			return err
		}

		// If no bytes where read return from loop.
		if err == io.EOF {
			break
		}

		if str != "commit;\n" {
			// If matched line is mysql comment.
			matched, err := regexp.Match(`^[-]{2}`, []byte(str))
			if err != nil {
				return err
			}
			// If matched continue.
			if matched {
				continue
			}

			if strings.TrimSpace(str) == "" {
				continue
			}

			// If matched to line ending with x characters to ;\n
			matched, err = regexp.Match(`^.{0,}(;\n)`, []byte(str))
			if err != nil {
				return err
			}
			if !matched {
				query += str
				continue
			}
			// If matched last regex append string to query.
			query += str

			// Trim white space from query.
			query = strings.TrimSpace(query)
			// Trim query ; suffix.
			query = strings.TrimSuffix(query, ";")

			// Execute query.
			res, err := tx.Exec(query)
			if err != nil {
				return err
			}

			_, err = res.RowsAffected()
			if err != nil {
				return err
			}
			// Empty query..
			query = ""
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("SQL commit error: %v\n", err.Error())
	}

	return nil
}
