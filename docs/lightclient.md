

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
* [type Broadcaster](#Broadcaster)
* [type Certifier](#Certifier)
* [type Checker](#Checker)
* [type KeyInfo](#KeyInfo)
* [type KeyInfos](#KeyInfos)
  * [func (k KeyInfos) Len() int](#KeyInfos.Len)
  * [func (k KeyInfos) Less(i, j int) bool](#KeyInfos.Less)
  * [func (k KeyInfos) Sort()](#KeyInfos.Sort)
  * [func (k KeyInfos) Swap(i, j int)](#KeyInfos.Swap)
* [type KeyManager](#KeyManager)
* [type Proof](#Proof)
* [type ProofReader](#ProofReader)
* [type Searcher](#Searcher)
* [type Signable](#Signable)
* [type SignableReader](#SignableReader)
* [type Signer](#Signer)
* [type Storage](#Storage)
* [type TmBroadcastResult](#TmBroadcastResult)
  * [func (r TmBroadcastResult) IsOk() bool](#TmBroadcastResult.IsOk)
* [type TmCode](#TmCode)
  * [func (c TmCode) IsOK() bool](#TmCode.IsOK)
* [type TmHeader](#TmHeader)
* [type TmQueryResult](#TmQueryResult)
* [type TmSignedHeader](#TmSignedHeader)
  * [func (sh TmSignedHeader) Height() uint64](#TmSignedHeader.Height)
* [type TmStatusResult](#TmStatusResult)
* [type TmValidator](#TmValidator)
* [type TmValidatorResult](#TmValidatorResult)
* [type TmVote](#TmVote)
* [type TmVotes](#TmVotes)
  * [func (v TmVotes) ForBlock(hash []byte) bool](#TmVotes.ForBlock)
* [type Value](#Value)
* [type ValueReader](#ValueReader)


#### <a name="pkg-files">Package files</a>
[docs.go](/src/github.com/tendermint/light-client/docs.go) [proofs.go](/src/github.com/tendermint/light-client/proofs.go) [rpc.go](/src/github.com/tendermint/light-client/rpc.go) [storage.go](/src/github.com/tendermint/light-client/storage.go) [transactions.go](/src/github.com/tendermint/light-client/transactions.go) [types.go](/src/github.com/tendermint/light-client/types.go) 






## <a name="Broadcaster">type</a> [Broadcaster](/src/target/rpc.go?s=101:360#L1)
``` go
type Broadcaster interface {
    // Broadcast sends into to the chain
    // We only implement BroadcastCommit for now, add others???
    // The return result cannot be fully trusted without downloading signed headers
    Broadcast(tx []byte) (TmBroadcastResult, error)
}
```
Broadcaster provides a way to send a signed transaction to a tendermint node










## <a name="Certifier">type</a> [Certifier](/src/target/proofs.go?s=893:958#L14)
``` go
type Certifier interface {
    Certify(block TmSignedHeader) error
}
```
Certifier checks the votes to make sure the block really is signed properly.
Certifier must know the current set of validitors by some other means.
TODO: some implementation to track the validator set (various algorithms)










## <a name="Checker">type</a> [Checker](/src/target/rpc.go?s=482:1020#L3)
``` go
type Checker interface {
    // Prove returns a merkle proof for the given key
    Prove(key []byte) (TmQueryResult, error)

    // SignedHeader gives us Header data along with the backing signatures,
    // so we can validate it externally (matching with the list of
    // known validators)
    SignedHeader(height uint64) (TmSignedHeader, error)

    // WaitForHeight is a useful helper to poll the server until the
    // data is ready for SignedHeader.  Returns nil when the data
    // is present, and error if it aborts.
    WaitForHeight(height uint64) error
}
```
Checker provides access to calls to get data from the tendermint core
and all cryptographic proof of its validity










## <a name="KeyInfo">type</a> [KeyInfo](/src/target/transactions.go?s=133:193#L1)
``` go
type KeyInfo struct {
    Name   string
    PubKey crypto.PubKey
}
```
KeyInfo is the public information about a key










## <a name="KeyInfos">type</a> [KeyInfos](/src/target/transactions.go?s=263:286#L6)
``` go
type KeyInfos []KeyInfo
```
KeyInfos is a wrapper to allows alphabetical sorting of the keys










### <a name="KeyInfos.Len">func</a> (KeyInfos) [Len](/src/target/transactions.go?s=288:315#L8)
``` go
func (k KeyInfos) Len() int
```



### <a name="KeyInfos.Less">func</a> (KeyInfos) [Less](/src/target/transactions.go?s=344:381#L9)
``` go
func (k KeyInfos) Less(i, j int) bool
```



### <a name="KeyInfos.Sort">func</a> (KeyInfos) [Sort](/src/target/transactions.go?s=481:505#L11)
``` go
func (k KeyInfos) Sort()
```



### <a name="KeyInfos.Swap">func</a> (KeyInfos) [Swap](/src/target/transactions.go?s=415:447#L10)
``` go
func (k KeyInfos) Swap(i, j int)
```



## <a name="KeyManager">type</a> [KeyManager](/src/target/transactions.go?s=1743:1956#L50)
``` go
type KeyManager interface {
    Create(name, passphrase string) error
    List() (KeyInfos, error)
    Get(name string) (KeyInfo, error)
    Update(name, oldpass, newpass string) error
    Delete(name, passphrase string) error
}
```
KeyManager allows simple CRUD on a keystore, as an aid to signing










## <a name="Proof">type</a> [Proof](/src/target/proofs.go?s=228:533#L1)
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










## <a name="ProofReader">type</a> [ProofReader](/src/target/proofs.go?s=591:660#L7)
``` go
type ProofReader interface {
    ReadProof(data []byte) (Proof, error)
}
```
ProofReader is an abstraction to let us parse proofs










## <a name="Searcher">type</a> [Searcher](/src/target/rpc.go?s=1022:1217#L18)
``` go
type Searcher interface {
    // Query gets data from the Blockchain state, possibly with a
    // complex path.  It doesn't worry about proofs
    Query(path string, data []byte) (TmQueryResult, error)
}
```









## <a name="Signable">type</a> [Signable](/src/target/transactions.go?s=683:1414#L19)
``` go
type Signable interface {
    // SignBytes is the immutable data, which needs to be signed
    SignBytes() []byte

    // Sign will add a signature and pubkey.
    //
    // Depending on the Signable, one may be able to call this multiple times for multisig
    // Returns error if called with invalid data or too many times
    Sign(pubkey crypto.PubKey, sig crypto.Signature) error

    // Signers will return the public key(s) that signed if the signature
    // is valid, or an error if there is any issue with the signature,
    // including if there are no signatures
    Signers() ([]crypto.PubKey, error)

    // TxBytes returns the transaction data as well as all signatures
    // It should return an error if Sign was never called
    TxBytes() ([]byte, error)
}
```
Signable represents any transaction we wish to send to tendermint core
These methods allow us to sign arbitrary Tx with the KeyStore










## <a name="SignableReader">type</a> [SignableReader](/src/target/transactions.go?s=1478:1556#L40)
``` go
type SignableReader interface {
    ReadSignable(data []byte) (Signable, error)
}
```
SignableReader is an abstraction to let us parse Signables










## <a name="Signer">type</a> [Signer](/src/target/transactions.go?s=1597:1672#L45)
``` go
type Signer interface {
    Sign(name, passphrase string, tx Signable) error
}
```
Signer allows one to use a keystore










## <a name="Storage">type</a> [Storage](/src/target/storage.go?s=149:322#L1)
``` go
type Storage interface {
    Put(name string, key []byte, info KeyInfo) error
    Get(name string) ([]byte, KeyInfo, error)
    List() ([]KeyInfo, error)
    Delete(name string) error
}
```
Storage has many implementation, based on security and sharing requirements
like disk-backed, mem-backed, vault, db, etc.










## <a name="TmBroadcastResult">type</a> [TmBroadcastResult](/src/target/types.go?s=203:338#L7)
``` go
type TmBroadcastResult struct {
    Code TmCode `json:"code"` // TODO: rethink this
    Data []byte `json:"data"`
    Log  string `json:"log"`
}
```









### <a name="TmBroadcastResult.IsOk">func</a> (TmBroadcastResult) [IsOk](/src/target/types.go?s=340:378#L13)
``` go
func (r TmBroadcastResult) IsOk() bool
```



## <a name="TmCode">type</a> [TmCode](/src/target/types.go?s=129:146#L1)
``` go
type TmCode int32
```
TmCode is a code from Tendermint










### <a name="TmCode.IsOK">func</a> (TmCode) [IsOK](/src/target/types.go?s=148:175#L3)
``` go
func (c TmCode) IsOK() bool
```



## <a name="TmHeader">type</a> [TmHeader](/src/target/types.go?s=1540:2158#L48)
``` go
type TmHeader struct {
    ChainID string    `json:"chain_id"`
    Height  uint64    `json:"height"`
    Time    time.Time `json:"time"` // or int64 nanoseconds????
    // NumTxs         int       `json:"num_txs"` // XXX: Can we get rid of this?
    LastBlockID    []byte `json:"last_block_id"`
    LastCommitHash []byte `json:"last_commit_hash"` // commit from validators from the last block
    DataHash       []byte `json:"data_hash"`        // transactions
    ValidatorsHash []byte `json:"validators_hash"`  // validators for the current block
    AppHash        []byte `json:"app_hash"`         // state after txs from the previous block
}
```
TmHeader is the info in block headers (from tendermint/types/block.go)










## <a name="TmQueryResult">type</a> [TmQueryResult](/src/target/types.go?s=598:803#L20)
``` go
type TmQueryResult struct {
    Code   TmCode `json:"code"`
    Key    []byte `json:"key"`
    Value  Value  `json:"value"`
    Proof  Proof  `json:"proof"`
    Height uint64 `json:"height"`
    Log    string `json:"log"`
}
```
TmQueryResult stores all info from a query or a proof
The Checker/Searcher is responsible for parsing the Value and Proof into
usable formats (not just opaque bytes) via their Reader










## <a name="TmSignedHeader">type</a> [TmSignedHeader](/src/target/types.go?s=1146:1321#L35)
``` go
type TmSignedHeader struct {
    Hash       []byte
    Header     TmHeader // contains height
    Votes      TmVotes
    LastHeight uint64 // the most recent block commited to the chain
}
```
TmSignedHeader returns the header info and the corresponding validator
signatures for it (even if those are on differrent rpc endpoints).

Checker is responsible for validating that the Hash matches the Header
which lets us return the Header in a non-binary format.
Votes must also be pre-verified, signatures checked, etc.










### <a name="TmSignedHeader.Height">func</a> (TmSignedHeader) [Height](/src/target/types.go?s=1395:1435#L43)
``` go
func (sh TmSignedHeader) Height() uint64
```
Height has been verified to be the same for the header and all votes




## <a name="TmStatusResult">type</a> [TmStatusResult](/src/target/types.go?s=3257:3507#L97)
``` go
type TmStatusResult struct {
    LatestBlockHash   []byte `json:"latest_block_hash"`
    LatestAppHash     []byte `json:"latest_app_hash"`
    LatestBlockHeight int    `json:"latest_block_height"`
    LatestBlockTime   int64  `json:"latest_block_time"` // nano
}
```









## <a name="TmValidator">type</a> [TmValidator](/src/target/types.go?s=3559:3748#L105)
``` go
type TmValidator struct {
    Address []byte `json:"address"`
    // PubKey  []byte `json:"pub_key"`
    PubKey      crypto.PubKey `json:"pub_key"`
    VotingPower int64         `json:"voting_power"`
}
```
TmValidator more or less from tendermint/types










## <a name="TmValidatorResult">type</a> [TmValidatorResult](/src/target/types.go?s=3750:3830#L112)
``` go
type TmValidatorResult struct {
    BlockHeight uint64
    Validators  []TmValidator
}
```









## <a name="TmVote">type</a> [TmVote](/src/target/types.go?s=2363:2898#L63)
``` go
type TmVote struct {
    // SignBytes are the cannonical bytes the signature refers to
    // This is verified in the Checker
    SignBytes []byte `json:"sign_bytes"`

    // Signature and ValidatorAddress represent who signed this
    // This information is not verified and must be validated
    // by the caller
    Signature        crypto.Signature `json:"signature"`
    ValidatorAddress []byte           `json:"validator_address"`

    // Height and BlockHash is embedded in TxBytes
    Height    uint64 `json:"height"`
    BlockHash []byte `json:"block_hash"`
}
```
TmVote must be verified by the Node implementation, this asserts a validly
signed precommit vote for the given Height and BlockHash.
The client can decide if these validators are to be trusted.










## <a name="TmVotes">type</a> [TmVotes](/src/target/types.go?s=2980:3001#L80)
``` go
type TmVotes []TmVote
```
TmVotes is a slice of TmVote structs, but let's add some control access here










### <a name="TmVotes.ForBlock">func</a> (TmVotes) [ForBlock](/src/target/types.go?s=3070:3113#L83)
``` go
func (v TmVotes) ForBlock(hash []byte) bool
```
ForBlock returns true only if all votes are for the given block




## <a name="Value">type</a> [Value](/src/target/rpc.go?s=1448:1488#L29)
``` go
type Value interface {
    Bytes() []byte
}
```
Value represents a database value and is generally a structure
that can be json serialized.  Bytes() is needed to get the original
data bytes for validation of proofs

TODO: add Fields() method to get field info???










## <a name="ValueReader">type</a> [ValueReader](/src/target/rpc.go?s=1567:1900#L34)
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
