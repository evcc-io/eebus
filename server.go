package eebus

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/evcc-io/eebus/cert"
	"github.com/evcc-io/eebus/ship"
	"github.com/evcc-io/eebus/util"
	"github.com/grandcat/zeroconf"
)

type Server struct {
	Log                    util.Logger
	Addr, Path             string
	ID, Brand, Model, Type string
	Register               bool
	Certificate            tls.Certificate
}

func (c *Server) Announce() (*zeroconf.Server, error) {
	ski, err := cert.SkiFromCert(c.Certificate)
	if err != nil {
		return nil, err
	}

	path := c.Path
	if path == "" {
		path = "/"
	}

	_, port, err := net.SplitHostPort(c.Addr)
	if err != nil {
		return nil, err
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	if c.Log != nil {
		c.Log.Printf("mDNS: announcing id: %s ski: %s", c.ID, ski)
	}

	server, err := zeroconf.Register(c.Model, ship.ZeroconfType, ship.ZeroconfDomain, portInt, []string{
		"txtvers=1",
		"path=" + path,
		"id=" + c.ID,
		"ski=" + ski,
		"brand=" + c.Brand,
		"model=" + c.Model,
		"type=" + c.Type,
		"register=" + fmt.Sprintf("%v", c.Register),
	}, nil)

	if err != nil {
		err = fmt.Errorf("mDNS: failed registering service: %w", err)
	}

	return server, err
}

func (c *Server) createVerifier(verifier func(*x509.Certificate) error) func(state tls.ConnectionState) error {
	return func(state tls.ConnectionState) error {
		if len(state.PeerCertificates) == 0 {
			return errors.New("missing client certificate")
		}

		cert := state.PeerCertificates[0]

		if len(cert.SubjectKeyId) == 0 {
			return errors.New("missing client ski")
		}

		return verifier(state.PeerCertificates[0])
	}
}

func (c *Server) Listen(handler http.Handler, verifier func(*x509.Certificate) error) error {
	s := &http.Server{
		Addr:    c.Addr,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{c.Certificate},
			ClientAuth:   tls.RequireAnyClientCert,
			CipherSuites: ship.CipherSuites,
		},
	}

	if verifier != nil {
		s.TLSConfig.VerifyConnection = c.createVerifier(verifier)
	}

	return s.ListenAndServeTLS("", "")
}
