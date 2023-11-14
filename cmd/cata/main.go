// Copyright (c) 2023 Justus Zorn

package main

import (
	"fmt"
	"os"

	"github.com/justusjz/cata/internal/parser"
)

func main() {
	fn, err := parser.Parse("test.cata")
	if err != nil {
		fmt.Printf("error: file 'test.cata' not found\n")
		os.Exit(1)
	}
	println(fn)
}
