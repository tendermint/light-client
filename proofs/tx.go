package proofs

import (
	"bytes"

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

	proof := TxProof{
		Height: uint64(block.Block.Height),
		Txs:    block.Block.Txs,
	}
	// sets the index if key is there, otherwise sets error
	err = proof.SetTx(key)
	return proof, err
}

func (t TxProver) Unmarshal(data []byte) (lc.Proof, error) {
	var proof TxProof
	var n int
	var err error
	wire.ReadBinary(&proof, bytes.NewBuffer(data), txLimit, &n, &err)
	return proof, errors.WithStack(err)
}

// TxProof checks ALL txs for one block... we need a better way!
type TxProof struct {
	Height uint64
	Index  int
	Txs    types.Txs
}

// Tx returns the selected transaction from the block
func (p TxProof) Tx() []byte {
	return p.Txs[p.Index]
}

// SetTx sets the index for the specified tx, returns an error if not present
func (p *TxProof) SetTx(tx []byte) error {
	for i, t := range p.Txs {
		if bytes.Equal(tx, t) {
			p.Index = i
			return nil
		}
	}
	return errors.Errorf("Specified Tx not found in block (numtxs = %d)", len(p.Txs))
}

func (p TxProof) BlockHeight() uint64 {
	return p.Height
}

func (p TxProof) Validate(check lc.Checkpoint) error {
	if uint64(check.Height()) != p.Height {
		return errors.Errorf("Trying to validate proof for block %d with header for block %d",
			p.Height, check.Height())
	}

	hash := p.Txs.Hash()
	if !bytes.Equal(check.Header.DataHash, hash) {
		return errors.Errorf("Hash mismatch: checkpoint = %X, proof = %X",
			check.Header.DataHash, hash)
	}

	return nil
}

func (p TxProof) Marshal() ([]byte, error) {
	data := wire.BinaryBytes(p)
	return data, nil
}

// TODO: one tx plus proof.... need changes in the
// func (c Checkpoint) CheckTx(tx types.Tx) error {
//   return nil
// }
