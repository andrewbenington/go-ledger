/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/andrewbenington/go-ledger/app"
	"github.com/andrewbenington/go-ledger/cmd"
)

func main() {
	if strings.HasPrefix(runtime.Version(), "go1.21.0") {
		fmt.Printf("Go 1.21.0 not supported due to a bug in the standard XML library. See https://github.com/golang/go/issues/61881")
		os.Exit(1)
	}
	if len(os.Args) == 1 {
		app.Start()
		os.Exit(0)
	}
	cmd.Execute()
}
