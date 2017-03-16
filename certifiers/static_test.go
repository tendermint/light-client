package certifiers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/certifiers"
	"github.com/tendermint/tendermint/types"
)

func TestStaticCert(t *testing.T) {
	// assert, require := assert.New(t), require.New(t)
	assert := assert.New(t)
	// require := require.New(t)

	keys := certifiers.GenValKeys(4)
	// 20, 30, 40, 50 - the first 3 don't have 2/3, the last 3 do!
	vals := keys.ToValidators(20, 10)
	// and a certifier based on our known set
	chainID := "test-static"
	cert := certifiers.NewStatic(chainID, vals)

	cases := []struct {
		keys        certifiers.ValKeys
		vals        []*types.Validator
		height      int
		first, last int  // who actually signs
		proper      bool // true -> expect no error
	}{
		// perfect, signed by everyone
		{keys, vals, 1, 0, len(keys), true},
		// skip little guy is okay
		{keys, vals, 2, 1, len(keys), true},
		// but not the big guy
		{keys, vals, 3, 0, len(keys) - 1, false},
		// even changing the power a little bit breaks the static validator
		// the sigs are enough, but the validator hash is unknown
		{keys, keys.ToValidators(20, 11), 4, 0, len(keys), false},
	}

	for _, tc := range cases {
		// let's make a header and sign it with everyone, must work
		header := certifiers.GenHeader(chainID, tc.height, nil, tc.vals, []byte("foo"))
		check := lc.Checkpoint{
			Header: header,
			Commit: tc.keys.SignHeader(header, tc.first, tc.last),
		}
		err := cert.Certify(check)
		if tc.proper {
			assert.Nil(err, "%+v", err)
		} else {
			assert.NotNil(err)
		}
	}

}
