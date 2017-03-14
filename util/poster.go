package util

import (
	keys "github.com/tendermint/go-keys"
	lc "github.com/tendermint/light-client"
)

// Poster combines KeyStore and Node to process a Signable and deliver it to tendermint
// returning the results from the tendermint node, once the transaction is processed.
//
// Only handles single signatures
type Poster struct {
	server lc.Broadcaster
	signer keys.Signer
}

func NewPoster(server lc.Broadcaster, signer keys.Signer) Poster {
	return Poster{server, signer}
}

// Post will sign the transaction with the given credentials and push it to
// the tendermint server
func (p Poster) Post(sign keys.Signable, keyname, passphrase string) (res lc.TmBroadcastResult, err error) {
	var signed []byte

	err = p.signer.Sign(keyname, passphrase, sign)
	if err != nil {
		return
	}

	signed, err = sign.TxBytes()
	if err != nil {
		return
	}

	res, err = p.server.Broadcast(signed)
	return
}
