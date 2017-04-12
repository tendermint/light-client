package proofs

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	data "github.com/tendermint/go-data"
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

	// load the seed as specified
	prover := p.Prover()
	h := viper.GetInt(heightFlag)

	proof, err := prover.Get(key, uint64(h))
	if err != nil {
		return err
	}

	// TODO: plenty!
	//
	// get the header and seed to validate it (from provider)
	// validate it
	// print out some data
	// store it
	fmt.Println("Got it")
	data, err := data.ToJSON(proof)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
