package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestReadMeasurementListData(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 3
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 3
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 4920
	var cmdClassifier CmdClassifierType = CmdClassifierTypeRead
	var ackRequest bool = true

	datagram := CmiDatagramType{
		Datagram: DatagramType{
			Header: HeaderType{
				SpecificationVersion: &specificationVersion,
				AddressSource: &FeatureAddressType{
					Device:  &deviceHems,
					Entity:  []AddressEntityType{addressSoureEntity1},
					Feature: &addressSourceFeature,
				},
				AddressDestination: &FeatureAddressType{
					Device:  &deviceEvse,
					Entity:  []AddressEntityType{addressDestinationEntity1, addressDestinationEntity1},
					Feature: &addressDestinationFeature,
				},
				MsgCounter:    &msgCounter,
				CmdClassifier: &cmdClassifier,
				AckRequest:    &ackRequest,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					MeasurementListData: &MeasurementListDataType{},
				}},
			},
		},
	}
	datagramJSON, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestReadMeasurementListData() error = %v", err)
	}
	expectedJSON := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":3}]},{"addressDestination":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":3}]},{"msgCounter":4920},{"cmdClassifier":"read"},{"ackRequest":true}]},{"payload":[{"cmd":[[{"measurementListData":[]}]]}]}]}`
	expectedDatagram := CmiDatagramType{}
	err = json.Unmarshal(json.RawMessage(expectedJSON), &expectedDatagram)
	if err != nil {
		t.Errorf("TestReadMeasurementListData() Unmarshal failed error = %v", err)
	} else {
		if string(datagramJSON) != expectedJSON {
			fmt.Println("EXPECTED:")
			fmt.Println(string(expectedJSON))
			fmt.Println("\nACTUAL:")
			fmt.Println(string(datagramJSON))

			t.Errorf("TestReadMeasurementListData() actual json string doesn't match expected result")
		}
	}
}

func TestMeasurementListDataResponse(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 3
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 3
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 235
	var msgCounterReference MsgCounterType = 4920
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	var measurementId1 MeasurementIdType = 1
	var valueType1 MeasurementValueTypeType = MeasurementValueTypeType(MeasurementValueTypeEnumTypeValue)
	var timestamp1 string = "2021-02-17T20:19:05.099Z"
	var valueNumber1 NumberType = 0
	var valueScale1 ScaleType = 0
	var valueSource1 MeasurementValueSourceType = MeasurementValueSourceType(MeasurementValueSourceEnumTypeMeasuredValue)

	var measurementId2 MeasurementIdType = 4
	var valueType2 MeasurementValueTypeType = MeasurementValueTypeType(MeasurementValueTypeEnumTypeValue)
	var timestamp2 string = "2021-02-17T20:19:05.099Z"
	var valueNumber2 NumberType = 0
	var valueScale2 ScaleType = 0
	var valueSource2 MeasurementValueSourceType = MeasurementValueSourceType(MeasurementValueSourceEnumTypeMeasuredValue)

	var measurementId3 MeasurementIdType = 7
	var valueType3 MeasurementValueTypeType = MeasurementValueTypeType(MeasurementValueTypeEnumTypeValue)
	var timestamp3 string = "2021-02-17T20:19:05.099Z"
	var valueNumber3 NumberType = 0
	var valueScale3 ScaleType = 0
	var valueSource3 MeasurementValueSourceType = MeasurementValueSourceType(MeasurementValueSourceEnumTypeMeasuredValue)

	datagram := CmiDatagramType{
		Datagram: DatagramType{
			Header: HeaderType{
				SpecificationVersion: &specificationVersion,
				AddressSource: &FeatureAddressType{
					Device:  &deviceEvse,
					Entity:  []AddressEntityType{addressSoureEntity1, addressSoureEntity1},
					Feature: &addressSourceFeature,
				},
				AddressDestination: &FeatureAddressType{
					Device:  &deviceHems,
					Entity:  []AddressEntityType{addressDestinationEntity1},
					Feature: &addressDestinationFeature,
				},
				MsgCounter:          &msgCounter,
				MsgCounterReference: &msgCounterReference,
				CmdClassifier:       &cmdClassifier,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					MeasurementListData: &MeasurementListDataType{
						MeasurementData: []MeasurementDataType{
							{
								MeasurementId: &measurementId1,
								ValueType:     &valueType1,
								Timestamp:     &timestamp1,
								Value: &ScaledNumberType{
									Number: &valueNumber1,
									Scale:  &valueScale1,
								},
								ValueSource: &valueSource1,
							},
							{
								MeasurementId: &measurementId2,
								ValueType:     &valueType2,
								Timestamp:     &timestamp2,
								Value: &ScaledNumberType{
									Number: &valueNumber2,
									Scale:  &valueScale2,
								},
								ValueSource: &valueSource2,
							},
							{
								MeasurementId: &measurementId3,
								ValueType:     &valueType3,
								Timestamp:     &timestamp3,
								Value: &ScaledNumberType{
									Number: &valueNumber3,
									Scale:  &valueScale3,
								},
								ValueSource: &valueSource3,
							},
						},
					},
				}},
			},
		},
	}
	datagramJSON, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestMeasurementListDataResponse() error = %v", err)
	}
	expectedJSON := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":3}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":3}]},{"msgCounter":235},{"msgCounterReference":4920},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"measurementListData":[{"measurementData":[[{"measurementId":1},{"valueType":"value"},{"timestamp":"2021-02-17T20:19:05.099Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":4},{"valueType":"value"},{"timestamp":"2021-02-17T20:19:05.099Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":7},{"valueType":"value"},{"timestamp":"2021-02-17T20:19:05.099Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}]]}]}]]}]}]}`
	expectedDatagram := CmiDatagramType{}
	err = json.Unmarshal(json.RawMessage(expectedJSON), &expectedDatagram)
	if err != nil {
		t.Errorf("TestMeasurementListDataResponse() Unmarshal failed error = %v", err)
	} else {
		if string(datagramJSON) != expectedJSON {
			fmt.Println("EXPECTED:")
			fmt.Println(string(expectedJSON))
			fmt.Println("\nACTUAL:")
			fmt.Println(string(datagramJSON))

			t.Errorf("TestMeasurementListDataResponse() actual json string doesn't match expected result")
		}
	}
}

func TestMeasurementDescriptionListDataResponse(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 3
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 3
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 230
	var msgCounterReference MsgCounterType = 4917
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	var measurementId1 MeasurementIdType = 1
	var measurementType1 MeasurementTypeType = MeasurementTypeType(MeasurementTypeEnumTypeCurrent)
	var commodityType1 CommodityTypeType = CommodityTypeType(CommodityTypeEnumTypeElectricity)
	var unit1 UnitOfMeasurementType = UnitOfMeasurementType(UnitOfMeasurementEnumTypeA)
	var scopeType1 ScopeTypeType = ScopeTypeType(ScopeTypeEnumTypeACCurrent)

	var measurementId2 MeasurementIdType = 4
	var measurementType2 MeasurementTypeType = MeasurementTypeType(MeasurementTypeEnumTypePower)
	var commodityType2 CommodityTypeType = CommodityTypeType(CommodityTypeEnumTypeElectricity)
	var unit2 UnitOfMeasurementType = UnitOfMeasurementType(UnitOfMeasurementEnumTypeW)
	var scopeType2 ScopeTypeType = ScopeTypeType(ScopeTypeEnumTypeACPower)

	var measurementId3 MeasurementIdType = 7
	var measurementType3 MeasurementTypeType = MeasurementTypeType(MeasurementTypeEnumTypeEnergy)
	var commodityType3 CommodityTypeType = CommodityTypeType(CommodityTypeEnumTypeElectricity)
	var unit3 UnitOfMeasurementType = UnitOfMeasurementType(UnitOfMeasurementEnumTypeWh)
	var scopeType3 ScopeTypeType = ScopeTypeType(ScopeTypeEnumTypeCharge)

	datagram := CmiDatagramType{
		Datagram: DatagramType{
			Header: HeaderType{
				SpecificationVersion: &specificationVersion,
				AddressSource: &FeatureAddressType{
					Device:  &deviceEvse,
					Entity:  []AddressEntityType{addressSoureEntity1, addressSoureEntity1},
					Feature: &addressSourceFeature,
				},
				AddressDestination: &FeatureAddressType{
					Device:  &deviceHems,
					Entity:  []AddressEntityType{addressDestinationEntity1},
					Feature: &addressDestinationFeature,
				},
				MsgCounter:          &msgCounter,
				MsgCounterReference: &msgCounterReference,
				CmdClassifier:       &cmdClassifier,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					MeasurementDescriptionListData: &MeasurementDescriptionListDataType{
						MeasurementDescriptionData: []MeasurementDescriptionDataType{
							{
								MeasurementId:   &measurementId1,
								MeasurementType: &measurementType1,
								CommodityType:   &commodityType1,
								Unit:            &unit1,
								ScopeType:       &scopeType1,
							},
							{
								MeasurementId:   &measurementId2,
								MeasurementType: &measurementType2,
								CommodityType:   &commodityType2,
								Unit:            &unit2,
								ScopeType:       &scopeType2,
							},
							{
								MeasurementId:   &measurementId3,
								MeasurementType: &measurementType3,
								CommodityType:   &commodityType3,
								Unit:            &unit3,
								ScopeType:       &scopeType3,
							},
						},
					},
				}},
			},
		},
	}
	datagramJSON, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestMeasurementDescriptionListDataResponse() error = %v", err)
	}
	expectedJSON := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":3}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":3}]},{"msgCounter":230},{"msgCounterReference":4917},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"measurementDescriptionListData":[{"measurementDescriptionData":[[{"measurementId":1},{"measurementType":"current"},{"commodityType":"electricity"},{"unit":"A"},{"scopeType":"acCurrent"}],[{"measurementId":4},{"measurementType":"power"},{"commodityType":"electricity"},{"unit":"W"},{"scopeType":"acPower"}],[{"measurementId":7},{"measurementType":"energy"},{"commodityType":"electricity"},{"unit":"Wh"},{"scopeType":"charge"}]]}]}]]}]}]}`
	expectedDatagram := CmiDatagramType{}
	err = json.Unmarshal(json.RawMessage(expectedJSON), &expectedDatagram)
	if err != nil {
		t.Errorf("TestMeasurementDescriptionListDataResponse() Unmarshal failed error = %v", err)
	} else {
		if string(datagramJSON) != expectedJSON {
			fmt.Println("EXPECTED:")
			fmt.Println(string(expectedJSON))
			fmt.Println("\nACTUAL:")
			fmt.Println(string(datagramJSON))

			t.Errorf("TestMeasurementDescriptionListDataResponse() actual json string doesn't match expected result")
		}
	}
}

func TestMeasurementListDataNotify(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.1.1"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 3
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 3
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 239
	var cmdClassifier CmdClassifierType = CmdClassifierTypeNotify

	var measurementId1 MeasurementIdType = 1
	var valueType1 MeasurementValueTypeType = MeasurementValueTypeType(MeasurementValueTypeEnumTypeValue)
	var timestamp1 string = "2021-02-17T20:19:06.188Z"
	var valueNumber1 NumberType = 0
	var valueScale1 ScaleType = 0
	var valueSource1 MeasurementValueSourceType = MeasurementValueSourceType(MeasurementValueSourceEnumTypeMeasuredValue)

	var measurementId2 MeasurementIdType = 4
	var valueType2 MeasurementValueTypeType = MeasurementValueTypeType(MeasurementValueTypeEnumTypeValue)
	var timestamp2 string = "2021-02-17T20:19:06.188Z"
	var valueNumber2 NumberType = 0
	var valueScale2 ScaleType = 0
	var valueSource2 MeasurementValueSourceType = MeasurementValueSourceType(MeasurementValueSourceEnumTypeMeasuredValue)

	var measurementId3 MeasurementIdType = 7
	var valueType3 MeasurementValueTypeType = MeasurementValueTypeType(MeasurementValueTypeEnumTypeValue)
	var timestamp3 string = "2021-02-17T20:19:06.188Z"
	var valueNumber3 NumberType = 0
	var valueScale3 ScaleType = 0
	var valueSource3 MeasurementValueSourceType = MeasurementValueSourceType(MeasurementValueSourceEnumTypeMeasuredValue)

	datagram := CmiDatagramType{
		Datagram: DatagramType{
			Header: HeaderType{
				SpecificationVersion: &specificationVersion,
				AddressSource: &FeatureAddressType{
					Device:  &deviceEvse,
					Entity:  []AddressEntityType{addressSoureEntity1, addressSoureEntity1},
					Feature: &addressSourceFeature,
				},
				AddressDestination: &FeatureAddressType{
					Device:  &deviceHems,
					Entity:  []AddressEntityType{addressDestinationEntity1},
					Feature: &addressDestinationFeature,
				},
				MsgCounter:    &msgCounter,
				CmdClassifier: &cmdClassifier,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					MeasurementListData: &MeasurementListDataType{
						MeasurementData: []MeasurementDataType{
							{
								MeasurementId: &measurementId1,
								ValueType:     &valueType1,
								Timestamp:     &timestamp1,
								Value: &ScaledNumberType{
									Number: &valueNumber1,
									Scale:  &valueScale1,
								},
								ValueSource: &valueSource1,
							},
							{
								MeasurementId: &measurementId2,
								ValueType:     &valueType2,
								Timestamp:     &timestamp2,
								Value: &ScaledNumberType{
									Number: &valueNumber2,
									Scale:  &valueScale2,
								},
								ValueSource: &valueSource2,
							},
							{
								MeasurementId: &measurementId3,
								ValueType:     &valueType3,
								Timestamp:     &timestamp3,
								Value: &ScaledNumberType{
									Number: &valueNumber3,
									Scale:  &valueScale3,
								},
								ValueSource: &valueSource3,
							}},
					},
				}},
			},
		},
	}
	datagramJSON, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestMeasurementListDataNotify() error = %v", err)
	}
	expectedJSON := `{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":3}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":3}]},{"msgCounter":239},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"measurementListData":[{"measurementData":[[{"measurementId":1},{"valueType":"value"},{"timestamp":"2021-02-17T20:19:06.188Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":4},{"valueType":"value"},{"timestamp":"2021-02-17T20:19:06.188Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":7},{"valueType":"value"},{"timestamp":"2021-02-17T20:19:06.188Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}]]}]}]]}]}]}`
	expectedDatagram := CmiDatagramType{}
	err = json.Unmarshal(json.RawMessage(expectedJSON), &expectedDatagram)
	if err != nil {
		t.Errorf("TestMeasurementListDataNotify() Unmarshal failed error = %v", err)
	} else {
		if string(datagramJSON) != expectedJSON {
			fmt.Println("EXPECTED:")
			fmt.Println(string(expectedJSON))
			fmt.Println("\nACTUAL:")
			fmt.Println(string(datagramJSON))

			t.Errorf("TestMeasurementListDataNotify() actual json string doesn't match expected result")
		}
	}
}
