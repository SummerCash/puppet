// Package cli defines helpful cli helper methods.
package cli

import (
	"github.com/urfave/cli"
)

/* BEGIN EXPORTED METHODS */

// SetupCreateCommand sets up the create CLI command.
func SetupCreateCommand(app *cli.App) {
	app.Commands = append(app.Commands, cli.Command{
		Name: "create", // Set name
		Aliases: []string{"new", "init"}, // Set aliases
		Usage: "create a new SummerCash network", // Set usage
		Action: createNetwork, // Set action
	})
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// createNetwork handles the create command.
func createNetwork(c *cli.Context) error {
	return nil // No error occurred, return nil
}

/* END INTERNAL METHODS */