package txs

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/light-client/commands"
)

const (
	NameFlag  = "name"
	InputFlag = "input"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               "tx",
	Short:             "Create and post transactions to the node",
	PersistentPreRunE: commands.RequireInit,
}

func init() {
	RootCmd.PersistentFlags().String(NameFlag, "", "name to sign the tx")
	RootCmd.PersistentFlags().String(InputFlag, "", "file with tx in json format")
}
