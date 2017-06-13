package proofs

import (
	"github.com/tendermint/tendermint/rpc/client"

	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
)

// GetProof performs the get command directly from the proof (not from the CLI)
func GetProof(node client.Client, prover lc.Prover, key []byte, height int) (proof lc.Proof, err error) {
	proof, err = prover.Get(key, uint64(height))
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
	client.WaitForHeight(node, ph, nil)
	commit, err := node.Commit(ph)
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
