package proofs

import (
	"github.com/spf13/cobra"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/proofs"
	"github.com/tendermint/tendermint/rpc/client"
)

var TxPresenters = proofs.NewPresenters()

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
	txProver := ProofCommander{
		ProverFunc: txProver,
		Presenters: TxPresenters,
	}
	txProver.Register(txCmd)
	RootCmd.AddCommand(txCmd)
}

func txProver(node client.Client) lc.Prover {
	return proofs.NewTxProver(node)
}
