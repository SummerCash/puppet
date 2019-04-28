// Package cli defines helpful cli helper methods.
package cli

import (
	"os"
	"github.com/SummerCash/puppet/common"

	"github.com/urfave/cli"

	"github.com/tcnksm/go-input"
)

// CLI defines a command-line-interface.
type CLI struct {
	App *cli.App // CLI app
	InputConfig *input.UI // Input config
}

/* BEGIN EXPORTED METHODS */

// NewCLI initializes a new cli application.
func NewCLI() *CLI {
	app := cli.NewApp() // Initialize CLI app

	app.Name = "puppet" // Set name
	app.Usage = "a visual CLI for managing, creating, and analyzing SummerCash networks." // Set description
	app.Version = "v0.1.0" // Set version
	app.EnableBashCompletion = true // Enable auto-completion

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "data-dir, data", // Set name
			Value: common.DataDir, // Set value
			Usage: "path to store network files in", // Set usage
			Destination: &common.DataDir, // Set destination
		},
		cli.StringFlag{
			Name: "network-name, network", // Set name
			Value: "main_net", // Set value
			Usage: "name to register network as", // Set usage
		},
		cli.StringFlag{
			Name: "config-path, config", // Set name
			Value: "", // Set value
			Usage: "existing network configuration to bootstrap network creation from", // Set usage
		},
		cli.StringFlag{
			Name: "genesis-path, genesis", // Set name
			Value: "", // Set value
			Usage: "file to bootstrap network configuration creation from; must contain supply, network id, and inflation rate definitions", // Set usage
		},
	}

	return &CLI{
		App: app, // Set app
		InputConfig: &input.UI{
			Writer: os.Stdout, // Set output
			Reader: os.Stdin, // Set input
		}, // Set config
	} // Return CLI
}

/* END EXPORTED METHODS */