package proxy_test

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/tendermint/light-client/cryptostore"
	"github.com/tendermint/light-client/proxy"
	"github.com/tendermint/light-client/storage/memstorage"
)

func TestKeyServer(t *testing.T) {
	// assert, require := assert.New(t), require.New(t)

	// make the storage with reasonable defaults
	cstore := cryptostore.New(
		cryptostore.GenSecp256k1,
		cryptostore.SecretBox,
		memstorage.New(),
	)

	// build your http server
	ks := proxy.NewKeyServer(cstore)
	r := mux.NewRouter()
	sk := r.PathPrefix("/keys").Subrouter()
	ks.Register(sk)

	// TODO: http test server and client

	// n1, n2, n3 := "personal", "business", "other"
	// p1, p2 := "1234", "really-secure!@#$"

	// list = 0
	// create n1, n2
	// list = 2 (in order)
	// get works
	// delete with proper key
	// update works with proper key

}
