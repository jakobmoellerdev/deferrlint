// Package analyzer provides the deferrlint static analysis for Go code.
// It checks for incorrect usage of defer statements, specifically assignments to error-typed variables
// that are not named return values in deferred functions.
//
// This analyzer is useful for catching potential bugs where a deferred function
// attempts to assign a value to an error variable that is expected to be a named return value,
// but isn't. This can help ensure that deferred functions behave as intended,
// and that error handling is consistent and predictable.
package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Analyzer reports when a deferred function assigns to an error-typed variable
// that is not a named return value.
var Analyzer = &analysis.Analyzer{
	Name: "deferrlint",
	Doc:  "reports when a deferred function assigns to an error-typed variable that is not a named return value",
	Run:  run,
}

// run executes the analysis pass for the deferrlint analyzer.
// It inspects Go files for deferred functions that assign to error-typed variables
// which are not named return values.
func run(pass *analysis.Pass) (interface{}, error) {
	errorType := types.Universe.Lookup("error").Type()

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				return true
			}

			namedErrorObjs := collectNamedErrorReturns(fn, pass, errorType)
			checkDeferAssignments(fn.Body, pass, namedErrorObjs, errorType)
			return true
		})
	}
	return nil, nil
}

// collectNamedErrorReturns returns the set of named return variables of type `error`.
// It inspects the function declaration and collects all named return variables
// that have the error type.
func collectNamedErrorReturns(fn *ast.FuncDecl, pass *analysis.Pass, errorType types.Type) map[*types.Var]bool {
	result := make(map[*types.Var]bool)
	if fn.Type.Results == nil {
		return result
	}

	for _, field := range fn.Type.Results.List {
		for _, name := range field.Names {
			if name == nil {
				continue
			}
			obj := pass.TypesInfo.ObjectOf(name)
			if v, ok := obj.(*types.Var); ok && types.Identical(pass.TypesInfo.TypeOf(name), errorType) {
				result[v] = true
			}
		}
	}
	return result
}

// checkDeferAssignments inspects deferred functions for assignments to error-typed vars
// and reports them if they are not named return variables.
func checkDeferAssignments(body *ast.BlockStmt, pass *analysis.Pass, namedErrors map[*types.Var]bool, errorType types.Type) {
	for _, stmt := range body.List {
		deferStmt, ok := stmt.(*ast.DeferStmt)
		if !ok {
			continue
		}

		fnLit, ok := deferStmt.Call.Fun.(*ast.FuncLit)
		if !ok {
			continue
		}

		ast.Inspect(fnLit.Body, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}

			for _, lhs := range assign.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}

				obj := pass.TypesInfo.ObjectOf(ident)
				typ := pass.TypesInfo.TypeOf(ident)
				if obj == nil || typ == nil || !types.Identical(typ, errorType) {
					continue
				}

				if _, ok := obj.(*types.Var); ok {
					if !namedErrors[obj.(*types.Var)] {
						pass.Reportf(ident.Pos(), "deferred function assigns to error-typed variable `%s`, but it is not a named return value", ident.Name)
					}
				}
			}
			return true
		})
	}
}
