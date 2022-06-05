package mdns

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/evcc-io/eebus/ship"
	"github.com/evcc-io/eebus/util"
	"github.com/gorilla/websocket"
	"github.com/grandcat/zeroconf"
	"github.com/mitchellh/mapstructure"
)

// ServiceDescription contains the ship service parameters
type ServiceDescription struct {
	Model, Brand string
	SKI          string
	Register     bool
	Path         string
	ID           string
}

// Service is the ship service
type Service struct {
	ServiceDescription
	URIs []string
	Conn *ship.Connector
}

// NewFromDNSEntry creates ship service from its DNS definition
func NewFromDNSEntry(zc *zeroconf.ServiceEntry) (*Service, error) {
	ss := Service{}

	txtM := make(map[string]interface{})
	for _, txtE := range zc.Text {
		split := strings.SplitN(txtE, "=", 2)
		if len(split) == 2 {
			txtM[split[0]] = split[1]
		}
	}

	decoderConfig := &mapstructure.DecoderConfig{
		Result:           &ss.ServiceDescription,
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err == nil {
		err = decoder.Decode(txtM)
		if err != nil {
			return &ss, err
		}
	}

	ss.URIs, err = URIsFromDNS(zc, ss.ServiceDescription.Path)

	return &ss, err
}

// URIsFromDNS returns the service URI and appends the path
func URIsFromDNS(zc *zeroconf.ServiceEntry, path string) ([]string, error) {
	var uris []string

	if len(zc.HostName) > 0 {
		uris = append(uris, createURI(zc.HostName, zc.Port)+path)
	}

	for _, address := range zc.AddrIPv4 {
		uris = append(uris, createURI(address.String(), zc.Port)+path)
	}

	for _, address := range zc.AddrIPv6 {
		uris = append(uris, createURI(address.String(), zc.Port)+path)
	}

	if len(uris) == 0 {
		return uris, errors.New("mDNS record doesn't contain an IP address")
	}

	return uris, nil
}

func createURI(host string, port int) string {
	return ship.Scheme + net.JoinHostPort(host, fmt.Sprintf("%d", port))
}

// WebsocketConnector is the connector used for establishing new websocket connections
var WebsocketConnector func(uri string) (*websocket.Conn, error)

// Connect connects to the service endpoint and performs handshake
func (ss *Service) Connect(log util.Logger, accessMethod string, cert tls.Certificate, closeHandler func(string)) (ship.Conn, error) {
	WebsocketConnector = ship.TLSConnection(cert)

	for _, uri := range ss.URIs {
		ws, err := WebsocketConnector(uri)
		if err != nil {
			log.Printf("Failed to connect to %s: %s\n", uri, err)
			continue
		}

		sc := &ship.Connector{
			Log:          log,
			Local:        ship.Service{Pin: "", Methods: accessMethod},
			Remote:       ship.Service{Pin: ""},
			CloseHandler: closeHandler,
			SKI:          ss.ServiceDescription.SKI,
		}

		conn, err := sc.Connect(ws)
		if err == nil {
			return conn, nil
		}
	}

	return nil, errors.New("cannot connect to any service endpoint")
}
