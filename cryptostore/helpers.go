package cryptostore

// import (
// 	crypto "github.com/tendermint/go-crypto"
// 	wire "github.com/tendermint/go-wire"
// )

// // TODO: is this necessary????
// func init() {
// 	registerGoCryptoGoWire()
// }

// type pubwire struct{ crypto.PubKey }
// type sigwire struct{ crypto.Signature }
// type message struct {
// 	Data      []byte
// 	Signature crypto.Signature
// 	PubKey    crypto.PubKey
// }

// // we must register these types here, to make sure they parse (maybe go-wire issue??)
// // TODO: fix go-wire, remove this code
// func registerGoCryptoGoWire() {
// 	wire.RegisterInterface(
// 		pubwire{},
// 		wire.ConcreteType{O: crypto.PubKeyEd25519{}, Byte: crypto.PubKeyTypeEd25519},
// 		wire.ConcreteType{O: crypto.PubKeySecp256k1{}, Byte: crypto.PubKeyTypeSecp256k1},
// 	)
// 	wire.RegisterInterface(
// 		sigwire{},
// 		wire.ConcreteType{O: crypto.SignatureEd25519{}, Byte: crypto.SignatureTypeEd25519},
// 		wire.ConcreteType{O: crypto.SignatureSecp256k1{}, Byte: crypto.SignatureTypeSecp256k1},
// 	)
// }
