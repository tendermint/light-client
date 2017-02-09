

# cryptostore
`import "github.com/tendermint/light-client/cryptostore"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
package cryptostore maintains everything needed for doing public-key signing and
key management in software, based on the go-crypto library from tendermint.

It is flexible, and allows the user to provide a key generation algorithm
(currently Ed25519 or Secp256k1), an encoder to passphrase-encrypt our keys
when storing them (currently SecretBox from NaCl), and a method to persist
the keys (currently FileStorage like ssh, or MemStorage for tests).
It should be relatively simple to write your own implementation of these
interfaces to match your specific security requirements.

Note that the private keys are never exposed outside the package, and the
interface of Manager could be implemented by an HSM in the future for
enhanced security.  It would require a completely different implementation
however.

This Manager aims to implement Signer and KeyManager interfaces, along
with some extensions to allow importing/exporting keys and updating the
passphrase.

Encoder and Generator implementations are currently in this package,
lightclient.Storage implementations exist as subpackages of
lightclient/storage




## <a name="pkg-index">Index</a>
* [type Encoder](#Encoder)
* [type GenFunc](#GenFunc)
  * [func (f GenFunc) Generate() crypto.PrivKey](#GenFunc.Generate)
* [type Generator](#Generator)
* [type Manager](#Manager)
  * [func New(gen Generator, coder Encoder, store lightclient.Storage) Manager](#New)
  * [func (s Manager) Create(name, passphrase string) error](#Manager.Create)
  * [func (s Manager) Delete(name, passphrase string) error](#Manager.Delete)
  * [func (s Manager) Export(name, oldpass, transferpass string) ([]byte, error)](#Manager.Export)
  * [func (s Manager) Get(name string) (lightclient.KeyInfo, error)](#Manager.Get)
  * [func (s Manager) Import(name, newpass, transferpass string, data []byte) error](#Manager.Import)
  * [func (s Manager) List() (lightclient.KeyInfos, error)](#Manager.List)
  * [func (s Manager) Sign(name, passphrase string, tx lightclient.Signable) error](#Manager.Sign)
  * [func (s Manager) Update(name, oldpass, newpass string) error](#Manager.Update)


#### <a name="pkg-files">Package files</a>
[docs.go](/src/github.com/tendermint/light-client/cryptostore/docs.go) [enc_storage.go](/src/github.com/tendermint/light-client/cryptostore/enc_storage.go) [encoder.go](/src/github.com/tendermint/light-client/cryptostore/encoder.go) [generator.go](/src/github.com/tendermint/light-client/cryptostore/generator.go) [holder.go](/src/github.com/tendermint/light-client/cryptostore/holder.go) 






## <a name="Encoder">type</a> [Encoder](/src/target/encoder.go?s=440:583#L8)
``` go
type Encoder interface {
    Encrypt(key crypto.PrivKey, pass string) ([]byte, error)
    Decrypt(data []byte, pass string) (crypto.PrivKey, error)
}
```
Encoder is used to encrypt any key with a passphrase for storage.

This should use a well-designed symetric encryption algorithm


``` go
var (
    // SecretBox uses the algorithm from NaCL to store secrets securely
    SecretBox Encoder = secretbox{}
    // Noop doesn't do any encryption, should only be used in test code
    Noop Encoder = noop{}
)
```









## <a name="GenFunc">type</a> [GenFunc](/src/target/generator.go?s=453:487#L8)
``` go
type GenFunc func() crypto.PrivKey
```
GenFunc is a helper to transform a function into a Generator










### <a name="GenFunc.Generate">func</a> (GenFunc) [Generate](/src/target/generator.go?s=489:531#L10)
``` go
func (f GenFunc) Generate() crypto.PrivKey
```



## <a name="Generator">type</a> [Generator](/src/target/generator.go?s=332:387#L3)
``` go
type Generator interface {
    Generate() crypto.PrivKey
}
```
Generator determines the type of private key the keystore creates


``` go
var (
    // GenEd25519 produces Ed25519 private keys
    GenEd25519 Generator = GenFunc(genEd25519)
    // GenSecp256k1 produces Secp256k1 private keys
    GenSecp256k1 Generator = GenFunc(genSecp256)
)
```









## <a name="Manager">type</a> [Manager](/src/target/holder.go?s=177:237#L1)
``` go
type Manager struct {
    // contains filtered or unexported fields
}
```
Manager combines encyption and storage implementation to provide
a full-featured key manager







### <a name="New">func</a> [New](/src/target/holder.go?s=239:312#L2)
``` go
func New(gen Generator, coder Encoder, store lightclient.Storage) Manager
```




### <a name="Manager.Create">func</a> (Manager) [Create](/src/target/holder.go?s=790:844#L24)
``` go
func (s Manager) Create(name, passphrase string) error
```
Create adds a new key to the storage engine, returning error if
another key already stored under this name




### <a name="Manager.Delete">func</a> (Manager) [Delete](/src/target/holder.go?s=2723:2777#L88)
``` go
func (s Manager) Delete(name, passphrase string) error
```
Delete removes key forever

TODO: make sure we have the passphrase before deleting it (for security)




### <a name="Manager.Export">func</a> (Manager) [Export](/src/target/holder.go?s=2008:2083#L63)
``` go
func (s Manager) Export(name, oldpass, transferpass string) ([]byte, error)
```
Export decodes the private key with the current password, encodes
it with a secure one-time password and generates a sequence that can be
Imported by another Manager

This is designed to copy from one device to another, or provide backups
during version updates.




### <a name="Manager.Get">func</a> (Manager) [Get](/src/target/holder.go?s=1188:1250#L38)
``` go
func (s Manager) Get(name string) (lightclient.KeyInfo, error)
```
Get returns the public information about one key




### <a name="Manager.Import">func</a> (Manager) [Import](/src/target/holder.go?s=2407:2485#L76)
``` go
func (s Manager) Import(name, newpass, transferpass string, data []byte) error
```
Import accepts bytes generated by Export along with the same transferpass
If they are valid, it stores the password under the given name with the
new passphrase.




### <a name="Manager.List">func</a> (Manager) [List](/src/target/holder.go?s=987:1040#L30)
``` go
func (s Manager) List() (lightclient.KeyInfos, error)
```
List loads the keys from the storage and enforces alphabetical order




### <a name="Manager.Sign">func</a> (Manager) [Sign](/src/target/holder.go?s=1487:1564#L47)
``` go
func (s Manager) Sign(name, passphrase string, tx lightclient.Signable) error
```
Sign will modify the Signable in order to attach a valid signature with
this public key

If no key for this name, or the passphrase doesn't match, returns an error




### <a name="Manager.Update">func</a> (Manager) [Update](/src/target/holder.go?s=3018:3078#L96)
``` go
func (s Manager) Update(name, oldpass, newpass string) error
```
Update changes the passphrase with which a already stored key is encoded.

oldpass must be the current passphrase used for encoding, newpass will be
the only valid passphrase from this time forward








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
