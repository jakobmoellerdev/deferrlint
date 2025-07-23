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
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Analyzer reports when a deferred function assigns to an error-typed variable
// that is not a named return value.
var Analyzer = &analysis.Analyzer{
	Name: "deferrlint",
	Doc:  "reports when a deferred function assigns to an error-typed variable that is not a named return value",
	Run:  run,
	URL:  "github.com/jakobmoellerdev/deferrlint",
}

// run executes the analysis pass for the deferrlint analyzer.
// It inspects Go files for deferred functions that assign to error-typed variables
// which are not named return values.
func run(pass *analysis.Pass) (interface{}, error) {
	errorType := types.Universe.Lookup("error").Type()

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}

			namedErrors := collectNamedErrorReturns(fn, pass, errorType)
			checkDeferredAssignments(pass, fn.Body, namedErrors, errorType)
		}
	}

	return nil, nil
}

func collectNamedErrorReturns(fn *ast.FuncDecl, pass *analysis.Pass, errorType types.Type) map[*types.Var]bool {
	namedErrors := make(map[*types.Var]bool)

	if fn.Type.Results == nil {
		return namedErrors
	}

	for _, field := range fn.Type.Results.List {
		for _, name := range field.Names {
			if name == nil {
				continue
			}
			obj := pass.TypesInfo.ObjectOf(name)
			if v, ok := obj.(*types.Var); ok && types.Identical(pass.TypesInfo.TypeOf(name), errorType) {
				namedErrors[v] = true
			}
		}
	}

	return namedErrors
}

func checkDeferredAssignments(pass *analysis.Pass, body *ast.BlockStmt, namedErrors map[*types.Var]bool, errorType types.Type) {
	for _, stmt := range body.List {
		deferStmt, ok := stmt.(*ast.DeferStmt)
		if !ok {
			continue
		}

		fnLit, ok := deferStmt.Call.Fun.(*ast.FuncLit)
		if !ok {
			continue
		}

		checkAssignmentsInDefer(pass, fnLit, namedErrors, errorType)
	}
}

func checkAssignmentsInDefer(pass *analysis.Pass, fnLit *ast.FuncLit, namedErrors map[*types.Var]bool, errorType types.Type) {
	for _, stmt := range fnLit.Body.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}

		for _, lhs := range assign.Lhs {
			reportIfBadAssignment(pass, lhs, assign, namedErrors, errorType)
		}
	}
}

func reportIfBadAssignment(pass *analysis.Pass, lhs ast.Expr, assign *ast.AssignStmt, namedErrors map[*types.Var]bool, errorType types.Type) {
	ident, ok := lhs.(*ast.Ident)
	if !ok || ident.Name == "_" {
		return
	}

	// Skip := definitions of new variables
	if assign.Tok == token.DEFINE {
		if _, isDef := pass.TypesInfo.Defs[ident]; isDef {
			return
		}
	}

	obj := pass.TypesInfo.ObjectOf(ident)
	typ := pass.TypesInfo.TypeOf(ident)
	if obj == nil || typ == nil || !types.Identical(typ, errorType) {
		return
	}

	if v, ok := obj.(*types.Var); ok && !namedErrors[v] {
		pass.Report(analysis.Diagnostic{
			Pos:     ident.Pos(),
			Message: fmt.Sprintf("deferred function assigns to error %q, which is not a named return – this assignment will not affect the function's return value", ident.Name),
			URL:     "github.com/jakobmoellerdev/deferrlint",
			SuggestedFixes: []analysis.SuggestedFix{
				{Message: fmt.Sprintf("Consider making %q a named return value", ident.Name)},
			},
		})
	}
}
