package proofs

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-wire/data"

	"github.com/tendermint/light-client/proofs"
)

const (
	heightFlag = "height"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "query",
	Short: "Get and store merkle proofs for blockchain data",
	Long: `Proofs allows you to validate data and merkle proofs.

These proofs tie the data to a checkpoint, which is managed by "seeds".
Here we can validate these proofs and import/export them to prove specific
data to other peers as needed.
`,
}

func init() {
	RootCmd.Flags().Int(heightFlag, 0, "Height to query (skip to use latest block)")
}

// ParseHexKey parses the key flag as hex and converts to bytes or returns error
// argname is used to customize the error message
func ParseHexKey(args []string, argname string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.Errorf("Missing required argument [%s]", argname)
	}
	if len(args) > 1 {
		return nil, errors.Errorf("Only accepts one argument [%s]", argname)
	}
	rawkey := args[0]
	if rawkey == "" {
		return nil, errors.Errorf("[%s] argument must be non-empty ", argname)
	}
	// with tx, we always just parse key as hex and use to lookup
	return proofs.ParseHexKey(rawkey)
}

func GetHeight() int {
	return viper.GetInt(heightFlag)
}

// OutputProof prints the proof to stdout
// reuse this for printing proofs and we should enhance this for text/json,
// better presentation of height
func OutputProof(info interface{}, height uint64) error {
	res, err := data.ToJSON(info)
	if err != nil {
		return err
	}

	// TODO: store the proof or do something more interesting than just printing
	fmt.Printf("Height: %d\n", height)
	fmt.Println(string(res))
	return nil
}
