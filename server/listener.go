package server

import (
	"errors"
	"net/http"

	"github.com/evcc-io/eebus/cert"
	"github.com/evcc-io/eebus/ship"
	"github.com/evcc-io/eebus/util"
	"github.com/gorilla/websocket"
)

type Listener struct {
	Log          util.Logger
	Handler      func(ski string, conn ship.Conn) error
	AccessMethod string
}

func (s *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Log == nil {
		s.Log = &util.NopLogger{}
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
		Subprotocols:    []string{ship.SubProtocol},
	}

	// upgrade
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Log.Println(err)
		return
	}

	// return and close connection
	if ws.Subprotocol() != ship.SubProtocol {
		s.Log.Println("protocol mismatch:", ws.Subprotocol())
		return
	}

	// ship
	shipSrv := &ship.Server{
		Log:    s.Log,
		Local:  ship.Service{Pin: "", Methods: s.AccessMethod},
		Remote: ship.Service{Pin: ""},
	}

	conn, err := shipSrv.Serve(ws)
	if err != nil {
		s.Log.Println(err)
		return
	}

	if s.Handler == nil {
		_ = conn.Close()
		err = errors.New("no handler")
	}

	if err == nil {
		var ski string
		if len(r.TLS.PeerCertificates) > 0 {
			ski, err = cert.SkiFromX509(r.TLS.PeerCertificates[0])
		}

		if err == nil {
			err = s.Handler(ski, conn)
		}
	}

	s.Log.Println("done:", err)
}
