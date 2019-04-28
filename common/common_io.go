// Package common defines common helper methods and variables.
package common

import (
	"path/filepath"
)

var (
	// DataDir is the global data directory definition.
	DataDir = GetDefaultDataPath() // Get default data path
)

/* BEGIN EXPORTED METHODS */

// GetDefaultDataPath gets the default data directory.
func GetDefaultDataPath() string {
	path, _ := filepath.Abs(filepath.FromSlash("./data")) // Get data path

	return path // Return path
}

/* END EXPORTED METHODS */