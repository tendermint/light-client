/*
package rpc provides higher-level functionality to interact with tendermint's
rpc server.

The basic client implementation is HTTPClient, which provides full, direct
access to the rpc server, providing type-safety and marshaling/unmarshaling,
but no additional functionality.

The more advanced client implementation is Node, which provides a few high-level
actions based upon HTTPClient and parsing, processing, and validating the
return values.  Node currently implements Broadcaster, Checker, and Searcher
interfaces from lightclient, which are the general high-level actions
one wants to perform on tendermint.

Higher-level functionality should be built upon Node, extending Node as needed,
or defining another type, but not directly on HTTPClient.  This package provides
the bridge between RPC calls and higher-level functionality.
*/
package rpc
