// Package cli defines helpful cli helper methods.
package cli

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"github.com/SummerCash/puppet/common"
	"github.com/tcnksm/go-input"
	"github.com/urfave/cli"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"

	walletAccounts "github.com/SummerCash/summercash-wallet-server/accounts"

	"github.com/SummerCash/go-summercash/accounts"
	summercashCommon "github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/crypto"
	"github.com/SummerCash/go-summercash/types"
)

/* BEGIN EXPORTED METHODS */

// SetupCreateCommand sets up the create CLI command.
func (app *CLI) SetupCreateCommand() {
	(*app).App.Commands = append((*app).App.Commands, cli.Command{
		Name:    "create",                          // Set name
		Aliases: []string{"new", "init"},           // Set aliases
		Usage:   "create a new SummerCash network", // Set usage
		Action:  app.createNetwork,                 // Set action
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "data-dir, data",                 // Set name
				Value:       common.DataDir,                   // Set value
				Usage:       "path to store network files in", // Set usage
				Destination: &common.DataDir,                  // Set destination
			},
			cli.StringFlag{
				Name:  "network-name, network",       // Set name
				Value: "main_net",                    // Set value
				Usage: "name to register network as", // Set usage
			},
			cli.StringFlag{
				Name:  "config-path, config",                                               // Set name
				Value: "",                                                                  // Set value
				Usage: "existing network configuration to bootstrap network creation from", // Set usage
			},
			cli.StringFlag{
				Name:  "genesis-path, genesis",                                                                                                 // Set name
				Value: "",                                                                                                                      // Set value
				Usage: "file to bootstrap network configuration creation from; can contain supply, network id, and inflation rate definitions", // Set usage
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
		var dataDir, err = app.InputConfig.Ask("Where would you like your new network to be stored?", &input.Options{
			Default:   common.GetDefaultDataPath(), // Set default
			Required:  false,                       // Make optional
			HideOrder: true,                        // Hide extra question
		})

		if err != nil { // Check for errors
			return err // Return found error
		}

		if dataDir == "\r" { // Check no value set
			dataDir = common.DataDir // Set to default
		}

		common.DataDir = dataDir // Set data dir
	}

	summercashCommon.DataDir = common.DataDir // Set smc data dir

	if configPath := c.String("config-path"); configPath != "" { // Check has existing configuration file.
		return constructNetwork(c, configPath) // Construct network
	}

	config, err := app.parseGenesisFile(c.String("genesis-path")) // Parse genesis file

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = config.WriteToMemory() // Write config to persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	return constructNetwork(c, c.String("data-path")) // Construct network
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

// parseGenesisFile parses a genesis file at a given genesisPath.
func (app *CLI) parseGenesisFile(genesisPath string) (*config.ChainConfig, error) {
	rawJSON := []byte("{}") // Init raw JSON buffer
	var err error           // Init error buffer

	if genesisPath != "" { // Check no genesis file
		rawJSON, err = ioutil.ReadFile(genesisPath) // Read genesis file

		if err != nil { // Check for errors
			return &config.ChainConfig{}, err // Return found error
		}
	}

	var readJSON map[string]interface{} // Init buffer

	err = json.Unmarshal(rawJSON, &readJSON) // Unmarshal to buffer

	if err != nil { // Check for errors
		return &config.ChainConfig{}, err // Return error
	}

	alloc := make(map[string]*big.Float) // Init alloc map

	allocAddresses := []summercashCommon.Address{} // Init alloc address buffer

	networkID := uint(0) // Init network ID buffer

	inflation := float64(0) // Init inflation rate buffer

	chainID := summercashCommon.Hash{} // Init hash buffer

	if readJSON["networkID"] != nil { // Check has network ID
		networkID = uint(readJSON["networkID"].(float64)) // Set network ID
	} else {
		networkIDString, err := app.InputConfig.Ask("What is this network's network ID?", &input.Options{
			Default:   "1",   // Set default
			Required:  false, // Make required
			HideOrder: true,  // Hide extra question
		})

		if networkIDString == "\r" { // Check no ID
			networkIDString = "1" // Set to default ID
		}

		networkIDString = strings.Replace(networkIDString, "\r", "", 1) // Remove \r

		networkIDInt, err := strconv.Atoi(networkIDString) // Parse network ID

		if err != nil { // Check for errors
			return &config.ChainConfig{}, err // Return error
		}

		networkID = uint(networkIDInt) // Convert to uint
	}

	if readJSON["alloc"] != nil { // Check has alloc
		alloc, allocAddresses, err = parseAlloc(readJSON) // Parse alloc

		if err != nil { // Check for errors
			return &config.ChainConfig{}, err // Return error
		}
	} else { // User has not specified alloc in genesis
		alloc, allocAddresses, err = app.requestAlloc(networkID) // Request alloc

		if err != nil { // Check for errors
			return &config.ChainConfig{}, err // Return error
		}
	}

	if readJSON["inflation"] != nil { // Check has inflation rate
		inflation = readJSON["inflation"].(float64) // Set inflation
	} else {
		inflationString, err := app.InputConfig.Ask("What will this network's inflation rate be?", &input.Options{
			Default:   "0.0", // Set default
			Required:  false, // Make required
			HideOrder: true,  // Hide extra question
		})

		inflation, err = strconv.ParseFloat(inflationString, 64) // Parse float

		if err != nil { // Check for errors
			return &config.ChainConfig{}, err // Return found error
		}
	}

	chainID = summercashCommon.NewHash(crypto.Sha3([]byte{byte(networkID)})) // Set chain ID

	chainVersion := config.Version // Set chain version

	return &config.ChainConfig{
		Alloc:          alloc,          // Set alloc
		AllocAddresses: allocAddresses, // Set alloc addresses
		NetworkID:      networkID,      // Set network ID
		InflationRate:  inflation,      // Set inflation
		ChainID:        chainID,        // Set chain ID
		ChainVersion:   chainVersion,   // Set chain version
	}, nil // Return chain config
}

func (app *CLI) requestAlloc(networkID uint) (map[string]*big.Float, []summercashCommon.Address, error) {
	alloc := make(map[string]*big.Float)           // Init alloc map
	allocAddresses := []summercashCommon.Address{} // Init alloc address buffer

	totalIssuanceString, err := app.InputConfig.Ask("How many coins would you like to issue?", &input.Options{
		Default:   "21000000", // Set default
		Required:  true,       // Make required
		HideOrder: true,       // Hide extra question
	})

	if totalIssuanceString == "\r" { // Check no value specified
		totalIssuanceString = "21000000" // Set default
	}

	totalIssuanceString = strings.Replace(totalIssuanceString, "\r", "", 1) // Remove \r

	totalIssuanceBigVal, _, _ := big.ParseFloat(totalIssuanceString, 10, 18, big.ToNearestEven) // Parse total issuance

	genesisAccount, err := newAccount(networkID) // Initialize genesis account

	if err != nil { // Check for errors
		return nil, []summercashCommon.Address{}, err // Return found error
	}

	alloc[genesisAccount.Address.String()] = totalIssuanceBigVal    // Set value
	allocAddresses = append(allocAddresses, genesisAccount.Address) // Append genesis account address

	shouldEnableFaucetString, err := app.InputConfig.Ask("Would you like to enable the SummerCash faucet?", &input.Options{
		Default:   "true", // Set default
		Required:  true,   // Make required
		HideOrder: true,   // Hide extra question
	})

	if shouldEnableFaucetString == "\r" { // Check no value specified
		shouldEnableFaucetString = "true" // Set should not
	}

	shouldEnableFaucetString = strings.Replace(shouldEnableFaucetString, "\r", "", 1) // Remove \r

	shouldEnableFaucet, err := strconv.ParseBool(shouldEnableFaucetString) // Parse should enable faucet

	if err != nil { // Check for errors
		return nil, []summercashCommon.Address{}, err // Return found error
	}

	if shouldEnableFaucet { // Check should enable faucet
		_, err = walletAccounts.OpenDB() // Open db, create faucet account

		if err != nil { // Check for errors
			return nil, []summercashCommon.Address{}, err // Return found error
		}
	}

	for x := 0; true; x++ { // Do until break
		message := "Would you like to add a genesis address (optional, press enter to skip)?" // Set default message

		if x > 0 { // Check multiple addresses
			message = "Would you like to add another genesis address (optional, press enter to skip)?" // Set message
		}

		additionalAddress, err := app.InputConfig.Ask(message, &input.Options{
			Default:   "",    // Set default
			Required:  false, // Make optional
			HideOrder: true,  // Hide extra question
		})

		if additionalAddress == "" { // Check no additional address
			break // Break
		}

		address, err := summercashCommon.StringToAddress(additionalAddress) // Parse string address

		if err != nil { // Check for errors
			return nil, []summercashCommon.Address{}, err // Return found error
		}

		additionalBalance, err := app.InputConfig.Ask("How much SummerCash would you like to give to this address?", &input.Options{
			Default:   "0",  // Set default
			Required:  true, // Make required
			HideOrder: true, // Hide extra question
		})

		additionalBalanceBigVal, _, _ := big.ParseFloat(additionalBalance, 64, 18, big.ToNearestEven) // Parse balance string val

		alloc[address.String()] = additionalBalanceBigVal // Set val
		allocAddresses = append(allocAddresses, address)  // Append alloc address
	}

	return alloc, allocAddresses, nil // No error occurred, return nil
}

// parseAlloc parses an alloc.
func parseAlloc(json map[string]interface{}) (map[string]*big.Float, []summercashCommon.Address, error) {
	alloc := make(map[string]*big.Float) // Init alloc map

	allocAddresses := []summercashCommon.Address{} // Init alloc address buffer

	for key, value := range json["alloc"].(map[string]interface{}) { // Iterate through genesis addresses
		floatVal, _, _ := big.ParseFloat(value.(map[string]interface{})["balance"].(string), 10, 350, big.ToNearestEven) // Parse float

		address, err := summercashCommon.StringToAddress(key) // Get address value

		if err != nil { // Check for errors
			return nil, []summercashCommon.Address{}, err // Return error
		}

		allocAddresses = append(allocAddresses, address) // Append address

		alloc[key] = floatVal // Set int val
	}

	return alloc, allocAddresses, nil // No error occurred, return nil
}

// newAccount initializes a new account, along with a chain.
func newAccount(networkID uint) (*accounts.Account, error) {
	account := &accounts.Account{
		Address: summercashCommon.Address{'\r'}, // Set mock address
	} // Init account buffer

	for bytes.Contains(account.Address.Bytes(), []byte{'\r'}) { // Generate accounts until valid
		privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

		if err != nil { // Check for errors
			return &accounts.Account{}, err // Return error
		}

		account, err = accounts.AccountFromKey(privateKey) // Generate account from key

		if err != nil { // Check for errors
			return &accounts.Account{}, err // Return error
		}
	}

	chain := &types.Chain{ // Init chain
		Account:      account.Address,
		Transactions: []*types.Transaction{},
		NetworkID:    networkID,
	}

	(*chain).ID = summercashCommon.NewHash(crypto.Sha3(chain.Bytes())) // Set ID

	return account, chain.WriteToMemory() // Write to memory
}

/* END INTERNAL METHODS */
