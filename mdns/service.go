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
	URI  string
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
	}

	baseURI, err := baseURIFromDNS(zc)
	if err == nil {
		ss.URI = baseURI + ss.ServiceDescription.Path
	}

	return &ss, err
}

// baseURIFromDNS returns the service URI
func baseURIFromDNS(zc *zeroconf.ServiceEntry) (string, error) {
	var address net.IP
	var hostname string
	var uri string

	if len(zc.HostName) > 0 {
		hostname = zc.HostName
	} else if len(zc.AddrIPv4) > 0 {
		address = zc.AddrIPv4[0]
	} else if len(zc.AddrIPv6) > 0 {
		address = zc.AddrIPv6[0]
	} else {
		return uri, errors.New("mDNS record doesn't contain an IP address")
	}
	if len(hostname) > 0 {
		uri = ship.Scheme + net.JoinHostPort(hostname, fmt.Sprintf("%d", zc.Port))
	} else {
		uri = ship.Scheme + net.JoinHostPort(address.String(), fmt.Sprintf("%d", zc.Port))
	}

	return uri, nil
}

// WebsocketConnector is the connector used for establishing new websocket connections
var WebsocketConnector func(uri string) (*websocket.Conn, error)

// Connect connects to the service endpoint and performs handshake
func (ss *Service) Connect(log util.Logger, accessMethod string, cert tls.Certificate, closeHandler func(string)) (ship.Conn, error) {
	WebsocketConnector = ship.TLSConnection(cert)

	ws, err := WebsocketConnector(ss.URI)
	if err != nil {
		return nil, err
	}

	sc := &ship.Connector{
		Log:          log,
		Local:        ship.Service{Pin: "", Methods: accessMethod},
		Remote:       ship.Service{Pin: ""},
		CloseHandler: closeHandler,
		SKI:          ss.ServiceDescription.SKI,
	}

	return sc.Connect(ws)
}
