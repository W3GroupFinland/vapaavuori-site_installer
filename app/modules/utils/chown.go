package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

func ChownRecursive(path string, username string, group string) error {
	log.Println(group)
	u, err := user.Lookup(username)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(u.Uid)
	gid, err := strconv.Atoi(u.Gid)

	if err != nil {
		msg := fmt.Sprintf("Chown: User %v doesn't exist.\n", username)
		return errors.New(msg)
	}

	rc := RecursiveChown{User: uid, Group: gid}
	log.Printf("Chown directory recursively %v for user %v.\n", path, username)
	err = rc.chownRecursive(path)
	return err
}

type RecursiveChown struct {
	User  int
	Group int
}

func (rc *RecursiveChown) chownRecursive(root string) error {
	err := filepath.Walk(root, rc.chownWalkFunc)
	if err != nil {
		return err
	}

	return nil
}

// Walk function to recursively copy files and directories in given path.
func (rc *RecursiveChown) chownWalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	err = os.Chown(path, rc.User, rc.Group)

	return err
}
