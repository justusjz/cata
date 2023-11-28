// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"
	"os"

	"github.com/justusjz/cata/internal/ast"
)

type genericInstance struct {
	ty   *ast.NamedType
	name string
}

type generator struct {
	body      *os.File
	header    *os.File
	indent    int
	structs   map[string]*ast.StructDecl
	instances []genericInstance
	counter   int
}

func (g *generator) writeIndent() {
	for i := 0; i < g.indent; i++ {
		fmt.Fprint(g.body, "\t")
	}
}

func GenModule(module *ast.Module, out string) error {
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
	g := generator{body: body, header: header, indent: 0, structs: map[string]*ast.StructDecl{}, instances: []genericInstance{}, counter: 0}
	for _, st := range module.Structs {
		// add struct to generator
		g.structs[st.Name.Ident] = st
	}
	for _, fn := range module.Fns {
		// generate functions
		g.genFn(fn)
	}
	body.Close()
	header.Close()
	return nil
}
