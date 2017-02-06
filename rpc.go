package lightclient

// Node represents someway to query a tendermint node for info
// Typically via RPC, but could be mocked or connected locally
// TODO: trim this down and distinguish from RPC a bit! (custom types)
type Node interface {
	// Broadcast sends into to the chain
	// We only implement BroadcastCommit for now, add others???
	// Or just register a callback?
	// The return result cannot be fully trusted without downloading signed headers
	Broadcast(tx []byte) (BroadcastResult, error)

	// Status and Validators give some info, nothing to be trusted though...
	// Unless we find that eg. the ValidatorResult matches the ValidatorHash
	// in a properly signed block header
	Status() (StatusResult, error)
	Validators() (ValidatorResult, error)

	// Query gets data from the Blockchain state, and can optionally provide
	// a Proof we can validate
	Query(path string, data []byte, prove bool) (QueryResult, error)

	// You need to check the Headers and Votes together to prove anything
	// is actually on the chain
	Headers(minHeight, maxHeight int) ([]BlockMeta, error)
	Votes(height int) (Votes, error)

	// TODO: let's make this reactive if possible
	// TODO: listen for a transaction being committed?
	// TODO: listen for a new block?
	// TODO: listen for change to a given key in the merkle store?

	// 	NetInfo() (*ctypes.ResultNetInfo, error)
	// 	DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error)
	// 	Genesis() (*ctypes.ResultGenesis, error)
	//  Block(height int) (*ctypes.ResultBlock, error)

	// 	ABCIInfo() (*ctypes.ResultABCIInfo, error)
}
