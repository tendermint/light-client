/*
package tests contain integration tests and helper functions for testing
the RPC interface

In particular, it allows us to spin up a tendermint node in process, with
a live RPC server, which we can use to verify our rpc calls.  It provides
all data structures, enabling us to do more complex tests (like node_test.go)
that introspect the blocks themselves to validate signatures and the like.

It currently only spins up one node, it would be interesting to expand it
to multiple nodes to see the real effects of validating partially signed
blocks.
*/
package basecoin_test

import (
	"os"
	"testing"

	"github.com/tendermint/basecoin/app"
	"github.com/tendermint/basecoin/plugins/counter"
	eyes "github.com/tendermint/merkleeyes/client"
	nm "github.com/tendermint/tendermint/node"
	rpctest "github.com/tendermint/tendermint/rpc/test"
)

// bcapp can be used in test cases directly,
// to SetOption as needed for preparing data
var bcapp *app.Basecoin

var node *nm.Node

const ChainID = "lc-test-chain-id"

func TestMain(m *testing.M) {
	// start a tendermint node (and basecoin) in the background to test against
	cli := eyes.NewLocalClient("", 100)
	bcapp = app.NewBasecoin(cli)
	bcapp.SetOption("base/chainID", ChainID)

	// add counter plugins
	cntr := counter.New()
	bcapp.RegisterPlugin(cntr)

	node = rpctest.StartTendermint(bcapp)
	code := m.Run()

	// and shut down proper at the end
	node.Stop()
	node.Wait()
	os.Exit(code)
}
