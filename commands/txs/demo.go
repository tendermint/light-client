package txs

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	wire "github.com/tendermint/go-wire"
	lightclient "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
)

var DemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demo tx creation",
	RunE:  runDemo,
}

const (
	UserFlag = "user"
	AgeFlag  = "age"
)

// do something like this in main.go to enable it
// txs.RootCmd.AddCommand(txs.DemoCmd)
func init() {
	DemoCmd.Flags().String(UserFlag, "", "username you want")
	DemoCmd.Flags().Int(AgeFlag, 0, "your age... for real like")
}

type DemoTx struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

func (d DemoTx) Bytes() []byte {
	return wire.BinaryBytes(&d)
}

// this is what we implement for a non-signable tx
var _ lightclient.Value = DemoTx{}

// runDemo is an example of how to make a tx
func runDemo(cmd *cobra.Command, args []string) error {
	var templ DemoTx
	tx, err := LoadJSON(&templ)
	if err != nil {
		return err
	}
	if tx == nil {
		// parse custom flags
		templ.User = viper.GetString(UserFlag)
		templ.Age = viper.GetInt(AgeFlag)
		if templ.Age < 18 {
			return errors.New("Sorry, dude, you're too young to blockchain!")
		}
		tx = templ
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
