package ship

// func (c *Client) dataHandshake() error {
// 	var specificationVersion model.SpecificationVersionType = "1.2.0"
// 	var device model.AddressDeviceType = "d:_i:3210_HEMS"
// 	var target model.AddressDeviceType = "19667_PorscheEVSE-00016544"

// 	var entity0 model.AddressEntityType = 0
// 	var feature0 model.AddressFeatureType = 0

// 	var msgCounter model.MsgCounterType = 5876
// 	var cmdClassifier model.CmdClassifierType = model.CmdClassifierTypeRead

// 	datagram := model.CmiDatagramType{
// 		Datagram: model.DatagramType{
// 			Header: model.HeaderType{
// 				SpecificationVersion: &specificationVersion,
// 				AddressSource: &model.FeatureAddressType{
// 					Device:  &device,
// 					Entity:  []model.AddressEntityType{entity0},
// 					Feature: &feature0,
// 				},
// 				AddressDestination: &model.FeatureAddressType{
// 					Entity:  []model.AddressEntityType{entity0},
// 					Feature: &feature0,
// 				},
// 				MsgCounter:    &msgCounter,
// 				CmdClassifier: &cmdClassifier,
// 			},
// 			Payload: model.PayloadType{
// 				Cmd: []model.CmdType{{
// 					NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
// 				}},
// 			},
// 		},
// 	}
// 	discovery, err := json.Marshal(datagram)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(string(discovery))

// 	hs := ship.CmiData{
// 		Data: ship.Data{
// 			Header: ship.HeaderType{
// 				ProtocolId: ship.ProtocolIdType(message.ProtocolID),
// 			},
// 			Payload: json.RawMessage(discovery),
// 		},
// 	}
// 	err = c.t.WriteJSON(message.CmiTypeData, hs)
// 	if err != nil {
// 		return err
// 	}

// 	smart := model.NetworkManagementFeatureSetTypeSmart

// 	msgCounter = 5881
// 	var msgCounterReference model.MsgCounterType = 1
// 	cmdClassifier = model.CmdClassifierTypeReply

// 	var entity1 model.AddressEntityType = 1
// 	var entity2 model.AddressEntityType = 2

// 	var feature1 model.AddressFeatureType = 1
// 	var feature2 model.AddressFeatureType = 2
// 	var feature3 model.AddressFeatureType = 3
// 	var feature4 model.AddressFeatureType = 4
// 	var feature5 model.AddressFeatureType = 5
// 	var feature7 model.AddressFeatureType = 7
// 	var feature8 model.AddressFeatureType = 8
// 	var feature9 model.AddressFeatureType = 9
// 	var feature10 model.AddressFeatureType = 10
// 	var feature11 model.AddressFeatureType = 11
// 	var feature12 model.AddressFeatureType = 12

// 	var cemEnergyGuardDescription model.DescriptionType = "CEM Energy Guard"
// 	var cemControllableSystemDescription model.DescriptionType = "CEM Controllable System"
// 	var deviceClassificationDescription model.DescriptionType = "Device Classification"
// 	var deviceDiagnosisDescription model.DescriptionType = "Device Diagnosis"
// 	var deviceDiagnosisServerDescription model.DescriptionType = "DeviceDiag"
// 	var measurementForClientDescription model.DescriptionType = "Measurement for client"
// 	var deviceConfigurationDescription model.DescriptionType = "Device Configuration"
// 	var loadControlDescription model.DescriptionType = "LoadControl client for CEM"
// 	var loadControlServerDescription model.DescriptionType = "Load Control"
// 	var evIdentificationDescription model.DescriptionType = "EV identification"
// 	var electricalConnectionDescription model.DescriptionType = "Electrical Connection"

// 	var featureRoleTypeSpecial model.RoleType = model.RoleTypeSpecial
// 	var featureRoleTypeServer model.RoleType = model.RoleTypeServer
// 	var featureRoleTypeClient model.RoleType = model.RoleTypeClient

// 	var specificationVersionDataType model.SpecificationVersionDataType = model.SpecificationVersionDataType(specificationVersion)
// 	var deviceTypeEMS model.DeviceTypeType = model.DeviceTypeType(model.DeviceTypeEnumTypeEnergyManagementSystem)

// 	var entityTypeDeviceInformation model.EntityTypeType = model.EntityTypeType(model.EntityTypeEnumTypeDeviceInformation)
// 	var entityTypeCem model.EntityTypeType = model.EntityTypeType(model.EntityTypeEnumTypeCEM)
// 	var entityTypeHeatpump model.EntityTypeType = model.EntityTypeType(model.EntityTypeEnumTypeHeatPumpAppliance)

// 	var featureTypeNodeMgmt model.FeatureTypeType = model.FeatureTypeType(model.FeatureTypeEnumTypeNodeManagement)
// 	var functionTypeNodeMgmtDetailedDiscovery model.FunctionType = model.FunctionType(model.FunctionEnumTypeNodeManagementDetailedDiscoveryData)
// 	var functionTypeNodeMgmtSubReqCall model.FunctionType = model.FunctionType(model.FunctionEnumTypeNodeManagementSubscriptionRequestCall)
// 	var functionTypeNodeMgmtBindReqCall model.FunctionType = model.FunctionType(model.FunctionEnumTypeNodeManagementBindingRequestCall)
// 	var functionTypeNodeMgmtSubDelCall model.FunctionType = model.FunctionType(model.FunctionEnumTypeNodeManagementSubscriptionDeleteCall)
// 	var functionTypeNodeMgmtBindDelCall model.FunctionType = model.FunctionType(model.FunctionEnumTypeNodeManagementBindingDeleteCall)
// 	var functionTypeNodeMgmtSubData model.FunctionType = model.FunctionType(model.FunctionEnumTypeNodeManagementSubscriptionData)
// 	var functionTypeNodeMgmtBindData model.FunctionType = model.FunctionType(model.FunctionEnumTypeNodeManagementBindingData)
// 	var functionTypeNodeMgmtUseCaseData model.FunctionType = model.FunctionType(model.FunctionEnumTypeNodeManagementUseCaseData)

// 	var featureTypeDeviceClassification model.FeatureTypeType = model.FeatureTypeType(model.FeatureTypeEnumTypeDeviceClassification)
// 	var functionDeviceClassificationManufacturerData model.FunctionType = model.FunctionType(model.FunctionEnumTypeDeviceClassificationManufacturerDdata)

// 	var featureTypeDeviceDiagnosis model.FeatureTypeType = model.FeatureTypeType(model.FeatureTypeEnumTypeDeviceDiagnosis)
// 	var featureTypeMeasurement model.FeatureTypeType = model.FeatureTypeType(model.FeatureTypeEnumTypeMeasurement)

// 	var featureTypeDeviceConfiguration model.FeatureTypeType = model.FeatureTypeType(model.FeatureTypeEnumTypeDeviceConfiguration)
// 	var functionDeviceConfigurationKeyValueDescriptionListData model.FunctionType = model.FunctionType(model.FunctionEnumTypeDeviceConfigurationKeyValueDescriptionListData)
// 	var functionDeviceConfigurationKeyValueListData model.FunctionType = model.FunctionType(model.FunctionEnumTypeDeviceConfigurationKeyValueListData)

// 	var functionDeviceDiagnosisStateData model.FunctionType = model.FunctionType(model.FunctionEnumTypeDeviceDiagnosisStateData)
// 	var functionDeviceDiagnosisHeartBeatData model.FunctionType = model.FunctionType(model.FunctionEnumTypeDeviceDiagnosisHeartBeatData)

// 	var featureTypeLoadControl model.FeatureTypeType = model.FeatureTypeType(model.FeatureTypeEnumTypeLoadControl)
// 	var functionLoadControlLimitDescriptionListData model.FunctionType = model.FunctionType(model.FunctionEnumTypeLoadControlLimitDescriptionListData)
// 	var functionLoadControlLimitListData model.FunctionType = model.FunctionType(model.FunctionEnumTypeLoadControlLimitListData)

// 	var featureTypeIdentification model.FeatureTypeType = model.FeatureTypeType(model.FeatureTypeEnumTypeIdentification)

// 	var featureTypeElectricalConnection model.FeatureTypeType = model.FeatureTypeType(model.FeatureTypeEnumTypeElectricalConnection)
// 	var functionElectricalConnectionDescriptionListData model.FunctionType = model.FunctionType(model.FunctionEnumTypeElectricalConnectionDescriptionListData)

// 	datagram = model.CmiDatagramType{
// 		Datagram: model.DatagramType{
// 			Header: model.HeaderType{
// 				SpecificationVersion: &specificationVersion,
// 				AddressSource: &model.FeatureAddressType{
// 					Feature: &feature0,
// 					Entity:  []model.AddressEntityType{entity0},
// 					Device:  &device,
// 				},
// 				AddressDestination: &model.FeatureAddressType{
// 					Feature: &feature0,
// 					Entity:  []model.AddressEntityType{entity0},
// 					Device:  &target,
// 				},
// 				MsgCounter:          &msgCounter,
// 				MsgCounterReference: &msgCounterReference,
// 				CmdClassifier:       &cmdClassifier,
// 			},
// 			Payload: model.PayloadType{
// 				Cmd: []model.CmdType{{
// 					NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{
// 						SpecificationVersionList: &model.NodeManagementSpecificationVersionListType{
// 							SpecificationVersion: []model.SpecificationVersionDataType{specificationVersionDataType},
// 						},
// 						DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
// 							Description: &model.NetworkManagementDeviceDescriptionDataType{
// 								DeviceAddress: &model.DeviceAddressType{
// 									Device: &device,
// 								},
// 								DeviceType:        &deviceTypeEMS,
// 								NetworkFeatureSet: &smart,
// 							},
// 						},
// 						EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{
// 							{
// 								Description: &model.NetworkManagementEntityDescriptionDataType{
// 									EntityAddress: &model.EntityAddressType{Entity: []model.AddressEntityType{entity0}},
// 									EntityType:    &entityTypeDeviceInformation,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementEntityDescriptionDataType{
// 									EntityAddress: &model.EntityAddressType{Entity: []model.AddressEntityType{entity1}},
// 									EntityType:    &entityTypeCem,
// 									Description:   &cemEnergyGuardDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementEntityDescriptionDataType{
// 									EntityAddress: &model.EntityAddressType{Entity: []model.AddressEntityType{entity2}},
// 									EntityType:    &entityTypeHeatpump,
// 									Description:   &cemControllableSystemDescription,
// 								},
// 							},
// 						},
// 						FeatureInformation: []model.NodeManagementDetailedDiscoveryFeatureInformationType{
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature0,
// 										Entity:  []model.AddressEntityType{entity0},
// 									},
// 									FeatureType: &featureTypeNodeMgmt,
// 									Role:        &featureRoleTypeSpecial,
// 									SupportedFunction: []model.FunctionPropertyType{
// 										{
// 											Function:           &functionTypeNodeMgmtDetailedDiscovery,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 										{
// 											Function:           &functionTypeNodeMgmtSubReqCall,
// 											PossibleOperations: &model.PossibleOperationsType{},
// 										},
// 										{
// 											Function:           &functionTypeNodeMgmtBindReqCall,
// 											PossibleOperations: &model.PossibleOperationsType{},
// 										},
// 										{
// 											Function:           &functionTypeNodeMgmtSubDelCall,
// 											PossibleOperations: &model.PossibleOperationsType{},
// 										},
// 										{
// 											Function:           &functionTypeNodeMgmtBindDelCall,
// 											PossibleOperations: &model.PossibleOperationsType{},
// 										},
// 										{
// 											Function:           &functionTypeNodeMgmtSubData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 										{
// 											Function:           &functionTypeNodeMgmtBindData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 										{
// 											Function:           &functionTypeNodeMgmtUseCaseData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 									},
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature1,
// 										Entity:  []model.AddressEntityType{entity0},
// 									},
// 									FeatureType: &featureTypeDeviceClassification,
// 									Role:        &featureRoleTypeServer,
// 									SupportedFunction: []model.FunctionPropertyType{
// 										{
// 											Function:           &functionDeviceClassificationManufacturerData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 									},
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature1,
// 										Entity:  []model.AddressEntityType{entity1},
// 									},
// 									FeatureType: &featureTypeDeviceClassification,
// 									Role:        &featureRoleTypeClient,
// 									Description: &deviceClassificationDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature2,
// 										Entity:  []model.AddressEntityType{entity1},
// 									},
// 									FeatureType: &featureTypeDeviceDiagnosis,
// 									Role:        &featureRoleTypeClient,
// 									Description: &deviceDiagnosisDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature3,
// 										Entity:  []model.AddressEntityType{entity1},
// 									},
// 									FeatureType: &featureTypeMeasurement,
// 									Role:        &featureRoleTypeClient,
// 									Description: &measurementForClientDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature4,
// 										Entity:  []model.AddressEntityType{entity1},
// 									},
// 									FeatureType: &featureTypeDeviceConfiguration,
// 									Role:        &featureRoleTypeClient,
// 									Description: &deviceConfigurationDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature5,
// 										Entity:  []model.AddressEntityType{entity1},
// 									},
// 									FeatureType: &featureTypeDeviceDiagnosis,
// 									Role:        &featureRoleTypeServer,
// 									SupportedFunction: []model.FunctionPropertyType{
// 										{
// 											Function:           &functionDeviceDiagnosisStateData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 										{
// 											Function:           &functionDeviceDiagnosisHeartBeatData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 									},
// 									Description: &deviceDiagnosisServerDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature7,
// 										Entity:  []model.AddressEntityType{entity1},
// 									},
// 									FeatureType: &featureTypeLoadControl,
// 									Role:        &featureRoleTypeClient,
// 									Description: &loadControlDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature8,
// 										Entity:  []model.AddressEntityType{entity1},
// 									},
// 									FeatureType: &featureTypeIdentification,
// 									Role:        &featureRoleTypeClient,
// 									Description: &evIdentificationDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature9,
// 										Entity:  []model.AddressEntityType{entity1},
// 									},
// 									FeatureType: &featureTypeElectricalConnection,
// 									Role:        &featureRoleTypeClient,
// 									Description: &electricalConnectionDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature10,
// 										Entity:  []model.AddressEntityType{entity2},
// 									},
// 									FeatureType: &featureTypeLoadControl,
// 									Role:        &featureRoleTypeServer,
// 									SupportedFunction: []model.FunctionPropertyType{
// 										{
// 											Function:           &functionLoadControlLimitDescriptionListData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 										{
// 											Function:           &functionLoadControlLimitListData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}, Write: &model.PossibleOperationsWriteType{}},
// 										},
// 									},
// 									Description: &loadControlServerDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature11,
// 										Entity:  []model.AddressEntityType{entity2},
// 									},
// 									FeatureType: &featureTypeElectricalConnection,
// 									Role:        &featureRoleTypeServer,
// 									SupportedFunction: []model.FunctionPropertyType{
// 										{
// 											Function:           &functionElectricalConnectionDescriptionListData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 									},
// 									Description: &electricalConnectionDescription,
// 								},
// 							},
// 							{
// 								Description: &model.NetworkManagementFeatureDescriptionDataType{
// 									FeatureAddress: &model.FeatureAddressType{
// 										Feature: &feature12,
// 										Entity:  []model.AddressEntityType{entity2},
// 									},
// 									FeatureType: &featureTypeDeviceConfiguration,
// 									Role:        &featureRoleTypeServer,
// 									SupportedFunction: []model.FunctionPropertyType{
// 										{
// 											Function:           &functionDeviceConfigurationKeyValueDescriptionListData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}},
// 										},
// 										{
// 											Function:           &functionDeviceConfigurationKeyValueListData,
// 											PossibleOperations: &model.PossibleOperationsType{Read: &model.PossibleOperationsReadType{}, Write: &model.PossibleOperationsWriteType{}},
// 										},
// 									},
// 									Description: &deviceConfigurationDescription,
// 								},
// 							},
// 						},
// 					},
// 				}},
// 			},
// 		},
// 	}

// 	services, err := json.Marshal(datagram)
// 	if err != nil {
// 		var servicesSent bool
// 		for err == nil {
// 			err = c.t.DataReceive()
// 			if err == nil && !servicesSent {
// 				hs := ship.CmiData{
// 					Data: ship.Data{
// 						Header: ship.HeaderType{
// 							ProtocolId: ship.ProtocolIdType(message.ProtocolID),
// 						},
// 						Payload: json.RawMessage(services),
// 					},
// 				}

// 				err = c.t.WriteJSON(message.CmiTypeData, hs)
// 				servicesSent = true
// 			}
// 		}
// 	}

// 	return err
// }
