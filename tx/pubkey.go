package tx

import crypto "github.com/tendermint/go-crypto"

type PubKey struct {
	crypto.PubKey
}

func (p *PubKey) UnmarshalJSON(b []byte) error {
	var data HexData
	err := data.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	p.PubKey, err = crypto.PubKeyFromBytes(data)
	return err
}

func (p PubKey) MarshalJSON() ([]byte, error) {
	var data []byte
	if p.PubKey != nil {
		data = p.PubKey.Bytes()
	}
	return HexData(data).MarshalJSON()
}
