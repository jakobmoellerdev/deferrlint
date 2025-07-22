// Package main provides the command-line interface for the deferrlint static analysis tool.
// Run this tool to check Go code for incorrect usage of defer statements and error assignments.
// See README.md for installation and usage instructions.

package main

import (
	"github.com/jakobmoellerdev/deferrlint/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
