package tests

/**
This file is base HEAVILY on tendermint/tendermint/rpc/tests/helpers.go
However, I wanted to use public variables, so this could be a basis
of tests in various packages.
**/

import (
	"fmt"

	logger "github.com/tendermint/go-logger"
	"github.com/tendermint/light-client/rpc"

	"github.com/tendermint/abci/example/dummy"
	abci "github.com/tendermint/abci/types"
	cfg "github.com/tendermint/go-config"
	"github.com/tendermint/tendermint/config/tendermint_test"
	nm "github.com/tendermint/tendermint/node"
)

var (
	config cfg.Config
)

const tmLogLevel = "error"

// GetConfig returns a config for the test cases as a singleton
func GetConfig() cfg.Config {
	if config == nil {
		config = tendermint_test.ResetConfig("rpc_test_client_test")
		// Shut up the logging
		logger.SetLogLevel(tmLogLevel)
	}
	return config
}

// GetClient gets a rpc client pointing to the test tendermint rpc
func GetClient() *rpc.HTTPClient {
	rpcAddr := GetConfig().GetString("rpc_laddr")
	return rpc.NewClient(rpcAddr, "/websocket")
}

// GetNodeClient gets a Node object pointing to this test tendermint rpc
func GetNode() rpc.Node {
	rpcAddr := GetConfig().GetString("rpc_laddr")
	chainID := GetConfig().GetString("chain_id")
	return rpc.NewNode(rpcAddr, chainID)
}

// StartTendermint starts a test tendermint server in a go routine and returns when it is initialized
// TODO: can one pass an Application in????
func StartTendermint() {
	// start a node
	fmt.Println("Start Tendermint")
	ready := make(chan struct{})

	app := dummy.NewDummyApplication()
	go NewTendermint(ready, app)
	<-ready
}

// NewTendermint creates a new tendermint server and sleeps forever
func NewTendermint(ready chan struct{}, app abci.Application) {
	// Create & start node
	config := GetConfig()
	node := nm.NewNodeDefault(config)
	// privValidator := ttypes.GenPrivValidator()
	// node := nm.NewNode(config, privValidator,
	// 	proxy.NewLocalClientCreator(app))

	// node.Start now does everything including the RPC server
	node.Start()
	ready <- struct{}{}

	// Sleep forever
	ch := make(chan struct{})
	<-ch
}
