

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

ReadSignableBinary() can be used by an ACBi app to deserialize a
signed, packed OneSig or MultiSig object.  You must write your
own SignableReader to translate json into a binary data blob and
wrap it with OneSig or MultiSig to provide an external interface
for the api proxy or language bindings.

Maybe this package should be moved out of lightclient, more to a repo
designed for application support, as that is the main usecase?




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func ReadSignableBinary(data []byte) (keys.Signable, error)](#ReadSignableBinary)
* [type B58Data](#B58Data)
  * [func (d B58Data) MarshalJSON() ([]byte, error)](#B58Data.MarshalJSON)
  * [func (d *B58Data) UnmarshalJSON(b []byte) (err error)](#B58Data.UnmarshalJSON)
* [type HexData](#HexData)
  * [func (h HexData) MarshalJSON() ([]byte, error)](#HexData.MarshalJSON)
  * [func (h *HexData) UnmarshalJSON(b []byte) (err error)](#HexData.UnmarshalJSON)
* [type JSONPubKey](#JSONPubKey)
  * [func (p JSONPubKey) MarshalJSON() ([]byte, error)](#JSONPubKey.MarshalJSON)
  * [func (p *JSONPubKey) UnmarshalJSON(b []byte) error](#JSONPubKey.UnmarshalJSON)
* [type MultiSig](#MultiSig)
  * [func LoadMulti(serialized []byte) (*MultiSig, error)](#LoadMulti)
  * [func NewMulti(data []byte) *MultiSig](#NewMulti)
  * [func (s *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#MultiSig.Sign)
  * [func (s *MultiSig) SignBytes() []byte](#MultiSig.SignBytes)
  * [func (s *MultiSig) Signers() ([]crypto.PubKey, error)](#MultiSig.Signers)
  * [func (s *MultiSig) TxBytes() ([]byte, error)](#MultiSig.TxBytes)
* [type OneSig](#OneSig)
  * [func Load(serialized []byte) (*OneSig, error)](#Load)
  * [func New(data []byte) *OneSig](#New)
  * [func (s *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#OneSig.Sign)
  * [func (s *OneSig) SignBytes() []byte](#OneSig.SignBytes)
  * [func (s *OneSig) Signers() ([]crypto.PubKey, error)](#OneSig.Signers)
  * [func (s *OneSig) TxBytes() ([]byte, error)](#OneSig.TxBytes)
* [type RawValue](#RawValue)
  * [func NewValue(val []byte) RawValue](#NewValue)
  * [func (v RawValue) Bytes() []byte](#RawValue.Bytes)


#### <a name="pkg-files">Package files</a>
[docs.go](/src/github.com/tendermint/light-client/tx/docs.go) [multi.go](/src/github.com/tendermint/light-client/tx/multi.go) [one.go](/src/github.com/tendermint/light-client/tx/one.go) [pubkey.go](/src/github.com/tendermint/light-client/tx/pubkey.go) [reader.go](/src/github.com/tendermint/light-client/tx/reader.go) [value.go](/src/github.com/tendermint/light-client/tx/value.go) 


## <a name="pkg-constants">Constants</a>
``` go
const RawValueType = "raw"
```



## <a name="ReadSignableBinary">func</a> [ReadSignableBinary](/src/target/reader.go?s=353:412#L13)
``` go
func ReadSignableBinary(data []byte) (keys.Signable, error)
```



## <a name="B58Data">type</a> [B58Data](/src/target/value.go?s=638:657#L22)
``` go
type B58Data []byte
```
B58Data let's us treat a byte slice as base58, like bitcoin addresses










### <a name="B58Data.MarshalJSON">func</a> (B58Data) [MarshalJSON](/src/target/value.go?s=895:941#L36)
``` go
func (d B58Data) MarshalJSON() ([]byte, error)
```



### <a name="B58Data.UnmarshalJSON">func</a> (\*B58Data) [UnmarshalJSON](/src/target/value.go?s=659:712#L24)
``` go
func (d *B58Data) UnmarshalJSON(b []byte) (err error)
```



## <a name="HexData">type</a> [HexData](/src/target/value.go?s=213:232#L2)
``` go
type HexData []byte
```
HexData let's us treat a byte slice as hex data, rather than default base64










### <a name="HexData.MarshalJSON">func</a> (HexData) [MarshalJSON](/src/target/value.go?s=461:507#L16)
``` go
func (h HexData) MarshalJSON() ([]byte, error)
```



### <a name="HexData.UnmarshalJSON">func</a> (\*HexData) [UnmarshalJSON](/src/target/value.go?s=234:287#L4)
``` go
func (h *HexData) UnmarshalJSON(b []byte) (err error)
```



## <a name="JSONPubKey">type</a> [JSONPubKey](/src/target/pubkey.go?s=61:102#L1)
``` go
type JSONPubKey struct {
    crypto.PubKey
}
```









### <a name="JSONPubKey.MarshalJSON">func</a> (JSONPubKey) [MarshalJSON](/src/target/pubkey.go?s=340:389#L10)
``` go
func (p JSONPubKey) MarshalJSON() ([]byte, error)
```



### <a name="JSONPubKey.UnmarshalJSON">func</a> (\*JSONPubKey) [UnmarshalJSON](/src/target/pubkey.go?s=145:195#L1)
``` go
func (p *JSONPubKey) UnmarshalJSON(b []byte) error
```
TODO: use B58Data instead of HexData?




## <a name="MultiSig">type</a> [MultiSig](/src/target/multi.go?s=357:409#L4)
``` go
type MultiSig struct {
    // contains filtered or unexported fields
}
```
MultiSig lets us wrap arbitrary data with a go-crypto signature

TODO: rethink how we want to integrate this with KeyStore so it makes
more sense (particularly the verify method)







### <a name="LoadMulti">func</a> [LoadMulti](/src/target/multi.go?s=554:606#L18)
``` go
func LoadMulti(serialized []byte) (*MultiSig, error)
```

### <a name="NewMulti">func</a> [NewMulti](/src/target/multi.go?s=482:518#L14)
``` go
func NewMulti(data []byte) *MultiSig
```




### <a name="MultiSig.Sign">func</a> (\*MultiSig) [Sign](/src/target/multi.go?s=1149:1222#L38)
``` go
func (s *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```
Sign will add a signature and pubkey.

Depending on the Signable, one may be able to call this multiple times for multisig
Returns error if called with invalid data or too many times




### <a name="MultiSig.SignBytes">func</a> (\*MultiSig) [SignBytes](/src/target/multi.go?s=897:934#L30)
``` go
func (s *MultiSig) SignBytes() []byte
```
SignBytes returns the original data passed into `NewSig`




### <a name="MultiSig.Signers">func</a> (\*MultiSig) [Signers](/src/target/multi.go?s=1593:1646#L52)
``` go
func (s *MultiSig) Signers() ([]crypto.PubKey, error)
```
Signers will return the public key(s) that signed if the signature
is valid, or an error if there is any issue with the signature,
including if there are no signatures




### <a name="MultiSig.TxBytes">func</a> (\*MultiSig) [TxBytes](/src/target/multi.go?s=2106:2150#L71)
``` go
func (s *MultiSig) TxBytes() ([]byte, error)
```
TxBytes serializes the Sig to send it to a tendermint app.
It returns an error if the Sig was never Signed.




## <a name="OneSig">type</a> [OneSig](/src/target/one.go?s=355:439#L4)
``` go
type OneSig struct {
    // contains filtered or unexported fields
}
```
OneSig lets us wrap arbitrary data with a go-crypto signature

TODO: rethink how we want to integrate this with KeyStore so it makes
more sense (particularly the verify method)







### <a name="Load">func</a> [Load](/src/target/one.go?s=504:549#L14)
``` go
func Load(serialized []byte) (*OneSig, error)
```

### <a name="New">func</a> [New](/src/target/one.go?s=441:470#L10)
``` go
func New(data []byte) *OneSig
```




### <a name="OneSig.Sign">func</a> (\*OneSig) [Sign](/src/target/one.go?s=1086:1157#L34)
``` go
func (s *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```
Sign will add a signature and pubkey.

Depending on the Signable, one may be able to call this multiple times for multisig
Returns error if called with invalid data or too many times




### <a name="OneSig.SignBytes">func</a> (\*OneSig) [SignBytes](/src/target/one.go?s=836:871#L26)
``` go
func (s *OneSig) SignBytes() []byte
```
SignBytes returns the original data passed into `NewSig`




### <a name="OneSig.Signers">func</a> (\*OneSig) [Signers](/src/target/one.go?s=1588:1639#L52)
``` go
func (s *OneSig) Signers() ([]crypto.PubKey, error)
```
Signers will return the public key(s) that signed if the signature
is valid, or an error if there is any issue with the signature,
including if there are no signatures




### <a name="OneSig.TxBytes">func</a> (\*OneSig) [TxBytes](/src/target/one.go?s=1980:2022#L66)
``` go
func (s *OneSig) TxBytes() ([]byte, error)
```
TxBytes serializes the Sig to send it to a tendermint app.
It returns an error if the Sig was never Signed.




## <a name="RawValue">type</a> [RawValue](/src/target/value.go?s=1019:1102#L43)
``` go
type RawValue struct {
    Type  string  `json:"type"`
    Value HexData `json:"value"`
}
```






### <a name="NewValue">func</a> [NewValue](/src/target/value.go?s=1104:1138#L48)
``` go
func NewValue(val []byte) RawValue
```




### <a name="RawValue.Bytes">func</a> (RawValue) [Bytes](/src/target/value.go?s=1211:1243#L55)
``` go
func (v RawValue) Bytes() []byte
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
