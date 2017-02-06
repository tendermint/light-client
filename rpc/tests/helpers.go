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

// GetClient gets a rpc client pointing to the test node
func GetClient() *rpc.HTTPClient {
	rpcAddr := GetConfig().GetString("rpc_laddr")
	return rpc.New(rpcAddr, "/websocket")
}

// StartNode starts a test node in a go routine and returns when it is initialized
// TODO: can one pass an Application in????
func StartNode() {
	// start a node
	fmt.Println("StartNode")
	ready := make(chan struct{})
	go NewNode(ready)
	<-ready
}

// NewNode creates a new node and sleeps forever
func NewNode(ready chan struct{}) {
	// Create & start node
	node := nm.NewNodeDefault(GetConfig())
	// node.Start now does everything including the RPC server
	node.Start()
	ready <- struct{}{}

	// Sleep forever
	ch := make(chan struct{})
	<-ch
}
