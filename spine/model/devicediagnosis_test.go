package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDeviceDiagnosisHeartbeatData(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 5
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 5
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 5971
	var cmdClassifier CmdClassifierType = CmdClassifierTypeNotify

	var heartBeatCounter uint64 = 2245
	var heartBeatTimeout string = "PT4S"

	datagram := CmiDatagramType{
		Datagram: DatagramType{
			Header: HeaderType{
				SpecificationVersion: &specificationVersion,
				AddressSource: &FeatureAddressType{
					Feature: &addressSourceFeature,
					Entity:  []AddressEntityType{addressSoureEntity1},
					Device:  &deviceHems,
				},
				AddressDestination: &FeatureAddressType{
					Feature: &addressDestinationFeature,
					Entity:  []AddressEntityType{addressDestinationEntity1},
					Device:  &deviceEvse,
				},
				MsgCounter:    &msgCounter,
				CmdClassifier: &cmdClassifier,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					DeviceDiagnosisHeartbeatData: &DeviceDiagnosisHeartbeatDataType{
						HeartbeatCounter: &heartBeatCounter,
						HeartbeatTimeout: &heartBeatTimeout,
					},
				}},
			},
		},
	}
	datagramJSON, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestDeviceDiagnosisHeartbeatData() error = %v", err)
	}

	expectedJSON := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":5}]},{"addressDestination":[{"device":"d:_i:3210_EVSE"},{"entity":[1]},{"feature":5}]},{"msgCounter":5971},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"deviceDiagnosisHeartbeatData":[{"heartbeatCounter":2245},{"heartbeatTimeout":"PT4S"}]}]]}]}]}`
	expectedDatagram := CmiDatagramType{}
	err = json.Unmarshal(json.RawMessage(expectedJSON), &expectedDatagram)
	if err != nil {
		t.Errorf("TestDeviceDiagnosisHeartbeatData() Unmarshal failed error = %v", err)
	} else {
		if string(datagramJSON) != expectedJSON {
			fmt.Println("EXPECTED:")
			fmt.Println(string(expectedJSON))
			fmt.Println("\nACTUAL:")
			fmt.Println(string(datagramJSON))

			t.Errorf("TestDeviceDiagnosisHeartbeatData() actual json string doesn't match expected result")
		}
	}
}
