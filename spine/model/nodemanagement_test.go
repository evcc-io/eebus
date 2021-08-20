package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestReadNodeManagementDetailedDiscoveryData(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var addressSourceFeature AddressFeatureType = 0
	var addressSoureEntity1 AddressEntityType = 0
	var addressDestinationFeature AddressFeatureType = 0
	var addressDestinationEntity1 AddressEntityType = 0
	var msgCounter MsgCounterType = 5876
	var cmdClassifier CmdClassifierType = CmdClassifierTypeRead

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
				},
				MsgCounter:    &msgCounter,
				CmdClassifier: &cmdClassifier,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					NodeManagementDetailedDiscoveryData: &NodeManagementDetailedDiscoveryDataType{},
				}},
			},
		},
	}
	datagramJSON, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestReadNodeManagementDetailedDiscoveryData() error = %v", err)
	}
	expectedJSON := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_HEMS"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"entity":[0]},{"feature":0}]},{"msgCounter":5876},{"cmdClassifier":"read"}]},{"payload":[{"cmd":[[{"nodeManagementDetailedDiscoveryData":[]}]]}]}]}`
	expectedDatagram := CmiDatagramType{}
	err = json.Unmarshal(json.RawMessage(expectedJSON), &expectedDatagram)
	if err != nil {
		t.Errorf("TestReadNodeManagementDetailedDiscoveryData() Unmarshal failed error = %v", err)
	} else {
		if string(datagramJSON) != expectedJSON {
			fmt.Println("EXPECTED:")
			fmt.Println(string(expectedJSON))
			fmt.Println("\nACTUAL:")
			fmt.Println(string(datagramJSON))

			t.Errorf("TestReadNodeManagementDetailedDiscoveryData() actual json string doesn't match expected result")
		}
	}
}

func TestWriteNodeManagementDetailedDiscoveryData(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.2.0"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 0
	var addressSoureEntity1 AddressEntityType = 0
	var addressDestinationFeature AddressFeatureType = 0
	var addressDestinationEntity1 AddressEntityType = 0
	var msgCounter MsgCounterType = 5881
	var msgCounterReference MsgCounterType = 1
	var cmdClassifier CmdClassifierType = CmdClassifierTypeReply

	smart := NetworkManagementFeatureSetTypeSmart
	var cemEnergyGuardDescription DescriptionType = "CEM Energy Guard"
	var cemControllableSystemDescription DescriptionType = "CEM Controllable System"
	var deviceClassificationDescription DescriptionType = "Device Classification"
	var deviceDiagnosisDescription DescriptionType = "Device Diagnosis"
	var deviceDiagnosisServerDescription DescriptionType = "DeviceDiag"
	var measurementForClientDescription DescriptionType = "Measurement for client"
	var deviceConfigurationDescription DescriptionType = "Device Configuration"
	var loadControlDescription DescriptionType = "LoadControl client for CEM"
	var loadControlServerDescription DescriptionType = "Load Control"
	var evIdentificationDescription DescriptionType = "EV identification"
	var electricalConnectionDescription DescriptionType = "Electrical Connection"

	var featureRoleTypeSpecial RoleType = RoleTypeSpecial
	var featureRoleTypeServer RoleType = RoleTypeServer
	var featureRoleTypeClient RoleType = RoleTypeClient

	var addressEntityType0 AddressEntityType = 0
	var addressEntityType1 AddressEntityType = 1
	var addressEntityType2 AddressEntityType = 2

	var feature0 AddressFeatureType = 0
	var feature1 AddressFeatureType = 1
	var feature2 AddressFeatureType = 2
	var feature3 AddressFeatureType = 3
	var feature4 AddressFeatureType = 4
	var feature5 AddressFeatureType = 5
	var feature7 AddressFeatureType = 7
	var feature8 AddressFeatureType = 8
	var feature9 AddressFeatureType = 9
	var feature10 AddressFeatureType = 10
	var feature11 AddressFeatureType = 11
	var feature12 AddressFeatureType = 12

	var featureTypeNodeMgmt FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeNodeManagement)
	var functionTypeNodeMgmtDetailedDiscovery FunctionType = FunctionType(FunctionEnumTypeNodeManagementDetailedDiscoveryData)
	var functionTypeNodeMgmtSubReqCall FunctionType = FunctionType(FunctionEnumTypeNodeManagementSubscriptionRequestCall)
	var functionTypeNodeMgmtBindReqCall FunctionType = FunctionType(FunctionEnumTypeNodeManagementBindingRequestCall)
	var functionTypeNodeMgmtSubDelCall FunctionType = FunctionType(FunctionEnumTypeNodeManagementSubscriptionDeleteCall)
	var functionTypeNodeMgmtBindDelCall FunctionType = FunctionType(FunctionEnumTypeNodeManagementBindingDeleteCall)
	var functionTypeNodeMgmtSubData FunctionType = FunctionType(FunctionEnumTypeNodeManagementSubscriptionData)
	var functionTypeNodeMgmtBindData FunctionType = FunctionType(FunctionEnumTypeNodeManagementBindingData)
	var functionTypeNodeMgmtUseCaseData FunctionType = FunctionType(FunctionEnumTypeNodeManagementUseCaseData)

	var featureTypeDeviceClassification FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeDeviceClassification)
	var functionTypeDevClassManuData FunctionType = FunctionType(FunctionEnumTypeDeviceClassificationManufacturerData)

	var featureTypeDeviceDiagnosis FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeDeviceDiagnosis)
	var functionTypeDevDiagStateData FunctionType = FunctionType(FunctionEnumTypeDeviceDiagnosisStateData)
	var functionTypeDevDiagHeartBeatData FunctionType = FunctionType(FunctionEnumTypeDeviceDiagnosisHeartbeatData)

	var featureTypeMeasurement FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeMeasurement)

	var featureTypeDeviceConfiguration FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeDeviceConfiguration)
	var functionTypeDevConfKeyValueDescListData FunctionType = FunctionType(FunctionEnumTypeDeviceConfigurationKeyValueDescriptionListData)
	var functionTypeDevConfigKeyValueListData FunctionType = FunctionType(FunctionEnumTypeDeviceConfigurationKeyValueListData)

	var featureTypeLoadControl FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeLoadControl)
	var functionTypeLoadControlLimitDescriptionListData = FunctionType(FunctionEnumTypeLoadControlLimitDescriptionListData)
	var functionTypeLoadControlLimitListData = FunctionType(FunctionEnumTypeLoadControlLimitListData)

	var featureTypeIdentification FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeIdentification)

	var featureTypeElectricalConnection FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeElectricalConnection)
	var functionTypeElConDescriptionListData FunctionType = FunctionType(FunctionEnumTypeElectricalConnectionDescriptionListData)

	var specificationVersionDataType SpecificationVersionDataType = SpecificationVersionDataType(specificationVersion)
	var deviceTypeEMS DeviceTypeType = DeviceTypeType(DeviceTypeEnumTypeEnergyManagementSystem)
	var entityTypeDeviceInformation EntityTypeType = EntityTypeType(EntityTypeEnumTypeDeviceInformation)
	var entityTypeCem EntityTypeType = EntityTypeType(EntityTypeEnumTypeCEM)
	var entityTypeHeatpump EntityTypeType = EntityTypeType(EntityTypeEnumTypeHeatPumpAppliance)

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
				MsgCounter:          &msgCounter,
				MsgCounterReference: &msgCounterReference,
				CmdClassifier:       &cmdClassifier,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					NodeManagementDetailedDiscoveryData: &NodeManagementDetailedDiscoveryDataType{
						SpecificationVersionList: &NodeManagementSpecificationVersionListType{
							SpecificationVersion: []SpecificationVersionDataType{specificationVersionDataType},
						},
						DeviceInformation: &NodeManagementDetailedDiscoveryDeviceInformationType{
							Description: &NetworkManagementDeviceDescriptionDataType{
								DeviceAddress: &DeviceAddressType{
									Device: &deviceHems,
								},
								DeviceType:        &deviceTypeEMS,
								NetworkFeatureSet: &smart,
							},
						},
						EntityInformation: []NodeManagementDetailedDiscoveryEntityInformationType{
							{
								Description: &NetworkManagementEntityDescriptionDataType{
									EntityAddress: &EntityAddressType{Entity: []AddressEntityType{addressEntityType0}},
									EntityType:    &entityTypeDeviceInformation,
								},
							},
							{
								Description: &NetworkManagementEntityDescriptionDataType{
									EntityAddress: &EntityAddressType{Entity: []AddressEntityType{addressEntityType1}},
									EntityType:    &entityTypeCem,
									Description:   &cemEnergyGuardDescription,
								},
							},
							{
								Description: &NetworkManagementEntityDescriptionDataType{
									EntityAddress: &EntityAddressType{Entity: []AddressEntityType{addressEntityType2}},
									EntityType:    &entityTypeHeatpump,
									Description:   &cemControllableSystemDescription,
								},
							},
						},
						FeatureInformation: []NodeManagementDetailedDiscoveryFeatureInformationType{
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature0,
										Entity:  []AddressEntityType{addressEntityType0},
									},
									FeatureType: &featureTypeNodeMgmt,
									Role:        &featureRoleTypeSpecial,
									SupportedFunction: []FunctionPropertyType{
										{
											Function:           &functionTypeNodeMgmtDetailedDiscovery,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
										{
											Function:           &functionTypeNodeMgmtSubReqCall,
											PossibleOperations: &PossibleOperationsType{},
										},
										{
											Function:           &functionTypeNodeMgmtBindReqCall,
											PossibleOperations: &PossibleOperationsType{},
										},
										{
											Function:           &functionTypeNodeMgmtSubDelCall,
											PossibleOperations: &PossibleOperationsType{},
										},
										{
											Function:           &functionTypeNodeMgmtBindDelCall,
											PossibleOperations: &PossibleOperationsType{},
										},
										{
											Function:           &functionTypeNodeMgmtSubData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
										{
											Function:           &functionTypeNodeMgmtBindData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
										{
											Function:           &functionTypeNodeMgmtUseCaseData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
									},
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature1,
										Entity:  []AddressEntityType{addressEntityType0},
									},
									FeatureType: &featureTypeDeviceClassification,
									Role:        &featureRoleTypeServer,
									SupportedFunction: []FunctionPropertyType{
										{
											Function:           &functionTypeDevClassManuData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
									},
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature1,
										Entity:  []AddressEntityType{addressEntityType1},
									},
									FeatureType: &featureTypeDeviceClassification,
									Role:        &featureRoleTypeClient,
									Description: &deviceClassificationDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature2,
										Entity:  []AddressEntityType{addressEntityType1},
									},
									FeatureType: &featureTypeDeviceDiagnosis,
									Role:        &featureRoleTypeClient,
									Description: &deviceDiagnosisDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature3,
										Entity:  []AddressEntityType{addressEntityType1},
									},
									FeatureType: &featureTypeMeasurement,
									Role:        &featureRoleTypeClient,
									Description: &measurementForClientDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature4,
										Entity:  []AddressEntityType{addressEntityType1},
									},
									FeatureType: &featureTypeDeviceConfiguration,
									Role:        &featureRoleTypeClient,
									Description: &deviceConfigurationDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature5,
										Entity:  []AddressEntityType{addressEntityType1},
									},
									FeatureType: &featureTypeDeviceDiagnosis,
									Role:        &featureRoleTypeServer,
									SupportedFunction: []FunctionPropertyType{
										{
											Function:           &functionTypeDevDiagStateData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
										{
											Function:           &functionTypeDevDiagHeartBeatData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
									},
									Description: &deviceDiagnosisServerDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature7,
										Entity:  []AddressEntityType{addressEntityType1},
									},
									FeatureType: &featureTypeLoadControl,
									Role:        &featureRoleTypeClient,
									Description: &loadControlDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature8,
										Entity:  []AddressEntityType{addressEntityType1},
									},
									FeatureType: &featureTypeIdentification,
									Role:        &featureRoleTypeClient,
									Description: &evIdentificationDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature9,
										Entity:  []AddressEntityType{addressEntityType1},
									},
									FeatureType: &featureTypeElectricalConnection,
									Role:        &featureRoleTypeClient,
									Description: &electricalConnectionDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature10,
										Entity:  []AddressEntityType{addressEntityType2},
									},
									FeatureType: &featureTypeLoadControl,
									Role:        &featureRoleTypeServer,
									SupportedFunction: []FunctionPropertyType{
										{
											Function:           &functionTypeLoadControlLimitDescriptionListData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
										{
											Function:           &functionTypeLoadControlLimitListData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}, Write: &PossibleOperationsWriteType{}},
										},
									},
									Description: &loadControlServerDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature11,
										Entity:  []AddressEntityType{addressEntityType2},
									},
									FeatureType: &featureTypeElectricalConnection,
									Role:        &featureRoleTypeServer,
									SupportedFunction: []FunctionPropertyType{
										{
											Function:           &functionTypeElConDescriptionListData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
									},
									Description: &electricalConnectionDescription,
								},
							},
							{
								Description: &NetworkManagementFeatureDescriptionDataType{
									FeatureAddress: &FeatureAddressType{
										Feature: &feature12,
										Entity:  []AddressEntityType{addressEntityType2},
									},
									FeatureType: &featureTypeDeviceConfiguration,
									Role:        &featureRoleTypeServer,
									SupportedFunction: []FunctionPropertyType{
										{
											Function:           &functionTypeDevConfKeyValueDescListData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
										},
										{
											Function:           &functionTypeDevConfigKeyValueListData,
											PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}, Write: &PossibleOperationsWriteType{}},
										},
									},
									Description: &deviceConfigurationDescription,
								},
							},
						},
					},
				}},
			},
		},
	}
	datagramJSON, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestWriteNodeManagementDetailedDiscoveryData() error = %v", err)
	}
	expectedJSON := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_HEMS"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"device":"d:_i:3210_EVSE"},{"entity":[0]},{"feature":0}]},{"msgCounter":5881},{"msgCounterReference":1},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"nodeManagementDetailedDiscoveryData":[{"specificationVersionList":[{"specificationVersion":["1.2.0"]}]},{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:3210_HEMS"}]},{"deviceType":"EnergyManagementSystem"},{"networkFeatureSet":"smart"}]}]},{"entityInformation":[[{"description":[{"entityAddress":[{"entity":[0]}]},{"entityType":"DeviceInformation"}]}],[{"description":[{"entityAddress":[{"entity":[1]}]},{"entityType":"CEM"},{"description":"CEM Energy Guard"}]}],[{"description":[{"entityAddress":[{"entity":[2]}]},{"entityType":"HeatPumpAppliance"},{"description":"CEM Controllable System"}]}]]},{"featureInformation":[[{"description":[{"featureAddress":[{"entity":[0]},{"feature":0}]},{"featureType":"NodeManagement"},{"role":"special"},{"supportedFunction":[[{"function":"nodeManagementDetailedDiscoveryData"},{"possibleOperations":[{"read":[]}]}],[{"function":"nodeManagementSubscriptionRequestCall"},{"possibleOperations":[]}],[{"function":"nodeManagementBindingRequestCall"},{"possibleOperations":[]}],[{"function":"nodeManagementSubscriptionDeleteCall"},{"possibleOperations":[]}],[{"function":"nodeManagementBindingDeleteCall"},{"possibleOperations":[]}],[{"function":"nodeManagementSubscriptionData"},{"possibleOperations":[{"read":[]}]}],[{"function":"nodeManagementBindingData"},{"possibleOperations":[{"read":[]}]}],[{"function":"nodeManagementUseCaseData"},{"possibleOperations":[{"read":[]}]}]]}]}],[{"description":[{"featureAddress":[{"entity":[0]},{"feature":1}]},{"featureType":"DeviceClassification"},{"role":"server"},{"supportedFunction":[[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]]}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":1}]},{"featureType":"DeviceClassification"},{"role":"client"},{"description":"Device Classification"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":2}]},{"featureType":"DeviceDiagnosis"},{"role":"client"},{"description":"Device Diagnosis"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":3}]},{"featureType":"Measurement"},{"role":"client"},{"description":"Measurement for client"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":4}]},{"featureType":"DeviceConfiguration"},{"role":"client"},{"description":"Device Configuration"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":5}]},{"featureType":"DeviceDiagnosis"},{"role":"server"},{"supportedFunction":[[{"function":"deviceDiagnosisStateData"},{"possibleOperations":[{"read":[]}]}],[{"function":"deviceDiagnosisHeartbeatData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"DeviceDiag"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":7}]},{"featureType":"LoadControl"},{"role":"client"},{"description":"LoadControl client for CEM"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":8}]},{"featureType":"Identification"},{"role":"client"},{"description":"EV identification"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":9}]},{"featureType":"ElectricalConnection"},{"role":"client"},{"description":"Electrical Connection"}]}],[{"description":[{"featureAddress":[{"entity":[2]},{"feature":10}]},{"featureType":"LoadControl"},{"role":"server"},{"supportedFunction":[[{"function":"loadControlLimitDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"loadControlLimitListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Load Control"}]}],[{"description":[{"featureAddress":[{"entity":[2]},{"feature":11}]},{"featureType":"ElectricalConnection"},{"role":"server"},{"supportedFunction":[[{"function":"electricalConnectionDescriptionListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Electrical Connection"}]}],[{"description":[{"featureAddress":[{"entity":[2]},{"feature":12}]},{"featureType":"DeviceConfiguration"},{"role":"server"},{"supportedFunction":[[{"function":"deviceConfigurationKeyValueDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"deviceConfigurationKeyValueListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Device Configuration"}]}]]}]}]]}]}]}`
	expectedDatagram := CmiDatagramType{}
	err = json.Unmarshal(json.RawMessage(expectedJSON), &expectedDatagram)
	if err != nil {
		t.Errorf("TestWriteNodeManagementDetailedDiscoveryData() Unmarshal failed error = %v", err)
	} else {
		if string(datagramJSON) != expectedJSON {
			fmt.Println("EXPECTED:")
			fmt.Println(string(expectedJSON))
			fmt.Println("\nACTUAL:")
			fmt.Println(string(datagramJSON))

			t.Errorf("TestWriteNodeManagementDetailedDiscoveryData() actual json string doesn't match expected result")
		}
	}
}

func TestNodeManagementSubscriptionRequestCall(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.1.1"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 5
	var addressSoureEntity1 AddressEntityType = 1
	var addressDestinationFeature AddressFeatureType = 0
	var addressDestinationEntity1 AddressEntityType = 0
	var msgCounter MsgCounterType = 80
	var cmdClassifier CmdClassifierType = CmdClassifierTypeCall
	var ackRequest bool = true

	var addressServerFeature AddressFeatureType = 5
	var addressServerEntity1 AddressEntityType = 1
	var serverFeatureType FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeDeviceDiagnosis)

	datagram := CmiDatagramType{
		Datagram: DatagramType{
			Header: HeaderType{
				SpecificationVersion: &specificationVersion,
				AddressSource: &FeatureAddressType{
					Feature: &addressSourceFeature,
					Entity:  []AddressEntityType{addressSoureEntity1},
					Device:  &deviceEvse,
				},
				AddressDestination: &FeatureAddressType{
					Feature: &addressDestinationFeature,
					Entity:  []AddressEntityType{addressDestinationEntity1},
					Device:  &deviceHems,
				},
				MsgCounter:    &msgCounter,
				CmdClassifier: &cmdClassifier,
				AckRequest:    &ackRequest,
			},
			Payload: PayloadType{
				Cmd: []CmdType{{
					NodeManagementSubscriptionRequestCall: &NodeManagementSubscriptionRequestCallType{
						SubscriptionRequest: &SubscriptionManagementRequestCallType{
							ClientAddress: &FeatureAddressType{
								Feature: &addressSourceFeature,
								Entity:  []AddressEntityType{addressSoureEntity1},
								Device:  &deviceEvse,
							},
							ServerAddress: &FeatureAddressType{
								Feature: &addressServerFeature,
								Entity:  []AddressEntityType{addressServerEntity1},
								Device:  &deviceHems,
							},
							ServerFeatureType: &serverFeatureType,
						},
					},
				}},
			},
		},
	}
	datagramJSON, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestNodeManagementSubscriptionRequestCall() error = %v", err)
	}

	expectedJSON := `{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1]},{"feature":5}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[0]},{"feature":0}]},{"msgCounter":80},{"cmdClassifier":"call"},{"ackRequest":true}]},{"payload":[{"cmd":[[{"nodeManagementSubscriptionRequestCall":[{"subscriptionRequest":[{"clientAddress":[{"device":"d:_i:3210_EVSE"},{"entity":[1]},{"feature":5}]},{"serverAddress":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":5}]},{"serverFeatureType":"DeviceDiagnosis"}]}]}]]}]}]}`
	expectedDatagram := CmiDatagramType{}
	err = json.Unmarshal(json.RawMessage(expectedJSON), &expectedDatagram)
	if err != nil {
		t.Errorf("TestNodeManagementSubscriptionRequestCall() Unmarshal failed error = %v", err)
	} else {
		if string(datagramJSON) != expectedJSON {
			fmt.Println("EXPECTED:")
			fmt.Println(string(expectedJSON))
			fmt.Println("\nACTUAL:")
			fmt.Println(string(datagramJSON))

			t.Errorf("TestNodeManagementSubscriptionRequestCall() actual json string doesn't match expected result")
		}
	}
}

func TestNodeManagementDetailedDiscoveryDataNotify(t *testing.T) {
	var specificationVersion SpecificationVersionType = "1.1.1"
	var deviceHems AddressDeviceType = "d:_i:3210_HEMS"
	var deviceEvse AddressDeviceType = "d:_i:3210_EVSE"
	var addressSourceFeature AddressFeatureType = 0
	var addressSoureEntity1 AddressEntityType = 0
	var addressDestinationFeature AddressFeatureType = 0
	var addressDestinationEntity1 AddressEntityType = 0
	var msgCounter MsgCounterType = 116
	var cmdClassifier CmdClassifierType = CmdClassifierTypeNotify

	var functionType FunctionType = FunctionType(FunctionEnumTypeNodeManagementDetailedDiscoveryData)

	var entityAddress1 AddressEntityType = 1
	var entityTypeEV EntityTypeType = EntityTypeType(EntityTypeEnumTypeEV)
	var entityLastStateChange NetworkManagementStateChangeType = NetworkManagementStateChangeTypeAdded
	var entityDescription DescriptionType = "Electric Vehicle"

	var featureRoleTypeServer RoleType = RoleTypeServer

	var featureEntityType1 AddressEntityType = 1

	var feature1 AddressFeatureType = 1
	var featureDescription1 DescriptionType = "Load Control"

	var feature2 AddressFeatureType = 2
	var featureDescription2 DescriptionType = "Electrical Connection"

	var feature3 AddressFeatureType = 3
	var featureSpecificUsage3 FeatureSpecificUsageType = FeatureSpecificUsageType(FeatureMeasurementSpecificUsageEnumTypeElectrical)
	var featureDescription3 DescriptionType = "Measurements"

	var feature5 AddressFeatureType = 5
	var featureDescription5 DescriptionType = "Device Configuration EV"

	var feature6 AddressFeatureType = 6
	var featureDescription6 DescriptionType = "Device Classification for EV"

	var feature7 AddressFeatureType = 7
	var featureDescription7 DescriptionType = "Time Series"

	var feature8 AddressFeatureType = 8
	var featureDescription8 DescriptionType = "Incentive Table"

	var feature9 AddressFeatureType = 9
	var featureDescription9 DescriptionType = "Device Diagnosis EV"

	var feature10 AddressFeatureType = 10
	var featureDescription10 DescriptionType = "Identification for EV"

	var featureTypeLoadControl FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeLoadControl)
	var featureTypeElectricalConnection FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeElectricalConnection)
	var featureTypeDeviceConfiguration FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeDeviceConfiguration)
	var featureTypeMeasurement FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeMeasurement)
	var featureTypeDeviceClassification FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeDeviceClassification)
	var featureTypeTimeSeries FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeTimeSeries)
	var featureTypeIncentiveTable FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeIncentiveTable)
	var featureTypeDeviceDiagnosis FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeDeviceDiagnosis)
	var featureTypeIdentification FeatureTypeType = FeatureTypeType(FeatureTypeEnumTypeIdentification)

	var functionLoadControlLimitDescriptionListData FunctionType = FunctionType(FunctionEnumTypeLoadControlLimitDescriptionListData)
	var functionLoadControlLimitListData FunctionType = FunctionType(FunctionEnumTypeLoadControlLimitListData)

	var functionElectricalConnectionParameterDescriptionListData FunctionType = FunctionType(FunctionEnumTypeElectricalConnectionParameterDescriptionListData)
	var functionElectricalConnectionDescriptionListData FunctionType = FunctionType(FunctionEnumTypeElectricalConnectionDescriptionListData)
	var functionElectricalConnectionPermittedValueSetListData FunctionType = FunctionType(FunctionEnumTypeElectricalConnectionPermittedValueSetListData)

	var functionMeasurementListData FunctionType = FunctionType(FunctionEnumTypeMeasurementListData)
	var functionMeasurementDescriptionListData FunctionType = FunctionType(FunctionEnumTypeMeasurementDescriptionListData)

	var functionDeviceConfigurationKeyValueDescriptionListData FunctionType = FunctionType(FunctionEnumTypeDeviceConfigurationKeyValueDescriptionListData)
	var functionDeviceConfigurationKeyValueListData FunctionType = FunctionType(FunctionEnumTypeDeviceConfigurationKeyValueListData)

	var functionDeviceClassificationManufacturerData FunctionType = FunctionType(FunctionEnumTypeDeviceClassificationManufacturerData)

	var functionTimeSeriesConstraintsListData FunctionType = FunctionType(FunctionEnumTypeTimeSeriesConstraintsListData)
	var functionTimeSeriesDescriptionListData FunctionType = FunctionType(FunctionEnumTypeTimeSeriesDescriptionListData)
	var functionTimeSeriesListData FunctionType = FunctionType(FunctionEnumTypeTimeSeriesListData)

	var functionIncentiveTableConstraintsData FunctionType = FunctionType(FunctionEnumTypeIncentiveTableConstraintsData)
	var functionIncentiveTableData FunctionType = FunctionType(FunctionEnumTypeIncentiveTableData)
	var functionIncentiveTableDescriptionData FunctionType = FunctionType(FunctionEnumTypeIncentiveTableDescriptionData)

	var functionDeviceDiagnosisStateData FunctionType = FunctionType(FunctionEnumTypeDeviceDiagnosisStateData)

	var functionIdentificationListData FunctionType = FunctionType(FunctionEnumTypeIdentificationListData)

	datagram := CmiDatagramType{
		Datagram: DatagramType{
			Header: HeaderType{
				SpecificationVersion: &specificationVersion,
				AddressSource: &FeatureAddressType{
					Feature: &addressSourceFeature,
					Entity:  []AddressEntityType{addressSoureEntity1},
					Device:  &deviceEvse,
				},
				AddressDestination: &FeatureAddressType{
					Feature: &addressDestinationFeature,
					Entity:  []AddressEntityType{addressDestinationEntity1},
					Device:  &deviceHems,
				},
				MsgCounter:    &msgCounter,
				CmdClassifier: &cmdClassifier,
			},
			Payload: PayloadType{
				Cmd: []CmdType{
					{
						Function: &functionType,
						Filter: []FilterType{
							{
								CmdControl: &CmdControlType{
									Partial: &ElementTagType{},
								},
							},
						},
						NodeManagementDetailedDiscoveryData: &NodeManagementDetailedDiscoveryDataType{
							DeviceInformation: &NodeManagementDetailedDiscoveryDeviceInformationType{
								Description: &NetworkManagementDeviceDescriptionDataType{
									DeviceAddress: &DeviceAddressType{
										Device: &deviceEvse,
									},
								},
							},
							EntityInformation: []NodeManagementDetailedDiscoveryEntityInformationType{
								{
									Description: &NetworkManagementEntityDescriptionDataType{
										EntityAddress:   &EntityAddressType{Entity: []AddressEntityType{entityAddress1, entityAddress1}},
										EntityType:      &entityTypeEV,
										LastStateChange: &entityLastStateChange,
										Description:     &entityDescription,
									},
								},
							},
							FeatureInformation: []NodeManagementDetailedDiscoveryFeatureInformationType{
								{
									Description: &NetworkManagementFeatureDescriptionDataType{
										FeatureAddress: &FeatureAddressType{
											Entity:  []AddressEntityType{featureEntityType1, featureEntityType1},
											Feature: &feature1,
										},
										FeatureType: &featureTypeLoadControl,
										Role:        &featureRoleTypeServer,
										SupportedFunction: []FunctionPropertyType{
											{
												Function:           &functionLoadControlLimitDescriptionListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
											{
												Function:           &functionLoadControlLimitListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}, Write: &PossibleOperationsWriteType{}},
											},
										},
										Description: &featureDescription1,
									},
								},
								{
									Description: &NetworkManagementFeatureDescriptionDataType{
										FeatureAddress: &FeatureAddressType{
											Entity:  []AddressEntityType{featureEntityType1, featureEntityType1},
											Feature: &feature2,
										},
										FeatureType: &featureTypeElectricalConnection,
										Role:        &featureRoleTypeServer,
										SupportedFunction: []FunctionPropertyType{
											{
												Function:           &functionElectricalConnectionParameterDescriptionListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
											{
												Function:           &functionElectricalConnectionDescriptionListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
											{
												Function:           &functionElectricalConnectionPermittedValueSetListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
										},
										Description: &featureDescription2,
									},
								},
								{
									Description: &NetworkManagementFeatureDescriptionDataType{
										FeatureAddress: &FeatureAddressType{
											Entity:  []AddressEntityType{featureEntityType1, featureEntityType1},
											Feature: &feature3,
										},
										FeatureType: &featureTypeMeasurement,
										SpecificUsage: []FeatureSpecificUsageType{
											featureSpecificUsage3,
										},
										Role: &featureRoleTypeServer,
										SupportedFunction: []FunctionPropertyType{
											{
												Function:           &functionMeasurementListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
											{
												Function:           &functionMeasurementDescriptionListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
										},
										Description: &featureDescription3,
									},
								},
								{
									Description: &NetworkManagementFeatureDescriptionDataType{
										FeatureAddress: &FeatureAddressType{
											Entity:  []AddressEntityType{featureEntityType1, featureEntityType1},
											Feature: &feature5,
										},
										FeatureType: &featureTypeDeviceConfiguration,
										Role:        &featureRoleTypeServer,
										SupportedFunction: []FunctionPropertyType{
											{
												Function:           &functionDeviceConfigurationKeyValueDescriptionListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
											{
												Function:           &functionDeviceConfigurationKeyValueListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
										},
										Description: &featureDescription5,
									},
								},
								{
									Description: &NetworkManagementFeatureDescriptionDataType{
										FeatureAddress: &FeatureAddressType{
											Entity:  []AddressEntityType{featureEntityType1, featureEntityType1},
											Feature: &feature6,
										},
										FeatureType: &featureTypeDeviceClassification,
										Role:        &featureRoleTypeServer,
										SupportedFunction: []FunctionPropertyType{
											{
												Function:           &functionDeviceClassificationManufacturerData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
										},
										Description: &featureDescription6,
									},
								},
								{
									Description: &NetworkManagementFeatureDescriptionDataType{
										FeatureAddress: &FeatureAddressType{
											Entity:  []AddressEntityType{featureEntityType1, featureEntityType1},
											Feature: &feature7,
										},
										FeatureType: &featureTypeTimeSeries,
										Role:        &featureRoleTypeServer,
										SupportedFunction: []FunctionPropertyType{
											{
												Function:           &functionTimeSeriesConstraintsListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
											{
												Function:           &functionTimeSeriesDescriptionListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
											{
												Function:           &functionTimeSeriesListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}, Write: &PossibleOperationsWriteType{}},
											},
										},
										Description: &featureDescription7,
									},
								},
								{
									Description: &NetworkManagementFeatureDescriptionDataType{
										FeatureAddress: &FeatureAddressType{
											Entity:  []AddressEntityType{featureEntityType1, featureEntityType1},
											Feature: &feature8,
										},
										FeatureType: &featureTypeIncentiveTable,
										Role:        &featureRoleTypeServer,
										SupportedFunction: []FunctionPropertyType{
											{
												Function:           &functionIncentiveTableConstraintsData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
											{
												Function:           &functionIncentiveTableData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}, Write: &PossibleOperationsWriteType{}},
											},
											{
												Function:           &functionIncentiveTableDescriptionData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}, Write: &PossibleOperationsWriteType{}},
											},
										},
										Description: &featureDescription8,
									},
								},
								{
									Description: &NetworkManagementFeatureDescriptionDataType{
										FeatureAddress: &FeatureAddressType{
											Entity:  []AddressEntityType{featureEntityType1, featureEntityType1},
											Feature: &feature9,
										},
										FeatureType: &featureTypeDeviceDiagnosis,
										Role:        &featureRoleTypeServer,
										SupportedFunction: []FunctionPropertyType{
											{
												Function:           &functionDeviceDiagnosisStateData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
										},
										Description: &featureDescription9,
									},
								},
								{
									Description: &NetworkManagementFeatureDescriptionDataType{
										FeatureAddress: &FeatureAddressType{
											Entity:  []AddressEntityType{featureEntityType1, featureEntityType1},
											Feature: &feature10,
										},
										FeatureType: &featureTypeIdentification,
										Role:        &featureRoleTypeServer,
										SupportedFunction: []FunctionPropertyType{
											{
												Function:           &functionIdentificationListData,
												PossibleOperations: &PossibleOperationsType{Read: &PossibleOperationsReadType{}},
											},
										},
										Description: &featureDescription10,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	datagramJSON, err := json.Marshal(datagram)
	if err != nil {
		t.Errorf("TestNodeManagementSubscriptionRequestCall() error = %v", err)
	}

	expectedJSON := `{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[0]},{"feature":0}]},{"msgCounter":116},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"function":"nodeManagementDetailedDiscoveryData"},{"filter":[[{"cmdControl":[{"partial":[]}]}]]},{"nodeManagementDetailedDiscoveryData":[{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:3210_EVSE"}]}]}]},{"entityInformation":[[{"description":[{"entityAddress":[{"entity":[1,1]}]},{"entityType":"EV"},{"lastStateChange":"added"},{"description":"Electric Vehicle"}]}]]},{"featureInformation":[[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":1}]},{"featureType":"LoadControl"},{"role":"server"},{"supportedFunction":[[{"function":"loadControlLimitDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"loadControlLimitListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Load Control"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":2}]},{"featureType":"ElectricalConnection"},{"role":"server"},{"supportedFunction":[[{"function":"electricalConnectionParameterDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionPermittedValueSetListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Electrical Connection"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":3}]},{"featureType":"Measurement"},{"specificUsage":["Electrical"]},{"role":"server"},{"supportedFunction":[[{"function":"measurementListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"measurementDescriptionListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Measurements"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":5}]},{"featureType":"DeviceConfiguration"},{"role":"server"},{"supportedFunction":[[{"function":"deviceConfigurationKeyValueDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"deviceConfigurationKeyValueListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Configuration EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":6}]},{"featureType":"DeviceClassification"},{"role":"server"},{"supportedFunction":[[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Classification for EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":7}]},{"featureType":"TimeSeries"},{"role":"server"},{"supportedFunction":[[{"function":"timeSeriesConstraintsListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"timeSeriesDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"timeSeriesListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Time Series"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":8}]},{"featureType":"IncentiveTable"},{"role":"server"},{"supportedFunction":[[{"function":"incentiveTableConstraintsData"},{"possibleOperations":[{"read":[]}]}],[{"function":"incentiveTableData"},{"possibleOperations":[{"read":[]},{"write":[]}]}],[{"function":"incentiveTableDescriptionData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Incentive Table"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":9}]},{"featureType":"DeviceDiagnosis"},{"role":"server"},{"supportedFunction":[[{"function":"deviceDiagnosisStateData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Diagnosis EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":10}]},{"featureType":"Identification"},{"role":"server"},{"supportedFunction":[[{"function":"identificationListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Identification for EV"}]}]]}]}]]}]}]}`
	expectedDatagram := CmiDatagramType{}
	err = json.Unmarshal(json.RawMessage(expectedJSON), &expectedDatagram)
	if err != nil {
		t.Errorf("TestNodeManagementSubscriptionRequestCall() Unmarshal failed error = %v", err)
	} else {
		if string(datagramJSON) != expectedJSON {
			fmt.Println("EXPECTED:")
			fmt.Println(string(expectedJSON))
			fmt.Println("\nACTUAL:")
			fmt.Println(string(datagramJSON))

			t.Errorf("TestNodeManagementSubscriptionRequestCall() actual json string doesn't match expected result")
		}
	}
}
