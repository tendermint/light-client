package proofs

import (
	"github.com/spf13/cobra"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/proofs"
	"github.com/tendermint/tendermint/rpc/tendermint/client"
)

var StatePresenters = proofs.NewPresenters()

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Handle proofs for state of abci app",
	Long: `Proofs allows you to validate abci state with merkle proofs.

These proofs tie the data to a checkpoint, which is managed by "seeds".
Here we can validate these proofs and import/export them to prove specific
data to other peers as needed.
`,
}

func init() {
	stateProver := ProofCommander{
		ProverFunc: stateProver,
		Presenters: StatePresenters,
	}
	stateProver.Register(stateCmd)
	RootCmd.AddCommand(stateCmd)
}

func stateProver(node client.Client) lc.Prover {
	return proofs.NewAppProver(node)
}
