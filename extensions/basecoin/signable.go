package basecoin

import lc "github.com/tendermint/light-client"

type BasecoinTx struct{}

// Turn json into a signable object
func (t BasecoinTx) ReadSignable(data []byte) (lc.Signable, error) {
	// TODO
	return nil, nil
}

func (t BasecoinTx) assertSignableReader() lc.SignableReader {
	return t
}
