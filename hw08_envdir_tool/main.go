package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("You must pass 2 or more arguments")
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Printf("Error parse env dir")
	}

	outputCode := RunCmd(os.Args[2:], env)
	os.Exit(outputCode)
}
