package ship

import "crypto/tls"

// protocol constants
const (
	Scheme         = "wss://"
	SubProtocol    = "ship"
	ZeroconfType   = "_ship._tcp"
	ZeroconfDomain = "local."
)

// CipherSuites are the SHIP cipher suites
var CipherSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
}
