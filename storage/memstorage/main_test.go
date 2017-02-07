package memstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	lightclient "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/mock"
)

func TestBasicCRUD(t *testing.T) {
	assert := assert.New(t)
	store := New()

	name := "foo"
	key := []byte("secret-key-here")
	pubkey := mock.PubKey{
		Val: []byte("pubkey"),
	}
	info := lightclient.KeyInfo{
		Name:   name,
		PubKey: pubkey,
	}

	// No data: Get and Delete return nothing
	_, _, err := store.Get(name)
	assert.NotNil(err)
	err = store.Delete(name)
	assert.NotNil(err)
	// List returns empty list
	l, err := store.List()
	assert.Nil(err)
	assert.Empty(l)

	// Putting the key in the  store must work
	err = store.Put(name, key, info)
	assert.Nil(err)
	// But a second time is a failure
	err = store.Put(name, key, info)
	assert.NotNil(err)

	// Now, we can get and list properly
	k, i, err := store.Get(name)
	assert.Nil(err)
	assert.Equal(key, k)
	assert.Equal(info, i)
	l, err = store.List()
	assert.Nil(err)
	assert.Equal(1, len(l))
	assert.Equal(info, l[0])

	// querying a non-existent key fails
	_, _, err = store.Get("badname")
	assert.NotNil(err)

	// We can only delete once
	err = store.Delete(name)
	assert.Nil(err)
	err = store.Delete(name)
	assert.NotNil(err)

	// and then Get and List don't work
	_, _, err = store.Get(name)
	assert.NotNil(err)
	// List returns empty list
	l, err = store.List()
	assert.Nil(err)
	assert.Empty(l)
}
