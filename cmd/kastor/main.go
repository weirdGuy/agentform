package main

import (
	"fmt"
	"os"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "kastor: %v\n", err)
		os.Exit(exitCode(err))
	}
}
