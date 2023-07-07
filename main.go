/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andrewbenington/go-ledger/cmd"
)

func main() {
	if len(os.Args) == 1 {
		GetCommands()
		os.Exit(0)
	}
	cmd.Execute()
}

func GetCommands() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter a command: ")
		args, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			continue
		}

		// Trim any leading/trailing white spaces and newlines
		args = strings.TrimSpace(args)

		// Check if the user wants to exit
		if args == "exit" {
			fmt.Println("Exiting...")
			break
		}
		cmd.ExecuteArgs(strings.Split(args, " "))

	}

	fmt.Println("Program ended.")
}
