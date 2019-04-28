// Package common defines common helper methods and variables.
package common

import (
	"path/filepath"
)

var (
	// DataDir is the global data directory definition.
	DataDir = getDefaultDataPath() // Get default data path
)

/* BEGIN EXPORTED METHODS */

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// getDefaultDataPath gets the default data directory.
func getDefaultDataPath() string {
	path, _ := filepath.Abs(filepath.FromSlash("./data")) // Get data path

	return path // Return path
}

/* END INTERNAL METHODS */