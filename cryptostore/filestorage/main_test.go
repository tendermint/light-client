package filestorage

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lightclient "github.com/tendermint/light-client"
)

func TestBasicCRUD(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "filestorage-test")
	assert.Nil(err)
	defer os.RemoveAll(dir)
	store := New(dir)

	name := "bar"
	key := []byte("secret-key-here")
	info := lightclient.KeyInfo{
		Name:    name,
		Address: []byte("addr"),
		PubKey:  []byte("pubkey"),
	}

	// No data: Get and Delete return nothing
	_, _, err = store.Get(name)
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

	// Now, we can get and list properly
	k, i, err := store.Get(name)
	assert.Nil(err)
	assert.Equal(key, k)
	assert.Equal(info, i)
	l, err = store.List()
	assert.Nil(err, "%+v", err)
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

func TestDirectoryHandling(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	// prepare a temp dir and make sure it is not there
	newDir := path.Join(os.TempDir(), "file-test-dir")
	_, err := os.Open(newDir)
	assert.True(os.IsNotExist(err))

	// create a new storage, and verify it creates the directory with good permissions
	New(newDir)
	defer os.RemoveAll(newDir)
	d, err := os.Open(newDir)
	require.Nil(err)
	defer d.Close()

	stat, err := d.Stat()
	require.Nil(err)
	assert.Equal(dirPerm, stat.Mode()&os.ModePerm)
}
