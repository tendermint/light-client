

# lightclient
`import "github.com/tendermint/light-client"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
package lightclient is a complete solution for integrating a light client with
tendermint.  It provides all common functionality that a client needs to create
and sign transactions, get and verify state, and synchronize with a tendermint node.
It is intended to expose this data both through golang interfaces, a local RPC server,
and language bindings.  You can find more info on the aims of this package in the
Readme: <a href="https://github.com/tendermint/light-client/blob/master/README.md">https://github.com/tendermint/light-client/blob/master/README.md</a>

The package layout attempts to expose common domain types in the
top-level with no other dependencies.  Main packages should select which
dependencies they wish to have and wire them together with common glue code
that only depends on the interface.
More info here: <a href="https://medium.com/%40benbjohnson/standard-package-layout-7cdbc8391fc1">https://medium.com/%40benbjohnson/standard-package-layout-7cdbc8391fc1</a>

The majority of the definitions here are interfaces, to be implemented in subpackages,
or data structures encapsulating tendermint return values.  All tendermint data
structures are prefixed by Tm to make the documentation clearer. The other data
structure defined on the top level is KeyInfo/KeyInfos to represent infomation
on the public keys stored in the KeyManager, along with logic to sort themselves.




## <a name="pkg-index">Index</a>
* [type ByteValue](#ByteValue)
  * [func (b ByteValue) Bytes() []byte](#ByteValue.Bytes)
* [type Certifier](#Certifier)
* [type Checkpoint](#Checkpoint)
  * [func NewCheckpoint(commit *rtypes.ResultCommit) Checkpoint](#NewCheckpoint)
  * [func (c Checkpoint) CheckAppState(k, v []byte, proof Proof) error](#Checkpoint.CheckAppState)
  * [func (c Checkpoint) CheckTxs(txs types.Txs) error](#Checkpoint.CheckTxs)
  * [func (c Checkpoint) CheckValidators(vals []*types.Validator) error](#Checkpoint.CheckValidators)
  * [func (c Checkpoint) Height() int](#Checkpoint.Height)
  * [func (c Checkpoint) ValidateBasic(chainID string) error](#Checkpoint.ValidateBasic)
* [type Poster](#Poster)
  * [func NewPoster(server client.ABCIClient, signer keys.Signer) Poster](#NewPoster)
  * [func (p Poster) Post(sign keys.Signable, keyname, passphrase string) (*ctypes.ResultBroadcastTxCommit, error)](#Poster.Post)
* [type Proof](#Proof)
* [type ProofReader](#ProofReader)
* [type SignableReader](#SignableReader)
* [type Value](#Value)
* [type ValueReader](#ValueReader)


#### <a name="pkg-files">Package files</a>
[checkpoint.go](/src/github.com/tendermint/light-client/checkpoint.go) [docs.go](/src/github.com/tendermint/light-client/docs.go) [poster.go](/src/github.com/tendermint/light-client/poster.go) [readers.go](/src/github.com/tendermint/light-client/readers.go) 






## <a name="ByteValue">type</a> [ByteValue](/src/target/readers.go?s=1303:1324#L35)
``` go
type ByteValue []byte
```









### <a name="ByteValue.Bytes">func</a> (ByteValue) [Bytes](/src/target/readers.go?s=1326:1359#L37)
``` go
func (b ByteValue) Bytes() []byte
```



## <a name="Certifier">type</a> [Certifier](/src/target/checkpoint.go?s=322:383#L3)
``` go
type Certifier interface {
    Certify(check Checkpoint) error
}
```
Certifier checks the votes to make sure the block really is signed properly.
Certifier must know the current set of validitors by some other means.










## <a name="Checkpoint">type</a> [Checkpoint](/src/target/checkpoint.go?s=748:850#L13)
``` go
type Checkpoint struct {
    Header *types.Header `json:"header"`
    Commit *types.Commit `json:"commit"`
}
```
Checkpoint is basically the rpc /commit response, but extended

This is the basepoint for proving anything on the blockchain. It contains
a signed header.  If the signatures are valid and > 2/3 of the known set,
we can store this checkpoint and use it to prove any number of aspects of
the system: such as txs, abci state, validator sets, etc...







### <a name="NewCheckpoint">func</a> [NewCheckpoint](/src/target/checkpoint.go?s=852:910#L18)
``` go
func NewCheckpoint(commit *rtypes.ResultCommit) Checkpoint
```




### <a name="Checkpoint.CheckAppState">func</a> (Checkpoint) [CheckAppState](/src/target/checkpoint.go?s=3349:3414#L103)
``` go
func (c Checkpoint) CheckAppState(k, v []byte, proof Proof) error
```
CheckAppState validates whether the key-value pair and merkle proof
can be verified with this Checkpoint.




### <a name="Checkpoint.CheckTxs">func</a> (Checkpoint) [CheckTxs](/src/target/checkpoint.go?s=2896:2945#L87)
``` go
func (c Checkpoint) CheckTxs(txs types.Txs) error
```
CheckTxs checks if the entire set of transactions for the block matches
the Checkpoint header.




### <a name="Checkpoint.CheckValidators">func</a> (Checkpoint) [CheckValidators](/src/target/checkpoint.go?s=2497:2563#L74)
``` go
func (c Checkpoint) CheckValidators(vals []*types.Validator) error
```
CheckValidators should only be used after you fully trust this checkpoint

It checks if these really are the validators authorized to sign the
checkpoint.




### <a name="Checkpoint.Height">func</a> (Checkpoint) [Height](/src/target/checkpoint.go?s=989:1021#L25)
``` go
func (c Checkpoint) Height() int
```



### <a name="Checkpoint.ValidateBasic">func</a> (Checkpoint) [ValidateBasic](/src/target/checkpoint.go?s=1321:1376#L34)
``` go
func (c Checkpoint) ValidateBasic(chainID string) error
```
ValidateBasic does basic consistency checks and makes sure the headers
and commits are all consistent and refer to our chain.

Make sure to use a Verifier to validate the signatures actually provide
a significantly strong proof for this header's validity.




## <a name="Poster">type</a> [Poster](/src/target/poster.go?s=387:455#L3)
``` go
type Poster struct {
    // contains filtered or unexported fields
}
```
Poster combines KeyStore and Node to process a Signable and deliver it to tendermint
returning the results from the tendermint node, once the transaction is processed.

Only handles single signatures







### <a name="NewPoster">func</a> [NewPoster](/src/target/poster.go?s=457:524#L8)
``` go
func NewPoster(server client.ABCIClient, signer keys.Signer) Poster
```




### <a name="Poster.Post">func</a> (Poster) [Post](/src/target/poster.go?s=662:771#L14)
``` go
func (p Poster) Post(sign keys.Signable, keyname, passphrase string) (*ctypes.ResultBroadcastTxCommit, error)
```
Post will sign the transaction with the given credentials and push it to
the tendermint server




## <a name="Proof">type</a> [Proof](/src/target/readers.go?s=320:625#L1)
``` go
type Proof interface {
    // Root returns the RootHash of the merkle tree used in the proof,
    // This is important for correlating it with a block header.
    Root() []byte

    // Verify returns true iff this proof validates this key and value belong
    // to the given root
    Verify(key, value, root []byte) bool
}
```
Proof is a generalization of merkle.IAVLProof and represents any
merkle proof that can validate a key-value pair back to a root hash.
TODO: someway to save/export a given proof for another client??










## <a name="ProofReader">type</a> [ProofReader](/src/target/readers.go?s=683:752#L12)
``` go
type ProofReader interface {
    ReadProof(data []byte) (Proof, error)
}
```
ProofReader is an abstraction to let us parse proofs


``` go
var MerkleReader ProofReader = merkleReader{}
```
MerkleReader is currently the only implementation of ProofReader,
using the IAVLProof from go-merkle










## <a name="SignableReader">type</a> [SignableReader](/src/target/readers.go?s=1786:1869#L49)
``` go
type SignableReader interface {
    ReadSignable(data []byte) (keys.Signable, error)
}
```









## <a name="Value">type</a> [Value](/src/target/readers.go?s=1261:1301#L31)
``` go
type Value interface {
    Bytes() []byte
}
```
Value represents a database value and is generally a structure
that can be json serialized.  Bytes() is needed to get the original
data bytes for validation of proofs

TODO: add Fields() method to get field info???










## <a name="ValueReader">type</a> [ValueReader](/src/target/readers.go?s=1451:1784#L40)
``` go
type ValueReader interface {
    // ReadValue accepts a key, value pair to decode.  The value bytes must be
    // retained in the returned Value implementation.
    //
    // key *may* be present and can be used as a hint of how to parse the data
    // when your application handles multiple formats
    ReadValue(key, value []byte) (Value, error)
}
```
ValueReader is an abstraction to let us parse application-specific values














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
