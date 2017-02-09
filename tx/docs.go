/*
package tx contains generic Signable implementations that can be used
by your application to handle authentication needs.

It currently supports transaction data as opaque bytes and either single
or multiple private key signatures using straightforward algorithms.
It currently does not support N-of-M key share signing of other more
complex algorithms (although it would be great to add them)

**TODO**
It also contains a SignableReader for deserialization of these two types.

Maybe this package should be moved out of lightclient, more to a repo
designed for application support, as that is the main usecase?
*/
package tx
