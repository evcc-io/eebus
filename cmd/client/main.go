package main

import (
	"context"
	"crypto/tls"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"time"

	"os"

	"github.com/evcc-io/eebus/app"
	certhelper "github.com/evcc-io/eebus/cert"
	"github.com/evcc-io/eebus/communication"
	"github.com/evcc-io/eebus/mdns"
	"github.com/evcc-io/eebus/server"
	"github.com/evcc-io/eebus/ship"
	"github.com/evcc-io/eebus/spine/model"
	"github.com/libp2p/zeroconf/v2"
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

func connectService(entry *zeroconf.ServiceEntry, id string, details communication.ManufacturerDetails, cert tls.Certificate) {
	svc, err := mdns.NewFromDNSEntry(entry)

	var conn ship.Conn
	if err == nil {
		log.Printf("%s: client connect", entry.HostName)
		conn, err = svc.Connect(log.Default(), id, cert, nil)
	}

	if err != nil {
		log.Printf("%s: client done: %v", entry.HostName, err)
		return
	}

	hems := app.HEMS(details)
	ctrl := communication.NewConnectionController(log.Default(), conn, hems)

	err = ctrl.Boot()
	if err != nil {
		log.Printf("%s: connection startup failed: ", err)
		return
	}
}

func discoverDNS(results <-chan *zeroconf.ServiceEntry, connector func(*zeroconf.ServiceEntry)) {
	for entry := range results {
		log.Println("mDNS:", entry.HostName, entry.AddrIPv4, entry.Text)

		for _, typ := range entry.Text {
			if strings.HasPrefix(typ, "type=") && typ == "type=EVSE" {
				connector(entry)
			}
		}
	}
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

	srv := &server.Server{
		Log:         log.Default(),
		Addr:        ":4712",
		Path:        "/ship/",
		Certificate: cert,
		ID:          id,
		Brand:       details.BrandName,
		Model:       details.DeviceCode,
		Type:        string(model.DeviceTypeEnumTypeEnergyManagementSystem),
		Register:    true,
	}

	// use announcements even in client example so remote party can see us
	zc, err := srv.Announce()
	if err != nil {
		panic(err)
	}
	defer zc.Shutdown()

	entries := make(chan *zeroconf.ServiceEntry)
	go discoverDNS(entries, func(entry *zeroconf.ServiceEntry) {
		connectService(entry, id, details, cert)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// discover all services on the network (e.g. _workstation._tcp)
	if err = zeroconf.Browse(ctx, ship.ZeroconfType, ship.ZeroconfDomain, entries); err != nil {
		panic(fmt.Errorf("failed to browse: %w", err))
	}

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
