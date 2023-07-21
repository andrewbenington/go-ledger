package app

import (
	"fmt"
	"os"
)

func log(content string, filepath string) {
	filePath := "go-ledger.log"

	// Open the file in append mode with write-only permissions
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
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
	log(content, "error.log")
}

func Log(content string, args ...any) {
	log(fmt.Sprintf(content, args...), "go-ledger.log")
}

func LogInterface(obj interface{}) {
	log(fmt.Sprintf("%+v", obj), "go-ledger.log")
}

func LogStack() {
	str := "Stack:"
	for _, cmd := range stack {
		str += " " + cmd.Name
	}
	log(str, "go-ledger.log")
}
