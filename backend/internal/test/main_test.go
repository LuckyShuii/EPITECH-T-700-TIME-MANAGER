package test

import (
	"os"
	"testing"
)

// Is called once before all tests
func TestMain(m *testing.M) {
	// Launch the test DB container
	SetupTestDB()

	// Run all tests
	code := m.Run()

	// Clean up the container at the end
	Cleanup()

	os.Exit(code)
}
