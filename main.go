/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/andrewbenington/go-ledger/app"
	"github.com/andrewbenington/go-ledger/cmd"
)

func main() {
	if len(os.Args) == 1 {
		app.Start()
		os.Exit(0)
	}
	cmd.Execute()
}
