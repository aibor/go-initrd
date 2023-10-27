package main

import (
	"fmt"
	"os"

	"github.com/aibor/initramfs"
)

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no init file given")
	}

	initFile := args[0]
	additionalFiles := args[1:]
	libSearchPath := os.Getenv("LD_LIBRARY_PATH")

	initRamFS := initramfs.New(initFile)
	if err := initRamFS.AddFiles(additionalFiles...); err != nil {
		return fmt.Errorf("add files: %v", err)
	}
	if err := initRamFS.ResolveLinkedLibs(libSearchPath); err != nil {
		return fmt.Errorf("add linked libs: %v", err)
	}
	if err := initRamFS.WriteCPIO(os.Stdout); err != nil {
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
