// Package common defines common helper methods and variables.
package common

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/SummerCash/go-summercash/common"
)

var (
	// DataDir is the global data directory definition.
	DataDir = GetDefaultDataPath() // Get default data path
)

/* BEGIN EXPORTED METHODS */

// GetDefaultDataPath gets the default data directory.
func GetDefaultDataPath() string {
	path, err := filepath.Abs(filepath.FromSlash("./data")) // Get data path

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	if !strings.Contains(path, "puppet") { // Check not in puppet dir
		user, err := user.Current() // Get current user

		if err != nil { // Check for errors
			panic(err) // Panic
		}

		puppetPath, _ := filepath.Abs(filepath.FromSlash(fmt.Sprintf("%s/puppet", user.HomeDir))) // Set data path

		path, _ = filepath.Abs(filepath.FromSlash(fmt.Sprintf("%s/puppet/data", user.HomeDir))) // Set data path

		err = common.CreateDirIfDoesNotExist(puppetPath) // Create puppet dir

		if err != nil { // Check for errors
			panic(err) // Panic
		}
	}

	return path // Return path
}

/* END EXPORTED METHODS */
