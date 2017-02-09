

# mock
`import "github.com/tendermint/light-client/mock"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
package mock contains various mock implementations of the lightclient
interfaces for use in testing.

This code must not depend on any other subpackages of lightclient to avoid
potential circular imports, and is designed to be imported by test code
anywhere.

If you are importing this from production code, please think twice.




## <a name="pkg-index">Index</a>
* [func Reader() lc.SignableReader](#Reader)
* [func ValueReader() lc.ValueReader](#ValueReader)
* [type ByteValue](#ByteValue)
  * [func (b ByteValue) Bytes() []byte](#ByteValue.Bytes)
* [type ByteValueReader](#ByteValueReader)
  * [func (b ByteValueReader) ReadValue(key, value []byte) (lc.Value, error)](#ByteValueReader.ReadValue)
* [type MultiSig](#MultiSig)
  * [func NewMultiSig(data []byte) *MultiSig](#NewMultiSig)
  * [func (m *MultiSig) Bytes() []byte](#MultiSig.Bytes)
  * [func (m *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#MultiSig.Sign)
  * [func (m *MultiSig) SignedBy() ([]crypto.PubKey, error)](#MultiSig.SignedBy)
  * [func (m *MultiSig) SignedBytes() ([]byte, error)](#MultiSig.SignedBytes)
* [type OneSig](#OneSig)
  * [func NewSig(data []byte) *OneSig](#NewSig)
  * [func (o *OneSig) Bytes() []byte](#OneSig.Bytes)
  * [func (o *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#OneSig.Sign)
  * [func (o *OneSig) SignedBy() ([]crypto.PubKey, error)](#OneSig.SignedBy)
  * [func (o *OneSig) SignedBytes() ([]byte, error)](#OneSig.SignedBytes)
* [type PubKey](#PubKey)
  * [func LoadPubKey(data []byte) (PubKey, error)](#LoadPubKey)
  * [func (p PubKey) Address() []byte](#PubKey.Address)
  * [func (p PubKey) Bytes() []byte](#PubKey.Bytes)
  * [func (p PubKey) Equals(pk crypto.PubKey) bool](#PubKey.Equals)
  * [func (p PubKey) KeyString() string](#PubKey.KeyString)
  * [func (p PubKey) VerifyBytes(msg []byte, sig crypto.Signature) bool](#PubKey.VerifyBytes)
* [type Signature](#Signature)
  * [func (s Signature) Bytes() []byte](#Signature.Bytes)
  * [func (s Signature) Equals(cs crypto.Signature) bool](#Signature.Equals)
  * [func (s Signature) IsZero() bool](#Signature.IsZero)
  * [func (s Signature) String() string](#Signature.String)


#### <a name="pkg-files">Package files</a>
[crypto.go](/src/github.com/tendermint/light-client/mock/crypto.go) [docs.go](/src/github.com/tendermint/light-client/mock/docs.go) [reader.go](/src/github.com/tendermint/light-client/mock/reader.go) [signable.go](/src/github.com/tendermint/light-client/mock/signable.go) [value.go](/src/github.com/tendermint/light-client/mock/value.go) 





## <a name="Reader">func</a> [Reader](/src/target/reader.go?s=531:562#L18)
``` go
func Reader() lc.SignableReader
```
Reader constructs a SignableReader that can parse OneSig and MultiSig

TODO: add some args to configure go-wire, rather than relying on init???7



## <a name="ValueReader">func</a> [ValueReader](/src/target/value.go?s=438:471#L9)
``` go
func ValueReader() lc.ValueReader
```
ValueReader returns a mock ValueReader for test cases




## <a name="ByteValue">type</a> [ByteValue](/src/target/value.go?s=244:265#L1)
``` go
type ByteValue []byte
```
ByteValue is a simple way to pass byte slices unparsed as Values
meant only for test cases, as usually you will want the ValueReader
to actually make sense of the structure










### <a name="ByteValue.Bytes">func</a> (ByteValue) [Bytes](/src/target/value.go?s=267:300#L1)
``` go
func (b ByteValue) Bytes() []byte
```



## <a name="ByteValueReader">type</a> [ByteValueReader](/src/target/value.go?s=665:694#L17)
``` go
type ByteValueReader struct{}
```
ByteValueReader is a simple implementation that just wraps the bytes
in ByteValue.

Intended for testing where there is no app-specific data structure










### <a name="ByteValueReader.ReadValue">func</a> (ByteValueReader) [ReadValue](/src/target/value.go?s=696:767#L19)
``` go
func (b ByteValueReader) ReadValue(key, value []byte) (lc.Value, error)
```



## <a name="MultiSig">type</a> [MultiSig](/src/target/signable.go?s=1214:1266#L45)
``` go
type MultiSig struct {
    Data []byte
    // contains filtered or unexported fields
}
```
MultiSig is a Signable implementation that can be used to
record the values and inspect them later.  It performs no validation.

It supports an arbitrary number of signatures







### <a name="NewMultiSig">func</a> [NewMultiSig](/src/target/signable.go?s=1339:1378#L55)
``` go
func NewMultiSig(data []byte) *MultiSig
```




### <a name="MultiSig.Bytes">func</a> (\*MultiSig) [Bytes](/src/target/signable.go?s=1486:1519#L63)
``` go
func (m *MultiSig) Bytes() []byte
```



### <a name="MultiSig.Sign">func</a> (\*MultiSig) [Sign](/src/target/signable.go?s=1540:1613#L67)
``` go
func (m *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```



### <a name="MultiSig.SignedBy">func</a> (\*MultiSig) [SignedBy](/src/target/signable.go?s=1685:1739#L73)
``` go
func (m *MultiSig) SignedBy() ([]crypto.PubKey, error)
```



### <a name="MultiSig.SignedBytes">func</a> (\*MultiSig) [SignedBytes](/src/target/signable.go?s=1940:1988#L84)
``` go
func (m *MultiSig) SignedBytes() ([]byte, error)
```



## <a name="OneSig">type</a> [OneSig](/src/target/signable.go?s=299:383#L3)
``` go
type OneSig struct {
    Data   []byte
    PubKey crypto.PubKey
    Sig    crypto.Signature
}
```
OneSig is a Signable implementation that can be used to
record the values and inspect them later.  It performs no validation.







### <a name="NewSig">func</a> [NewSig](/src/target/signable.go?s=385:417#L9)
``` go
func NewSig(data []byte) *OneSig
```




### <a name="OneSig.Bytes">func</a> (\*OneSig) [Bytes](/src/target/signable.go?s=521:552#L17)
``` go
func (o *OneSig) Bytes() []byte
```



### <a name="OneSig.Sign">func</a> (\*OneSig) [Sign](/src/target/signable.go?s=573:644#L21)
``` go
func (o *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```



### <a name="OneSig.SignedBy">func</a> (\*OneSig) [SignedBy](/src/target/signable.go?s=764:816#L30)
``` go
func (o *OneSig) SignedBy() ([]crypto.PubKey, error)
```



### <a name="OneSig.SignedBytes">func</a> (\*OneSig) [SignedBytes](/src/target/signable.go?s=934:980#L37)
``` go
func (o *OneSig) SignedBytes() ([]byte, error)
```



## <a name="PubKey">type</a> [PubKey](/src/target/crypto.go?s=195:229#L2)
``` go
type PubKey struct {
    Val []byte
}
```
PubKey lets us wrap some bytes to provide crypto PubKey for those
methods that just act on it.







### <a name="LoadPubKey">func</a> [LoadPubKey](/src/target/crypto.go?s=792:836#L35)
``` go
func LoadPubKey(data []byte) (PubKey, error)
```




### <a name="PubKey.Address">func</a> (PubKey) [Address](/src/target/crypto.go?s=400:432#L15)
``` go
func (p PubKey) Address() []byte
```
Address is just the pubkey with some constant prepended




### <a name="PubKey.Bytes">func</a> (PubKey) [Bytes](/src/target/crypto.go?s=291:321#L10)
``` go
func (p PubKey) Bytes() []byte
```



### <a name="PubKey.Equals">func</a> (PubKey) [Equals](/src/target/crypto.go?s=648:693#L28)
``` go
func (p PubKey) Equals(pk crypto.PubKey) bool
```



### <a name="PubKey.KeyString">func</a> (PubKey) [KeyString](/src/target/crypto.go?s=480:514#L19)
``` go
func (p PubKey) KeyString() string
```



### <a name="PubKey.VerifyBytes">func</a> (PubKey) [VerifyBytes](/src/target/crypto.go?s=554:620#L23)
``` go
func (p PubKey) VerifyBytes(msg []byte, sig crypto.Signature) bool
```



## <a name="Signature">type</a> [Signature](/src/target/crypto.go?s=980:1017#L41)
``` go
type Signature struct {
    Val []byte
}
```
Signature lets us wrap some bytes to provide crypto signature for those
methods that just act on it.










### <a name="Signature.Bytes">func</a> (Signature) [Bytes](/src/target/crypto.go?s=1085:1118#L49)
``` go
func (s Signature) Bytes() []byte
```



### <a name="Signature.Equals">func</a> (Signature) [Equals](/src/target/crypto.go?s=1274:1325#L61)
``` go
func (s Signature) Equals(cs crypto.Signature) bool
```



### <a name="Signature.IsZero">func</a> (Signature) [IsZero](/src/target/crypto.go?s=1138:1170#L53)
``` go
func (s Signature) IsZero() bool
```



### <a name="Signature.String">func</a> (Signature) [String](/src/target/crypto.go?s=1200:1234#L57)
``` go
func (s Signature) String() string
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
