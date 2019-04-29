// Package cli defines helpful cli helper methods.
package cli

import (
	"os"

	"github.com/urfave/cli"
	"github.com/tcnksm/go-input"
)

// CLI defines a command-line-interface.
type CLI struct {
	App         *cli.App  // CLI app
	InputConfig *input.UI // Input config
}

/* BEGIN EXPORTED METHODS */

// NewCLI initializes a new cli application.
func NewCLI() *CLI {
	app := cli.NewApp() // Initialize CLI app

	app.Name = "puppet"                                                                   // Set name
	app.Usage = "a visual CLI for managing, creating, and analyzing SummerCash networks." // Set description
	app.Version = "v0.1.3"                                                                // Set version
	app.EnableBashCompletion = true                                                       // Enable auto-completion

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "node-port, port",                                    // Set name
			Value: 3033,                                                 // Set value
			Usage: "port to use for p2p communications (if applicable)", // Set usage
		},
	}

	return &CLI{
		App: app, // Set app
		InputConfig: &input.UI{
			Writer: os.Stdout, // Set output
			Reader: os.Stdin,  // Set input
		}, // Set config
	} // Return CLI
}

/* END EXPORTED METHODS */
