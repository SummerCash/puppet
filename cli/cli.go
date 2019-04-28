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
		},
	}

	return app // Return app
}

/* END EXPORTED METHODS */