package proofs

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	data "github.com/tendermint/go-data"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
)

func (p ProofCommander) GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a proof from the tendermint node",
		RunE:  p.doGet,
	}
	cmd.Flags().Int(heightFlag, 0, "Height to query (skip to use latest block)")
	return cmd
}

func (p ProofCommander) doGet(cmd *cobra.Command, args []string) error {
	key, err := getHexArg(args)
	if err != nil {
		return err
	}

	// instantiate the prover instance and get a proof from the server
	p.Init()
	h := viper.GetInt(heightFlag)
	proof, err := p.Get(key, uint64(h))
	if err != nil {
		return err
	}

	// here is the certifier, root of all knowledge
	cert, err := commands.GetCertifier()
	if err != nil {
		return err
	}

	// get and validate a signed header for this proof

	// FIXME: cannot use cert.GetByHeight for now, as it also requires
	// Validators and will fail on querying tendermint for non-current height.
	// When this is supported, we should use it instead...
	commit, err := p.node.Commit(h)
	if err != nil {
		return err
	}
	check := lc.Checkpoint{commit.Header, commit.Commit}
	err = cert.Certify(check)
	if err != nil {
		return err
	}

	// validate the proof against the certified header to ensure data integrity
	err = proof.Validate(check)
	if err != nil {
		return err
	}

	// TODO: store the proof or do something more interesting than just printing
	fmt.Println("Your data is 100% certified:")
	data, err := data.ToJSON(proof)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
