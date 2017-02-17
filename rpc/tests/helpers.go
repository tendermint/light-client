package tests

/**
This file is base HEAVILY on tendermint/tendermint/rpc/tests/helpers.go
However, I wanted to use public variables, so this could be a basis
of tests in various packages.
**/

import (
	"fmt"

	logger "github.com/tendermint/go-logger"
	"github.com/tendermint/light-client/mock"
	"github.com/tendermint/light-client/rpc"

	abci "github.com/tendermint/abci/types"
	cfg "github.com/tendermint/go-config"
	"github.com/tendermint/tendermint/config/tendermint_test"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
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
	return rpc.NewNode(rpcAddr, chainID, mock.ValueReader())
}

// StartTendermint starts a test tendermint server in a go routine and returns when it is initialized
// TODO: can one pass an Application in????
func StartTendermint(app abci.Application) *nm.Node {
	// start a node
	fmt.Println("Starting Tendermint...")

	node := NewTendermint(app)
	fmt.Println("Tendermint running!")
	return node
}

// NewTendermint creates a new tendermint server and sleeps forever
func NewTendermint(app abci.Application) *nm.Node {
	// Create & start node
	config := GetConfig()
	privValidatorFile := config.GetString("priv_validator_file")
	privValidator := types.LoadOrGenPrivValidator(privValidatorFile)
	papp := proxy.NewLocalClientCreator(app)
	node := nm.NewNode(config, privValidator, papp)

	// node.Start now does everything including the RPC server
	node.Start()
	return node
}
