package ship

import (
	"bytes"
	"fmt"
	"time"

	"github.com/evcc-io/eebus/ship/message"
	"github.com/evcc-io/eebus/ship/ship"
	"github.com/evcc-io/eebus/ship/transport"
	"github.com/evcc-io/eebus/util"
	"github.com/gorilla/websocket"
)

// Connector is the ship client connector
type Connector struct {
	Log          util.Logger
	Local        Service
	Remote       Service
	CloseHandler func(string)
	SKI          string

	// mux    sync.Mutex
	closedHandlerInvoked bool
}

// init creates the connection
func (c *Connector) init(t *transport.Transport) error {
	init := []byte{message.CmiTypeInit, 0x00}

	// CMI_STATE_CLIENT_SEND
	if err := t.WriteBinary(init); err != nil {
		return err
	}

	timer := time.NewTimer(message.CmiTimeout)

	// CMI_STATE_CLIENT_WAIT
	msg, err := t.ReadBinary(timer.C)
	if err != nil {
		return err
	}

	// CMI_STATE_CLIENT_EVALUATE
	if !bytes.Equal(init, msg) {
		return fmt.Errorf("init: invalid response")
	}

	return nil
}

func (c *Connector) protocolHandshake(t *transport.Transport) error {
	hs := ship.CmiMessageProtocolHandshake{
		MessageProtocolHandshake: ship.MessageProtocolHandshake{
			HandshakeType: ship.ProtocolHandshakeTypeTypeAnnouncemax,
			Version:       ship.Version{Major: 1, Minor: 0},
			Formats: ship.MessageProtocolFormatsType{
				Format: []ship.MessageProtocolFormatType{ship.ProtocolHandshakeFormatJSON},
			},
		},
	}
	if err := t.WriteJSON(message.CmiTypeControl, hs); err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	// receive server selection and send selection back to server
	err := t.HandshakeReceiveSelect()
	if err == nil {
		hs.MessageProtocolHandshake.HandshakeType = ship.ProtocolHandshakeTypeTypeSelect
		err = t.WriteJSON(message.CmiTypeControl, hs)
	}

	return err
}

// // Close performs ordered close of client connection
// func (c *Connector) Close(t *transport.Transport) error {
// c.mux.Lock()
// defer c.mux.Unlock()

// 	if c.closed {
// 		return os.ErrClosed
// 	}

// 	c.closed = true

// 	// stop readPump
// 	// defer close(c.closeC)

// 	return t.Close()
// }

// Connect performs the client connection handshake
func (c *Connector) Connect(conn *websocket.Conn) (Conn, error) {
	t := transport.New(c.Log, conn)
	t.CloseHandler = c.TransportClosed

	if err := c.init(t); err != nil {
		return nil, err
	}

	err := t.Hello()
	if err == nil {
		err = c.protocolHandshake(t)
	}

	if err == nil {
		err = t.PinState(
			ship.PinValueType(c.Local.Pin),
			ship.PinValueType(c.Remote.Pin),
		)
	}

	if err == nil {
		c.Remote.Methods, err = t.AccessMethodsRequest(c.Local.Methods)
	}

	// close connection if handshake or hello fails
	if err != nil {
		_ = t.Close()
		c.TransportClosed()
	}

	shipConn := &connection{t: t}

	return shipConn, err
}

// TransportClosed handles a closed transport conncection
func (c *Connector) TransportClosed() {
	if c.CloseHandler != nil && !c.closedHandlerInvoked {
		// make sure the close handler is only invoked once for this connection
		c.closedHandlerInvoked = true
		c.CloseHandler(c.SKI)
	}
}

// Connect performs the client connection handshake
func Connect(conn *websocket.Conn) (Conn, error) {
	c := &Connector{Log: &util.NopLogger{}}
	return c.Connect(conn)
}
