package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "testdata/", nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("could not parse: %s", err)
	}

	// just print sources
	for _, pkg := range pkgs {
		log.Printf("package=%s", pkg.Name)
		for name, f := range pkg.Files {
			log.Printf("file=%s", name)
			if err := printer.Fprint(os.Stdout, fset, f); err != nil {
				log.Printf("could not print the source: %s", err)
			}
		}
	}

	// print AST nodes
	for _, pkg := range pkgs {
		log.Printf("package=%s", pkg.Name)
		for name, f := range pkg.Files {
			log.Printf("file=%s", name)

			ast.Inspect(f, func(node ast.Node) bool {
				//if err := ast.Print(fset, node); err != nil {
				//	log.Printf("could not print the node: %s", err)
				//}

				switch node.(type) {
				case *ast.ImportSpec:
					// imports
					imp := node.(*ast.ImportSpec)
					log.Printf("import %s as %s", imp.Path.Value, imp.Name)

				case *ast.CallExpr:
					call := node.(*ast.CallExpr)
					switch call.Fun.(type) {
					case *ast.SelectorExpr:
						// package function call
						fun := call.Fun.(*ast.SelectorExpr)
						log.Printf("call %s.%s with %d arg(s)", fun.X, fun.Sel, len(call.Args))
					case *ast.Ident:
						// local function call
						ident := call.Fun.(*ast.Ident)
						log.Printf("call %s with %d arg(s)", ident.Name, len(call.Args))
					default:
						log.Printf("call %T", call.Fun)
					}
				}
				return true
			})
		}
	}
}
