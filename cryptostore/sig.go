package cryptostore

import (
	"github.com/pkg/errors"
	wire "github.com/tendermint/go-wire"
	lightclient "github.com/tendermint/light-client"
)

func NewSig(data []byte) *Sig {
	return &Sig{data: data}
}

// Sig lets us wrap arbitrary data with a go-crypto signature
//
// TODO: rethink how we want to integrate this with KeyStore so it makes
// more sense (particularly the verify method)
type Sig struct {
	data   []byte
	sig    []byte
	pubkey []byte
}

// assertSignable is just to make sure we stay in sync with the Signable interface
func (s *Sig) assertSignable() lightclient.Signable {
	return s
}

// Bytes returns the original data passed into `NewSig`
func (s *Sig) Bytes() []byte {
	return s.data
}

// Sign attaches information from `KeySigner.Signature`, which should be
// a properly encoded public key signature
func (s *Sig) Sign(addr, pubkey, sig []byte) error {
	if len(sig) == 0 || len(pubkey) == 0 {
		return errors.New("Signature or Key missing")
	}
	if s.sig != nil {
		return errors.New("Transaction can only be signed once")
	}
	s.sig = sig
	s.pubkey = pubkey
	return nil
}

// Signed serializes the Sig to send it to a tendermint app.
// It returns an error if the Sig was never Signed.
func (s *Sig) Signed() ([]byte, error) {
	if s.sig == nil {
		return nil, errors.New("Transaction was never signed")
	}
	return wire.BinaryBytes(s), nil
}

// TODO: how do we verify this??? need some function to deserialize and verify the sigs!

// // Validate will deserialize the contained action, and validate the signature or return an error
// func (tx SignedAction) Validate() (ValidatedAction, error) {
// 	res := ValidatedAction{
// 		SignedAction: tx,
// 	}
// 	valid := tx.Signer.VerifyBytes(tx.ActionData, tx.Signature)
// 	if !valid {
// 		return res, errors.New("Invalid signature")
// 	}

// 	var err error
// 	res.action, err = ActionFromBytes(tx.ActionData)
// 	if err == nil {
// 		res.valid = true
// 	}
// 	return res, err
// }

// // SignAction will serialize the action and sign it with your key
// func SignAction(action Action, privKey crypto.PrivKey) (res SignedAction, err error) {
// 	res.ActionData, err = ActionToBytes(action)
// 	if err != nil {
// 		return res, err
// 	}
// 	res.Signature = privKey.Sign(res.ActionData)
// 	res.Signer = privKey.PubKey()
// 	return res, nil
// }
