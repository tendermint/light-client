

# basecoin
`import "github.com/tendermint/light-client/extensions/basecoin"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [type Account](#Account)
  * [func (a Account) Bytes() []byte](#Account.Bytes)
* [type AccountData](#AccountData)
* [type AppDataReader](#AppDataReader)
* [type AppTx](#AppTx)
  * [func (a *AppTx) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#AppTx.Sign)
  * [func (a *AppTx) SignBytes() []byte](#AppTx.SignBytes)
  * [func (a *AppTx) Signers() ([]crypto.PubKey, error)](#AppTx.Signers)
  * [func (a *AppTx) TxBytes() ([]byte, error)](#AppTx.TxBytes)
* [type BasecoinTx](#BasecoinTx)
  * [func NewBasecoinTx(chainID string) BasecoinTx](#NewBasecoinTx)
  * [func (t BasecoinTx) ReadSignable(data []byte) (keys.Signable, error)](#BasecoinTx.ReadSignable)
  * [func (t BasecoinTx) RegisterParser(name, txType string, reader TxReader)](#BasecoinTx.RegisterParser)
* [type BasecoinValues](#BasecoinValues)
  * [func NewBasecoinValues() *BasecoinValues](#NewBasecoinValues)
  * [func (t *BasecoinValues) ReadValue(key, value []byte) (lc.Value, error)](#BasecoinValues.ReadValue)
  * [func (t *BasecoinValues) RegisterPlugin(reader lc.ValueReader)](#BasecoinValues.RegisterPlugin)
* [type SendTx](#SendTx)
  * [func (s *SendTx) Sign(pubkey crypto.PubKey, sig crypto.Signature) error](#SendTx.Sign)
  * [func (s *SendTx) SignBytes() []byte](#SendTx.SignBytes)
  * [func (s *SendTx) Signers() ([]crypto.PubKey, error)](#SendTx.Signers)
  * [func (s *SendTx) TxBytes() ([]byte, error)](#SendTx.TxBytes)
* [type TxReader](#TxReader)
* [type TxType](#TxType)


#### <a name="pkg-files">Package files</a>
[apptx.go](/src/github.com/tendermint/light-client/extensions/basecoin/apptx.go) [parse.go](/src/github.com/tendermint/light-client/extensions/basecoin/parse.go) [sendtx.go](/src/github.com/tendermint/light-client/extensions/basecoin/sendtx.go) [signable.go](/src/github.com/tendermint/light-client/extensions/basecoin/signable.go) [values.go](/src/github.com/tendermint/light-client/extensions/basecoin/values.go) 


## <a name="pkg-constants">Constants</a>
``` go
const AccountType = "account"
```




## <a name="Account">type</a> [Account](/src/target/values.go?s=1587:1733#L55)
``` go
type Account struct {
    Type  string      `json:"type"`
    Value AccountData `json:"value"` // TODO: custom encoding?
    // contains filtered or unexported fields
}
```









### <a name="Account.Bytes">func</a> (Account) [Bytes](/src/target/values.go?s=1930:1961#L67)
``` go
func (a Account) Bytes() []byte
```



## <a name="AccountData">type</a> [AccountData](/src/target/values.go?s=1735:1928#L61)
``` go
type AccountData struct {
    PubKey   crypto.PubKeyS `json:"pub_key,omitempty"` // May be empty, if not known.
    Sequence int            `json:"sequence"`
    Balance  bc.Coins       `json:"coins"`
}
```









## <a name="AppDataReader">type</a> [AppDataReader](/src/target/signable.go?s=514:587#L19)
``` go
type AppDataReader func(name, txType string, json []byte) ([]byte, error)
```
AppDataReader takes a plugin name and txType and some json
and serializes it into a binary format










## <a name="AppTx">type</a> [AppTx](/src/target/apptx.go?s=216:295#L1)
``` go
type AppTx struct {
    Tx *bc.AppTx
    // contains filtered or unexported fields
}
```









### <a name="AppTx.Sign">func</a> (\*AppTx) [Sign](/src/target/apptx.go?s=873:943#L29)
``` go
func (a *AppTx) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```
Sign will add a signature and pubkey.

Depending on the Signable, one may be able to call this multiple times for multisig
Returns error if called with invalid data or too many times




### <a name="AppTx.SignBytes">func</a> (\*AppTx) [SignBytes](/src/target/apptx.go?s=605:639#L21)
``` go
func (a *AppTx) SignBytes() []byte
```
SignBytes returned the unsigned bytes, needing a signature




### <a name="AppTx.Signers">func</a> (\*AppTx) [Signers](/src/target/apptx.go?s=1250:1300#L41)
``` go
func (a *AppTx) Signers() ([]crypto.PubKey, error)
```
Signers will return the public key(s) that signed if the signature
is valid, or an error if there is any issue with the signature,
including if there are no signatures




### <a name="AppTx.TxBytes">func</a> (\*AppTx) [TxBytes](/src/target/apptx.go?s=1541:1582#L50)
``` go
func (a *AppTx) TxBytes() ([]byte, error)
```
TxBytes returns the transaction data as well as all signatures
It should return an error if Sign was never called




## <a name="BasecoinTx">type</a> [BasecoinTx](/src/target/signable.go?s=209:280#L5)
``` go
type BasecoinTx struct {
    // contains filtered or unexported fields
}
```






### <a name="NewBasecoinTx">func</a> [NewBasecoinTx](/src/target/signable.go?s=282:327#L10)
``` go
func NewBasecoinTx(chainID string) BasecoinTx
```




### <a name="BasecoinTx.ReadSignable">func</a> (BasecoinTx) [ReadSignable](/src/target/signable.go?s=738:806#L25)
``` go
func (t BasecoinTx) ReadSignable(data []byte) (keys.Signable, error)
```
Turn json into a signable object




### <a name="BasecoinTx.RegisterParser">func</a> (BasecoinTx) [RegisterParser](/src/target/signable.go?s=1130:1202#L40)
``` go
func (t BasecoinTx) RegisterParser(name, txType string, reader TxReader)
```



## <a name="BasecoinValues">type</a> [BasecoinValues](/src/target/values.go?s=229:285#L3)
``` go
type BasecoinValues struct {
    // contains filtered or unexported fields
}
```






### <a name="NewBasecoinValues">func</a> [NewBasecoinValues](/src/target/values.go?s=287:327#L7)
``` go
func NewBasecoinValues() *BasecoinValues
```




### <a name="BasecoinValues.ReadValue">func</a> (\*BasecoinValues) [ReadValue](/src/target/values.go?s=454:525#L14)
``` go
func (t *BasecoinValues) ReadValue(key, value []byte) (lc.Value, error)
```
Turn merkle binary into a json-able struct




### <a name="BasecoinValues.RegisterPlugin">func</a> (\*BasecoinValues) [RegisterPlugin](/src/target/values.go?s=810:872#L28)
``` go
func (t *BasecoinValues) RegisterPlugin(reader lc.ValueReader)
```



## <a name="SendTx">type</a> [SendTx](/src/target/sendtx.go?s=216:299#L1)
``` go
type SendTx struct {
    Tx *bc.SendTx
    // contains filtered or unexported fields
}
```









### <a name="SendTx.Sign">func</a> (\*SendTx) [Sign](/src/target/sendtx.go?s=873:944#L29)
``` go
func (s *SendTx) Sign(pubkey crypto.PubKey, sig crypto.Signature) error
```
Sign will add a signature and pubkey.

Depending on the Signable, one may be able to call this multiple times for multisig
Returns error if called with invalid data or too many times




### <a name="SendTx.SignBytes">func</a> (\*SendTx) [SignBytes](/src/target/sendtx.go?s=604:639#L21)
``` go
func (s *SendTx) SignBytes() []byte
```
SignBytes returned the unsigned bytes, needing a signature




### <a name="SendTx.Signers">func</a> (\*SendTx) [Signers](/src/target/sendtx.go?s=1323:1374#L42)
``` go
func (s *SendTx) Signers() ([]crypto.PubKey, error)
```
Signers will return the public key(s) that signed if the signature
is valid, or an error if there is any issue with the signature,
including if there are no signatures




### <a name="SendTx.TxBytes">func</a> (\*SendTx) [TxBytes](/src/target/sendtx.go?s=1604:1646#L51)
``` go
func (s *SendTx) TxBytes() ([]byte, error)
```
TxBytes returns the transaction data as well as all signatures
It should return an error if Sign was never called




## <a name="TxReader">type</a> [TxReader](/src/target/signable.go?s=653:700#L22)
``` go
type TxReader func(json []byte) ([]byte, error)
```
TxReader handles parsing and serializing one particular type










## <a name="TxType">type</a> [TxType](/src/target/signable.go?s=111:207#L1)
``` go
type TxType struct {
    Type string           `json:"type"`
    Data *json.RawMessage `json:"data"`
}
```













- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
