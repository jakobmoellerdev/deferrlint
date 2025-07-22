package main

import (
	"errors"
)

func main() {
	if err := deferNotOk(); err != nil {
		print(err.Error())
	}
}

func deferNotOk() error {
	var err error
	defer func() {
		err = errors.Join(err, errors.New("test")) // want "deferred function assigns to error-typed variable `err`, but it is not a named return value"
	}()

	return err
}
