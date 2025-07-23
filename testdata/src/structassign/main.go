package main

import (
	"errors"
)

func main() {
	if err := structField(); err != nil {
		print(err.Error())
	}
}

type wrapper struct {
	err error
}

func structField() error {
	w := wrapper{}
	defer func() {
		w.err = errors.New("ok") // no lint
	}()
	return nil
}
