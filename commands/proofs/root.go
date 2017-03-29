package proofs

import (
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

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Handle proofs for state of abci app",
	Long: `Proofs allows you to validate abci state with merkle proofs.

These proofs tie the data to a checkpoint, which is managed by "seeds".
Here we can validate these proofs and import/export them to prove specific
data to other peers as needed.
`,
}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Handle proofs of commited txs",
	Long: `Proofs allows you to validate abci state with merkle proofs.

These proofs tie the data to a checkpoint, which is managed by "seeds".
Here we can validate these proofs and import/export them to prove specific
data to other peers as needed.
`,
}

func init() {
	RootCmd.AddCommand(stateCmd)
	RootCmd.AddCommand(txCmd)
}

type Prover interface {
	Get(key []byte, h int) (Proof, error)
	Unmarshal([]byte) (Proof, error)
}

type Proof interface {
	BlockHeight() int
	Validate(lc.Checkpoint) error // Make sure the checkpoint is validated and proper height
	Marshal() ([]byte, error)
}
