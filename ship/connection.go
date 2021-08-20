package ship

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/evcc-io/eebus/ship/message"
	"github.com/evcc-io/eebus/ship/ship"
	"github.com/evcc-io/eebus/ship/transport"
	"github.com/gorilla/websocket"
)

// TLSConnection creates an encrypted websocket connection
func TLSConnection(cert tls.Certificate) func(uri string) (*websocket.Conn, error) {
	return func(uri string) (*websocket.Conn, error) {
		dialer := &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
				CipherSuites:       CipherSuites,
			},
			Subprotocols: []string{SubProtocol},
		}

		conn, _, err := dialer.Dial(uri, nil)

		return conn, err
	}
}

var ErrInvalidMessageType = errors.New("invalid message type")

type Conn interface {
	Read() (json.RawMessage, error)
	Write(json.RawMessage) error
	Close() error
	IsConnectionClosed() bool
}

var _ Conn = (*connection)(nil)

type connection struct {
	t *transport.Transport
}

func (c *connection) IsConnectionClosed() bool {
	return c.t.IsConnectionClosed()
}

func (c *connection) Read() (json.RawMessage, error) {
	timer := time.NewTimer(600 * time.Second)

	msg, err := c.t.ReadMessage(timer.C)
	if err != nil {
		return nil, err
	}

	switch typed := msg.(type) {
	case ship.Data:
		return typed.Payload, nil

	case ship.ConnectionClose:
		return nil, c.t.AcceptClose()

	default:
		err = ErrInvalidMessageType
	}

	return nil, err
}

func (c *connection) Write(payload json.RawMessage) error {
	hs := ship.CmiData{
		Data: ship.Data{
			Header: ship.HeaderType{
				ProtocolId: ship.ProtocolIdType(message.ProtocolID),
			},
			Payload: payload,
		},
	}

	return c.t.WriteJSON(message.CmiTypeData, &hs)
}

func (c *connection) Close() error {
	return c.t.Close()
}
