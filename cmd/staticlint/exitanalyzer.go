// Package main - checks for os.Exit in main function
package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// ErrMainExit - new analyzer for os.Exit in main
var ErrMainExit = &analysis.Analyzer{
	Name: "mainExit",
	Doc:  "checks for os.Exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	expr := func(x *ast.ExprStmt) {
		// проверяем, что выражение представляет собой вызов функции,
		if callExpr, ok := x.X.(*ast.CallExpr); ok {
			{
				if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if selectorExpr.Sel.Name == "Exit" && selectorExpr.X.(*ast.Ident).Name == "os" {
						pass.Reportf(x.Pos(), "Exit is forbidden")
					}
				}
			}
		}
	}
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			if file.Name.Name == "main" {
				for _, i := range file.Decls {
					fn, ok := i.(*ast.FuncDecl)
					if !ok {
						continue
					}
					if fn.Name.Name == "main" {
						switch x := node.(type) {
						case *ast.ExprStmt: // выражение
							expr(x)
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
