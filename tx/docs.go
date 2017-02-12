/*
package tx contains generic Signable implementations that can be used
by your application to handle authentication needs.

It currently supports transaction data as opaque bytes and either single
or multiple private key signatures using straightforward algorithms.
It currently does not support N-of-M key share signing of other more
complex algorithms (although it would be great to add them)

ReadSignableBinary() can be used by an ACBi app to deserialize a
signed, packed OneSig or MultiSig object.  You must write your
own SignableReader to translate json into a binary data blob and
wrap it with OneSig or MultiSig to provide an external interface
for the api proxy or language bindings.

Maybe this package should be moved out of lightclient, more to a repo
designed for application support, as that is the main usecase?
*/
package tx
