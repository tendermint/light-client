/*
Package certifiers allows you to validate headers
without a full node.
The purpose here is to provide all security algorithms in
order to enable the quick and efficient construction of
light clients.

Commits

There are two main data structures that we pass around - Commit
and FullCommit. Both of them mirror what information is
exposed in tendermint rpc.

Commit is a block header along with enough validator signatures
to prove its validity (> 2/3 of the voting power). A FullCommit
is a Commit along with the full validator set. When the
validator set doesn't change, the Commit is enough, but since
the block header only has a hash, we need the FullCommit to
follow any changes to the validator set.

Certifiers

A Certifier validates a new Commit given the currently known
state. There are three different types of Certifiers exposed,
each one building on the last one, with additional complexity.

Static - given the validator set upon initialization. Verifies
all signatures against that set and if the validator set
changes, it will reject all headers.

Dynamic - This wraps Static and has the same Certify
method. However, it adds an Update method, which can be called
with a FullCommit when the validator set changes. If it can
prove this is a valid transition, it will update the validator
set.

Inquiring - this wraps Dynamic and implements an auto-update
strategy on top of the Dynamic update. If a call to
Certify fails as the validator set has changed, then it
attempts to find a FullCommit and Update to that header.
To get these FullCommits, it makes use of a Provider.

Providers

A Provider allows us to store and retrieve the FullCommits,
to provide memory to the Inquiring Certifier.

NewMemStoreProvider - in-memory cache.

files.NewProvider - disk backed storage.

client.NewHTTPProvider - query tendermint rpc.

NewCacheProvider - combine multiple providers.

The suggested use for local light clients is
client.NewHTTPProvider for getting new data (Source),
and NewCacheProvider(NewMemStoreProvider(),
files.NewProvider()) to store confirmed headers (Trusted)
*/
package certifiers
