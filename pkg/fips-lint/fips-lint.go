package fipslint

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"

	"github.com/fatih/color"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "fips140",
	Doc:  "checks fips 140-2 compliance",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {

			callExpr, ok := n.(*ast.CallExpr)
			if ok {

				// Prting the call expression
				// fmt.Printf("CallExpr: %s\n", render(pass.Fset, callExpr))

				// Check if the function call is crypto/md5.New()
				if funSel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					//if funSel.Sel.Name == "New" {
					//					fmt.Printf("!! CallExpr: %s\n", render(pass.Fset, callExpr))

					if iden, ok := funSel.X.(*ast.Ident); ok {
						// Check if the package name is crypto/md5
						if obj := pass.TypesInfo.ObjectOf(iden); obj != nil {
							if pkgName, ok := obj.(*types.PkgName); ok {
								// fmt.Printf("!! CallExpr: %s at line:%v\n", render(pass.Fset, callExpr), pass.Fset.Position(callExpr.Pos()).Filename)
								importedPkg := pkgName.Imported().Path()
								funcName := funSel.Sel.Name
								message := fmt.Sprintf("Package: %-20s function name: %-20s file %s (line:%d)",
									importedPkg,
									funcName,
									pass.Fset.Position(callExpr.Pos()).Filename,
									pass.Fset.Position(callExpr.Pos()).Line)
								if isBannedFunction(function{importedPkg, funcName}) {
									color.Red(message)
								} else {
									color.Green(message)
								}

							}

							if pkg := obj.Pkg(); pkg != nil && pkg.Path() == "crypto/md5" {
								fmt.Printf("crypto/md5.New() function call found at line: %d\n", pass.Fset.Position(callExpr.Pos()).Line)
							}
						}
					}
					//}
				}
			}

			// check if this node is a function call

			be, ok := n.(*ast.BinaryExpr)
			if !ok {
				return true
			}

			if be.Op != token.ADD {
				return true
			}

			if _, ok := be.X.(*ast.BasicLit); !ok {
				return true
			}

			if _, ok := be.Y.(*ast.BasicLit); !ok {
				return true
			}

			isInteger := func(expr ast.Expr) bool {
				t := pass.TypesInfo.TypeOf(expr)
				if t == nil {
					return false
				}

				bt, ok := t.Underlying().(*types.Basic)
				if !ok {
					return false
				}

				if (bt.Info() & types.IsInteger) == 0 {
					return false
				}

				return true
			}

			// check that both left and right hand side are integers
			if !isInteger(be.X) || !isInteger(be.Y) {
				return true
			}

			pass.Reportf(be.Pos(), "integer addition found %q",
				render(pass.Fset, be))
			return true
		})
	}

	return nil, nil
}

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
 
func isBannedFunction(f function) bool{
	for _, excludedFunction := range excludedFunctions {
		if f.name == excludedFunction.name && f.pkg == excludedFunction.pkg {
			return true
		}
	}
	return false
}