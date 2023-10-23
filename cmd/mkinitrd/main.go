package main

import (
	"fmt"
	"os"

	"github.com/aibor/go-initrd"
)

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no init file given")
	}

	initFile := args[0]
	additionalFiles := args[1:]
	libSearchPath := os.Getenv("GOINITRDLIBPATH")

	initRD := initrd.New(initFile)
	if err := initRD.AddFiles(additionalFiles...); err != nil {
		return fmt.Errorf("add files: %v", err)
	}
	if err := initRD.ResolveLinkedLibs(libSearchPath); err != nil {
		return fmt.Errorf("add linked libs: %v", err)
	}
	if err := initRD.WriteCPIO(os.Stdout); err != nil {
		return fmt.Errorf("write: %v", err)
	}

	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
