package main

import (
	"os"
)

func main() {
	if len(os.Args) > 1 {
		// Build a queue for all arguments
		cmdQue := make([]string, 0)
		for _, cmd := range os.Args[1:] {
			cmdQue = append(cmdQue, cmd)
		}
		handleCmd(cmdQue)
	}
}
