

# proxy
`import "github.com/tendermint/light-client/proxy"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
package proxy provides http handlers to construct a proxy server
for key management, transaction signing, and query validation.

Please read the README and godoc to see how to
configure the server for your application.




## <a name="pkg-index">Index</a>
* [func RegisterDefault(r *mux.Router, keys KeyStore, node rpc.Node, txReader lc.SignableReader, valReader lc.ValueReader)](#RegisterDefault)
* [type KeyServer](#KeyServer)
  * [func NewKeyServer(manager lc.KeyManager) KeyServer](#NewKeyServer)
  * [func (k KeyServer) DeleteKey(w http.ResponseWriter, r *http.Request)](#KeyServer.DeleteKey)
  * [func (k KeyServer) GenerateKey(w http.ResponseWriter, r *http.Request)](#KeyServer.GenerateKey)
  * [func (k KeyServer) GetKey(w http.ResponseWriter, r *http.Request)](#KeyServer.GetKey)
  * [func (k KeyServer) ListKeys(w http.ResponseWriter, r *http.Request)](#KeyServer.ListKeys)
  * [func (k KeyServer) Register(r *mux.Router)](#KeyServer.Register)
  * [func (k KeyServer) UpdateKey(w http.ResponseWriter, r *http.Request)](#KeyServer.UpdateKey)
* [type KeyStore](#KeyStore)
* [type TxSigner](#TxSigner)
  * [func NewTxSigner(server lc.Broadcaster, signer lc.Signer, reader lc.SignableReader) TxSigner](#NewTxSigner)
  * [func (t TxSigner) PostTransaction(w http.ResponseWriter, r *http.Request)](#TxSigner.PostTransaction)
  * [func (t TxSigner) Register(r *mux.Router)](#TxSigner.Register)
* [type Viewer](#Viewer)
  * [func NewViewer(checker lc.Checker, searcher lc.Searcher, cert lc.Certifier) Viewer](#NewViewer)
  * [func (v Viewer) ProveData(w http.ResponseWriter, r *http.Request)](#Viewer.ProveData)
  * [func (v Viewer) QueryData(w http.ResponseWriter, r *http.Request)](#Viewer.QueryData)
  * [func (v Viewer) Register(r *mux.Router)](#Viewer.Register)


#### <a name="pkg-files">Package files</a>
[helpers.go](/src/github.com/tendermint/light-client/proxy/helpers.go) [keys.go](/src/github.com/tendermint/light-client/proxy/keys.go) [tx.go](/src/github.com/tendermint/light-client/proxy/tx.go) [valid.go](/src/github.com/tendermint/light-client/proxy/valid.go) [viewer.go](/src/github.com/tendermint/light-client/proxy/viewer.go) 





## <a name="RegisterDefault">func</a> [RegisterDefault](/src/target/helpers.go?s=781:901#L25)
``` go
func RegisterDefault(r *mux.Router, keys KeyStore, node rpc.Node,
    txReader lc.SignableReader, valReader lc.ValueReader)
```
RegisterDefault constructs all components and wires them up under
standard routes.

TODO: something more intelligent for getting validators,
this is pretty insecure right now




## <a name="KeyServer">type</a> [KeyServer](/src/target/keys.go?s=215:263#L4)
``` go
type KeyServer struct {
    // contains filtered or unexported fields
}
```






### <a name="NewKeyServer">func</a> [NewKeyServer](/src/target/keys.go?s=265:315#L8)
``` go
func NewKeyServer(manager lc.KeyManager) KeyServer
```




### <a name="KeyServer.DeleteKey">func</a> (KeyServer) [DeleteKey](/src/target/keys.go?s=1821:1889#L94)
``` go
func (k KeyServer) DeleteKey(w http.ResponseWriter, r *http.Request)
```



### <a name="KeyServer.GenerateKey">func</a> (KeyServer) [GenerateKey](/src/target/keys.go?s=363:433#L14)
``` go
func (k KeyServer) GenerateKey(w http.ResponseWriter, r *http.Request)
```



### <a name="KeyServer.GetKey">func</a> (KeyServer) [GetKey](/src/target/keys.go?s=789:854#L38)
``` go
func (k KeyServer) GetKey(w http.ResponseWriter, r *http.Request)
```



### <a name="KeyServer.ListKeys">func</a> (KeyServer) [ListKeys](/src/target/keys.go?s=1035:1102#L51)
``` go
func (k KeyServer) ListKeys(w http.ResponseWriter, r *http.Request)
```



### <a name="KeyServer.Register">func</a> (KeyServer) [Register](/src/target/keys.go?s=2392:2434#L122)
``` go
func (k KeyServer) Register(r *mux.Router)
```



### <a name="KeyServer.UpdateKey">func</a> (KeyServer) [UpdateKey](/src/target/keys.go?s=1241:1309#L63)
``` go
func (k KeyServer) UpdateKey(w http.ResponseWriter, r *http.Request)
```



## <a name="KeyStore">type</a> [KeyStore](/src/target/helpers.go?s=537:590#L15)
``` go
type KeyStore interface {
    lc.KeyManager
    lc.Signer
}
```
KeyStore is implemented by cryptostore.Manager










## <a name="TxSigner">type</a> [TxSigner](/src/target/tx.go?s=200:256#L2)
``` go
type TxSigner struct {
    lc.SignableReader
    util.Poster
}
```






### <a name="NewTxSigner">func</a> [NewTxSigner](/src/target/tx.go?s=258:351#L7)
``` go
func NewTxSigner(server lc.Broadcaster, signer lc.Signer,
    reader lc.SignableReader) TxSigner
```




### <a name="TxSigner.PostTransaction">func</a> (TxSigner) [PostTransaction](/src/target/tx.go?s=455:528#L16)
``` go
func (t TxSigner) PostTransaction(w http.ResponseWriter, r *http.Request)
```



### <a name="TxSigner.Register">func</a> (TxSigner) [Register](/src/target/tx.go?s=887:928#L40)
``` go
func (t TxSigner) Register(r *mux.Router)
```



## <a name="Viewer">type</a> [Viewer](/src/target/viewer.go?s=243:320#L5)
``` go
type Viewer struct {
    lc.Checker
    lc.Searcher
    lc.ValueReader
    util.Auditor
}
```






### <a name="NewViewer">func</a> [NewViewer](/src/target/viewer.go?s=322:405#L12)
``` go
func NewViewer(checker lc.Checker, searcher lc.Searcher,
    cert lc.Certifier) Viewer
```




### <a name="Viewer.ProveData">func</a> (Viewer) [ProveData](/src/target/viewer.go?s=1072:1137#L51)
``` go
func (v Viewer) ProveData(w http.ResponseWriter, r *http.Request)
```



### <a name="Viewer.QueryData">func</a> (Viewer) [QueryData](/src/target/viewer.go?s=508:573#L21)
``` go
func (v Viewer) QueryData(w http.ResponseWriter, r *http.Request)
```



### <a name="Viewer.Register">func</a> (Viewer) [Register](/src/target/viewer.go?s=2061:2100#L102)
``` go
func (v Viewer) Register(r *mux.Router)
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
