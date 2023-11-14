// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"
	"os"

	"github.com/justusjz/cata/internal/ast"
	"github.com/justusjz/cata/internal/scanner"
)

type generator struct {
	body    *os.File
	header  *os.File
	scanner *scanner.Scanner
	indent  int
	scope   *scope
}

func (g *generator) diagnose(pos scanner.Pos, format string, a ...any) {
	g.scanner.Diagnose(pos, format, a...)
}

func (g *generator) writeIndent() {
	for i := 0; i < g.indent; i++ {
		fmt.Fprint(g.body, " ")
	}
}

func Gen(fnDecl *ast.FnDecl, out string) error {
	body, err := os.Create(out + ".c")
	if err != nil {
		return err
	}
	header, err := os.Create(out + ".h")
	if err != nil {
		return err
	}
	fmt.Fprint(header, "#include <stdint.h>\n\n")
	fmt.Fprintf(body, "#include \"%s.h\"\n\n", out)
	g := generator{body: body, header: header, scanner: fnDecl.Scanner, indent: 0, scope: newGlobalScope()}
	g.genFnDecl(fnDecl)
	body.Close()
	header.Close()
	return nil
}
