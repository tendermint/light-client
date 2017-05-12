package proxy

import (
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/tendermint/rpc/client"
	rpc "github.com/tendermint/tendermint/rpc/lib/server"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Run proxy server, verifying tendermint rpc",
	Long: `This node will run a secure proxy to a tendermint rpc server.

All calls that can be tracked back to a block header by a proof
will be verified before passing them back to the caller. Other that
that it will present the same interface as a full tendermint node,
just with added trust and running locally.`,
	RunE:         runProxy,
	SilenceUsage: true,
}

const (
	bindFlag   = "serve"
	wsEndpoint = "/websocket"
)

func init() {
	RootCmd.Flags().String(bindFlag, ":8888", "Serve the proxy on the given port")
}

func runProxy(cmd *cobra.Command, args []string) error {
	// First, connect a client
	c := client.NewHTTP(viper.GetString(commands.NodeFlag), "/websocket")
	r := routes(c)

	// build the handler...
	mux := http.NewServeMux()
	rpc.RegisterRPCFuncs(mux, r)
	wm := rpc.NewWebsocketManager(r, nil)
	// wm.SetLogger(log.TestingLogger())
	mux.HandleFunc(wsEndpoint, wm.WebsocketHandler)

	// TODO: pass in a proper logger
	_, err := rpc.StartHTTPServer(viper.GetString(bindFlag), mux)
	if err != nil {
		return nil
	}
	// uhh... better way?
	time.Sleep(1000 * time.Minute)

	return err
}

// First step, proxy with no checks....
func routes(c *client.HTTP) map[string]*rpc.RPCFunc {
	// // subscribe/unsubscribe are reserved for websocket events.
	//  "subscribe":   rpc.NewWSRPCFunc(Subscribe, "event"),
	//  "unsubscribe": rpc.NewWSRPCFunc(Unsubscribe, "event"),

	return map[string]*rpc.RPCFunc{
		// info API
		"status":     rpc.NewRPCFunc(c.Status, ""),
		"net_info":   rpc.NewRPCFunc(c.NetInfo, ""),
		"blockchain": rpc.NewRPCFunc(c.BlockchainInfo, "minHeight,maxHeight"),
		"genesis":    rpc.NewRPCFunc(c.Genesis, ""),
		"block":      rpc.NewRPCFunc(c.Block, "height"),
		"commit":     rpc.NewRPCFunc(c.Commit, "height"),
		"tx":         rpc.NewRPCFunc(c.Tx, "hash,prove"),
		"validators": rpc.NewRPCFunc(c.Validators, ""),

		// broadcast API
		"broadcast_tx_commit": rpc.NewRPCFunc(c.BroadcastTxCommit, "tx"),
		"broadcast_tx_sync":   rpc.NewRPCFunc(c.BroadcastTxSync, "tx"),
		"broadcast_tx_async":  rpc.NewRPCFunc(c.BroadcastTxAsync, "tx"),

		// abci API
		"abci_query": rpc.NewRPCFunc(c.ABCIQuery, "path,data,prove"),
		"abci_info":  rpc.NewRPCFunc(c.ABCIInfo, ""),
	}
}

// Step two, add some checks....
