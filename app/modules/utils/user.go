package utils

import (
	"os/user"
)

func UserExists(username string) bool {
	_, err := user.Lookup(username)

	if err != nil {
		return false
	}

	return true
}
