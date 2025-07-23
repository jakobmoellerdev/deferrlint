package main

import (
	"errors"
)

func main() {
	loopDefer()
}

func loopDefer() {
	for i := 0; i < 1; i++ {
		var err error
		defer func() {
			err = errors.New("inside loop") // want
		}()
		_ = err
	}
}
