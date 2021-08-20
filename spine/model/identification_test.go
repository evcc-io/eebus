package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIdentificationListDataResponse(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 10
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 8
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 247
	var msgCounterReference MsgCounterType = 4949
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	var identificationId IdentificationIdType = 0
	var identificationType IdentificationTypeType = IdentificationTypeType(IdentificationTypeEnumTypeEui48)
	var identificationValue IdentificationValueType = ""

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
					IdentificationListData: &IdentificationListDataType{
						IdentificationData: []IdentificationDataType{
							{
								IdentificationId:    &identificationId,
								IdentificationType:  &identificationType,
								IdentificationValue: &identificationValue,
							},
						},
					},
				}},
			},
		},
	}
	json, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestIdentificationListDataResponse() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":10}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":8}]},{"msgCounter":247},{"msgCounterReference":4949},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"identificationListData":[{"identificationData":[[{"identificationId":0},{"identificationType":"eui48"},{"identificationValue":""}]]}]}]]}]}]}`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestIdentificationListDataResponse() actual json string doesn't match expected result")
	}
}
