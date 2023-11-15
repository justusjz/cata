// Copyright (c) 2023 Justus Zorn

package main

import (
	"github.com/justusjz/cata/internal/gen"
	"github.com/justusjz/cata/internal/parser"
)

func main() {
	module := parser.Parse("test.cata")
	gen.Gen(module, "out")
}
