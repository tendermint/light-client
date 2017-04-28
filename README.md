# Tendermint Light Client

Once you have built your amazing new ABCi app, possibly with the help of the [Basecoin framework](https://github.com/tendermint/basecoin/blob/develop/README.md) and the [example apps](https://github.com/tendermint/basecoin-examples/blob/master/README.md), you now want to make some sort of client to access it.

Basecoin comes with a [nice simple cli](https://github.com/tendermint/basecoin-examples/blob/master/tutorial.md), that is nice for testing and developing your application, but is probably not the first thing you would hand off to the future users of your blockchain.  You want something pretty, something like a web page or mobile app.  But where do you start?  The [Tendermint RPC](https://tendermint.com/docs/internals/rpc) is documented and a good start, but there are style plenty of opaque hex strings (byte slices) returned that may need go code to decipher.  And how do your properly sign that Basecoin AppTx anyway?

If you're still with me, then you should take a deeper look at this repo.  The purpose here is to build a helper library to perform most common actions one would want to do with a client, make it extensible to easily support custom transaction and data types, and provide bindings to other languages.

## Important Note

The current code will work with Tendermint 0.9 and Basecoin 0.3.1.

This is an incomplete but usable state. The main purpose here is validating headers along with signatures (here refered to as seeds). And some basic checks for abci and tx proofs.  To make this better, the cli needs to be more aware of the actual data structures used in the particular abci app, and expose them.  This is an area of development.

The develop branch is tracking Tendermint 0.9.1 release and a number of refactors for internal libraries, and there should be another release soon.  There will also be some better api for this version.

However, a number of desired features require some breaking changes to the Tendermint RPC itself, which are planned for 0.10, and that version of the light-client should be more secure in the face of malicious nodes.

### Let's go already

 1. Compile the code with `make install`
 2. Run a Basecoin 0.3.1 instance on some machine (or better yet, a cluster)
 3. Initialize the local client:
    * Run `basecli init --chain_id <ID> --node <host>:<port>`
    * This will ask you to confirm the validator set of the running cluster and verify the chain id is correct, check this step.
    * You must use `--force-reset` to overwrite this dir
    * You can also use `-r` or `--root` to store the chain config in a custom dir (and support two chains at once)
4. After some time (and possible validator set changes), update the the current chain state securely
    * Run `basecli seeds update`
    * Run `basecli seeds show --height X` to show the closest seed to that height if available
5. Try `basecli --help` to see what is available and try it out.

### More reasons

If I develop a desktop/mobile client I don't want either:

* Sync the entire chain on my device (be a non-validating node) or
* Blindly trust whichever node I communicate with to be honest

One goal of this project is to provide a library that pulls together all the crypto and algorithms, so given a relatively recent (< unbonding period) known validator set, one can get indisputable proof that data is in the chain (current state) or detect if the node is lying to the client.

Tendermint RPC exposes a lot of info, but a malicious node could return any data it wants to queries, or even to block headers, even making up fake signatures from non-existent validators to justify it.  This is a lot of logic to get right, and I want to make a small, easy to use library, that does this for you, so people can just build nice UI.

I refer to the tendermint consensus engine and rpc as a `node`, the abci app as an `app` (which implicitly runs in a trusted environment with a node), and any user-interface that is external to the validator network as a `client`.

These external clients who have no strong trust relationship with any node, just the validator set as a whole. Beyond a nice mobile or desktop application, the cosmos hub is another important example of a `client`, that needs undeniable proof without syncing the full chain.

## Code Documentation

I try to provide extensive documentation in the code, to make it as easy to use as possible.  I highly recommend running `godoc -http :6060` locally and browsing to [interactive documentation](http://localhost:6060/pkg/github.com/tendermint/light-client/). If you have not downloaded the code locally, you can also browse the [generated godoc](./docs) thanks to the excellent [godoc2md](https://github.com/davecheney/godoc2md) tool

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

This code is now a separate repo called [go-keys](https://github.com/tendermint/go-keys) and is embedded as a subcommand in `basecli keys`. Try that with `list` and `new` to see info.  Also, note the `-o json` command to see a machine readable format with more info.

The general concept (create, list, sign, import, export...) was inspired by [Ethereum Key Management](https://github.com/ethereum/go-ethereum/wiki/Managing-Your-Accounts).  The code and architecture was developed completely independently (I didn't even look at their code, so as not to possibly violate the GPLv3 license).

### Tracking Validators

Unless you want to blindly trust the node you talk with, you need to trace every response back to a hash in a block header and validate the commit signatures of that block header match the proper validator set.  If there is a contant validator set, you store it locally upon initialization of the client, and check against that every time.

Once there is a dynamic validator set, the issue of verifying a block becomes a bit more tricky. There is background information in a [github issue](https://github.com/tendermint/tendermint/issues/377), and the [concept of validators](https://tendermint.com/docs/internals/validators).

I refer to a complete proof at one height as a seed (block header, block commit signatures, validator set).  All code to validate these seeds and to use these seeds to validate other headers can be found in the [certifiers package](https://github.com/tendermint/light-client/tree/master/certifiers).

In short, if there is a block at height H with a known (trusted) validator set V, and another block at height H' (H' > H) with validator set V' != V, then we want a way to safely update it. First, get the new (unconfirmed) validator set V' and verify H' is internally consistent and properly signed by this V'. Assuming it is a valid block, we check that at least 2/3 of the validators in V signed it, meaning it would also be valid under our old assumptions.  That should be enough, but we can also check that the V counts for at least 2/3 of the total votes in H' for extra safety (we can have a discussion if this is strictly required). If we can verify all this, then we can accept H' and V' as valid and use that to validate all blocks X > H'.

If we cannot update directly from H -> H' because there was too much changes to the validator set, then we can look for some Hm (H < Hm < H') with a validator set Vm.  Then we try to update H -> Hm and Hm -> H' in two separate steps.  If one of these steps doesn't work, then we continue bisecting, until we eventually have to externally validate the valdiator set changes at every block.

There is only one problem now... it is impossible to get old validator sets from the tendermint RPC API.  You can currently copy these seeds from one light client to another, with full verification.  Thus, if client A cannot update from height 800 to 1200, but client B has the full seed for height 1000, he can export it and client A can import that seed, doing the full verification procedure above to advance from 800 to 1000, and then move up to the current block height.

Since we never trust any server in this protocol, only the signatures themselves, it doesn't matter if the seed comes from a (possibly malicious) node or a (possibly malicious) user.  We can accept it or reject it only based on our trusted validator set and cryptographic proofs. This makes it extremely important to verify that you have the proper validator set when initializing the client, as that is the root of all trust.

Or course, this assumes that the known block is within the unbonding period to avoid the "nothing at stake" problem. If you haven't seen the state in a few months, you will need to manually verify the new validator set hash using off-chain means (the same as getting the initial hash).

### Viewing Data

When querying data, we often get binary data back from the server.  We need a way to unpack this data (using domain knowledge of the application's data format) and return it as JSON (or generic dictionary).  Something like how the basecoin cli [queries the account](https://github.com/tendermint/basecoin/blob/develop/cmd/basecoin/commands/utils.go#L59-L81), and then [renders it as json](https://github.com/tendermint/basecoin/blob/develop/cmd/basecoin/commands/query.go#L118-L124)

The data must be returned from the app as bytes that match the merkle proof, thus it is the responsibility of this library to parse it.  Since this is application-specific domain knowledge, we cannot program this, but rather allow the application designer to provide us this information in a special `ValueReader` interface, which knows how to read the application-specific values stored in the merkle tree, and convert them into a struct that can be json-encoded, or otherwise transformed for other binding.

If you just want to pass unparsed bytes as hex-data around, use `Bytes` from `go-data` to store the data, which serializes in hex by default and fulfills the `lightclient.Value` interface.

### Verifying Proofs

Beyond simply querying data from a blockchain, we often want **undeniable, cryptographic proof** of its validity.  This is the reason for exposing merkle proofs as first class objects in the new query [request](https://github.com/tendermint/abci/blob/develop/types/types.pb.go#L718-L723) and [response](https://github.com/tendermint/abci/blob/develop/types/types.pb.go#L1413-L1421).  However, this Proof byte slice, still generally requires go code to [parse and validate](https://github.com/tendermint/go-merkle/blob/develop/iavl_proof.go#L14-L42).

It is important to provide access to this essential functionality in a light client library, so we can provide this security to any UI we wish to build. Once we have a header we have properly certified by the above mechanisms, then we can accept merkle proofs for any data that leads to any of the root hashes in the block header (currently, this is meaningfully the apphash (for state), datahash (for txs in that block), and the validatorhash (used by the certifier)).

This is still a work in progress (will be enhanced soon), but you can view the initial code in the [proofs package](https://github.com/tendermint/light-client/tree/master/proofs).

## Extensibility

Of course, every application has its own transaction types, its own way of signing them, and its own data structures for which it wishes to present proofs. While we can work with interfaces that know how to serialize themselves (like `Signable`), this does not tell us how to *deserialize* the objects.  In fact, it is impossible to attach this information to an interface (as we allow many potentially unknown concrete types). We will allow each client to configure this app-specific parsing and display logic as plugins, which are registered in `package main`, thus allowing one to easily create light-client customized for their app, the way we can easily create specialized installs of basecoin with a few custom plugins configured.

## References

Some other projects that may inspire this:

* [Project Trillian](https://github.com/google/trillian) - Verifiable data structures (Apache 2.0)

