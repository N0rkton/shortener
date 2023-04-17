package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var ErrMainExit = &analysis.Analyzer{
	Name: "mainExit",
	Doc:  "check for os.Exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		// проверяем, что выражение представляет собой вызов функции,
		// у которой возвращаемая ошибка никак не обрабатывается
		if _, ok := x.X.(*ast.CallExpr); ok {
			{
				if x.X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "Exit is forbidden")
				}
			}
		}
	}
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {

			case *ast.ExprStmt: // выражение
				expr(x)
			}
			return true
		})
	}
	return nil, nil
}
