package transport

import (
	"errors"
	"time"

	"github.com/evcc-io/eebus/ship/message"
	"github.com/evcc-io/eebus/ship/ship"
)

// HandshakeReceiveSelect receives handshake
func (c *Transport) HandshakeReceiveSelect() error {
	timer := time.NewTimer(CmiReadWriteTimeout)
	msg, err := c.ReadMessage(timer.C)
	if err != nil {
		return err
	}

	switch typed := msg.(type) {
	case ship.MessageProtocolHandshake:
		if typed.HandshakeType != ship.ProtocolHandshakeTypeTypeSelect || !typed.Formats.IsSupported(ship.ProtocolHandshakeFormatJSON) {
			_ = c.WriteJSON(message.CmiTypeControl, ship.CmiMessageProtocolHandshakeError{
				MessageProtocolHandshakeError: ship.MessageProtocolHandshakeError{
					Error: "2", // TODO
				}})

			err = errors.New("handshake: invalid format")
		}

		return err

	case ship.ConnectionClose:
		err = errors.New("handshake: remote closed")

	default:
		err = errors.New("handshake: invalid type")
	}

	return err
}
