# Tendermint Light Client

Once you have built your amazing new ABCi app, possibly with the help of the [basecoin framework](https://github.com/tendermint/basecoin/blob/develop/README.md) and the [example apps](https://github.com/tendermint/basecoin-examples/blob/master/README.md), you now want to make some sort of client to access it.

Basecoin comes with a [nice simple cli](https://github.com/tendermint/basecoin-examples/blob/master/tutorial.md), that is nice for testing and developing your application, but is probably not the first thing you would hand off to the future users of your blockchain.  You want something pretty, something like a web page or mobile app.  But where do you start?  The [tendermint rpc](https://tendermint.com/docs/internals/rpc) is documented and a good start, but there are style plenty of opaque hex strings (byte slices) returned that may need go code to decipher.  And how do your properly sign that basecoin AppTx anyway?

If you're still with me, then you should take a deeper look at this repo.  The purpose here is to build a helper library to perform most common actions one would want to do with a client, make it extensible to easily support custom transaction and data types, and provide bindings to other languages.

## Important Note

The current code will work with tendermint 0.9 and basecoin 0.2.

However, I realized this was much more difficult than it needed to be, and ended up integrating an even more powerful rpc client in tendermint itself, as well as a more flexible json un/marshaler through a new repository `go-data`, along with changes to `go-crypto` and `basecoin` to support it.

So, feel free to use and read this code, but be aware that I will do a major rewrite to make use of those other libraries, so don't code to much on the current API.

### More reasons

If I develop a desktop/mobile client I don't want either:

* Sync the entire chain on my device (be a non-validating node) or
* Blindly trust whichever node I communicate with to be honest

One goal of this project is to provide a library that pulls together all the crypto and algorithms, so given a relatively recent (< unbonding period) known validator set, one can get indisputable proof that data is in the chain (current state) or detect if the node is lying to the client.

Tendermint RPC exposes a lot of info, but a malicious node could return any data it wants to queries, or even to block headers, even making up fake signatures from non-existent validators to justify it.  This is a lot of logic to get right, and I want to make a small, easy to use library, that does this for you, so people can just build nice UI.

I refer to the tendermint consensus engine and rpc as a `node`, the abci app as an `app` (which implicitly runs in a trusted environment with a node), and any user-interface that is external to the validator network as a `client`.

These external clients who have no strong trust relationship with any node, just the validator set as a whole. Beyond a nice mobile or desktop application, the cosmos hub is another important example of a `client`, that needs undeniable proof without syncing the full chain.

## Code Documentation

I try to provide extensive documentation in the code, to make it as easy to use as possible.  I highly recommend running `godoc -http:6060` locally and browsing to [interactive documentation](http://localhost:6060/pkg/github.com/tendermint/light-client/). If you have not downloaded the code locally, you can also browse the [generated godoc](./docs) thanks to the excellent [godoc2md](https://github.com/davecheney/godoc2md) tool

## Bindings

First, the library will provide a nice API to call directly from other programs written in go and thus integrate nicely with headless clients, and provide an easy way to extend this functionality via a different interface.

Second, it will include a proxy web server with a simple JSON REST API that you can run locally and verify and sign all interaction with a blockchain. This can be connected over unix sockets (more secure) or local TCP port (to easily expose tendermint from a webapp - be careful about CORS for security). This is primarily intended for webapp/javascript development, but anyone else who feels running a separate binary and making REST calls is easier than compiling against a go library.

(It may include gRPC bindings to the proxy as well, but those are of questionable use, as javascript clients cannot call gRPC, and native apps would likely use other bindings.  Note: actually gRPC support from browser is an [much discussed proposal](https://github.com/grpc/grpc/issues/8682) with it's own [private repo](https://github.com/grpc/grpc-web).  Maybe in some months this is possible).

Third, it will expose a subset of this functionality through a simple cli, inspired by the basecoin cli.  This could be used for development, or embedding in shell scripts (simple integration tests?).

The next usage would be building [gomobile bindings](https://github.com/golang/go/wiki/Mobile) for Android and iOS allow mobile developers to integrate a tendermint app as easily as any other web service.

Finally, we could export a nice `.so` file with a simple C ABI using [-buildmode=c-shared](https://golang.org/cmd/go/#hdr-Description_of_build_modes).  From this point, we could link it with a [C/C++ desktop app](http://stackoverflow.com/questions/12066279/using-c-libraries-for-c-programs), produce [python bindings](https://blog.filippo.io/building-python-modules-with-go-1-5/), call from Java [via JNI](https://blog.dogan.io/2015/08/15/java-jni-jnr-go/), even call it from [erlang](http://andrealeopardi.com/posts/using-c-from-elixir-with-nifs/) if that's what makes you happy.

## Functionality

### Key Management

We need to manage private keys locally, store them securely (passphrase protected), sign transactions, and display their addresses (for receiving transactions).

This code is in the [cryptostore directory](./cryptostore). It uses a composable architecture to allow you to customize the type of key (currently Ed25519 or Secp256k1), what symetric encryption algorithm we use to passphrase-encode the key for storage, and where we store the key (currently in-memory or on disk, could be extend to eg. vault, etcd, db...)

Please take a look at the godoc for this package, as care was taken to make it approachable. Note that you can find the [storage options](./storage) in their own package.  They can be used to store eg. validator lists as well.

The general concept (create, list, sign, import, export...) was inspired by [Ethereum Key Management](https://github.com/ethereum/go-ethereum/wiki/Managing-Your-Accounts).  The code and architecture was developed completely independently (I didn't even look at their code, so as not to possibly violate the GPLv3 license).

### Creating Transactions

If the server is writen in go, especially if it is based on basecoin, generating the transaction (with `go-wire`) and signing with `go-crypto` is very hard to reliably do from any language other than go.  This library will produce an interface, where the caller can simply provide the data to the transaction, as well as a keyname and passphrase, and the library will generate a byte array (or hex/base64 string) with the properly encoded and signed transaction. If running the proxy server, we will also post it directly to the blockchain for you.

This should be written in a way that it is easy to add custom transaction encodings to a custom build of this library without forking the codebase (just importing it and passing some initialization info).

We extracted these ideas and present the results in three interfaces:

* [Signable](./transactions.go#L28-L48), which can be implemented by any transaction
* [Signer](./transactions.go#L50-L54), which is responsible for attaching signatures to the `Signable` and is implemented by [cryptostore.Manager](./cryptostore/holder.go#L9-L14)
* [Poster](./transactions.go#L56-L62), which pulls together a `Signer` and `Broadcaster`, so we can `Post` the `Signable` directly to tendermint in one shot.

The infrastructure is in place, it is just up to an app to create transactions that implement the `Signable` interface, to take advantage of the lightclient. We provide various implementations that can simply be used by applications that don't have special requirements:

* `tx.OneSig` - supports one signature using go-crypto (`tx.New(data)`)
* `tx.MultiSig` - supports multi-sig using go-crypto (`tx.NewMulti(data)`)
* `mock.OneSig` - records a single signature for use in testing
* `mock.MultiSig` - records multi-sig for use in testing

**TODO** update basecoin transactions to support the `Signable` interface


### RPC Wrapper

First, we created a [simple interface](./rpc) to call the tendermint RPC, to avoid a lot of boilerplate casting and marshaling of data types when we call the RPC. This is a literal client of the existing tendermint RPC, and will track the most recent version of tendermint rpc (and multiple versions once 1.0 is released). The main advantage over using `github.com/tendermint/go-rpc/client` directly is that we handle casting the types and following the rpc interfaces, allowing you to just call simple, type-safe methods.

Secondly, we create two abstract interfaces `Broadcaster` and `Checker` representing the needs of a light client, either sending info to tendermint, or getting and validating a key-value pair.  These interfaces (and a but more) are implemented by [rpc.Node](https://github.com/tendermint/light-client/blob/develop/rpc/node.go), which takes the results from `rpc.HTTPClient` and does some validation and other processing on them.  This is responsible for all parsing of tendermint structures, as well as app-specific data structures.  To that end, it must be configurable to allow custom deserializing of `Proof` and `Value`.

### Viewing Data

When querying data, we often get binary data back from the server.  We need a way to unpack this data (using domain knowledge of the application's data format) and return it as JSON (or generic dictionary).  Something like how the basecoin cli [queries the account](https://github.com/tendermint/basecoin/blob/develop/cmd/basecoin/commands/utils.go#L59-L81), and then [renders it as json](https://github.com/tendermint/basecoin/blob/develop/cmd/basecoin/commands/query.go#L118-L124)

The data must be returned from the app as bytes that match the merkle proof, thus it is the responsibility of this library to parse it.  Since this is application-specific domain knowledge, we cannot program this, but rather allow the application designer to provide us this information in a special `ValueReader` interface, which knows how to read the application-specific values stored in the merkle tree, and convert them into a struct that can be json-encoded, or otherwise transformed for other binding.

There is a simple implementation `mock.ByteValueReader` for testing, which just wraps the `[]byte` into a `Value` without doing any parsing.  For integration with your application, please provide your application-specific logic.

### Verifying Proofs

Beyond simply querying data from a blockchain, we often want **undeniable, cryptographic proof** of its validity.  This is the reason for exposing merkle proofs as first class objects in the new query [request](https://github.com/tendermint/abci/blob/develop/types/types.pb.go#L718-L723) and [response](https://github.com/tendermint/abci/blob/develop/types/types.pb.go#L1413-L1421).  However, this Proof byte slice, still generally requires go code to [parse and validate](https://github.com/tendermint/go-merkle/blob/develop/iavl_proof.go#L14-L42).

It is important to provide access to this essential functionality in a light client library, so we can provide this security to any UI we wish to build. We provide a few methods to allow this to work:

* There are two important methods to get information, implemented by `rpc.Node`:
  * `Checker.Prove` method will parse the data into a format that can be used to prove a key-value slice leads to a root hash
  * `Checker.SignedHeader` retrieves a block header, containing the AppHash, as well as the Precommit signatures that prove which validators signed off on this commit
* `Certifier.Certify` uses some out-of-band knowledge of the validator set to check that these signatures are sufficient proof (> 2/3 votes)
  * `rpc.StaticCertifier` is a simple implementation, using a static validator set, for test code or very simple apps
  * **TODO**: special `TrackingCertifier`, which is seeded with some validator set and does occasional queries on the blockchain to update it safely.

Note: here is some info on the [block structure](https://tendermint.com/docs/internals/block-structure) we parse to get this data

### Tracking Validators

Once there is a dynamic validator set, the issue of verifying a block becomes a bit more tricky.  This is still a work in progress, but so far I have yet to see an app that dynamically updates the validator set. In any case "coming soon".s

There is background information in a [github issue](https://github.com/tendermint/tendermint/issues/377), and the [concept of validators](https://tendermint.com/docs/internals/validators)

**TODO**: Link to Bucky's document about this algorithm

**TODO**: Implement

## Extensibility

Of course, every application has its own transaction types, its own way of signing them, and its own data structures for which it wishes to present proofs. While we can work with interfaces that know how to serialize themselves (like `Signable`), this does not tell us how to *deserialize* the objects.  In fact, it is impossible to attach this information to an interface (as we allow many potentially unknown concrete types). We use two ways for this in the code:

###Hard code dependencies

Some things are hard coded for simplicity, or as there are no other options available.  I tried to abstract `crypto.PubKey` and `crypto.Signature` and quickly came to the issue that interface methods applying to interfaces cannot be implemented by other interfaces with the same footprint. I'll explain later, but for technical reasons it was quite difficult not to just use go-crypto, and anyway, this is supposed to be a general purpose implementation of any cryptographic algorithm, so I left it as a direct dependency.  If you wish to add new algorithms, you have to add them to go-crypto (by forking).

###External configuration

Other concepts need maximum flexibility, as there are many options, most specific to the application itself. For this we need to provide our own "Read" functions somehow.  My approach is to define interfaces for the Readers and pass them into the constructors of any struct that will need to deserialize bytes. Maybe you have another approach, like dynamically registering with `init` in your package, as is done to extend the basecoin commands. This is open to discussion as to which approach is the simplest to use and most maintainable.

* `Proof` is loaded by `rpc.Node` when running queries.  `rpc.Node` uses a `rpc.ProofReader` to define this behavior, which is set to `MerkleReader` in the constructor.  If needed, this can be made more configurable.
* `Signable` will need to be passed into the program somehow, and as most of the bindings (json, jni, etc) don't handle go structs so well, we provide a `SignableReader` that accepts json as `[]byte` and returns a concrete implementation of a `Signable`. This is app-specific, and since the byte layout doesn't matter until we sign the struct, we can expose a simple interface to pass in unsigned transaction. This would be mainly used by the bindings on input.
* `Value` - when we get data back from the server, it is just a bunch of bytes, exactly as stored in the app, which is important to be able to validate the proofs.  However, once we validate it, most clients would rather just have a struct or json object and some stamp that it was, in fact, cryptographically proven. Thus, we add a `ValueReader`, like the `SignableReader`, but this is executed by `rpc.Node` upon receiving the response to a query. Since an app can (and usually does) support multiple data formats, we also provide the key to help the Reader decipher the bytes.

###Interfaces referring to interfaces

TODO: why `lightclient.PubKey` implementing `crypto.PubKey` without importing `crypto` is impossible. digging deep with go interfaces.

## References

Some other projects that may inspire this:

* [Project Trillian](https://github.com/google/trillian) - Verifiable data structures (Apache 2.0)

