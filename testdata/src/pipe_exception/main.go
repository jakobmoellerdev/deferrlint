package main

import (
	"errors"
	"io"
)

func main() {
	if err := deferNotOkByDefaultButHasAnException(); err != nil {
		print(err.Error())
	}
}

func deferNotOkByDefaultButHasAnException() error {
	_, w := io.Pipe()
	var err error
	defer func() {
		_ = w.CloseWithError(err)
	}()

	defer func() {
		err = errors.Join(err, errors.New("test")) //nolint:deferrlint
	}()

	return err
}
