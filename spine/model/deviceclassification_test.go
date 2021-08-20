package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDeviceClassificationManufacturerData(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 6
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 1
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 194
	var msgCounterReference MsgCounterType = 4890
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	var deviceName DeviceClassificationStringType = ""
	var deviceCode DeviceClassificationStringType = ""
	var brandName DeviceClassificationStringType = ""

	datagram := CmiDatagramType{
		Datagram: DatagramType{
			Header: HeaderType{
				SpecificationVersion: &specificationVersion,
				AddressSource: &FeatureAddressType{
					Feature: &addressSourceFeature,
					Entity:  []AddressEntityType{addressSoureEntity1, addressSoureEntity1},
					Device:  &deviceEvse,
				},
				AddressDestination: &FeatureAddressType{
					Feature: &addressDestinationFeature,
					Entity:  []AddressEntityType{addressDestinationEntity1},
					Device:  &deviceHems,
				},
				MsgCounter:          &msgCounter,
				MsgCounterReference: &msgCounterReference,
				CmdClassifier:       &cmdClassifier,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					DeviceClassificationManufacturerData: &DeviceClassificationManufacturerDataType{
						DeviceName:  &deviceName,
						DeviceCode:  &deviceCode,
						BrandName:   &brandName,
						PowerSource: string(PowerSourceEnumTypeMains3Phase),
					},
				}},
			},
		},
	}
	json, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestDeviceClassificationManufacturerData() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":6}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":1}]},{"msgCounter":194},{"msgCounterReference":4890},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"deviceClassificationManufacturerData":[{"deviceName":""},{"deviceCode":""},{"brandName":""},{"powerSource":"mains3Phase"}]}]]}]}]}`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestDeviceClassificationManufacturerData() actual json string doesn't match expected result")
	}
}
