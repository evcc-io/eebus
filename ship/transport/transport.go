package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/evcc-io/eebus/ship/message"
	"github.com/evcc-io/eebus/util"
	"github.com/gorilla/websocket"
)

// CmiReadWriteTimeout timeout
const CmiReadWriteTimeout = 10 * time.Second

// ErrTimeout is the timeout error
var ErrTimeout = errors.New("timeout")

// Transport is the physical transport layer
type Transport struct {
	conn   *websocket.Conn
	logger util.Logger

	recv    chan []byte
	recvErr chan error
	send    chan []byte
	sendErr chan error
	closeC  chan struct{}

	CloseHandler func()
	isClosed     bool
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// New creates SHIP transport on given websocket connection
func New(log util.Logger, conn *websocket.Conn) *Transport {
	t := &Transport{
		conn:    conn,
		logger:  log,
		send:    make(chan []byte, 1),
		recv:    make(chan []byte, 1),
		sendErr: make(chan error, 1),
		recvErr: make(chan error, 1),
		closeC:  make(chan struct{}),
	}

	go t.readPump()
	go t.writePump()

	return t
}

func (c *Transport) IsConnectionClosed() bool {
	return c.isClosed
}

func (c *Transport) handleConnectionClose() {
	c.isClosed = true
	if c.CloseHandler != nil {
		c.CloseHandler()
	}
}

func (c *Transport) log() util.Logger {
	if c.logger == nil {
		return &util.NopLogger{}
	}
	return c.logger
}

func (c *Transport) readPump() {
	defer func() {
		c.conn.Close()
		c.handleConnectionClose()
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-c.closeC:
			return

		default:
			typ, b, err := c.conn.ReadMessage()
			if err == nil {
				if len(b) > 2 {
					b = bytes.TrimSuffix(b, []byte{0x00}) // workaround
					c.log().Println("recv:", string(b[1:]))
				}

				if typ != websocket.BinaryMessage {
					err = fmt.Errorf("invalid message type: %d", typ)
				}
			}

			if err == nil {
				c.recv <- b
			} else {
				c.recvErr <- err
			}
		}
	}
}

// ReadBinary reads binary message
func (c *Transport) ReadBinary(timerC <-chan time.Time) ([]byte, error) {
	select {
	case <-timerC:
		c.handleConnectionClose()
		return nil, ErrTimeout

	case <-c.closeC:
		c.handleConnectionClose()
		return nil, net.ErrClosed

	case b := <-c.recv:
		return b, nil

	case err := <-c.recvErr:
		return nil, err
	}
}

// ReadMessage reads JSON message
func (c *Transport) ReadMessage(timerC <-chan time.Time) (interface{}, error) {
	select {
	case <-timerC:
		c.handleConnectionClose()
		return nil, ErrTimeout

	case <-c.closeC:
		c.handleConnectionClose()
		return nil, net.ErrClosed

	case b := <-c.recv:
		if len(b) < 2 {
			return nil, errors.New("invalid length")
		}
		if b[0] < 1 {
			return nil, errors.New("invalid phase")
		}

		return message.Decode(b[1:])

	case err := <-c.recvErr:
		return nil, err
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Transport) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
		c.handleConnectionClose()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if len(msg) > 2 {
				c.log().Println("send:", string(msg[1:]))
			}

			err := c.conn.WriteMessage(websocket.BinaryMessage, msg)
			c.sendErr <- err

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WriteBinary writes binary message
func (c *Transport) WriteBinary(msg []byte) error {
	if c.isClosed {
		return errors.New("cannot write to closed connection")
	}

	c.send <- msg

	timer := time.NewTimer(10 * time.Second)
	select {
	case <-timer.C:
		c.handleConnectionClose()
		return ErrTimeout

	case <-c.closeC:
		c.handleConnectionClose()
		return net.ErrClosed

	case err := <-c.sendErr:
		return err
	}
}

// WriteJSON writes JSON message
func (c *Transport) WriteJSON(typ byte, jsonMsg interface{}) error {
	msg, err := json.Marshal(jsonMsg)
	if err != nil {
		return err
	}

	// add header
	b := bytes.NewBuffer([]byte{typ})
	if _, err = b.Write(msg); err == nil {
		err = c.WriteBinary(b.Bytes())
	}

	return err
}
