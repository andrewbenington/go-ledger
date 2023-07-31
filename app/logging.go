package app

import (
	"fmt"
	"os"
)

func logToFile(content string, filepath string) {
	if filepath == "" {
		filepath = "go-ledger.log"
	}

	// Open the file in append mode with write-only permissions
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	// Write the content to the end of the file
	_, err = file.WriteString(content + "\n")
	if err != nil {
		return
	}
}

func LogErr(content string) {
	logToFile(content, "error.log")
}

func Log(content string, args ...any) {
	logToFile(fmt.Sprintf(content, args...), "go-ledger.log")
}

func LogInterface(obj interface{}) {
	logToFile(fmt.Sprintf("%+v", obj), "go-ledger.log")
}

func LogStack() {
	str := "Stack:"
	for _, cmd := range stack {
		str += " " + cmd.Name
	}
	logToFile(str, "go-ledger.log")
}
