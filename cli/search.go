// Package cli defines helpful cli helper methods.
package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/kyokomi/emoji"
	"github.com/tcnksm/go-input"
	i "github.com/tockins/interact"
	"github.com/urfave/cli"

	summercashCommon "github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
	"github.com/SummerCash/puppet/common"
)

/* BEGIN EXPORTED METHODS */

// SetupSearchCommand sets up the search CLI command.
func (app *CLI) SetupSearchCommand() {
	(*app).App.Commands = append((*app).App.Commands, cli.Command{
		Name:    "search",                                                  // Set name
		Aliases: []string{"analyze", "query", "s"},                         // Set aliases
		Usage:   "search the SummerCash blockmesh for a particular phrase", // Set usage
		Action:  app.searchBlockmesh,                                       // Set action
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "data-dir, data",       // Set name
				Value:       common.DataDir,         // Set value
				Usage:       "path start search in", // Set usage
				Destination: &common.DataDir,        // Set destination
			},
			cli.StringFlag{
				Name:  "search-term, term",                   // Set name
				Usage: "term to search for in the blockmesh", // Set usage
			},
			cli.StringSliceFlag{
				Name:  "search-chains, chains",    // Set name
				Usage: "blockchains to search in", // Set usage
			},
		},
	})
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// searchBlockmesh handles the search command.
func (app *CLI) searchBlockmesh(c *cli.Context) error {
	summercashCommon.Silent = true // Silence logsconfigPath

	summercashCommon.DataDir = common.DataDir // Set smc data dir

	var err error // Init error buffer

	searchTerm := c.String("search-term") // Get search term

	if searchTermArg := c.Args().Get(0); searchTermArg != "" { // Check has search term arg
		searchTerm = searchTermArg // Set search term
	}

	searchChains := c.StringSlice("search-chains") // Get search chains

	if len(searchChains) == 0 { // Check search chains not provided
		for x := 1; true; x++ {
			searchChain := c.Args().Get(x) // Get arg

			if searchChain == "" {
				break // Break
			}

			searchChains = append(searchChains, searchChain) // Append search chain
		}
	}

	if searchTerm == "" { // Check search term not specified
		searchTerm, err = app.InputConfig.Ask("What term would you like to search for?", &input.Options{
			Required:  true, // Make optional
			HideOrder: true, // Hide extra question
		})

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	if len(searchChains) == 0 { // Check no search chains specified
		searchChains, err = app.requestSearchChains(c) // Request search chains

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	results, files, err := searchBlockmesh(searchChains, searchTerm) // Search

	if err != nil { // Check for errors
		return err // Return found error
	}

	if len(results) == 0 { // Check no results
		color.Red(fmt.Sprintf("No results were found in any of %d files matching your query for %s.", len(files), searchTerm)) // Log error

		return nil // No error occurred, return nil
	}

	green := color.New(color.FgGreen).PrintfFunc() // Init green

	green("All done! Found %d results in %d files matching your query for %s. Which result would you like to show?", len(results), len(files), searchTerm) // Print

	for {
		resultSelector, err := app.InputConfig.Ask("", &input.Options{
			Required:    true, // Make required
			HideOrder:   true, // Hide extra question
			HideDefault: true, // Hide default
		})

		if err != nil { // Check for errors
			break // Break
		}

		if resultSelector == "" || resultSelector == "\r" { // Check no value
			return nil // No error occurred, return nil
		}

		resultSelector = strings.Replace(resultSelector, "\r", "", 1) // Remove \r

		switch resultSelector {
		case "no":
			return nil // Finished
		case "yes":
			print("Which result would you like to show?") // Print query
			continue                                      // Continue
		}

		if !strings.Contains(resultSelector, "/") && !strings.Contains(resultSelector, ".") { // Check is not file selector
			intVal, err := strconv.Atoi(resultSelector) // Parse

			if err != nil { // Check for errors
				return err // Return found error
			}

			fmt.Println(results[intVal]) // Log result
			fmt.Println(files[intVal])   // Log filename
		} else {
			for i := 0; i < len(files); i++ { // Iterate through files
				if files[i] == resultSelector { // Check is file
					fmt.Println(results[i]) // Log result
					fmt.Println(files[i])   // Log filename
				}
			}
		}

		print("\nAre there any other results you would like to look at?") // Log query
	}

	return nil // No error occurred, return nil
}

// requestSearchChains requests the list of search chains from the user.
func (app *CLI) requestSearchChains(c *cli.Context) ([]string, error) {
	var searchChains []string // Init search chains buffer
	var err error             // Init error buffer

	shouldGetSearchChains := false // Init buffer

	i.Run(&i.Interact{
		Questions: []*i.Question{
			{
				Quest: i.Quest{
					Msg: "Would you like to search the entire blockmesh, or multiple chains?",
					Choices: i.Choices{
						Alternatives: []i.Choice{
							{
								Text:     "blockmesh",
								Response: false,
							},
							{
								Text:     "chain",
								Response: true,
							},
						},
					},
				},
				Action: func(c i.Context) interface{} {
					shouldGetSearchChains, err = c.Ans().Bool() // Set should get search chains

					return nil
				},
			},
		},
	})

	if err != nil { // Check for errors
		return []string{}, err // Return found error
	}

	if shouldGetSearchChains { // Check must get search chains
		searchChainsString, err := app.InputConfig.Ask("What chains would you like to search?", &input.Options{
			Required:  true, // Make optional
			HideOrder: true, // Hide extra question
		})

		if err != nil { // Check for errors
			return []string{}, err // Return found error
		}

		searchChains = strings.Split(searchChainsString, ", ") // Split
	}

	return searchChains, nil // Return search chains
}

// searchBlockmesh searches the blockmesh for a particular string.
// If the search term is found, the particular resources in which it is contained are returned as JSON string values,
// followed by their file names.
func searchBlockmesh(searchChains []string, searchTerm string) ([]string, []string, error) {
	w := wow.New(os.Stdout, spin.Get(spin.Dots), emoji.Sprintf(":mag: Searching the blockmesh...")) // Init logger

	w.Start() // Start spinner

	defer w.Stop() // Stop spinner

	localSearchChains := searchChains // Init local search chains buffer

	var results []string     // Init search results buffer
	var resultFiles []string // Init result files buffer
	var err error            // Init error buffer

	if len(searchChains) == 0 { // Check no search chains
		localSearchChains, err = types.GetAllLocalizedChains() // Get all local chains

		if err != nil { // Check for errors
			return []string{}, []string{}, err // Return found error
		}
	}

	for _, chainName := range localSearchChains { // Iterate through search chains
		address, err := summercashCommon.StringToAddress(chainName) // Parse string addr

		if err != nil { // Check for errors
			return []string{}, []string{}, err // Return found error
		}

		chain, err := types.ReadChainFromMemory(address) // Read chain

		if err != nil { // Check for errors
			return []string{}, []string{}, err // Return found error
		}

		if chainName == searchTerm { // Check direct match
			results = append(results, chain.String())                                                                                         // Append chain string
			resultFiles = append(resultFiles, filepath.FromSlash(fmt.Sprintf("%s/db/chain/chain_%s.json", common.DataDir, address.String()))) // Append result file

			continue // Continue
		}

		for _, transaction := range chain.Transactions { // Iterate through transactions
			if transaction != nil && transaction.Hash != nil { // Check is transaction
				if transaction.Hash.String() == searchTerm || (transaction.Sender != nil && transaction.Sender.String() == searchTerm) || (transaction.Recipient != nil && transaction.Recipient.String() == searchTerm) || bytes.Contains(transaction.Payload, []byte(searchTerm)) || (transaction.Amount != nil && strings.Contains(transaction.Amount.String(), searchTerm)) { // Check match
					results = append(results, transaction.String())                                                                                   // Append transaction
					resultFiles = append(resultFiles, filepath.FromSlash(fmt.Sprintf("%s/db/chain/chain_%s.json", common.DataDir, address.String()))) // Append result file

					continue // Continue
				}
			}
		}
	}

	return results, resultFiles, nil // Return results
}

/* BEGIN INTERNAL METHODS */
