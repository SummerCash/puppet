// Package cli defines helpful cli helper methods.
package cli

import (
	"github.com/SummerCash/puppet/common"
	"github.com/urfave/cli"
	"github.com/tcnksm/go-input"
)

/* BEGIN EXPORTED METHODS */

// SetupCreateCommand sets up the create CLI command.
func (app *CLI) SetupCreateCommand() {
	(*app).App.Commands = append((*app).App.Commands, cli.Command{
		Name: "create", // Set name
		Aliases: []string{"new", "init"}, // Set aliases
		Usage: "create a new SummerCash network", // Set usage
		Action: app.createNetwork, // Set action
	})
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// createNetwork handles the create command.
func (app *CLI) createNetwork(c *cli.Context) error {
	var err error // Init error buffer

	if c.GlobalString("data-dir") == common.GetDefaultDataPath() { // Check data directory not specified
		common.DataDir, err = app.InputConfig.Ask("Where would you like your new network to be stored?", &input.Options{
			Default: common.GetDefaultDataPath(), // Set default
			Required: false, // Make optional
			HideOrder: true, // Hide extra question
		})

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	return nil // No error occurred, return nil
}

/* END INTERNAL METHODS */