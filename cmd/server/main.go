package main

import (
	"crypto/tls"
	"crypto/x509/pkix"
	"log"
	"os/signal"

	"os"

	"github.com/evcc-io/eebus/app"
	certhelper "github.com/evcc-io/eebus/cert"
	"github.com/evcc-io/eebus/communication"
	"github.com/evcc-io/eebus/server"
	"github.com/evcc-io/eebus/ship"
	"github.com/evcc-io/eebus/spine/model"
	"github.com/evcc-io/eebus/util"
)

const (
	certFile = "evcc.crt"
	keyFile  = "evcc.key"
)

func certificate(details communication.ManufacturerDetails) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		if os.IsNotExist(err) {
			subject := pkix.Name{
				CommonName:   details.DeviceCode,
				Country:      []string{"DE"},
				Organization: []string{details.BrandName},
			}

			if cert, err = certhelper.CreateCertificate(true, subject); err == nil {
				err = certhelper.SaveX509KeyPair(certFile, keyFile, cert)
			}
		}

		if err != nil {
			panic(err)
		}
	}

	return cert
}

func main() {
	details := communication.ManufacturerDetails{
		BrandName:     "EVCC",
		DeviceName:    "EVCC",
		DeviceCode:    "EVCC_HEMS_01",
		DeviceAddress: "EVCC_HEMS",
	}

	cert := certificate(details)

	id := server.UniqueID{Prefix: details.BrandName}.String()
	log := log.New(&util.LogWriter{os.Stdout, "2006/01/02 15:04:05 "}, "[server] ", 0)

	srv := &server.Server{
		Log:         log,
		Addr:        ":4712",
		Path:        "/ship/",
		Certificate: cert,
		ID:          id,
		Brand:       details.BrandName,
		Model:       details.DeviceCode,
		Type:        string(model.DeviceTypeEnumTypeEnergyManagementSystem),
		Register:    true,
	}

	zc, err := srv.Announce()
	if err != nil {
		panic(err)
	}
	defer zc.Shutdown()

	hems := app.HEMS(details)

	ln := &server.Listener{
		Log:          log,
		AccessMethod: id,
		Handler: func(ski string, conn ship.Conn) error {
			ctrl := communication.NewConnectionController(log, conn, hems)
			return ctrl.Boot()
		},
	}

	go srv.Listen(ln, nil)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			log.Println("mDNS: shutdown")
			zc.Shutdown()
			os.Exit(0)
		}
	}()

	select {}
}
