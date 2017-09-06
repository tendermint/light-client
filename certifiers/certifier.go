package certifiers

import (
	lc "github.com/tendermint/light-client"
)

// Certifier is the interface all certifiers implement.
type Certifier interface {
	Certify(check lc.Checkpoint) error
}
