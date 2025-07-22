# deferrlint

A static analysis tool for Go that checks for incorrect usage of `defer` statements, specifically assignments to error-typed variables that are not named return values in deferred functions.

## Features
- Detects deferred functions that assign to error variables not declared as named return values.
- Helps avoid the common pitfall of using `defer` incorrectly with errors, which can lead to unexpected behavior in error handling.

## Installation

You can build deferrlint from source and install it locally:

_Make sure you have [Go](https://go.dev/doc/install) and [Task](https://taskfile.dev/installation/) installed on your machine._

```sh
git clone https://github.com/yourusername/deferrlint.git
cd deferrlint
task install
```

## Usage

Run deferrlint on your Go codebase:

```sh
go run main.go ./...
```

Or, if you built the binary with `task install`, you can run it directly:

```sh
deferrlint ./...
```

You can verify `deferrlint` is installed correctly by running:

```shell
deferrlint ./testdata/src/fail/...
```

## Example

```go
package main

import (
	"errors"
)

func foo() (err error) {
	defer func() {
		err = errors.New("foo") // OK: err is a named return value
	}()
	
	return err
}

func bar() error {
	var err error
	defer func() {
		err = errors.New("bar") // Not OK: err is not a named return value
	}()
	return err
}
```

`deferrlint` will report the assignment in `bar` as a potential issue.

## Contributing
Pull requests and issues are welcome!

## License

See the [LICENSE](LICENSE) file for details.

