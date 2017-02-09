/*
package lightclient is a complete solution for integrating a light client with
tendermint.  It provides all common functionality that a client needs to create
and sign transactions, get and verify state, and synchronize with a tendermint node.
It is intended to expose this data both through golang interfaces, a local RPC server,
and language bindings.  You can find more info on the aims of this package in the
Readme: https://github.com/tendermint/light-client/blob/master/README.md

The package layout attempts to expose common domain types in the
top-level with no other dependencies.  Main packages should select which
dependencies they wish to have and wire them together with common glue code
that only depends on the interface.
More info here: https://medium.com/%40benbjohnson/standard-package-layout-7cdbc8391fc1

The majority of the definitions here are interfaces, to be implemented in subpackages,
or data structures encapsulating tendermint return values.  All tendermint data
structures are prefixed by Tm to make the documentation clearer. The other data
structure defined on the top level is KeyInfo/KeyInfos to represent infomation
on the public keys stored in the KeyManager, along with logic to sort themselves.
*/
package lightclient
