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

	for _, pkg := range pkgs {
		log.Printf("package %s", pkg.ID)
		log.Printf("files %+v", pkg.CompiledGoFiles)
		for _, syntax := range pkg.Syntax {
			log.Printf("file=%s", syntax.Name)
			if err := printer.Fprint(os.Stdout, pkg.Fset, syntax); err != nil {
				log.Printf("could not print the source: %s", err)
			}
		}
	}

	// print AST nodes
	for _, pkg := range pkgs {
		for _, syntax := range pkg.Syntax {
			ast.Inspect(syntax, func(node ast.Node) bool {
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
