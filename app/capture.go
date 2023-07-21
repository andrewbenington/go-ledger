package app

import (
	"bufio"
	"os"
)

func captureFunctionErrors(f func() error) error {
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer func() {
		w.Close()
		os.Stderr = old
	}()
	os.Stderr = w

	running := true
	go func() {
		Log("starting error log goroutine")
		reader := bufio.NewReader(r)
		for running {
			line, _ := reader.ReadString('\n')
			LogErr(line)
		}
		Log("ending error log goroutine")
	}()

	err = f()

	running = false
	return err
}
