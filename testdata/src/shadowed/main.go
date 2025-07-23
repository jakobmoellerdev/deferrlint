package main

import (
	"errors"
)

func main() {
	if err := shadowedErr(); err != nil {
		print(err.Error())
	}
}

func shadowedErr() error {
	defer func() {
		err := errors.New("shadowed") // no lint (new variable)
		_ = err
	}()
	return nil
}
