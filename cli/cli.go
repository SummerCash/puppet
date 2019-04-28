// Package cli defines helpful cli helper methods.
package cli

import (
	"github.com/SummerCash/puppet/common"

	"github.com/urfave/cli"
)

/* BEGIN EXPORTED METHODS */

// NewCLIApp initializes a new cli application.
func NewCLIApp() *cli.App {
	app := cli.NewApp() // Initialize CLI app

	app.Name = "puppet" // Set name
	app.Usage = "a visual CLI for managing, creating, and analyzing SummerCash networks." // Set description
	app.Version = "v0.1.0" // Set version

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "data-dir", // Set name
			Value: common.DataDir, // Set value
			Usage: "path to store network files in", // Set usage
			Destination: &common.DataDir, // Set destination
		},
		cli.StringFlag{
			Name: "network-name", // Set name
			Value: "main_net", // Set value
			Usage: "name to register network as", // Set usage
		},
		cli.StringFlag{
			Name: "config-path", // Set name
			Value: "", // Set value
			Usage: "existing network configuration to bootstrap network creation from", // Set usage
		},
		cli.StringFlag{
			Name: "genesis-path", // Set name
			Value: "", // Set value
			Usage: "file to bootstrap network configuration creation from; must contain supply, network id, and inflation rate definitions", // Set usage
		},
	}

	return app // Return app
}

/* END EXPORTED METHODS */