package ship

import (
	"errors"
	"fmt"
)

// UniqueID creates ship ID with given prefix and salted by provided protectedID
// The protectedID is basically the same as `machineid.ProtectedID(appId)` above
// but provided by the system using this library
func UniqueIDWithProtectedID(prefix, protectedID string) (string, error) {
	if len(protectedID) == 0 {
		return "", errors.New("A generated machine specific protectedID needs to be provided")
	}

	return fmt.Sprintf("%s-%0x", prefix, protectedID[:8]), nil
}
