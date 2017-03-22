

# tx
`import "github.com/tendermint/light-client/tx"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
package tx contains generic Signable implementations that can be used
by your application or tests to handle authentication needs.

It currently supports transaction data as opaque bytes and either single
or multiple private key signatures using straightforward algorithms.
It currently does not support N-of-M key share signing of other more
complex algorithms (although it would be great to add them)

TODO: This package should be moved into go-keys/go-crypto




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type MultiSig](#MultiSig)
  * [func (s *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#MultiSig.Sign)
  * [func (s *MultiSig) SignBytes() []byte](#MultiSig.SignBytes)
  * [func (s *MultiSig) Signers() ([]crypto.PubKey, error)](#MultiSig.Signers)
* [type OneSig](#OneSig)
  * [func (s *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#OneSig.Sign)
  * [func (s *OneSig) SignBytes() []byte](#OneSig.SignBytes)
  * [func (s *OneSig) Signers() ([]crypto.PubKey, error)](#OneSig.Signers)
* [type Sig](#Sig)
  * [func New(data []byte) Sig](#New)
  * [func NewMulti(data []byte) Sig](#NewMulti)
  * [func (s Sig) TxBytes() ([]byte, error)](#Sig.TxBytes)
* [type SigInner](#SigInner)


#### <a name="pkg-files">Package files</a>
[docs.go](/src/github.com/tendermint/light-client/tx/docs.go) [multi.go](/src/github.com/tendermint/light-client/tx/multi.go) [one.go](/src/github.com/tendermint/light-client/tx/one.go) [reader.go](/src/github.com/tendermint/light-client/tx/reader.go) 



## <a name="pkg-variables">Variables</a>
``` go
var TxMapper data.Mapper
```



## <a name="MultiSig">type</a> [MultiSig](/src/target/multi.go?s=281:333#L2)
``` go
type MultiSig struct {
    // contains filtered or unexported fields
}
```
MultiSig lets us wrap arbitrary data with a go-crypto signature

TODO: rethink how we want to integrate this with KeyStore so it makes
more sense (particularly the verify method)










### <a name="MultiSig.Sign">func</a> (\*MultiSig) [Sign](/src/target/multi.go?s=819:892#L27)
``` go
func (s *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```
Sign will add a signature and pubkey.

Depending on the Signable, one may be able to call this multiple times for multisig
Returns error if called with invalid data or too many times




### <a name="MultiSig.SignBytes">func</a> (\*MultiSig) [SignBytes](/src/target/multi.go?s=567:604#L19)
``` go
func (s *MultiSig) SignBytes() []byte
```
SignBytes returns the original data passed into `NewSig`




### <a name="MultiSig.Signers">func</a> (\*MultiSig) [Signers](/src/target/multi.go?s=1263:1316#L41)
``` go
func (s *MultiSig) Signers() ([]crypto.PubKey, error)
```
Signers will return the public key(s) that signed if the signature
is valid, or an error if there is any issue with the signature,
including if there are no signatures




## <a name="OneSig">type</a> [OneSig](/src/target/one.go?s=279:322#L2)
``` go
type OneSig struct {
    // contains filtered or unexported fields
}
```
OneSig lets us wrap arbitrary data with a go-crypto signature

TODO: rethink how we want to integrate this with KeyStore so it makes
more sense (particularly the verify method)










### <a name="OneSig.Sign">func</a> (\*OneSig) [Sign](/src/target/one.go?s=726:797#L22)
``` go
func (s *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```
Sign will add a signature and pubkey.

Depending on the Signable, one may be able to call this multiple times for multisig
Returns error if called with invalid data or too many times




### <a name="OneSig.SignBytes">func</a> (\*OneSig) [SignBytes](/src/target/one.go?s=476:511#L14)
``` go
func (s *OneSig) SignBytes() []byte
```
SignBytes returns the original data passed into `NewSig`




### <a name="OneSig.Signers">func</a> (\*OneSig) [Signers](/src/target/one.go?s=1227:1278#L39)
``` go
func (s *OneSig) Signers() ([]crypto.PubKey, error)
```
Signers will return the public key(s) that signed if the signature
is valid, or an error if there is any issue with the signature,
including if there are no signatures




## <a name="Sig">type</a> [Sig](/src/target/reader.go?s=671:700#L22)
``` go
type Sig struct {
    SigInner
}
```
Sig is what is exported, and handles serialization







### <a name="New">func</a> [New](/src/target/one.go?s=352:377#L9)
``` go
func New(data []byte) Sig
```

### <a name="NewMulti">func</a> [NewMulti](/src/target/multi.go?s=436:466#L14)
``` go
func NewMulti(data []byte) Sig
```




### <a name="Sig.TxBytes">func</a> (Sig) [TxBytes](/src/target/reader.go?s=702:740#L26)
``` go
func (s Sig) TxBytes() ([]byte, error)
```



## <a name="SigInner">type</a> [SigInner](/src/target/reader.go?s=476:615#L15)
``` go
type SigInner interface {
    SignBytes() []byte
    Sign(pubkey crypto.PubKey, sig crypto.Signature) error
    Signers() ([]crypto.PubKey, error)
}
```













- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
