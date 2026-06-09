// Package id mints opaque identifiers for new domain objects — needed
// wherever a use case creates a Monitor or Incident before persistence
// assigns one.
package id

import (
	"crypto/rand"
	"encoding/hex"
)

type Generator func() string

func Random() Generator {
	return func() string {
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		return hex.EncodeToString(b)
	}
}
