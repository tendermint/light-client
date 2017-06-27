package rpc

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tendermint/go-wire/data"
	"github.com/tendermint/tendermint/rpc/client"

	certclient "github.com/tendermint/light-client/certifiers/client"
	"github.com/tendermint/light-client/commands"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Query the tendermint rpc, validating everything with a proof",
}

func init() {
	RootCmd.AddCommand(
		statusCmd,
		genesisCmd,
		validatorsCmd,
	)
}

func getSecureNode() (client.Client, error) {
	// First, connect a client
	c := commands.GetNode()
	cert, err := commands.GetCertifier()
	if err != nil {
		return nil, err
	}
	sc := certclient.Wrap(c, cert)
	return sc, nil
}

// printResult just writes the struct to the console, returns an error if it can't
func printResult(res interface{}) error {
	// TODO: handle text mode
	// switch viper.Get(cli.OutputFlag) {
	// case "text":
	// case "json":
	json, err := data.ToJSON(res)
	if err != nil {
		return err
	}
	fmt.Println(string(json))
	return nil
}

// // First step, proxy with no checks....
// func routes(c client.Client) map[string]*rpc.RPCFunc {

// 	return map[string]*rpc.RPCFunc{
// 		// Subscribe/unsubscribe are reserved for websocket events.
// 		// We can just use the core tendermint impl, which uses the
// 		// EventSwitch we registered in NewWebsocketManager above
// 		"subscribe":   rpc.NewWSRPCFunc(core.Subscribe, "event"),
// 		"unsubscribe": rpc.NewWSRPCFunc(core.Unsubscribe, "event"),

// 		// info API
// 		"status":     rpc.NewRPCFunc(c.Status, ""),
// 		"blockchain": rpc.NewRPCFunc(c.BlockchainInfo, "minHeight,maxHeight"),
// 		"genesis":    rpc.NewRPCFunc(c.Genesis, ""),
// 		"block":      rpc.NewRPCFunc(c.Block, "height"),
// 		"commit":     rpc.NewRPCFunc(c.Commit, "height"),
// 		"tx":         rpc.NewRPCFunc(c.Tx, "hash,prove"),
// 		"validators": rpc.NewRPCFunc(c.Validators, ""),

// 		// broadcast API
// 		"broadcast_tx_commit": rpc.NewRPCFunc(c.BroadcastTxCommit, "tx"),
// 		"broadcast_tx_sync":   rpc.NewRPCFunc(c.BroadcastTxSync, "tx"),
// 		"broadcast_tx_async":  rpc.NewRPCFunc(c.BroadcastTxAsync, "tx"),

// 		// abci API
// 		"abci_query": rpc.NewRPCFunc(c.ABCIQuery, "path,data,prove"),
// 		"abci_info":  rpc.NewRPCFunc(c.ABCIInfo, ""),
// 	}
// }
