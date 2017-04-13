package proofs

import (
	"github.com/pkg/errors"
	wire "github.com/tendermint/go-wire"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

var _ lc.Prover = TxProver{}
var _ lc.Proof = TxProof{}

// we store up to 10MB as a proof, as we need an entire block! right now
const txLimit = 10 * 1000 * 1000

// TxProver provides positive proofs of key-value pairs in the abciapp.
//
// TODO: also support negative proofs (this key is not set)
type TxProver struct {
	node client.Client
}

func NewTxProver(node client.Client) TxProver {
	return TxProver{node: node}
}

// Get tries to download a merkle hash for app state on this key from
// the tendermint node.
func (t TxProver) Get(key []byte, h uint64) (lc.Proof, error) {
	block, err := t.node.Block(int(h))
	if err != nil {
		return nil, err
	}

	// find the desired tx in the block
	txs := block.Block.Txs
	idx := txs.Index(key)
	if idx == -1 {
		return nil, errors.Errorf("Specified Tx not found in block (numtxs = %d)", len(txs))
	}
	// and build a proof for lighter storage
	proof := TxProof{
		Height: uint64(block.Block.Height),
		Proof:  txs.Proof(idx),
	}
	return proof, err
}

func (t TxProver) Unmarshal(data []byte) (pr lc.Proof, err error) {
	var proof TxProof
	err = errors.WithStack(wire.ReadBinaryBytes(data, &proof))
	return proof, err
}

// TxProof checks ALL txs for one block... we need a better way!
type TxProof struct {
	Height uint64
	Proof  types.TxProof
}

func (p *TxProof) Tx() types.Tx {
	return p.Proof.Data
}

func (p TxProof) BlockHeight() uint64 {
	return p.Height
}

func (p TxProof) Validate(check lc.Checkpoint) error {
	if uint64(check.Height()) != p.Height {
		return errors.Errorf("Trying to validate proof for block %d with header for block %d",
			p.Height, check.Height())
	}
	return p.Proof.Validate(check.Header.DataHash)
}

func (p TxProof) Marshal() ([]byte, error) {
	data := wire.BinaryBytes(p)
	return data, nil
}
