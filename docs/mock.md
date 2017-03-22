

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
  * [func (m *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#MultiSig.Sign)
  * [func (m *MultiSig) SignBytes() []byte](#MultiSig.SignBytes)
  * [func (m *MultiSig) Signers() ([]crypto.PubKey, error)](#MultiSig.Signers)
  * [func (m *MultiSig) TxBytes() ([]byte, error)](#MultiSig.TxBytes)
* [type OneSig](#OneSig)
  * [func NewSig(data []byte) *OneSig](#NewSig)
  * [func (o *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#OneSig.Sign)
  * [func (o *OneSig) SignBytes() []byte](#OneSig.SignBytes)
  * [func (o *OneSig) Signers() ([]crypto.PubKey, error)](#OneSig.Signers)
  * [func (o *OneSig) TxBytes() ([]byte, error)](#OneSig.TxBytes)


#### <a name="pkg-files">Package files</a>
[docs.go](/src/github.com/tendermint/light-client/mock/docs.go) [reader.go](/src/github.com/tendermint/light-client/mock/reader.go) [signable.go](/src/github.com/tendermint/light-client/mock/signable.go) [value.go](/src/github.com/tendermint/light-client/mock/value.go) 





## <a name="Reader">func</a> [Reader](/src/target/reader.go?s=571:602#L19)
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



## <a name="MultiSig">type</a> [MultiSig](/src/target/signable.go?s=1194:1246#L45)
``` go
type MultiSig struct {
    Data []byte
    // contains filtered or unexported fields
}
```
MultiSig is a Signable implementation that can be used to
record the values and inspect them later.  It performs no validation.

It supports an arbitrary number of signatures







### <a name="NewMultiSig">func</a> [NewMultiSig](/src/target/signable.go?s=1319:1358#L55)
``` go
func NewMultiSig(data []byte) *MultiSig
```




### <a name="MultiSig.Sign">func</a> (\*MultiSig) [Sign](/src/target/signable.go?s=1517:1590#L67)
``` go
func (m *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```



### <a name="MultiSig.SignBytes">func</a> (\*MultiSig) [SignBytes](/src/target/signable.go?s=1459:1496#L63)
``` go
func (m *MultiSig) SignBytes() []byte
```



### <a name="MultiSig.Signers">func</a> (\*MultiSig) [Signers](/src/target/signable.go?s=1662:1715#L73)
``` go
func (m *MultiSig) Signers() ([]crypto.PubKey, error)
```



### <a name="MultiSig.TxBytes">func</a> (\*MultiSig) [TxBytes](/src/target/signable.go?s=1916:1960#L84)
``` go
func (m *MultiSig) TxBytes() ([]byte, error)
```



## <a name="OneSig">type</a> [OneSig](/src/target/signable.go?s=287:371#L3)
``` go
type OneSig struct {
    Data   []byte
    PubKey crypto.PubKey
    Sig    crypto.Signature
}
```
OneSig is a Signable implementation that can be used to
record the values and inspect them later.  It performs no validation.







### <a name="NewSig">func</a> [NewSig](/src/target/signable.go?s=373:405#L9)
``` go
func NewSig(data []byte) *OneSig
```




### <a name="OneSig.Sign">func</a> (\*OneSig) [Sign](/src/target/signable.go?s=558:629#L21)
``` go
func (o *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```



### <a name="OneSig.SignBytes">func</a> (\*OneSig) [SignBytes](/src/target/signable.go?s=502:537#L17)
``` go
func (o *OneSig) SignBytes() []byte
```



### <a name="OneSig.Signers">func</a> (\*OneSig) [Signers](/src/target/signable.go?s=749:800#L30)
``` go
func (o *OneSig) Signers() ([]crypto.PubKey, error)
```



### <a name="OneSig.TxBytes">func</a> (\*OneSig) [TxBytes](/src/target/signable.go?s=918:960#L37)
``` go
func (o *OneSig) TxBytes() ([]byte, error)
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
