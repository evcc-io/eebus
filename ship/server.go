package ship

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/evcc-io/eebus/ship/message"
	"github.com/evcc-io/eebus/ship/ship"
	"github.com/evcc-io/eebus/ship/transport"
	"github.com/evcc-io/eebus/util"
	"github.com/gorilla/websocket"
)

// Server is the SHIP server
type Server struct {
	Log    util.Logger
	Local  Service
	Remote Service
}

// Init creates the connection
func (c *Server) init(t *transport.Transport) error {
	timer := time.NewTimer(message.CmiTimeout)

	// CMI_STATE_SERVER_WAIT
	msg, err := t.ReadBinary(timer.C)
	if err != nil {
		return err
	}

	// CMI_STATE_SERVER_EVALUATE
	init := []byte{message.CmiTypeInit, 0x00}
	if !bytes.Equal(init, msg) {
		return fmt.Errorf("init: invalid response")
	}

	return t.WriteBinary(init)
}

func (c *Server) protocolHandshake(t *transport.Transport) error {
	timer := time.NewTimer(transport.CmiReadWriteTimeout)
	msg, err := t.ReadMessage(timer.C)
	if err != nil {
		if errors.Is(err, transport.ErrTimeout) {
			_ = t.WriteJSON(message.CmiTypeControl, ship.CmiMessageProtocolHandshakeError{
				MessageProtocolHandshakeError: ship.MessageProtocolHandshakeError{
					Error: "2", // TODO
				}})
		}

		return err
	}

	switch typed := msg.(type) {
	case ship.MessageProtocolHandshake:
		if typed.HandshakeType != ship.ProtocolHandshakeTypeTypeAnnouncemax || !typed.Formats.IsSupported(ship.ProtocolHandshakeFormatJSON) {
			msg := ship.CmiMessageProtocolHandshakeError{
				MessageProtocolHandshakeError: ship.MessageProtocolHandshakeError{
					Error: "2", // TODO
				},
			}

			_ = t.WriteJSON(message.CmiTypeControl, msg)
			err = errors.New("handshake: invalid response")
			break
		}

		// send selection to client
		typed.HandshakeType = ship.ProtocolHandshakeTypeTypeSelect
		err = t.WriteJSON(message.CmiTypeControl, ship.CmiMessageProtocolHandshake{
			MessageProtocolHandshake: typed,
		})

	default:
		return fmt.Errorf("handshake: invalid type")
	}

	// receive selection back from client
	if err == nil {
		err = t.HandshakeReceiveSelect()
	}

	return err
}

// Close performs ordered close of server connection
// func (c *Server) Close() error {
// 	return t.Close()
// }

// Serve performs the server connection handshake
func (c *Server) Serve(conn *websocket.Conn) (Conn, error) {
	t := transport.New(c.Log, conn)

	if err := c.init(t); err != nil {
		return nil, err
	}

	// CMI_STATE_DATA_PREPARATION
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
	}

	shipConn := &connection{t: t}

	return shipConn, err
}
