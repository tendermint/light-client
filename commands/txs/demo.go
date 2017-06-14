package txs

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/go-wire/data"
	"github.com/tendermint/light-client/commands"
)

var DemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demo tx creation",
	RunE:  runDemo,
}

// do something like this in main.go to enable it
// txs.RootCmd.AddCommand(txs.DemoCmd)

// runDemo is not used!  This is just to serve as a demo
// of how to construct your command
func runDemo(cmd *cobra.Command, args []string) error {
	tx, err := LoadJSON(map[string]interface{}{})
	if err != nil {
		return err
	}
	if tx == nil {
		// custom flag parsing...
		tx = data.Bytes([]byte("foo"))
	}

	// TODO: add this pubkey to the loaded tx somehow
	// pubkey := GetSigner()

	packet, err := Sign(tx)
	if err != nil {
		return err
	}

	// post the bytes
	node := commands.GetNode()
	bres, err := node.BroadcastTxCommit(packet)
	if err != nil {
		return err
	}

	// output them
	return OutputTx(bres)
}
