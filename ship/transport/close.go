package transport

import (
	"errors"
	"time"

	"github.com/evcc-io/eebus/ship/message"
	"github.com/evcc-io/eebus/ship/ship"
)

// AcceptClose accepts connection close
func (c *Transport) AcceptClose() error {
	err := c.WriteJSON(message.CmiTypeEnd, ship.CmiConnectionClose{
		ConnectionClose: ship.ConnectionClose{
			Phase: ship.ConnectionClosePhaseTypeConfirm,
		},
	})

	// stop read/write pump
	if !c.isChannelClosed() {
		close(c.closeC)
	}
	c.conn.Close()
	c.handleConnectionClose()
	if c.CloseHandler != nil {
		c.CloseHandler()
	}

	return err
}

// Close closes the connection
func (c *Transport) Close() error {
	err := c.WriteJSON(message.CmiTypeEnd, ship.CmiConnectionClose{
		ConnectionClose: ship.ConnectionClose{
			Phase: ship.ConnectionClosePhaseTypeAnnounce,
			// MaxTime: int(ship.CmiCloseTimeout / time.Millisecond),
		},
	})

	timer := time.NewTimer(message.CmiCloseTimeout)
	for err == nil {
		var msg interface{}
		msg, err = c.ReadMessage(timer.C)
		if err != nil {
			break
		}

		if typed, ok := msg.(ship.ConnectionClose); ok && typed.Phase == ship.ConnectionClosePhaseTypeConfirm {
			break
		}

		err = errors.New("close: invalid response")
	}

	// stop read/write pump
	if !c.isChannelClosed() {
		close(c.closeC)
	}
	c.conn.Close()
	c.handleConnectionClose()
	if c.CloseHandler != nil {
		c.CloseHandler()
	}

	return err
}

func (c *Transport) isChannelClosed() bool {
	select {
	case <-c.closeC:
		return false
	default:
		return true
	}
}
