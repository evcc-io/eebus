package transport

import (
	"errors"
	"time"

	"github.com/evcc-io/eebus/ship/message"
	"github.com/evcc-io/eebus/ship/ship"
)

// read pin requirements
func (c *Transport) readPinState() (ship.ConnectionPinState, error) {
	timer := time.NewTimer(CmiReadWriteTimeout)
	msg, err := c.ReadMessage(timer.C)

	switch typed := msg.(type) {
	case ship.ConnectionPinState:
		return typed, err

	default:
		if err == nil {
			err = errors.New("pin: invalid type")
		}

		return ship.ConnectionPinState{}, err
	}
}

const (
	pinReceived = 1 << iota
	pinSent

	pinCompleted = pinReceived | pinSent
)

// PinState handles pin exchange
func (c *Transport) PinState(local, remote ship.PinValueType) error {
	pinState := ship.ConnectionPinState{
		PinState: ship.PinStateTypeNone,
	}

	var status int
	if local != "" {
		ok := ship.PinInputPermissionTypeOk
		pinState.PinState = ship.PinStateTypeRequired
		pinState.InputPermission = &ok
	} else {
		// always received if not necessary
		status |= pinReceived
	}

	err := c.WriteJSON(message.CmiTypeControl, ship.CmiConnectionPinState{
		ConnectionPinState: pinState,
	})

	timer := time.NewTimer(10 * time.Second)
	for err == nil && status != pinCompleted {
		var msg interface{}
		msg, err = c.ReadMessage(timer.C)
		if err != nil {
			break
		}

		switch typed := msg.(type) {
		// local pin
		case ship.ConnectionPinInput:
			// signal error to client
			if typed.Pin != local {
				err = c.WriteJSON(message.CmiTypeControl, ship.CmiConnectionPinError{
					ConnectionPinError: ship.ConnectionPinError{
						Error: "1", // TODO
					},
				})
			}

			status |= pinReceived

		// remote pin
		case ship.ConnectionPinState:
			if typed.PinState == ship.PinStateTypeOptional || typed.PinState == ship.PinStateTypeRequired {
				if remote != "" {
					err = c.WriteJSON(message.CmiTypeControl, ship.CmiConnectionPinInput{
						ConnectionPinInput: ship.ConnectionPinInput{
							Pin: remote,
						},
					})
				} else {
					err = errors.New("pin: remote pin required")
				}
			}

			status |= pinSent

		case ship.ConnectionPinError:
			err = errors.New("pin: remote pin mismatched")

		case ship.ConnectionClose:
			err = errors.New("pin: remote closed")

		default:
			err = errors.New("pin: invalid type")
		}
	}

	return err
}
