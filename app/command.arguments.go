package app

import (
	"errors"
	"os"
)

func (a *Application) ParseCommandLineArgs() {
	args := os.Args[1:]

	a.Arguments = make(map[string]string)
	for _, argument := range args {
		bytes := []byte(argument)
		var (
			key      []byte
			value    []byte
			keyEnded bool
		)
		for _, b := range bytes {
			if !keyEnded && b == 61 {
				keyEnded = true
				continue
			}

			if !keyEnded {
				key = append(key, b)
			} else {
				value = append(value, b)
			}
		}

		a.Arguments[string(key)] = string(value)
	}
}

func (a *Application) GetCommandArg(key string) (string, error) {
	if val, ok := a.Arguments[key]; ok {
		return val, nil
	}

	return "", errors.New("Command line argument key not found.")
}
