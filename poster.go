package lightclient

import (
	keys "github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Poster combines KeyStore and Node to process a Signable and deliver it to tendermint
// returning the results from the tendermint node, once the transaction is processed.
//
// Only handles single signatures
type Poster struct {
	server client.ABCIClient
	signer keys.Signer
}

func NewPoster(server client.ABCIClient, signer keys.Signer) Poster {
	return Poster{server, signer}
}

// Post will sign the transaction with the given credentials and push it to
// the tendermint server
func (p Poster) Post(sign keys.Signable, keyname, passphrase string) (*ctypes.ResultBroadcastTxCommit, error) {
	var signed []byte

	err := p.signer.Sign(keyname, passphrase, sign)
	if err != nil {
		return nil, err
	}

	signed, err = sign.TxBytes()
	if err != nil {
		return nil, err
	}

	return p.server.BroadcastTxCommit(signed)
}
