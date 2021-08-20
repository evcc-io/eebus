package transport

import (
	"errors"
	"time"

	"github.com/evcc-io/eebus/ship/ship"
)

// DataReceive receives handshake
func (c *Transport) DataReceive() error {
	timer := time.NewTimer(CmiReadWriteTimeout)
	msg, err := c.ReadMessage(timer.C)
	if err != nil {
		return err
	}

	switch typed := msg.(type) {
	case ship.Data:
		_ = typed
		return nil

	case ship.ConnectionClose:
		err = errors.New("data: remote closed")

	default:
		err = errors.New("data: invalid type")
	}

	return err
}
