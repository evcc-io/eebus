package ship

import (
	"fmt"

	"github.com/panta/machineid"
)

// UniqueID creates ship ID with given prefix and salted by appID
func UniqueID(prefix, appID string) (string, error) {
	s, err := machineid.ProtectedID(appID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%0x", prefix, s[:8]), nil
}
