package proofs

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	data "github.com/tendermint/go-data"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/tendermint/rpc/client"
)

// GetCmd creates the get command for a proof
func (p ProofCommander) GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Get a proof from the tendermint node",
		RunE:         p.getCmd,
		SilenceUsage: true,
	}
	cmd.Flags().Int(heightFlag, 0, "Height to query (skip to use latest block)")
	cmd.Flags().String(appFlag, "raw", "App to use to interpret data")
	cmd.Flags().String(keyFlag, "", "Key to query on")
	return cmd
}

func (p ProofCommander) getCmd(cmd *cobra.Command, args []string) error {
	app := viper.GetString(appFlag)

	rawkey := viper.GetString(keyFlag)
	if rawkey == "" {
		return errors.New("missing required flag: --" + keyFlag)
	}

	height := viper.GetInt(heightFlag)

	pres, err := p.Lookup(app)
	if err != nil {
		return err
	}

	// prepare the query in an app-dependent manner
	key, err := pres.MakeKey(rawkey)
	if err != nil {
		return err
	}

	//get the proof
	proof, err := p.GetProof(key, height)
	if err != nil {
		return err
	}

	info, err := pres.ParseData(proof.Data())
	if err != nil {
		return err
	}

	data, err := data.ToJSON(info)
	if err != nil {
		return err
	}

	// TODO: store the proof or do something more interesting than just printing
	fmt.Printf("Height: %d\n", proof.BlockHeight())
	fmt.Println(string(data))
	return nil
}

// GetProof performs the get command directly from the proof (not from the CLI)
func (p ProofCommander) GetProof(key []byte, height int) (proof lc.Proof, err error) {

	// instantiate the prover instance and get a proof from the server
	p.Init()
	proof, err = p.Get(key, uint64(height))
	if err != nil {
		return
	}
	ph := int(proof.BlockHeight())
	// here is the certifier, root of all knowledge
	cert, err := commands.GetCertifier()
	if err != nil {
		return
	}

	// get and validate a signed header for this proof

	// FIXME: cannot use cert.GetByHeight for now, as it also requires
	// Validators and will fail on querying tendermint for non-current height.
	// When this is supported, we should use it instead...
	client.WaitForHeight(p.node, ph, nil)
	commit, err := p.node.Commit(ph)
	if err != nil {
		return
	}
	check := lc.Checkpoint{
		Header: commit.Header,
		Commit: commit.Commit,
	}
	err = cert.Certify(check)
	if err != nil {
		return
	}

	// validate the proof against the certified header to ensure data integrity
	err = proof.Validate(check)
	if err != nil {
		return
	}

	return proof, err
}
