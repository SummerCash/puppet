// Package cli defines helpful cli helper methods.
package cli

import (
	"github.com/SummerCash/puppet/common"
	"github.com/urfave/cli"
	"strings"

	"github.com/tcnksm/go-input"

	i "github.com/tockins/interact"

	summercashCommon "github.com/SummerCash/go-summercash/common"
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

// searchBlockmesh handles the search command.
func (app *CLI) searchBlockmesh(c *cli.Context) error {
	summercashCommon.Silent = true // Silence logsconfigPath

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
			return err // Return found error
		}

		if shouldGetSearchChains { // Check must get search chains
			searchChainsString, err := app.InputConfig.Ask("What chains would you like to search?", &input.Options{
				Required:  true, // Make optional
				HideOrder: true, // Hide extra question
			})

			if err != nil { // Check for errors
				return err // Return found error
			}

			searchChains = strings.Split(searchChainsString, ", ") // Split
		}
	}

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */
