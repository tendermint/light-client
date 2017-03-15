package certifiers

import (
	merkle "github.com/tendermint/go-merkle"
	lc "github.com/tendermint/light-client"
)

/*** TODO: where does this go??? ***/

// MerkleReader is currently the only implementation of ProofReader,
// using the IAVLProof from go-merkle
var MerkleReader lc.ProofReader = merkleReader{}

type merkleReader struct{}

func (p merkleReader) ReadProof(data []byte) (lc.Proof, error) {
	return merkle.ReadProof(data)
}
