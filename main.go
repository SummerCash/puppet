// Package main defines the puppet entry point.
package main

import (
	"github.com/urfave/cli"
	"os"
	"log"
)

// main is the puppet entry function.
func main() {
	err := cli.NewApp().Run(os.Args) // Initialize CLI app

	if err != nil { // Check for errors
		log.Fatal(err) // Panic
	}
}