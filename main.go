// Package main defines the puppet entry point.
package main

import (
	"log"
	"os"

	"github.com/SummerCash/puppet/cli"
)

// main is the puppet entry function.
func main() {
	app := cli.NewCLIApp() // Initialize CLI app

	cli.SetupCreateCommand(app) // Setup create command

	err := app.App.Run(os.Args) // Initialize CLI app

	if err != nil { // Check for errors
		log.Fatal(err) // Panic
	}
}
