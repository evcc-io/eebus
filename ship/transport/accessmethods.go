package transport

import (
	"errors"
	"time"

	"github.com/evcc-io/eebus/ship/message"
	"github.com/evcc-io/eebus/ship/ship"
)

// AccessMethodsRequest sends access methods request and processes answer
func (c *Transport) AccessMethodsRequest(methods string) (string, error) {
	err := c.WriteJSON(message.CmiTypeControl, ship.CmiAccessMethodsRequest{
		AccessMethodsRequest: ship.AccessMethodsRequest{},
	})

	for err == nil {
		timer := time.NewTimer(CmiReadWriteTimeout)

		var msg interface{}
		msg, err = c.ReadMessage(timer.C)
		if err != nil {
			break
		}

		switch typed := msg.(type) {
		case ship.AccessMethods:
			// access methods received
			return typed.Id, nil

		case ship.AccessMethodsRequest:
			err = c.WriteJSON(message.CmiTypeControl, ship.CmiAccessMethods{
				AccessMethods: ship.AccessMethods{Id: methods},
			})

		default:
			err = errors.New("access methods: invalid type")
		}
	}

	return "", err
}
