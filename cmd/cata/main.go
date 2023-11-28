// Copyright (c) 2023 Justus Zorn

package main

import (
	"os"
	"os/exec"

	"github.com/justusjz/cata/internal/check"
	"github.com/justusjz/cata/internal/gen"
	"github.com/justusjz/cata/internal/parser"
)

func main() {
	module := parser.Parse("test.cata")
	check.CheckModule(module)
	gen.GenModule(module, "test.cata")
	cmd := exec.Command("tcc", "test.cata.c", "catalib.c")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return
	}
	cmd = exec.Command("./test.cata.exe")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
