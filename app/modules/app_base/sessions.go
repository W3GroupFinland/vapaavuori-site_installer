package app_base

import (
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
	"reflect"
)

type SessionKey struct {
	*sessions.Session
}

func (a *AppBase) InitSessions(secret string) {
	a.Sessions = sessions.NewCookieStore([]byte(secret))
}

func (a *AppBase) SetSessionKey(key string, data string, rw http.ResponseWriter, req *http.Request) {
	if data != "" {
		session, _ := a.Sessions.Get(req, key)
		session.Values[key] = data
		session.Save(req, rw)
	}
}

func (a *AppBase) InvalidateSessionKey(key string, rw http.ResponseWriter, req *http.Request) {
	http.SetCookie(rw, &http.Cookie{Name: key, MaxAge: -1, Path: "/"})
}

func (a *AppBase) GetSessionKey(key string, rw http.ResponseWriter, req *http.Request) (*SessionKey, error) {
	session, err := a.Sessions.Get(req, key)

	return &SessionKey{session}, err
}

func (sk *SessionKey) ToString(key string) (string, error) {
	var (
		i      interface{}
		exists bool
	)

	if i, exists = sk.Values[key]; !exists {
		return "", errors.New("Session value doesn't exist.")
	}

	typ := reflect.TypeOf(i)

	if typ.Name() != "string" {
		return "", errors.New("Failed to convert session value to string.")
	}

	val := reflect.ValueOf(i)

	return val.String(), nil
}
