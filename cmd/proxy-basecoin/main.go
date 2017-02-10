/*
proxy-basecoin is an example script that sets up a proxy
that knows about the basic basecoin types (sendtx and accounts).

It can be extended to add support for basecoin plugins,
or as a source of inspiration to configure a proxy for your
own app.

If you run the basecoin demo app locally (from the data dir),
try the following:

proxy-basecoin -chain=test_chain_id -rpc=localhost:46657

curl http://localhost:8108/keys/
curl -XPOST http://localhost:8108/keys/ -d '{"name": "john", "passphrase": "1234567890"}'
curl http://localhost:8108/keys/john

# -> at this point, grab that address, but it in the genesis for
# basecoin, so you are rich, and restart the basecoin server ;)

## TODO: working examples here

# query no data
curl http://localhost:8108/query/store/01234567

# 626173652f612f <- this is the magic base/a/ prefix for accounts in hex
# 1B1BE55F969F54064628A63B9559E7C21C925165 <- address with coins
when will this work????

# failing proof
curl http://localhost:8108/proof/01234567

# post a tx (not yet implemented)
curl -XPOST http://localhost:8108/txs/ -d \
  '{"name": "john", "passphrase": "1234567890", \
  "data": {"key": "value"}}'

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tendermint/light-client/cryptostore"
	"github.com/tendermint/light-client/extensions/basecoin"
	"github.com/tendermint/light-client/proxy"
	"github.com/tendermint/light-client/rpc"
	"github.com/tendermint/light-client/storage/filestorage"
)

var (
	port    = flag.Int("port", 8108, "port for proxy server")
	rpcAddr = flag.String("rpc", "", "url for tendermint rpc server")
	chainID = flag.String("chain", "", "id of the blockchain")
	keydir  = flag.String("keydir", ".keys", "directory to store private keys")
)

// TODO: add cors and unix-socket support like over in verifier
func main() {
	flag.Parse()

	// TODO: make these actually do something
	vr := basecoin.BasecoinValues{}
	sr := basecoin.BasecoinTx{}

	if *chainID == "" {
		fmt.Println("You must specify -chain with the chain_id")
		return
	}
	if *rpcAddr == "" {
		fmt.Println("You must specify -rpc with the location of a tendermint node")
		return
	}

	// set up all the pieces based on config
	r := mux.NewRouter()
	cstore := cryptostore.New(
		cryptostore.GenEd25519,
		cryptostore.SecretBox,
		filestorage.New(*keydir),
	)
	node := rpc.NewNode(*rpcAddr, *chainID, vr)
	proxy.RegisterDefault(r, cstore, node, sr, vr)

	// TODO: add some awesome middlewares...

	// now, just run the server and bind to localhost for security
	bind := fmt.Sprintf("127.0.0.1:%d", *port)
	log.Fatal(http.ListenAndServe(bind, r))

}
