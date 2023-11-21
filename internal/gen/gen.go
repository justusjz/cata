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

func extractParamTypes(params []ast.Param) []ast.TypeNode {
	types := []ast.TypeNode{}
	for _, param := range params {
		types = append(types, param.Type)
	}
	return types
}

func Gen(module *ast.Module, out string) error {
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
	g := generator{body: body, header: header, scanner: nil, indent: 0, scope: newGlobalScope()}
	for _, st := range module.Structs {
		// add structs to scope
		g.scanner = st.Scanner
		if g.scope.findType(st.Name.Ident) != nil {
			g.diagnose(st.Name.Pos, "duplicate identifier '%s'", st.Name.Ident)
		}
		g.scope.addType(st.Name.Ident, &scopeType{decl: st})
	}
	for _, fn := range module.Fns {
		// add functions to scope
		g.scanner = fn.Scanner
		fnParams := extractParamTypes(fn.Params)
		fnTy := &ast.FnType{Params: fnParams, ReturnType: fn.ReturnType}
		if g.scope.findVar(fn.Name.Ident) != nil {
			g.diagnose(fn.Name.Pos, "duplicate identifier '%s'", fn.Name.Ident)
		}
		g.scope.addVar(fn.Name.Ident, &scopeVar{ty: fnTy, mut: false})
	}
	for _, st := range module.Structs {
		// generate structs, necessary here for checking unused structs
		g.scanner = st.Scanner
		g.genStructDecl(st)
	}
	for _, fn := range module.Fns {
		// generate functions
		g.scanner = fn.Scanner
		g.genFnDecl(fn)
	}
	body.Close()
	header.Close()
	return nil
}
