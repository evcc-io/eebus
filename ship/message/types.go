package message

import "time"

const (
	CmiTypeInit    byte = 0
	CmiTypeControl byte = 1
	CmiTypeData    byte = 2
	CmiTypeEnd     byte = 3

	CmiTimeout                  = 60 * time.Second
	CmiHelloProlongationTimeout = 30 * time.Second
	CmiCloseTimeout             = 100 * time.Millisecond

	ProtocolID = "ee1.0"
)
