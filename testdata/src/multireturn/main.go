package main

import (
	"errors"
)

func main() {
	if _, err := multiReturns(); err != nil {
		print(err.Error())
	}
	if _, err := wrongNamed(); err != nil {
		print(err.Error())
	}
}

func multiReturns() (val int, err error) {
	defer func() {
		err = errors.Join(err, errors.New("test")) // no lint
	}()
	return 42, err
}

func wrongNamed() (int, error) {
	var err error
	defer func() {
		err = errors.Join(err, errors.New("test")) // want "deferred function assigns to error \"err\", which is not a named return – this assignment will not affect the function's return value"
	}()
	return 0, err
}
