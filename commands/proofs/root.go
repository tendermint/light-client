package proofs

import (
	"encoding/hex"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/tendermint/rpc/client"
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
	node client.Client
	lc.Prover
	ProverFunc func(client.Client) lc.Prover
}

// Init uses configuration info to create a network connection
// as well as initializing the prover
func (p ProofCommander) Init() {
	endpoint := viper.GetString(commands.NodeFlag)
	p.node = client.NewHTTP(endpoint, "/websockets")
	p.Prover = p.ProverFunc(p.node)
}

func (p ProofCommander) Register(parent *cobra.Command) {
	// we add each subcommand here, so we can register the
	// ProofCommander in one swoop
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
