package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestElectricalConnectionDescriptionListDataRead(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 9
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 2
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 4955
	var cmdClassifier CmdClassifierType = CmdClassifierTypeRead
	var ackRequest bool = true

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
					ElectricalConnectionDescriptionListData: &ElectricalConnectionDescriptionListDataType{},
				}},
			},
		},
	}
	json, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestElectricalConnectionDescriptionListDataRead() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":9}]},{"addressDestination":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":2}]},{"msgCounter":4955},{"cmdClassifier":"read"},{"ackRequest":true}]},{"payload":[{"cmd":[[{"electricalConnectionDescriptionListData":[]}]]}]}]}`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestElectricalConnectionDescriptionListDataRead() actual json string doesn't match expected result")
	}
}

func TestElectricalConnectionDescriptionListDataResponse(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 2
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 9
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 249
	var msgCounterReference MsgCounterType = 4955
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	var electricalConnectionId ElectricalConnectionIdType = 0
	var powerSupplyType ElectricalConnectionVoltageTypeType = ElectricalConnectionVoltageTypeType(ElectricalConnectionVoltageTypeEnumTypeAc)
	var acConnectedPhases uint = 1
	var positiveEnergyDirection EnergyDirectionType = EnergyDirectionType(EnergyDirectionEnumTypeConsume)

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
					ElectricalConnectionDescriptionListData: &ElectricalConnectionDescriptionListDataType{
						ElectricalConnectionDescriptionData: []ElectricalConnectionDescriptionDataType{
							{
								ElectricalConnectionId:  &electricalConnectionId,
								PowerSupplyType:         &powerSupplyType,
								AcConnectedPhases:       &acConnectedPhases,
								PositiveEnergyDirection: &positiveEnergyDirection,
							},
						}},
				}},
			},
		},
	}
	json, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestElectricalConnectionDescriptionListDataResponse() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":2}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":249},{"msgCounterReference":4955},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"electricalConnectionDescriptionListData":[{"electricalConnectionDescriptionData":[[{"electricalConnectionId":0},{"powerSupplyType":"ac"},{"acConnectedPhases":1},{"positiveEnergyDirection":"consume"}]]}]}]]}]}]}`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestElectricalConnectionDescriptionListDataResponse() actual json string doesn't match expected result")
	}
}

func TestElectricalConnectionParameterDescriptionListDataResponse(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 2
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 9
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 251
	var msgCounterReference MsgCounterType = 4958
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	var electricalConnectionId1 ElectricalConnectionIdType = 0
	var parameterId1 ElectricalConnectionParameterIdType = 1
	var measurementId1 MeasurementIdType = 1
	var voltageType1 ElectricalConnectionVoltageTypeType = ElectricalConnectionVoltageTypeType(ElectricalConnectionVoltageTypeEnumTypeAc)
	var acMeasuredPhases1 ElectricalConnectionPhaseNameType = ElectricalConnectionPhaseNameType(ElectricalConnectionPhaseNameEnumTypeA)
	var acMeasuredInReferenceTo1 ElectricalConnectionPhaseNameType = ElectricalConnectionPhaseNameType(ElectricalConnectionPhaseNameEnumTypeNeutral)
	var acMeasurementType1 ElectricalConnectionAcMeasurementTypeType = ElectricalConnectionAcMeasurementTypeType(ElectricalConnectionAcMeasurementTypeEnumTypeReal)
	var acMeasurementVariant1 ElectricalConnectionMeasurandVariantType = ElectricalConnectionMeasurandVariantType(ElectricalConnectionMeasurandVariantEnumTypeRms)
	var electricalConnectionId2 ElectricalConnectionIdType = 0
	var parameterId2 ElectricalConnectionParameterIdType = 2
	var measurementId2 MeasurementIdType = 4
	var voltageType2 ElectricalConnectionVoltageTypeType = ElectricalConnectionVoltageTypeType(ElectricalConnectionVoltageTypeEnumTypeAc)
	var acMeasuredPhases2 ElectricalConnectionPhaseNameType = ElectricalConnectionPhaseNameType(ElectricalConnectionPhaseNameEnumTypeA)
	var acMeasuredInReferenceTo2 ElectricalConnectionPhaseNameType = ElectricalConnectionPhaseNameType(ElectricalConnectionPhaseNameEnumTypeNeutral)
	var acMeasurementType2 ElectricalConnectionAcMeasurementTypeType = ElectricalConnectionAcMeasurementTypeType(ElectricalConnectionAcMeasurementTypeEnumTypeReal)
	var acMeasurementVariant2 ElectricalConnectionMeasurandVariantType = ElectricalConnectionMeasurandVariantType(ElectricalConnectionMeasurandVariantEnumTypeRms)
	var electricalConnectionId3 ElectricalConnectionIdType = 0
	var parameterId3 ElectricalConnectionParameterIdType = 3
	var measurementId3 MeasurementIdType = 7
	var voltageType3 ElectricalConnectionVoltageTypeType = ElectricalConnectionVoltageTypeType(ElectricalConnectionVoltageTypeEnumTypeAc)
	var acMeasuredPhases3 ElectricalConnectionPhaseNameType = ElectricalConnectionPhaseNameType(ElectricalConnectionPhaseNameEnumTypeA)
	var acMeasuredInReferenceTo3 ElectricalConnectionPhaseNameType = ElectricalConnectionPhaseNameType(ElectricalConnectionPhaseNameEnumTypeNeutral)
	var acMeasurementType3 ElectricalConnectionAcMeasurementTypeType = ElectricalConnectionAcMeasurementTypeType(ElectricalConnectionAcMeasurementTypeEnumTypeReal)
	var acMeasurementVariant3 ElectricalConnectionMeasurandVariantType = ElectricalConnectionMeasurandVariantType(ElectricalConnectionMeasurandVariantEnumTypeRms)
	var electricalConnectionId4 ElectricalConnectionIdType = 0
	var parameterId4 ElectricalConnectionParameterIdType = 8
	var acMeasuredPhases4 ElectricalConnectionPhaseNameType = ElectricalConnectionPhaseNameType(ElectricalConnectionPhaseNameEnumTypeA)
	var scopeType4 ScopeTypeType = ScopeTypeType(ScopeTypeEnumTypeACPowerTotal)

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
					ElectricalConnectionParameterDescriptionListData: &ElectricalConnectionParameterDescriptionListDataType{
						ElectricalConnectionParameterDescriptionData: []ElectricalConnectionParameterDescriptionDataType{
							{
								ElectricalConnectionId:  &electricalConnectionId1,
								ParameterId:             &parameterId1,
								MeasurementId:           &measurementId1,
								VoltageType:             &voltageType1,
								AcMeasuredPhases:        &acMeasuredPhases1,
								AcMeasuredInReferenceTo: &acMeasuredInReferenceTo1,
								AcMeasurementType:       &acMeasurementType1,
								AcMeasurementVariant:    &acMeasurementVariant1,
							},
							{
								ElectricalConnectionId:  &electricalConnectionId2,
								ParameterId:             &parameterId2,
								MeasurementId:           &measurementId2,
								VoltageType:             &voltageType2,
								AcMeasuredPhases:        &acMeasuredPhases2,
								AcMeasuredInReferenceTo: &acMeasuredInReferenceTo2,
								AcMeasurementType:       &acMeasurementType2,
								AcMeasurementVariant:    &acMeasurementVariant2,
							},
							{
								ElectricalConnectionId:  &electricalConnectionId3,
								ParameterId:             &parameterId3,
								MeasurementId:           &measurementId3,
								VoltageType:             &voltageType3,
								AcMeasuredPhases:        &acMeasuredPhases3,
								AcMeasuredInReferenceTo: &acMeasuredInReferenceTo3,
								AcMeasurementType:       &acMeasurementType3,
								AcMeasurementVariant:    &acMeasurementVariant3,
							},
							{
								ElectricalConnectionId: &electricalConnectionId4,
								ParameterId:            &parameterId4,
								AcMeasuredPhases:       &acMeasuredPhases4,
								ScopeType:              &scopeType4,
							},
						}},
				}},
			},
		},
	}
	json, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestElectricalConnectionParameterDescriptionListDataResponse() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":2}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":251},{"msgCounterReference":4958},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"electricalConnectionParameterDescriptionListData":[{"electricalConnectionParameterDescriptionData":[[{"electricalConnectionId":0},{"parameterId":1},{"measurementId":1},{"voltageType":"ac"},{"acMeasuredPhases":"a"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":2},{"measurementId":4},{"voltageType":"ac"},{"acMeasuredPhases":"a"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":3},{"measurementId":7},{"voltageType":"ac"},{"acMeasuredPhases":"a"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":8},{"acMeasuredPhases":"a"},{"scopeType":"acPowerTotal"}]]}]}]]}]}]}`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestElectricalConnectionParameterDescriptionListDataResponse() actual json string doesn't match expected result")
	}
}

func TestElectricalConnectionPermittedValueSetListDataResponse(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 2
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 9
	var addressDestinationEntity1 AddressEntityType = 1
	var msgCounter MsgCounterType = 253
	var msgCounterReference MsgCounterType = 4961
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	var electricalConnectionId1 ElectricalConnectionIdType = 0
	var parameterId1 ElectricalConnectionParameterIdType = 1
	var valueNumber1 NumberType = 0
	var valueScale1 ScaleType = 0
	var rangeMinNumber1 NumberType = 6
	var rangeMinScale1 ScaleType = 0
	var rangeMaxNumber1 NumberType = 10
	var rangeMaxScale1 ScaleType = 0
	var electricalConnectionId2 ElectricalConnectionIdType = 0
	var parameterId2 ElectricalConnectionParameterIdType = 8
	var valueNumber2 NumberType = 0
	var valueScale2 ScaleType = 0
	var rangeMinNumber2 NumberType = 1384200
	var rangeMinScale2 ScaleType = -3
	var rangeMaxNumber2 NumberType = 2307
	var rangeMaxScale2 ScaleType = 0

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
					ElectricalConnectionPermittedValueSetListData: &ElectricalConnectionPermittedValueSetListDataType{
						ElectricalConnectionPermittedValueSetData: []ElectricalConnectionPermittedValueSetDataType{
							{
								ElectricalConnectionId: &electricalConnectionId1,
								ParameterId:            &parameterId1,
								PermittedValueSet: []ScaledNumberSetType{
									{
										Value: []ScaledNumberType{
											{
												Number: &valueNumber1,
												Scale:  &valueScale1,
											},
										},
										Range: []ScaledNumberRangeType{
											{
												Min: &ScaledNumberType{
													Number: &rangeMinNumber1,
													Scale:  &rangeMinScale1,
												},
												Max: &ScaledNumberType{
													Number: &rangeMaxNumber1,
													Scale:  &rangeMaxScale1,
												},
											},
										},
									},
								},
							},
							{
								ElectricalConnectionId: &electricalConnectionId2,
								ParameterId:            &parameterId2,
								PermittedValueSet: []ScaledNumberSetType{
									{
										Value: []ScaledNumberType{
											{
												Number: &valueNumber2,
												Scale:  &valueScale2,
											},
										},
										Range: []ScaledNumberRangeType{
											{
												Min: &ScaledNumberType{
													Number: &rangeMinNumber2,
													Scale:  &rangeMinScale2,
												},
												Max: &ScaledNumberType{
													Number: &rangeMaxNumber2,
													Scale:  &rangeMaxScale2,
												},
											},
										},
									},
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
		t.Errorf("TestElectricalConnectionPermittedValueSetListDataResponse() error = %v", err)
	}
	jsonString := string(json)

	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":2}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":253},{"msgCounterReference":4961},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"electricalConnectionPermittedValueSetListData":[{"electricalConnectionPermittedValueSetData":[[{"electricalConnectionId":0},{"parameterId":1},{"permittedValueSet":[[{"value":[[{"number":0},{"scale":0}]]},{"range":[[{"min":[{"number":6},{"scale":0}]},{"max":[{"number":10},{"scale":0}]}]]}]]}],[{"electricalConnectionId":0},{"parameterId":8},{"permittedValueSet":[[{"value":[[{"number":0},{"scale":0}]]},{"range":[[{"min":[{"number":1384200},{"scale":-3}]},{"max":[{"number":2307},{"scale":0}]}]]}]]}]]}]}]]}]}]}`
	if jsonString != jsonTest {
		fmt.Println("EXPECTED:")
		fmt.Println(string(jsonTest))
		fmt.Println("\nACTUAL:")
		fmt.Println(string(jsonString))

		t.Errorf("TestElectricalConnectionPermittedValueSetListDataResponse() actual json string doesn't match expected result")
	}
}
