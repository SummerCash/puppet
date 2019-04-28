// Package common defines common helper methods and variables.
package common

import (
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestGetDefaultDataPath tests the functionality of the GetDefaultDataPath() method.
func TestGetDefaultDataPath(t *testing.T) {
	dataPath := GetDefaultDataPath() // Get default data dir

	if dataPath == "" { // Check is not defined
		t.Fatal("default data path must be defined") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
