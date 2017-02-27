

# rpc
`import "github.com/tendermint/light-client/rpc"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
package rpc provides higher-level functionality to interact with tendermint's
rpc server.

The basic client implementation is HTTPClient, which provides full, direct
access to the rpc server, providing type-safety and marshaling/unmarshaling,
but no additional functionality.

The more advanced client implementation is Node, which provides a few high-level
actions based upon HTTPClient and parsing, processing, and validating the
return values.  Node currently implements Broadcaster, Checker, and Searcher
interfaces from lightclient, which are the general high-level actions
one wants to perform on tendermint.

Higher-level functionality should be built upon Node, extending Node as needed,
or defining another type, but not directly on HTTPClient.  This package provides
the bridge between RPC calls and higher-level functionality.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func ImportSignedHeader(chainID string, data []byte) (lc.TmSignedHeader, error)](#ImportSignedHeader)
* [type HTTPClient](#HTTPClient)
  * [func NewClient(remote, wsEndpoint string) *HTTPClient](#NewClient)
  * [func (c *HTTPClient) ABCIInfo() (*ctypes.ResultABCIInfo, error)](#HTTPClient.ABCIInfo)
  * [func (c *HTTPClient) ABCIQuery(path string, data []byte, prove bool) (*ctypes.ResultABCIQuery, error)](#HTTPClient.ABCIQuery)
  * [func (c *HTTPClient) Block(height int) (*ctypes.ResultBlock, error)](#HTTPClient.Block)
  * [func (c *HTTPClient) BlockchainInfo(minHeight, maxHeight int) (*ctypes.ResultBlockchainInfo, error)](#HTTPClient.BlockchainInfo)
  * [func (c *HTTPClient) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error)](#HTTPClient.BroadcastTxAsync)
  * [func (c *HTTPClient) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error)](#HTTPClient.BroadcastTxCommit)
  * [func (c *HTTPClient) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error)](#HTTPClient.BroadcastTxSync)
  * [func (c *HTTPClient) DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error)](#HTTPClient.DialSeeds)
  * [func (c *HTTPClient) Genesis() (*ctypes.ResultGenesis, error)](#HTTPClient.Genesis)
  * [func (c *HTTPClient) GetEventChannels() (chan json.RawMessage, chan error)](#HTTPClient.GetEventChannels)
  * [func (c *HTTPClient) NetInfo() (*ctypes.ResultNetInfo, error)](#HTTPClient.NetInfo)
  * [func (c *HTTPClient) StartWebsocket() error](#HTTPClient.StartWebsocket)
  * [func (c *HTTPClient) Status() (*ctypes.ResultStatus, error)](#HTTPClient.Status)
  * [func (c *HTTPClient) StopWebsocket()](#HTTPClient.StopWebsocket)
  * [func (c *HTTPClient) Subscribe(event string) error](#HTTPClient.Subscribe)
  * [func (c *HTTPClient) Unsubscribe(event string) error](#HTTPClient.Unsubscribe)
  * [func (c *HTTPClient) Validators() (*ctypes.ResultValidators, error)](#HTTPClient.Validators)
* [type Node](#Node)
  * [func NewNode(rpcAddr, chainID string, valuer lc.ValueReader) Node](#NewNode)
  * [func (n Node) Broadcast(tx []byte) (res lc.TmBroadcastResult, err error)](#Node.Broadcast)
  * [func (n Node) ExportSignedHeader(height uint64) ([]byte, error)](#Node.ExportSignedHeader)
  * [func (n Node) ImportSignedHeader(data []byte) (lc.TmSignedHeader, error)](#Node.ImportSignedHeader)
  * [func (n Node) Prove(key []byte) (lc.TmQueryResult, error)](#Node.Prove)
  * [func (n Node) Query(path string, data []byte) (lc.TmQueryResult, error)](#Node.Query)
  * [func (n Node) SignedHeader(height uint64) (lc.TmSignedHeader, error)](#Node.SignedHeader)
  * [func (n Node) Validators() (lc.TmValidatorResult, error)](#Node.Validators)
  * [func (n Node) WaitForHeight(height uint64) error](#Node.WaitForHeight)
* [type StaticCertifier](#StaticCertifier)
  * [func (c StaticCertifier) Certify(block lc.TmSignedHeader) error](#StaticCertifier.Certify)


#### <a name="pkg-files">Package files</a>
[certifier.go](/src/github.com/tendermint/light-client/rpc/certifier.go) [client.go](/src/github.com/tendermint/light-client/rpc/client.go) [docs.go](/src/github.com/tendermint/light-client/rpc/docs.go) [node.go](/src/github.com/tendermint/light-client/rpc/node.go) 



## <a name="pkg-variables">Variables</a>
``` go
var MerkleReader lc.ProofReader = merkleReader{}
```
MerkleReader is currently the only implementation of ProofReader,
using the IAVLProof from go-merkle



## <a name="ImportSignedHeader">func</a> [ImportSignedHeader](/src/target/node.go?s=5014:5093#L182)
``` go
func ImportSignedHeader(chainID string, data []byte) (lc.TmSignedHeader, error)
```
ImportSignedHeader takes serialized data from ExportSignedHeader
and verifies and processes it the same as SignedHeader.

Use this where you have no Node (rpcclient), but you still need the
chainID.

The result can be used just as the result from SignedHeader, and
passed to a Certifier




## <a name="HTTPClient">type</a> [HTTPClient](/src/target/client.go?s=207:332#L2)
``` go
type HTTPClient struct {
    // contains filtered or unexported fields
}
```






### <a name="NewClient">func</a> [NewClient](/src/target/client.go?s=334:387#L9)
``` go
func NewClient(remote, wsEndpoint string) *HTTPClient
```




### <a name="HTTPClient.ABCIInfo">func</a> (\*HTTPClient) [ABCIInfo](/src/target/client.go?s=824:887#L27)
``` go
func (c *HTTPClient) ABCIInfo() (*ctypes.ResultABCIInfo, error)
```



### <a name="HTTPClient.ABCIQuery">func</a> (\*HTTPClient) [ABCIQuery](/src/target/client.go?s=1102:1203#L36)
``` go
func (c *HTTPClient) ABCIQuery(path string, data []byte, prove bool) (*ctypes.ResultABCIQuery, error)
```



### <a name="HTTPClient.Block">func</a> (\*HTTPClient) [Block](/src/target/client.go?s=3411:3478#L103)
``` go
func (c *HTTPClient) Block(height int) (*ctypes.ResultBlock, error)
```



### <a name="HTTPClient.BlockchainInfo">func</a> (\*HTTPClient) [BlockchainInfo](/src/target/client.go?s=2804:2903#L85)
``` go
func (c *HTTPClient) BlockchainInfo(minHeight, maxHeight int) (*ctypes.ResultBlockchainInfo, error)
```



### <a name="HTTPClient.BroadcastTxAsync">func</a> (\*HTTPClient) [BroadcastTxAsync](/src/target/client.go?s=1585:1676#L49)
``` go
func (c *HTTPClient) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error)
```



### <a name="HTTPClient.BroadcastTxCommit">func</a> (\*HTTPClient) [BroadcastTxCommit](/src/target/client.go?s=1438:1530#L45)
``` go
func (c *HTTPClient) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error)
```



### <a name="HTTPClient.BroadcastTxSync">func</a> (\*HTTPClient) [BroadcastTxSync](/src/target/client.go?s=1730:1820#L53)
``` go
func (c *HTTPClient) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error)
```



### <a name="HTTPClient.DialSeeds">func</a> (\*HTTPClient) [DialSeeds](/src/target/client.go?s=2449:2528#L75)
``` go
func (c *HTTPClient) DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error)
```



### <a name="HTTPClient.Genesis">func</a> (\*HTTPClient) [Genesis](/src/target/client.go?s=3151:3212#L94)
``` go
func (c *HTTPClient) Genesis() (*ctypes.ResultGenesis, error)
```



### <a name="HTTPClient.GetEventChannels">func</a> (\*HTTPClient) [GetEventChannels](/src/target/client.go?s=4554:4628#L146)
``` go
func (c *HTTPClient) GetEventChannels() (chan json.RawMessage, chan error)
```
GetEventChannels returns the results and error channel from the websocket




### <a name="HTTPClient.NetInfo">func</a> (\*HTTPClient) [NetInfo](/src/target/client.go?s=2188:2249#L66)
``` go
func (c *HTTPClient) NetInfo() (*ctypes.ResultNetInfo, error)
```



### <a name="HTTPClient.StartWebsocket">func</a> (\*HTTPClient) [StartWebsocket](/src/target/client.go?s=4102:4145#L125)
``` go
func (c *HTTPClient) StartWebsocket() error
```
StartWebsocket starts up a websocket and a listener goroutine
if already started, do nothing




### <a name="HTTPClient.Status">func</a> (\*HTTPClient) [Status](/src/target/client.go?s=509:568#L17)
``` go
func (c *HTTPClient) Status() (*ctypes.ResultStatus, error)
```



### <a name="HTTPClient.StopWebsocket">func</a> (\*HTTPClient) [StopWebsocket](/src/target/client.go?s=4387:4423#L138)
``` go
func (c *HTTPClient) StopWebsocket()
```
StopWebsocket stops the websocket connection




### <a name="HTTPClient.Subscribe">func</a> (\*HTTPClient) [Subscribe](/src/target/client.go?s=4711:4761#L153)
``` go
func (c *HTTPClient) Subscribe(event string) error
```



### <a name="HTTPClient.Unsubscribe">func</a> (\*HTTPClient) [Unsubscribe](/src/target/client.go?s=4823:4875#L157)
``` go
func (c *HTTPClient) Unsubscribe(event string) error
```



### <a name="HTTPClient.Validators">func</a> (\*HTTPClient) [Validators](/src/target/client.go?s=3689:3756#L112)
``` go
func (c *HTTPClient) Validators() (*ctypes.ResultValidators, error)
```



## <a name="Node">type</a> [Node](/src/target/node.go?s=296:437#L5)
``` go
type Node struct {
    lc.ProofReader
    lc.ValueReader
    // contains filtered or unexported fields
}
```






### <a name="NewNode">func</a> [NewNode](/src/target/node.go?s=723:788#L23)
``` go
func NewNode(rpcAddr, chainID string, valuer lc.ValueReader) Node
```




### <a name="Node.Broadcast">func</a> (Node) [Broadcast](/src/target/node.go?s=1317:1389#L49)
``` go
func (n Node) Broadcast(tx []byte) (res lc.TmBroadcastResult, err error)
```
Broadcast sends the transaction to a tendermint node and waits until
it is committed.

If it failed on CheckTx, we return the result of CheckTx, otherwise
we return the result of DeliverTx




### <a name="Node.ExportSignedHeader">func</a> (Node) [ExportSignedHeader](/src/target/node.go?s=3733:3796#L137)
``` go
func (n Node) ExportSignedHeader(height uint64) ([]byte, error)
```
ExportSignedHeader downloads and verifies the same info as
SignedHeader, but returns a serialized version of the proof to
be stored for later use.

The result should be considered opaque bytes, but can be passed into
ImportSignedHeader to get data ready for a Certifier




### <a name="Node.ImportSignedHeader">func</a> (Node) [ImportSignedHeader](/src/target/node.go?s=4583:4655#L170)
``` go
func (n Node) ImportSignedHeader(data []byte) (lc.TmSignedHeader, error)
```
ImportSignedHeader takes serialized data from ExportSignedHeader
and verifies and processes it the same as SignedHeader.

This is just a convenience wrapper around the same-named function
that passes in the chainID from the node.

The result can be used just as the result from SignedHeader, and
passed to a Certifier




### <a name="Node.Prove">func</a> (Node) [Prove](/src/target/node.go?s=2006:2063#L76)
``` go
func (n Node) Prove(key []byte) (lc.TmQueryResult, error)
```
Prove returns a merkle proof for the given key




### <a name="Node.Query">func</a> (Node) [Query](/src/target/node.go?s=1800:1871#L70)
``` go
func (n Node) Query(path string, data []byte) (lc.TmQueryResult, error)
```
Query gets data from the Blockchain state, possibly with a
complex path.  It doesn't worry about proofs




### <a name="Node.SignedHeader">func</a> (Node) [SignedHeader](/src/target/node.go?s=2940:3008#L108)
``` go
func (n Node) SignedHeader(height uint64) (lc.TmSignedHeader, error)
```
SignedHeader gives us Header data along with the backing signatures

It is also responsible for making the blockhash and header match,
and that all votes are valid pre-commit votes for this block
It does not check if the keys signing the votes are actual validators




### <a name="Node.Validators">func</a> (Node) [Validators](/src/target/node.go?s=6456:6512#L228)
``` go
func (n Node) Validators() (lc.TmValidatorResult, error)
```
UNSAFE - use wisely

Validators returns the current set of validators from the
node we call.  There is no guarantee it is correct.

This is intended for use in test cases, or to query many nodes
to find consensus before trusting it.




### <a name="Node.WaitForHeight">func</a> (Node) [WaitForHeight](/src/target/node.go?s=5600:5648#L198)
``` go
func (n Node) WaitForHeight(height uint64) error
```
Wait for height will poll status at reasonable intervals until
we can safely call SignedHeader at the given block height.
This means that both the block header itself, as well as all
validator signatures are available

In this current implementation, we must wait until height+1,
as the signatures are in the following block.




## <a name="StaticCertifier">type</a> [StaticCertifier](/src/target/certifier.go?s=339:393#L5)
``` go
type StaticCertifier struct {
    Vals []lc.TmValidator
}
```
StaticCertifier assumes a static set of validators, set on
initilization and checks against them.

Good for testing or really simple chains.  You will want a
better implementation when the validator set can actually change.










### <a name="StaticCertifier.Certify">func</a> (StaticCertifier) [Certify](/src/target/certifier.go?s=466:529#L13)
``` go
func (c StaticCertifier) Certify(block lc.TmSignedHeader) error
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
