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
				if err := ast.Print(fset, node); err != nil {
					log.Printf("could not print the node: %s", err)
				}
				return true
			})
		}
	}
}
