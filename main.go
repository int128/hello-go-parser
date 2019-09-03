package main

import (
	"go/ast"
	"go/printer"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

func main() {
	config := &packages.Config{
		Mode: packages.NeedCompiledGoFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(config, os.Args[1:]...)
	if err != nil {
		log.Fatalf("could not load packages: %s", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		log.Fatalf("error occurred")
	}

	// print the AST
	for _, pkg := range pkgs {
		if err := ast.Print(pkg.Fset, pkg); err != nil {
			log.Printf("could not print the AST: %s", err)
		}
	}

	// print the function calls
	for _, pkg := range pkgs {
		for _, syntax := range pkg.Syntax {
			ast.Inspect(syntax, func(node ast.Node) bool {
				switch node.(type) {
				case *ast.ImportSpec:
					// imports
					imp := node.(*ast.ImportSpec)
					log.Printf("import %s as %s", imp.Path.Value, imp.Name)

				case *ast.CallExpr:
					call := node.(*ast.CallExpr)
					switch fun := call.Fun.(type) {
					case *ast.SelectorExpr:
						switch x := fun.X.(type) {
						case *ast.Ident:
							// print the function call
							log.Printf("call %s.%s with %d arg(s)", x, fun.Sel, len(call.Args))

							// mutate the function call
							if x.Name == "errors" {
								x.Name = "xerrors"
							}
						default:
							log.Printf("unknown type of call.Fun.X: %T", fun)
						}
					default:
						log.Printf("unknown type of call.Fun: %T", fun)
					}
				}
				return true
			})
		}
	}

	// print the sources
	for _, pkg := range pkgs {
		for _, syntax := range pkg.Syntax {
			log.Printf("file=%s", syntax.Name)
			if err := printer.Fprint(os.Stdout, pkg.Fset, syntax); err != nil {
				log.Printf("could not print the source: %s", err)
			}
		}
	}
}
