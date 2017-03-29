package proofs

import (
	"encoding/hex"
	"errors"

	"github.com/spf13/cobra"
	lc "github.com/tendermint/light-client"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "proof",
	Short: "Get and store merkle proofs for blockchain data",
	Long: `Proofs allows you to validate data and merkle proofs.

These proofs tie the data to a checkpoint, which is managed by "seeds".
Here we can validate these proofs and import/export them to prove specific
data to other peers as needed.
`,
}

type ProofCommander struct {
	Prover func() lc.Prover
}

func (p ProofCommander) Register(parent *cobra.Command) {
	parent.AddCommand(p.GetCmd())
}

const (
	heightFlag = "height"
)

func getHexArg(args []string) ([]byte, error) {
	if len(args) != 1 || len(args[0]) == 0 {
		return nil, errors.New("You must provide exactly one arg in hex")
	}
	bytes, err := hex.DecodeString(args[0])
	return bytes, err
}
