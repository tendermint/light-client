package txs

import "github.com/spf13/cobra"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tx",
	Short: "Create and post transactions to the node",
}

const (
	NameFlag  = "name"
	InputFlag = "input"
)

func init() {
	RootCmd.PersistentFlags().String(NameFlag, "", "name to sign the tx")
	RootCmd.PersistentFlags().String(InputFlag, "", "file with tx in json format")
}

func Register(name string, maker ReaderMaker) {
	p := NewPoster(name, maker)
	RootCmd.AddCommand(p.CreateCommand())
}
