package proofs

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/light-client/proofs"
	"github.com/tendermint/tendermint/rpc/client"
)

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
		Prover: txProver,
	}
	txProver.Register(txCmd)
	RootCmd.AddCommand(txCmd)
}

func txProver() lc.Prover {
	endpoint := viper.GetString(commands.NodeFlag)
	node := client.NewHTTP(endpoint, "/websockets")
	return proofs.NewTxProver(node)
}
