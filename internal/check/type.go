// Copyright (c) 2023 Justus Zorn

package check

import (
	"github.com/justusjz/cata/internal/ast"
)

func TypeEqual(left ast.Type, right ast.Type) bool {
	if l, ok := left.(*ast.NamedType); ok {
		if r, ok := right.(*ast.NamedType); ok {
			if l.Name.Ident != r.Name.Ident {
				return false
			}
			if len(l.Args) != len(r.Args) {
				// different number of generic arguments, shouldn't happen
				panic("different number of generic arguments")
			}
			for i := 0; i < len(l.Args); i++ {
				if !TypeEqual(l.Args[i], r.Args[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (c *checker) checkType(ty *ast.NamedType) {
	rty := c.scope.findType(ty.Name.Ident)
	if rty == nil {
		c.diagnose(ty.At(), "undefined type '%s'", ty.Name.Ident)
	}
	// check generic arguments
	if len(ty.Args) != len(rty.Params) {
		c.diagnose(ty.Name.Pos, "expected %d generic arguments, but got %d", len(rty.Params), len(ty.Args))
	}
	for _, arg := range ty.Args {
		c.checkType(arg)
	}
	// if the type is a struct, check whether it is correct
	if st, ok := c.structs[ty.Name.Ident]; ok {
		c.checkStruct(st)
	}
}

func ResolveGenericType(ty *ast.NamedType, env map[string]*ast.NamedType) *ast.NamedType {
	if aty, ok := env[ty.Name.Ident]; ok {
		// type is a generic type
		if len(ty.Args) > 0 {
			panic("cannot instantiate generic parameter")
		}
		return aty
	} else {
		// type is not a generic type
		args := []*ast.NamedType{}
		for _, arg := range ty.Args {
			// resolve its generic arguments the same way
			args = append(args, ResolveGenericType(arg, env))
		}
		return &ast.NamedType{Name: ty.Name, Args: args}
	}
}

func GetGenericEnv(st *ast.StructDecl, args []*ast.NamedType) map[string]*ast.NamedType {
	if len(st.Params) != len(args) {
		panic("incorrect generic instantiation")
	}
	env := map[string]*ast.NamedType{}
	for i := 0; i < len(st.Params); i++ {
		env[st.Params[i].Ident] = args[i]
	}
	return env
}

func (c *checker) getField(ty ast.Type, name string) ast.Type {
	// find field if ty is of struct type
	if nty, ok := ty.(*ast.NamedType); ok {
		if st, ok := c.structs[nty.Name.Ident]; ok {
			for _, field := range st.Fields {
				if field.Name.Ident == name {
					// generic environment of struct
					env := GetGenericEnv(st, nty.Args)
					return ResolveGenericType(field.Type, env)
				}
			}
		}
	}
	return nil
}
