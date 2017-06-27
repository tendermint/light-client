package rpc

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/light-client/commands"
)

const (
	FlagHeight = "height"
	FlagMax    = "max"
	FlagMin    = "min"
)

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Get a validated block at the given height",
	RunE:  commands.RequireInit(runBlock),
}

func init() {
	blockCmd.Flags().Int(FlagHeight, 0, "block height")
}

func runBlock(cmd *cobra.Command, args []string) error {
	c := commands.GetNode()
	h := viper.GetInt(FlagHeight)
	block, err := c.Block(h)
	if err != nil {
		return err
	}
	return printResult(block)
}
