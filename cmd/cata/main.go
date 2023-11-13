// Copyright (c) 2023 Justus Zorn

package main

import (
	"fmt"
	"os"

	"github.com/justusjz/cata/internal/scanner"
)

func main() {
	s, err := scanner.New("test.cata")
	if err != nil {
		fmt.Printf("error: file 'test.cata' not found\n")
		os.Exit(1)
	}
	s.Diagnose(s.Pos(), "some error message")
}
