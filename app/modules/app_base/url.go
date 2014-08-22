package app_base

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Url struct {
	Url         *url.URL
	QueryValues url.Values
}

type QueryValues struct {
	Param  string
	Values []string
}

func RequestURIQuery(rawUrl string) (*Url, error) {
	p, err := url.ParseRequestURI(rawUrl)

	if err != nil {
		return &Url{}, err
	}

	u := Url{Url: p, QueryValues: p.Query()}

	return &u, nil
}

func (u *Url) GetParam(param string) (*QueryValues, error) {
	var (
		values []string
		exists bool
	)
	if values, exists = u.QueryValues[param]; !exists {
		msg := fmt.Sprintf("Query parameter %v doesn't exist.", param)
		return &QueryValues{}, errors.New(msg)
	}

	return &QueryValues{
		Values: values,
		Param:  param,
	}, nil
}

func (q *QueryValues) All() []string {
	return q.Values
}

func (q *QueryValues) GetKeyError(key int) error {
	if len(q.Values) < key {
		errMsg := fmt.Sprintf("Param %v value at %v doesn't exist.\n", q.Param, key)
		return errors.New(errMsg)
	}

	return nil
}

func (q *QueryValues) GetKey(key int) (string, error) {
	if err := q.GetKeyError(key); err != nil {
		return "", err
	}

	return q.Values[key], nil
}

func (q *QueryValues) GetKeyInt(key int) (int, error) {
	if err := q.GetKeyError(key); err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(q.Values[key])
	if err != nil {
		return val, err
	}

	return val, err
}

func (q *QueryValues) GetKeyInt64(key int) (int64, error) {
	val, err := q.GetKeyInt(key)
	if err != nil {
		return int64(val), err
	}

	return int64(val), err
}
