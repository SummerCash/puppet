// Package cli defines helpful cli helper methods.
package cli

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"github.com/SummerCash/puppet/common"
	"github.com/urfave/cli"
	"github.com/tcnksm/go-input"

	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/types"
)

/* BEGIN EXPORTED METHODS */

// SetupCreateCommand sets up the create CLI command.
func (app *CLI) SetupCreateCommand() {
	(*app).App.Commands = append((*app).App.Commands, cli.Command{
		Name: "create", // Set name
		Aliases: []string{"new", "init"}, // Set aliases
		Usage: "create a new SummerCash network", // Set usage
		Action: app.createNetwork, // Set action
		Flags: []cli.Flag{
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
		},
	})
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// createNetwork handles the create command.
func (app *CLI) createNetwork(c *cli.Context) error {
	var err error // Init error buffer

	if c.String("data-dir") == common.GetDefaultDataPath() { // Check data directory not specified
		common.DataDir, err = app.InputConfig.Ask("Where would you like your new network to be stored?", &input.Options{
			Default: common.GetDefaultDataPath(), // Set default
			Required: false, // Make optional
			HideOrder: true, // Hide extra question
		})

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	if configPath := c.String("config-path"); configPath != "" { // Check has existing configuration file.
		return constructNetwork(c, configPath) // Construct network
	}

	return nil // No error occurred, return nil
}

// constructNetwork constructs a network, assuming all configs have been set.
func constructNetwork(c *cli.Context, dataPath string) error {
	config, err := readChainConfig(dataPath) // Read chain config

	if err != nil { // Check for errors
		return err // Return found error
	}

	genesisAddress := config.AllocAddresses[0] // Get genesis address

	chain, err := types.NewChain(genesisAddress) // Initialize chain

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = chain.WriteToMemory() // Write chain to persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	_, err = chain.MakeGenesis(config) // Make genesis

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// readChainConfig reads a chain config at a given configPath.
func readChainConfig(configPath string) (*config.ChainConfig, error) {
	path, _ := filepath.Abs(filepath.FromSlash(configPath)) // Get direct path

	data, err := ioutil.ReadFile(path) // Read file

	if err != nil { // Check for errors
		return &config.ChainConfig{}, err // Return error
	}

	buffer := &config.ChainConfig{} // Initialize buffer

	err = json.Unmarshal(data, buffer) // Read json into buffer

	if err != nil { // Check for errors
		return &config.ChainConfig{}, err // Return error
	}

	return buffer, nil // No error occurred, return read config
}

/* END INTERNAL METHODS */