

# tx
`import "github.com/tendermint/light-client/tx"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
package tx contains generic Signable implementations that can be used
by your application to handle authentication needs.

It currently supports transaction data as opaque bytes and either single
or multiple private key signatures using straightforward algorithms.
It currently does not support N-of-M key share signing of other more
complex algorithms (although it would be great to add them)

Reader() will return a SignableReader suitable for deserialization
of these two types.

Maybe this package should be moved out of lightclient, more to a repo
designed for application support, as that is the main usecase?




## <a name="pkg-index">Index</a>
* [func Reader() lc.SignableReader](#Reader)
* [type MultiSig](#MultiSig)
  * [func LoadMulti(serialized []byte) (*MultiSig, error)](#LoadMulti)
  * [func NewMulti(data []byte) *MultiSig](#NewMulti)
  * [func (s *MultiSig) Bytes() []byte](#MultiSig.Bytes)
  * [func (s *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#MultiSig.Sign)
  * [func (s *MultiSig) SignedBy() ([]crypto.PubKey, error)](#MultiSig.SignedBy)
  * [func (s *MultiSig) SignedBytes() ([]byte, error)](#MultiSig.SignedBytes)
* [type OneSig](#OneSig)
  * [func Load(serialized []byte) (*OneSig, error)](#Load)
  * [func New(data []byte) *OneSig](#New)
  * [func (s *OneSig) Bytes() []byte](#OneSig.Bytes)
  * [func (s *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#OneSig.Sign)
  * [func (s *OneSig) SignedBy() ([]crypto.PubKey, error)](#OneSig.SignedBy)
  * [func (s *OneSig) SignedBytes() ([]byte, error)](#OneSig.SignedBytes)


#### <a name="pkg-files">Package files</a>
[docs.go](/src/github.com/tendermint/light-client/tx/docs.go) [multi.go](/src/github.com/tendermint/light-client/tx/multi.go) [one.go](/src/github.com/tendermint/light-client/tx/one.go) [reader.go](/src/github.com/tendermint/light-client/tx/reader.go) 





## <a name="Reader">func</a> [Reader](/src/target/reader.go?s=529:560#L18)
``` go
func Reader() lc.SignableReader
```
Reader constructs a SignableReader that can parse OneSig and MultiSig

TODO: add some args to configure go-wire, rather than relying on init???7




## <a name="MultiSig">type</a> [MultiSig](/src/target/multi.go?s=369:421#L4)
``` go
type MultiSig struct {
    // contains filtered or unexported fields
}
```
MultiSig lets us wrap arbitrary data with a go-crypto signature

TODO: rethink how we want to integrate this with KeyStore so it makes
more sense (particularly the verify method)







### <a name="LoadMulti">func</a> [LoadMulti](/src/target/multi.go?s=566:618#L18)
``` go
func LoadMulti(serialized []byte) (*MultiSig, error)
```

### <a name="NewMulti">func</a> [NewMulti](/src/target/multi.go?s=494:530#L14)
``` go
func NewMulti(data []byte) *MultiSig
```




### <a name="MultiSig.Bytes">func</a> (\*MultiSig) [Bytes](/src/target/multi.go?s=912:945#L30)
``` go
func (s *MultiSig) Bytes() []byte
```
Bytes returns the original data passed into `NewSig`




### <a name="MultiSig.Sign">func</a> (\*MultiSig) [Sign](/src/target/multi.go?s=1160:1233#L38)
``` go
func (s *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```
Sign will add a signature and pubkey.

Depending on the Signable, one may be able to call this multiple times for multisig
Returns error if called with invalid data or too many times




### <a name="MultiSig.SignedBy">func</a> (\*MultiSig) [SignedBy](/src/target/multi.go?s=1605:1659#L52)
``` go
func (s *MultiSig) SignedBy() ([]crypto.PubKey, error)
```
SignedBy will return the public key(s) that signed if the signature
is valid, or an error if there is any issue with the signature,
including if there are no signatures




### <a name="MultiSig.SignedBytes">func</a> (\*MultiSig) [SignedBytes](/src/target/multi.go?s=2123:2171#L71)
``` go
func (s *MultiSig) SignedBytes() ([]byte, error)
```
SignedBytes serializes the Sig to send it to a tendermint app.
It returns an error if the Sig was never Signed.




## <a name="OneSig">type</a> [OneSig](/src/target/one.go?s=367:451#L4)
``` go
type OneSig struct {
    // contains filtered or unexported fields
}
```
OneSig lets us wrap arbitrary data with a go-crypto signature

TODO: rethink how we want to integrate this with KeyStore so it makes
more sense (particularly the verify method)







### <a name="Load">func</a> [Load](/src/target/one.go?s=516:561#L14)
``` go
func Load(serialized []byte) (*OneSig, error)
```

### <a name="New">func</a> [New](/src/target/one.go?s=453:482#L10)
``` go
func New(data []byte) *OneSig
```




### <a name="OneSig.Bytes">func</a> (\*OneSig) [Bytes](/src/target/one.go?s=851:882#L26)
``` go
func (s *OneSig) Bytes() []byte
```
Bytes returns the original data passed into `NewSig`




### <a name="OneSig.Sign">func</a> (\*OneSig) [Sign](/src/target/one.go?s=1097:1168#L34)
``` go
func (s *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```
Sign will add a signature and pubkey.

Depending on the Signable, one may be able to call this multiple times for multisig
Returns error if called with invalid data or too many times




### <a name="OneSig.SignedBy">func</a> (\*OneSig) [SignedBy](/src/target/one.go?s=1600:1652#L52)
``` go
func (s *OneSig) SignedBy() ([]crypto.PubKey, error)
```
SignedBy will return the public key(s) that signed if the signature
is valid, or an error if there is any issue with the signature,
including if there are no signatures




### <a name="OneSig.SignedBytes">func</a> (\*OneSig) [SignedBytes](/src/target/one.go?s=1997:2043#L66)
``` go
func (s *OneSig) SignedBytes() ([]byte, error)
```
SignedBytes serializes the Sig to send it to a tendermint app.
It returns an error if the Sig was never Signed.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
