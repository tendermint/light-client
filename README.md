# Tendermint Light Client

Once you have built your amazing new ABCi app, possibly with the help of the [basecoin framework](https://github.com/tendermint/basecoin/blob/develop/README.md) and the [example apps](https://github.com/tendermint/basecoin-examples/blob/master/README.md), you now want to make some sort of client to access it.

Basecoin comes with a [nice simple cli](https://github.com/tendermint/basecoin-examples/blob/master/tutorial.md), that is nice for testing and developing your application, but is probably not the first thing you would hand off to the future users of your blockchain.  You want something pretty, something like a web page or mobile app.  But where do you start?  The [tendermint rpc](https://tendermint.com/docs/internals/rpc) is documented and a good start, but there are style plenty of opaque hex strings (byte slices) returned that may need go code to decipher.  And how do your properly sign that basecoin AppTx anyway?

If you're still with me, then you should take a deeper look at this repo.  The purpose here is to build a helper library to perform most common actions one would want to do with a client, make it extensible to easily support custom transaction and data types, and provide bindings to other languages.

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

Please take a look at the godoc for this package, as care was taken to make it approachable.

The general concept (create, list, sign, import, export...) was inspired by [Ethereum Key Management](https://github.com/ethereum/go-ethereum/wiki/Managing-Your-Accounts).  The code and architecture was developed completely independently (I didn't even look at their code, so as not to possibly violate the GPLv3 license).

### Creating Transactions

If the server is writen in go, especially if it is based on basecoin, generating the transaction (with `go-wire`) and signing with `go-crypto` is very hard to reliably do from any language other than go.  This library will produce an interface, where the caller can simply provide the data to the transaction, as well as a keyname and passphrase, and the library will generate a byte array (or hex/base64 string) with the properly encoded and signed transaction. If running the proxy server, we will also post it directly to the blockchain for you.

This should be written in a way that it is easy to add custom transaction encodings to a custom build of this library without forking the codebase (just importing it and passing some initialization info).  Much the

The beginning of this work is in [Signable and Poster](./transactions.go#L21-L43), and the example [transaction wrappers](./sign). Other applications should create concrete transactions that implement `Signable`, and then we just need to expose this via cli or proxy.

**Needs Work**


### RPC Wrapper

First, we create a [nicer interface](./rpc) to call the tendermint RPC, to avoid a lot of boilerplate casting and marshaling of data types when we call the RPC. This is a literal client of the existing tendermint RPC, and will track the most recent version of tendermint rpc (and multiple versions once 1.0 is released).

**Needs work**

Secondly, we create an abstract interface `Node` representing the needs of a light client, which takes the results from tendermint rpc and does some validation and other processing on them.  This doesn't pull in any other dependencies and is the interface that is used internally in the package for all code that needs to interact with the tendermint server.

**TODO**

### Viewing Data

When querying data, we often get binary data back from the server.  We need a way to unpack this data (using domain knowledge of the application's data format) and return it as JSON (or generic dictionary).  Something like how the basecoin cli [queries the account](https://github.com/tendermint/basecoin/blob/develop/cmd/basecoin/commands/utils.go#L59-L81), and then [renders it as json](https://github.com/tendermint/basecoin/blob/develop/cmd/basecoin/commands/query.go#L118-L124)

**TODO**

### Verifying Proofs

Beyond simply querying data from a blockchain, we often want **undeniable, cryptographic proof** of its validity.  This is the reason for exposing merkle proofs as first class objects in the new query [request](https://github.com/tendermint/abci/blob/develop/types/types.pb.go#L718-L723) and [response](https://github.com/tendermint/abci/blob/develop/types/types.pb.go#L1413-L1421).  However, this Proof byte slice, still generally requires go code to [parse and validate](https://github.com/tendermint/go-merkle/blob/develop/iavl_proof.go#L14-L42).

It is important to provide access to this essential functionality in a light client library, so we can provide this security to any UI we wish to build.

**TODO**

### Tracking Validators

A very important part of verifying a proof is to ensure that it's roothash was committed to a valid block.  We can get the proof and a block header, and validate they match.  But then we need to prove that the block is properly signed by a validator set.  The first step is to verify the signatures, the second that these public keys actually have the authority to sign blocks (not just some random keys).

If a chain has a static validator set, we just need to feed in this set upon initialization, and can verify a block header with a simple 2/3 check. If it can update, we need an easy way to track them, without having to process every block header.

There is background information in a [github issue](https://github.com/tendermint/tendermint/issues/377), the [block structure](https://tendermint.com/docs/internals/block-structure), and the [concept of validators](https://tendermint.com/docs/internals/validators)

**TODO**: Link to Bucky's document about this algorithm

**TODO**: Implement

## References

Some other projects that may inspire this:

* [Project Trillian](https://github.com/google/trillian) - Verifiable data structures (Apache 2.0)

