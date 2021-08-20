package transport

import (
	"errors"
	"time"

	"github.com/evcc-io/eebus/ship/message"
	"github.com/evcc-io/eebus/ship/ship"
)

// Hello is the common hello exchange
func (c *Transport) Hello() error {
	waitMs := uint(message.CmiTimeout / time.Millisecond)

	// SME_HELLO_STATE_READY_INIT
	err := c.WriteJSON(message.CmiTypeControl, ship.CmiConnectionHello{
		ConnectionHello: ship.ConnectionHello{
			Phase:   ship.ConnectionHelloPhaseTypeReady,
			Waiting: &waitMs,
		},
	})

	timer := time.NewTimer(message.CmiTimeout)
	for err == nil {
		// SME_HELLO_STATE_READY_LISTEN
		var msg interface{}
		msg, err = c.ReadMessage(timer.C)
		if err != nil {
			if errors.Is(err, ErrTimeout) {
				// SME_HELLO_STATE_READY_TIMEOUT
				_ = c.WriteJSON(message.CmiTypeControl, ship.CmiConnectionHello{
					ConnectionHello: ship.ConnectionHello{
						Phase: ship.ConnectionHelloPhaseTypeAborted,
					},
				})
			}

			return err
		}

		switch hello := msg.(type) {
		case ship.ConnectionHello:
			switch hello.Phase {
			case ship.ConnectionHelloPhaseTypeReady:
				// HELLO_OK
				return nil

			case ship.ConnectionHelloPhaseTypeAborted:
				err = errors.New("hello: aborted")

			case ship.ConnectionHelloPhaseTypePending:
				if hello.ProlongationRequest != nil && *hello.ProlongationRequest {
					timer = time.NewTimer(message.CmiHelloProlongationTimeout)
				}
			}

		case ship.ConnectionClose:
			err = errors.New("hello: remote closed")

		default:
			err = errors.New("hello: invalid type")
		}
	}

	return err
}
