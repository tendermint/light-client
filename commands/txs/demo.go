package txs

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	wire "github.com/tendermint/go-wire"
	lightclient "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
)

/*** this is how to build a command ***/

var DemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demo tx creation",
	RunE:  commands.RequireInit(runDemo),
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

// runDemo is an example of how to make a tx
func runDemo(cmd *cobra.Command, args []string) error {
	templ := new(DemoTx)

	// load data from json or flags
	found, err := LoadJSON(templ)
	if err != nil {
		return err
	}
	if !found {
		// parse custom flags
		templ.User = viper.GetString(UserFlag)
		templ.Age = viper.GetInt(AgeFlag)
	}

	// TODO: add this pubkey to the loaded tx somehow
	// pubkey := GetSigner()

	// Sign if needed and post.  This it the work-horse
	bres, err := SignAndPostTx(templ)
	if err != nil {
		return err
	}

	// output result
	return OutputTx(bres)
}

/*** this is the tx struct ***/

type DemoTx struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

func (d DemoTx) Bytes() []byte {
	return wire.BinaryBytes(&d)
}

func (d DemoTx) ValidateBasic() error {
	// validate both inputs here...
	if d.Age < 18 {
		return errors.New("Sorry, dude, you're too young to blockchain!")
	}
	return nil
}

// this is what we implement for a non-signable tx
var _ lightclient.Value = DemoTx{}
