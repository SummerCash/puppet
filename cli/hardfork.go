// Package cli defines helpful cli helper methods.
package cli

import (
	"strconv"
	"strings"

	"github.com/urfave/cli"

	summercashCommon "github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/puppet/common"
)

/* BEGIN EXPORTED METHODS */

// SetupHardforkCommand sets up the hardfork CLI command.
func (app *CLI) SetupHardforkCommand() {
	(*app).App.Commands = append((*app).App.Commands, cli.Command{
		Name:    "hardfork",                      // Set name
		Aliases: []string{"fork", "f"},           // Set aliases
		Usage:   "fork the SummerCash blockmesh", // Set usage
		Action:  app.forkBlockmesh,               // Set action
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "data-dir, data",       // Set name
				Value:       common.DataDir,         // Set value
				Usage:       "path start search in", // Set usage
				Destination: &common.DataDir,        // Set destination
			},
		},
	})
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// forkBlockmesh handles the fork command.
func (app *CLI) forkBlockmesh(c *cli.Context) error {
	summercashCommon.Silent = true // Silence logsconfigPath

	summercashCommon.DataDir = common.DataDir // Set smc data dir

	config, err := config.ReadChainConfigFromMemory() // Read config from persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	previousVersion := config.ChainVersion // Get old version

	err = config.UpdateChainVersion() // Update

	if err != nil { // Check for errors
		return err // Return found error
	}

	if previousVersion == config.ChainVersion { // Check has not upgraded
		splitVersion := strings.Split(config.ChainVersion, ".") // Split version

		parsed, err := strconv.Atoi(splitVersion[2]) // Parse version

		if err != nil { // Check for errors
			return err // Return found error
		}

		parsed++ // Increment

		config.ChainVersion = strings.Join(splitVersion[:2], ".") + "." + strconv.Itoa(parsed) // Set incremented
	}

	return nil // No error occurred, return nil
}

/* END INTERNAL METHODS */
