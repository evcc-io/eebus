package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoadControlLimitListDataWrite(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 7
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 1
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 5014
	var cmdClassifier CmdClassifierType = CmdClassifierTypeWrite
	var ackRequest bool = true

	var limitId1 LoadControlLimitIdType = 1
	var isLimitActive1 bool = true
	var valueNumber1 NumberType = 0
	var valueScale1 ScaleType = 0
	var limitId2 LoadControlLimitIdType = 2
	var isLimitActive2 bool = true
	var valueNumber2 NumberType = 0
	var valueScale2 ScaleType = 0

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
					Entity:  []AddressEntityType{addressDestinationEntity1, addressDestinationEntity1},
					Device:  &deviceEvse,
				},
				MsgCounter:    &msgCounter,
				CmdClassifier: &cmdClassifier,
				AckRequest:    &ackRequest,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					LoadControlLimitListData: &LoadControlLimitListDataType{
						LoadControlLimitData: []LoadControlLimitDataType{
							{
								LimitId:       &limitId1,
								IsLimitActive: &isLimitActive1,
								Value: &ScaledNumberType{
									Number: &valueNumber1,
									Scale:  &valueScale1,
								},
							},
							{
								LimitId:       &limitId2,
								IsLimitActive: &isLimitActive2,
								Value: &ScaledNumberType{
									Number: &valueNumber2,
									Scale:  &valueScale2,
								},
							},
						},
					},
				}},
			},
		},
	}
	json, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestLoadControlLimitListDataWrite() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":7}]},{"addressDestination":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":1}]},{"msgCounter":5014},{"cmdClassifier":"write"},{"ackRequest":true}]},{"payload":[{"cmd":[[{"loadControlLimitListData":[{"loadControlLimitData":[[{"limitId":1},{"isLimitActive":true},{"value":[{"number":0},{"scale":0}]}],[{"limitId":2},{"isLimitActive":true},{"value":[{"number":0},{"scale":0}]}]]}]}]]}]}]}`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestLoadControlLimitListDataWrite() actual json string doesn't match expected result")
	}
}
