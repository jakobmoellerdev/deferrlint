package main

import (
	"errors"
)

func main() {
	if err := deferOk(); err != nil {
		print(err.Error())
	}
}

func deferOk() (err error) {
	defer func() {
		err = errors.Join(err, errors.New("test"))
	}()

	return err
}
