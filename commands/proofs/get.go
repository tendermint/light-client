package proofs

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	data "github.com/tendermint/go-wire/data"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/tendermint/rpc/tendermint/client"
)

func (p ProofCommander) GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a proof from the tendermint node",
		RunE:  p.doGet,
	}
	cmd.Flags().Int(heightFlag, 0, "Height to query (skip to use latest block)")
	cmd.Flags().String(appFlag, "raw", "App to use to interpret data")
	cmd.Flags().String(keyFlag, "", "Key to query on")
	return cmd
}

func (p ProofCommander) doGet(cmd *cobra.Command, args []string) error {
	app := viper.GetString(appFlag)
	pres, err := p.Lookup(app)
	if err != nil {
		return err
	}

	rawkey := viper.GetString(keyFlag)
	if rawkey == "" {
		return errors.New("missing required flag: --" + keyFlag)
	}

	// prepare the query in an app-dependent manner
	key, err := pres.MakeKey(rawkey)
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
	ph := int(proof.BlockHeight())

	// here is the certifier, root of all knowledge
	cert, err := commands.GetCertifier()
	if err != nil {
		return err
	}

	// get and validate a signed header for this proof

	// FIXME: cannot use cert.GetByHeight for now, as it also requires
	// Validators and will fail on querying tendermint for non-current height.
	// When this is supported, we should use it instead...
	client.WaitForHeight(p.node, ph, nil)
	commit, err := p.node.Commit(ph)
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
	fmt.Printf("Height: %d\n", proof.BlockHeight())
	info, err := pres.ParseData(proof.Data())
	if err != nil {
		return err
	}
	data, err := data.ToJSON(info)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
