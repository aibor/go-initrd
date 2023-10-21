package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/aibor/go-initrd"
)

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no init file given")
	}

	initFile := args[0]
	additionalFiles := args[1:]
	i := initrd.New(initFile, additionalFiles...)

	searchPathString := strings.TrimSpace(os.Getenv("GOINITRDSEARCH"))
	searchPaths := strings.Split(searchPathString, ":")
	searchPaths = slices.DeleteFunc(searchPaths, func(e string) bool { return e == "" })

	resolver := initrd.NewELFLibResolver(searchPaths...)
	if err := i.ResolveLinkedLibs(resolver); err != nil {
		return fmt.Errorf("resolve: %v", err)
	}

	writer := initrd.NewCPIOWriter(os.Stdout)
	defer writer.Close()
	if err := i.WriteTo(writer); err != nil {
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
