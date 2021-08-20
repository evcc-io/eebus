package message

import (
	"encoding/json"
	"errors"

	"github.com/evcc-io/eebus/ship/ship"
)

func Decode(b []byte) (interface{}, error) {
	var sum map[string]json.RawMessage

	if err := json.Unmarshal(b, &sum); err != nil {
		return nil, err
	}

	var typ string
	var raw json.RawMessage
	for k, v := range sum {
		typ = k
		raw = v
	}

	switch typ {
	case "accessMethods":
		res := []ship.AccessMethods{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ship.AccessMethods{}, nil

	case "accessMethodsRequest":
		res := []ship.AccessMethodsRequest{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ship.AccessMethodsRequest{}, nil

	case "connectionPinState":
		res := []ship.ConnectionPinState{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ship.ConnectionPinState{}, nil

	case "connectionPinInput":
		res := []ship.ConnectionPinInput{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ship.ConnectionPinInput{}, nil

	case "connectionPinError":
		res := []ship.ConnectionPinError{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ship.ConnectionPinError{}, nil

	case "connectionHello":
		res := []ship.ConnectionHello{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ship.ConnectionHello{}, nil

	case "connectionClose":
		res := []ship.ConnectionClose{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ship.ConnectionClose{}, nil

	case "messageProtocolHandshake":
		res := ship.MessageProtocolHandshake{}
		err := json.Unmarshal(raw, &res)
		return res, err

	case "data":
		res := ship.Data{}
		err := json.Unmarshal(raw, &res)
		return res, err

	default:
		return nil, errors.New("invalid type")
	}
}
