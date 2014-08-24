package controllers

import (
	"database/sql"
	"errors"
	"github.com/tuomasvapaavuori/site_installer/app/models"
	a "github.com/tuomasvapaavuori/site_installer/app/modules/app_base"
	"log"
	"net/http"
)

type User struct {
	Base *a.AppBase
}

func (c *User) Init() {
}

func (c *User) ControllerName() string {
	return "app.controllers.web.user"
}

func (c *User) LoggedInAcl(rw http.ResponseWriter, r *http.Request) bool {
	if _, ok := c.Current(rw, r); !ok {
		return false
	}

	return true
}

func (c *User) UpdateHandler(rw http.ResponseWriter, r *http.Request) {
	_, ok := c.Current(rw, r)
	if !ok {
		http.Error(rw, c.Base.Http.Error(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
}

func (c *User) GetHandler(rw http.ResponseWriter, r *http.Request) {
	user, ok := c.Current(rw, r)
	if !ok {
		http.Error(rw, c.Base.Http.Error(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	encoder := c.Base.JSON.NewEncoder(rw)
	encoder.Encode(models.UserSend{
		Username: user.Username,
		Mail:     user.Mail,
		Status:   user.Status,
	})
}

func (c *User) Current(rw http.ResponseWriter, r *http.Request) (*models.User, bool) {
	key := "client-logged"
	session, err := c.Base.GetSessionKey(key, rw, r)

	if err != nil {
		return &models.User{}, false
	}

	var (
		value string
	)

	value, err = session.ToString(key)

	if err != nil {
		return &models.User{}, false
	}

	user, err := c.Load(value)

	if err != nil {
		log.Printf("Error loading user %v.", value)
		return &models.User{}, false
	}

	return user, true
}

func (c *User) Load(username string) (*models.User, error) {
	q := "SELECT u.id, u.username, u.mail, u.password, u.status FROM user u "
	q += "WHERE u.username = ?"
	row := c.Base.DataStore.DB.QueryRow(q, username)

	user := models.User{}
	err := row.Scan(
		&user.Uid,
		&user.Username,
		&user.Mail,
		&user.Password,
		&user.Status,
	)

	if err == sql.ErrNoRows {
		return &user, models.NoUserFoundError
	}

	return &user, err
}

func (c *User) Create(username string, mail string, password string, status bool) (*models.User, error) {
	exists, err := c.UsernameExists(username)
	if err != nil {
		return &models.User{}, err
	}

	if exists {
		return &models.User{}, errors.New("Username exists already.")
	}

	exists, err = c.MailExists(mail)
	if err != nil {
		return &models.User{}, err
	}

	if exists {
		return &models.User{}, errors.New("User mail exists already.")
	}

	q := "INSERT INTO user (username, mail, password, status) VALUES(?, ?, ?, ?)"
	res, err := c.Base.DataStore.DB.Exec(q,
		username,
		mail,
		password,
		status,
	)

	if err != nil {
		return &models.User{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return &models.User{}, err
	}

	return &models.User{Uid: id,
		Username: username,
		Password: password,
		Status:   status}, err
}

func (c *User) UsernameExists(username string) (bool, error) {
	q := "SELECT id FROM user u WHERE u.username = ?"

	var id int64
	row := c.Base.DataStore.DB.QueryRow(q, username)
	err := row.Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (c *User) MailExists(mail string) (bool, error) {
	q := "SELECT id FROM user u WHERE u.mail = ?"

	var id int64
	row := c.Base.DataStore.DB.QueryRow(q, mail)
	err := row.Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (c *User) LoadUserWithId(uid int64) (*models.User, error) {
	var username string
	q := "SELECT u.username FROM user u WHERE u.id = ?"
	row := c.Base.DataStore.DB.QueryRow(q, uid)

	err := row.Scan(&username)
	if err != nil {
		return &models.User{}, err
	}

	user, _ := c.Load(username)

	return user, nil
}

func (c *User) Update(user *models.User) error {
	_, err := c.Load(user.Username)
	if err != nil {
		return err
	}

	_, err = c.Base.DataStore.DB.Exec("UPDATE user SET user.mail=?, user.password=? WHERE user.username=?",
		user.Mail,
		user.Password,
		user.Username)

	return err
}

func (c *User) Delete(username string) error {
	_, err := c.Base.DataStore.DB.Exec("DELETE user WHERE user.username=?", username)

	return err
}
