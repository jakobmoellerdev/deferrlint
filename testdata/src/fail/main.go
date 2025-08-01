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
		err = errors.Join(err, errors.New("test")) // want "deferred function assigns to error \"err\", which is not a named return – this assignment will not affect the function's return value"
	}()

	return err
}
