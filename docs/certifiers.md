

# certifiers
`import "github.com/tendermint/light-client/certifiers"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [func GenHeader(chainID string, height int, txs types.Txs, vals []*types.Validator, appHash []byte) *types.Header](#GenHeader)
* [func NoPathFound(err error) bool](#NoPathFound)
* [func SeedNotFound(err error) bool](#SeedNotFound)
* [func TooMuchChange(err error) bool](#TooMuchChange)
* [func ValidatorsChanged(err error) bool](#ValidatorsChanged)
* [func VerifyCommitAny(old, cur *types.ValidatorSet, chainID string, blockID types.BlockID, height int, commit *types.Commit) error](#VerifyCommitAny)
* [type CacheProvider](#CacheProvider)
  * [func NewCacheProvider(providers ...Provider) CacheProvider](#NewCacheProvider)
  * [func (c CacheProvider) GetByHash(hash []byte) (s Seed, err error)](#CacheProvider.GetByHash)
  * [func (c CacheProvider) GetByHeight(h int) (s Seed, err error)](#CacheProvider.GetByHeight)
  * [func (c CacheProvider) StoreSeed(seed Seed) (err error)](#CacheProvider.StoreSeed)
* [type DynamicCertifier](#DynamicCertifier)
  * [func NewDynamic(chainID string, vals []*types.Validator) *DynamicCertifier](#NewDynamic)
  * [func (c *DynamicCertifier) Certify(check lc.Checkpoint) error](#DynamicCertifier.Certify)
  * [func (c *DynamicCertifier) Update(check lc.Checkpoint, vals []*types.Validator) error](#DynamicCertifier.Update)
* [type InquiringCertifier](#InquiringCertifier)
  * [func NewInquiring(chainID string, vals []*types.Validator, provider Provider) *InquiringCertifier](#NewInquiring)
  * [func (c *InquiringCertifier) Certify(check lc.Checkpoint) error](#InquiringCertifier.Certify)
  * [func (c *InquiringCertifier) Update(check lc.Checkpoint, vals []*types.Validator) error](#InquiringCertifier.Update)
* [type MemStoreProvider](#MemStoreProvider)
  * [func NewMemStoreProvider() *MemStoreProvider](#NewMemStoreProvider)
  * [func (m *MemStoreProvider) GetByHash(hash []byte) (Seed, error)](#MemStoreProvider.GetByHash)
  * [func (m *MemStoreProvider) GetByHeight(h int) (Seed, error)](#MemStoreProvider.GetByHeight)
  * [func (m *MemStoreProvider) StoreSeed(seed Seed) error](#MemStoreProvider.StoreSeed)
* [type MissingProvider](#MissingProvider)
  * [func NewMissingProvider() MissingProvider](#NewMissingProvider)
  * [func (_ MissingProvider) GetByHash(_ []byte) (Seed, error)](#MissingProvider.GetByHash)
  * [func (_ MissingProvider) GetByHeight(_ int) (Seed, error)](#MissingProvider.GetByHeight)
  * [func (_ MissingProvider) StoreSeed(_ Seed) error](#MissingProvider.StoreSeed)
* [type Provider](#Provider)
* [type Seed](#Seed)
  * [func (s Seed) Hash() []byte](#Seed.Hash)
  * [func (s Seed) Height() int](#Seed.Height)
* [type Seeds](#Seeds)
  * [func (s Seeds) Len() int](#Seeds.Len)
  * [func (s Seeds) Less(i, j int) bool](#Seeds.Less)
  * [func (s Seeds) Swap(i, j int)](#Seeds.Swap)
* [type StaticCertifier](#StaticCertifier)
  * [func NewStatic(chainID string, vals []*types.Validator) *StaticCertifier](#NewStatic)
  * [func (c *StaticCertifier) Certify(check lc.Checkpoint) error](#StaticCertifier.Certify)
  * [func (c *StaticCertifier) Hash() []byte](#StaticCertifier.Hash)
* [type ValKeys](#ValKeys)
  * [func GenSecValKeys(n int) ValKeys](#GenSecValKeys)
  * [func GenValKeys(n int) ValKeys](#GenValKeys)
  * [func (v ValKeys) Change(i int) ValKeys](#ValKeys.Change)
  * [func (v ValKeys) Extend(n int) ValKeys](#ValKeys.Extend)
  * [func (v ValKeys) ExtendSec(n int) ValKeys](#ValKeys.ExtendSec)
  * [func (v ValKeys) GenCheckpoint(chainID string, height int, txs types.Txs, vals []*types.Validator, appHash []byte, first, last int) lc.Checkpoint](#ValKeys.GenCheckpoint)
  * [func (v ValKeys) SignHeader(header *types.Header, first, last int) *types.Commit](#ValKeys.SignHeader)
  * [func (v ValKeys) ToValidators(init, inc int64) []*types.Validator](#ValKeys.ToValidators)


#### <a name="pkg-files">Package files</a>
[dynamic.go](/src/github.com/tendermint/light-client/certifiers/dynamic.go) [helper.go](/src/github.com/tendermint/light-client/certifiers/helper.go) [inquirer.go](/src/github.com/tendermint/light-client/certifiers/inquirer.go) [provider.go](/src/github.com/tendermint/light-client/certifiers/provider.go) [static.go](/src/github.com/tendermint/light-client/certifiers/static.go) 





## <a name="GenHeader">func</a> [GenHeader](/src/target/helper.go?s=2716:2829#L92)
``` go
func GenHeader(chainID string, height int, txs types.Txs,
    vals []*types.Validator, appHash []byte) *types.Header
```


## <a name="NoPathFound">func</a> [NoPathFound](/src/target/inquirer.go?s=360:392#L7)
``` go
func NoPathFound(err error) bool
```
NoPathFound asserts whether an error is due to no path of
validators in provider from where we are to where we want to be



## <a name="SeedNotFound">func</a> [SeedNotFound](/src/target/provider.go?s=315:348#L8)
``` go
func SeedNotFound(err error) bool
```
SeedNotFound asserts whether an error is due to missing data



## <a name="TooMuchChange">func</a> [TooMuchChange](/src/target/dynamic.go?s=423:457#L9)
``` go
func TooMuchChange(err error) bool
```
TooMuchChange asserts whether and error is due to too much change
between these validators sets



## <a name="ValidatorsChanged">func</a> [ValidatorsChanged](/src/target/static.go?s=342:380#L6)
``` go
func ValidatorsChanged(err error) bool
```
ValidatorsChanged asserts whether and error is due
to a differing validator set



## <a name="VerifyCommitAny">func</a> [VerifyCommitAny](/src/target/dynamic.go?s=2992:3122#L93)
``` go
func VerifyCommitAny(old, cur *types.ValidatorSet, chainID string,
    blockID types.BlockID, height int, commit *types.Commit) error
```
VerifyCommitAny will check to see if the set would
be valid with a different validator set.

old is the validator set that we know
* over 2/3 of the power in old signed this block

cur is the validator set that signed this block
* only votes from old are sufficient for 2/3 majority


	in the new set as well

That means that:
* 10% of the valset can't just declare themselves kings
* If the validator set is 3x old size, we need more proof to trust

*** TODO: move this.
It belongs in tendermint/types/validator_set.go: VerifyCommitAny




## <a name="CacheProvider">type</a> [CacheProvider](/src/target/provider.go?s=1583:1634#L51)
``` go
type CacheProvider struct {
    Providers []Provider
}
```
CacheProvider allows you to place one or more caches in front of a source
Provider.  It runs through them in order until a match is found.
So you can keep a local cache, and check with the network if
no data is there.







### <a name="NewCacheProvider">func</a> [NewCacheProvider](/src/target/provider.go?s=1636:1694#L55)
``` go
func NewCacheProvider(providers ...Provider) CacheProvider
```




### <a name="CacheProvider.GetByHash">func</a> (CacheProvider) [GetByHash](/src/target/provider.go?s=2208:2273#L84)
``` go
func (c CacheProvider) GetByHash(hash []byte) (s Seed, err error)
```



### <a name="CacheProvider.GetByHeight">func</a> (CacheProvider) [GetByHeight](/src/target/provider.go?s=2031:2092#L74)
``` go
func (c CacheProvider) GetByHeight(h int) (s Seed, err error)
```



### <a name="CacheProvider.StoreSeed">func</a> (CacheProvider) [StoreSeed](/src/target/provider.go?s=1864:1919#L64)
``` go
func (c CacheProvider) StoreSeed(seed Seed) (err error)
```
StoreSeed tries to add the seed to all providers.

Aborts on first error it encounters (closest provider)




## <a name="DynamicCertifier">type</a> [DynamicCertifier](/src/target/dynamic.go?s=781:858#L18)
``` go
type DynamicCertifier struct {
    Cert       *StaticCertifier
    LastHeight int
}
```
DynamicCertifier uses a StaticCertifier to evaluate the checkpoint
but allows for a change, if we present enough proof

TODO: do we keep a long history so we can use our memory to validate
checkpoints from previously valid validator sets????







### <a name="NewDynamic">func</a> [NewDynamic](/src/target/dynamic.go?s=860:934#L23)
``` go
func NewDynamic(chainID string, vals []*types.Validator) *DynamicCertifier
```




### <a name="DynamicCertifier.Certify">func</a> (\*DynamicCertifier) [Certify](/src/target/dynamic.go?s=1129:1190#L35)
``` go
func (c *DynamicCertifier) Certify(check lc.Checkpoint) error
```
Certify handles this with




### <a name="DynamicCertifier.Update">func</a> (\*DynamicCertifier) [Update](/src/target/dynamic.go?s=1526:1611#L48)
``` go
func (c *DynamicCertifier) Update(check lc.Checkpoint, vals []*types.Validator) error
```
Update will verify if this is a valid change and update
the certifying validator set if safe to do so.

Returns an error if update is impossible (invalid proof or TooMuchChange)




## <a name="InquiringCertifier">type</a> [InquiringCertifier](/src/target/inquirer.go?s=458:526#L11)
``` go
type InquiringCertifier struct {
    Cert *DynamicCertifier
    Provider
}
```






### <a name="NewInquiring">func</a> [NewInquiring](/src/target/inquirer.go?s=528:625#L16)
``` go
func NewInquiring(chainID string, vals []*types.Validator, provider Provider) *InquiringCertifier
```




### <a name="InquiringCertifier.Certify">func</a> (\*InquiringCertifier) [Certify](/src/target/inquirer.go?s=724:787#L23)
``` go
func (c *InquiringCertifier) Certify(check lc.Checkpoint) error
```



### <a name="InquiringCertifier.Update">func</a> (\*InquiringCertifier) [Update](/src/target/inquirer.go?s=983:1070#L35)
``` go
func (c *InquiringCertifier) Update(check lc.Checkpoint, vals []*types.Validator) error
```



## <a name="MemStoreProvider">type</a> [MemStoreProvider](/src/target/provider.go?s=2390:2595#L94)
``` go
type MemStoreProvider struct {
    // contains filtered or unexported fields
}
```






### <a name="NewMemStoreProvider">func</a> [NewMemStoreProvider](/src/target/provider.go?s=2597:2641#L101)
``` go
func NewMemStoreProvider() *MemStoreProvider
```




### <a name="MemStoreProvider.GetByHash">func</a> (\*MemStoreProvider) [GetByHash](/src/target/provider.go?s=3450:3513#L138)
``` go
func (m *MemStoreProvider) GetByHash(hash []byte) (Seed, error)
```



### <a name="MemStoreProvider.GetByHeight">func</a> (\*MemStoreProvider) [GetByHeight](/src/target/provider.go?s=3188:3247#L127)
``` go
func (m *MemStoreProvider) GetByHeight(h int) (Seed, error)
```



### <a name="MemStoreProvider.StoreSeed">func</a> (\*MemStoreProvider) [StoreSeed](/src/target/provider.go?s=2825:2878#L112)
``` go
func (m *MemStoreProvider) StoreSeed(seed Seed) error
```



## <a name="MissingProvider">type</a> [MissingProvider](/src/target/provider.go?s=3734:3763#L149)
``` go
type MissingProvider struct{}
```
MissingProvider doens't store anything, always a miss
Designed as a mock for testing







### <a name="NewMissingProvider">func</a> [NewMissingProvider](/src/target/provider.go?s=3765:3806#L151)
``` go
func NewMissingProvider() MissingProvider
```




### <a name="MissingProvider.GetByHash">func</a> (MissingProvider) [GetByHash](/src/target/provider.go?s=4014:4072#L159)
``` go
func (_ MissingProvider) GetByHash(_ []byte) (Seed, error)
```



### <a name="MissingProvider.GetByHeight">func</a> (MissingProvider) [GetByHeight](/src/target/provider.go?s=3902:3959#L156)
``` go
func (_ MissingProvider) GetByHeight(_ int) (Seed, error)
```



### <a name="MissingProvider.StoreSeed">func</a> (MissingProvider) [StoreSeed](/src/target/provider.go?s=3838:3886#L155)
``` go
func (_ MissingProvider) StoreSeed(_ Seed) error
```



## <a name="Provider">type</a> [Provider](/src/target/provider.go?s=1097:1351#L39)
``` go
type Provider interface {
    StoreSeed(seed Seed) error
    // GetByHeight returns the closest seed at with height <= h
    GetByHeight(h int) (Seed, error)
    // GetByHash returns a seed exactly matching this validator hash
    GetByHash(hash []byte) (Seed, error)
}
```
Provider is used to get more validators by other means

TODO: Also FileStoreProvider, NodeProvider, ...










## <a name="Seed">type</a> [Seed](/src/target/provider.go?s=576:642#L15)
``` go
type Seed struct {
    lc.Checkpoint
    Validators []*types.Validator
}
```
Seed is a checkpoint and the actual validator set, the base info you
need to update to a given point, assuming knowledge of some previous
validator set










### <a name="Seed.Hash">func</a> (Seed) [Hash](/src/target/provider.go?s=706:733#L24)
``` go
func (s Seed) Hash() []byte
```



### <a name="Seed.Height">func</a> (Seed) [Height](/src/target/provider.go?s=644:670#L20)
``` go
func (s Seed) Height() int
```



## <a name="Seeds">type</a> [Seeds](/src/target/provider.go?s=782:799#L28)
``` go
type Seeds []Seed
```









### <a name="Seeds.Len">func</a> (Seeds) [Len](/src/target/provider.go?s=801:825#L30)
``` go
func (s Seeds) Len() int
```



### <a name="Seeds.Less">func</a> (Seeds) [Less](/src/target/provider.go?s=907:941#L32)
``` go
func (s Seeds) Less(i, j int) bool
```



### <a name="Seeds.Swap">func</a> (Seeds) [Swap](/src/target/provider.go?s=849:878#L31)
``` go
func (s Seeds) Swap(i, j int)
```



## <a name="StaticCertifier">type</a> [StaticCertifier](/src/target/static.go?s=690:782#L15)
``` go
type StaticCertifier struct {
    ChainID string
    VSet    *types.ValidatorSet
    // contains filtered or unexported fields
}
```
StaticCertifier assumes a static set of validators, set on
initilization and checks against them.

Good for testing or really simple chains.  You will want a
better implementation when the validator set can actually change.







### <a name="NewStatic">func</a> [NewStatic](/src/target/static.go?s=784:856#L21)
``` go
func NewStatic(chainID string, vals []*types.Validator) *StaticCertifier
```




### <a name="StaticCertifier.Certify">func</a> (\*StaticCertifier) [Certify](/src/target/static.go?s=1137:1197#L39)
``` go
func (c *StaticCertifier) Certify(check lc.Checkpoint) error
```



### <a name="StaticCertifier.Hash">func</a> (\*StaticCertifier) [Hash](/src/target/static.go?s=951:990#L28)
``` go
func (c *StaticCertifier) Hash() []byte
```



## <a name="ValKeys">type</a> [ValKeys](/src/target/helper.go?s=216:245#L2)
``` go
type ValKeys []crypto.PrivKey
```
we use this to simulate signing with many keys







### <a name="GenSecValKeys">func</a> [GenSecValKeys](/src/target/helper.go?s=861:894#L28)
``` go
func GenSecValKeys(n int) ValKeys
```
GenSecValKeys produces an array of secp256k1 private keys to generate commits


### <a name="GenValKeys">func</a> [GenValKeys](/src/target/helper.go?s=315:345#L5)
``` go
func GenValKeys(n int) ValKeys
```
GenValKeys produces an array of private keys to generate commits





### <a name="ValKeys.Change">func</a> (ValKeys) [Change](/src/target/helper.go?s=489:527#L14)
``` go
func (v ValKeys) Change(i int) ValKeys
```
Change replaces the key at index i




### <a name="ValKeys.Extend">func</a> (ValKeys) [Extend](/src/target/helper.go?s=684:722#L22)
``` go
func (v ValKeys) Extend(n int) ValKeys
```
Extend adds n more keys (to remove, just take a slice)




### <a name="ValKeys.ExtendSec">func</a> (ValKeys) [ExtendSec](/src/target/helper.go?s=1070:1111#L37)
``` go
func (v ValKeys) ExtendSec(n int) ValKeys
```
Extend adds n more secp256k1 keys (to remove, just take a slice)




### <a name="ValKeys.GenCheckpoint">func</a> (ValKeys) [GenCheckpoint](/src/target/helper.go?s=3177:3323#L109)
``` go
func (v ValKeys) GenCheckpoint(chainID string, height int, txs types.Txs,
    vals []*types.Validator, appHash []byte, first, last int) lc.Checkpoint
```
GenCheckpoint calls GenHeader and SignHeader and combines them into a Checkpoint




### <a name="ValKeys.SignHeader">func</a> (ValKeys) [SignHeader](/src/target/helper.go?s=1710:1790#L55)
``` go
func (v ValKeys) SignHeader(header *types.Header, first, last int) *types.Commit
```
SignHeader properly signs the header with all keys from first to last exclusive




### <a name="ValKeys.ToValidators">func</a> (ValKeys) [ToValidators](/src/target/helper.go?s=1416:1481#L46)
``` go
func (v ValKeys) ToValidators(init, inc int64) []*types.Validator
```
ToValidators produces a list of validators from the set of keys
The first key has weight `init` and it increases by `inc` every step
so we can have all the same weight, or a simple linear distribution
(should be enough for testing)








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
