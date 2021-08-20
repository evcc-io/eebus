package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDeviceConfigurationKeyValueDescriptionListData(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 5
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 4
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 237
	var msgCounterReference MsgCounterType = 4926
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	var keyID1 DeviceConfigurationKeyIdType = 1
	var keyName1 string = string(DeviceConfigurationKeyNameEnumTypeAsymmetricChargingSupported)
	var keyID2 DeviceConfigurationKeyIdType = 2
	var keyName2 string = string(DeviceConfigurationKeyNameEnumTypeCommunicationsStandard)
	var typeBool DeviceConfigurationKeyValueTypeType = DeviceConfigurationKeyValueTypeTypeBoolean
	var typeString DeviceConfigurationKeyValueTypeType = DeviceConfigurationKeyValueTypeTypeString

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
					DeviceConfigurationKeyValueDescriptionListData: &DeviceConfigurationKeyValueDescriptionListDataType{
						DeviceConfigurationKeyValueDescriptionData: []DeviceConfigurationKeyValueDescriptionDataType{
							{
								KeyId:     &keyID1,
								KeyName:   &keyName1,
								ValueType: &typeBool,
							},
							{
								KeyId:     &keyID2,
								KeyName:   &keyName2,
								ValueType: &typeString,
							},
						},
					},
				}},
			},
		},
	}
	json, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestDeviceConfigurationKeyValueDescriptionListData() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":5}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":4}]},{"msgCounter":237},{"msgCounterReference":4926},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"deviceConfigurationKeyValueDescriptionListData":[{"deviceConfigurationKeyValueDescriptionData":[[{"keyId":1},{"keyName":"asymmetricChargingSupported"},{"valueType":"boolean"}],[{"keyId":2},{"keyName":"communicationsStandard"},{"valueType":"string"}]]}]}]]}]}]}`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestDeviceConfigurationKeyValueDescriptionListData() actual json string doesn't match expected result")
	}
}

func TestDeviceConfigurationKeyValueListData(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 5
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 4
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 241
	var msgCounterReference MsgCounterType = 4931
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	var keyID1 DeviceConfigurationKeyIdType = 1
	var keyID2 DeviceConfigurationKeyIdType = 2
	var keyID1Value bool = false
	var keyID2Value DeviceConfigurationKeyValueStringType = "iec61851"

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
					DeviceConfigurationKeyValueListData: &DeviceConfigurationKeyValueListDataType{
						DeviceConfigurationKeyValueData: []DeviceConfigurationKeyValueDataType{
							{
								KeyId: &keyID1,
								Value: &DeviceConfigurationKeyValueValueType{
									Boolean: &keyID1Value,
								},
							},
							{
								KeyId: &keyID2,
								Value: &DeviceConfigurationKeyValueValueType{
									String: &keyID2Value,
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
		t.Errorf("TestDeviceConfigurationKeyValueListData() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":5}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":4}]},{"msgCounter":241},{"msgCounterReference":4931},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"deviceConfigurationKeyValueListData":[{"deviceConfigurationKeyValueData":[[{"keyId":1},{"value":[{"boolean":false}]}],[{"keyId":2},{"value":[{"string":"iec61851"}]}]]}]}]]}]}]}`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestDeviceConfigurationKeyValueListData() actual json string doesn't match expected result")
	}
}
