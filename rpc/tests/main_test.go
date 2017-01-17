package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// start a tendermint node (and dummy app) in the background to test against
	StartNode()
	os.Exit(m.Run())
}
