package cryptostore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"
)

func TestSortKeys(t *testing.T) {
	assert := assert.New(t)

	gen := GenEd25519.Generate
	assert.NotEqual(gen(), gen())

	// alphabetical order is n3, n1, n2
	n1, n2, n3 := "john", "mike", "alice"
	infos := SortKeys{
		info(n1, gen()),
		info(n2, gen()),
		info(n3, gen()),
	}

	// make sure they are initialized unsorted
	assert.Equal(n1, infos[0].Name)
	assert.Equal(n2, infos[1].Name)
	assert.Equal(n3, infos[2].Name)

	// now they are sorted
	infos.Sort()
	assert.Equal(n3, infos[0].Name)
	assert.Equal(n1, infos[1].Name)
	assert.Equal(n2, infos[2].Name)

	// make sure info put some real data there...
	assert.NotEmpty(infos[0].PubKey)
	assert.NotEqual(infos[0].PubKey, infos[1].PubKey)
	assert.NotEqual(infos[0].Address, infos[1].Address)

	// and make sure the pubkey is really something we can use
	_, err := crypto.PubKeyFromBytes(infos[2].PubKey)
	assert.Nil(err, "%+v", err)
}
