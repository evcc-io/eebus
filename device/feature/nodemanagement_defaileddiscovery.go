package feature

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/evcc-io/eebus/device"
	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

func (f *NodeManagement) readDetailedDiscoveryData(ctrl spine.Context, data model.NodeManagementDetailedDiscoveryDataType) error {
	localDevice := f.Entity.GetDevice()

	var entityInformation []model.NodeManagementDetailedDiscoveryEntityInformationType
	var featureInformation []model.NodeManagementDetailedDiscoveryFeatureInformationType

	for _, e := range localDevice.GetEntities() {
		entityInformation = append(entityInformation, *e.Information())

		for _, f := range e.GetFeatures() {
			featureInformation = append(featureInformation, *f.Information())
		}
	}

	res := model.CmdType{
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{
			SpecificationVersionList: &model.NodeManagementSpecificationVersionListType{
				SpecificationVersion: []model.SpecificationVersionDataType{device.SpecificationVersion},
			},
			DeviceInformation:  localDevice.Information(),
			EntityInformation:  entityInformation,
			FeatureInformation: featureInformation,
		},
	}

	return ctrl.Reply(model.CmdClassifierTypeReply, res)
}

func (f *NodeManagement) checkEntityInformation(entity model.NodeManagementDetailedDiscoveryEntityInformationType) error {
	description := entity.Description
	if description == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description")
	}

	if description.EntityAddress == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress")
	}

	if description.EntityAddress.Entity == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress.Entity")
	}

	return nil
}

func (f *NodeManagement) addEntityAndFeatures(ctrl spine.Context, remoteDevice spine.Device, data model.NodeManagementDetailedDiscoveryDataType) error {
	for _, ei := range data.EntityInformation {
		if err := f.checkEntityInformation(ei); err != nil {
			return err
		}

		entityAddress := ei.Description.EntityAddress.Entity

		entity := remoteDevice.Entity(entityAddress)
		if entity == nil {
			entity = spine.UnmarshalEntity(remoteDevice.GetAddress(), ei)
		}

		for _, fi := range data.FeatureInformation {
			if reflect.DeepEqual(fi.Description.FeatureAddress.Entity, entityAddress) {
				if f := spine.UnmarshalFeature(fi); f != nil {
					entity.Add(f)
				}
			}
		}

		remoteDevice.Add(entity)

		if err := f.announceFeatureDiscovery(ctrl, entity); err != nil {
			return err
		}
	}

	return nil
}

func (f *NodeManagement) replyDetailedDiscoveryData(ctrl spine.Context, data model.NodeManagementDetailedDiscoveryDataType) error {
	remoteDevice := ctrl.GetDevice()

	deviceDescription := data.DeviceInformation.Description
	if deviceDescription == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid DeviceInformation.Description")
	}

	// create remote device if not existing
	// TODO handle partial updates
	if remoteDevice == nil {
		remoteDevice = spine.UnmarshalDevice(data)
	}

	if err := f.addEntityAndFeatures(ctrl, remoteDevice, data); err != nil {
		return err
	}

	ctrl.SetDevice(remoteDevice)
	ctrl.UpdateDevice(model.NetworkManagementStateChangeTypeAdded)

	return nil
}

func (f *NodeManagement) notifyDetailedDiscoveryData(ctrl spine.Context, data model.NodeManagementDetailedDiscoveryDataType, isPartialForCmd bool) error {
	// is this a partial request?
	if !isPartialForCmd {
		return errors.New("the received NodeManagementDetailedDiscovery.notify dataset should be partial")
	}

	if data.EntityInformation == nil || len(data.EntityInformation) == 0 || data.EntityInformation[0].Description == nil || data.EntityInformation[0].Description.LastStateChange == nil {
		return errors.New("the received NodeManagementDetailedDiscovery.notify dataset is incomplete")
	}

	lastStateChange := *data.EntityInformation[0].Description.LastStateChange
	remoteDevice := ctrl.GetDevice()

	// addition exmaple:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[0]},{"feature":0}]},{"msgCounter":6023},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"function":"nodeManagementDetailedDiscoveryData"},{"filter":[[{"cmdControl":[{"partial":[]}]}]]},{"nodeManagementDetailedDiscoveryData":[{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:19667_PorscheEVSE-00016544"}]}]}]},{"entityInformation":[[{"description":[{"entityAddress":[{"entity":[1,1]}]},{"entityType":"EV"},{"lastStateChange":"added"},{"description":"Electric Vehicle"}]}]]},{"featureInformation":[[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":1}]},{"featureType":"LoadControl"},{"role":"server"},{"supportedFunction":[[{"function":"loadControlLimitDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"loadControlLimitListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Load Control"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":2}]},{"featureType":"ElectricalConnection"},{"role":"server"},{"supportedFunction":[[{"function":"electricalConnectionParameterDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionPermittedValueSetListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Electrical Connection"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":3}]},{"featureType":"Measurement"},{"specificUsage":["Electrical"]},{"role":"server"},{"supportedFunction":[[{"function":"measurementListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"measurementDescriptionListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Measurements"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":5}]},{"featureType":"DeviceConfiguration"},{"role":"server"},{"supportedFunction":[[{"function":"deviceConfigurationKeyValueDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"deviceConfigurationKeyValueListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Configuration EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":6}]},{"featureType":"DeviceClassification"},{"role":"server"},{"supportedFunction":[[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Classification for EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":7}]},{"featureType":"TimeSeries"},{"role":"server"},{"supportedFunction":[[{"function":"timeSeriesConstraintsListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"timeSeriesDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"timeSeriesListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Time Series"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":8}]},{"featureType":"IncentiveTable"},{"role":"server"},{"supportedFunction":[[{"function":"incentiveTableConstraintsData"},{"possibleOperations":[{"read":[]}]}],[{"function":"incentiveTableData"},{"possibleOperations":[{"read":[]},{"write":[]}]}],[{"function":"incentiveTableDescriptionData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Incentive Table"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":9}]},{"featureType":"DeviceDiagnosis"},{"role":"server"},{"supportedFunction":[[{"function":"deviceDiagnosisStateData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Diagnosis EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":10}]},{"featureType":"Identification"},{"role":"server"},{"supportedFunction":[[{"function":"identificationListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Identification for EV"}]}]]}]}]]}]}]}}]}
	// {"cmd":[[
	// 	{"function":"nodeManagementDetailedDiscoveryData"},
	// 	{"filter":[[{"cmdControl":[{"partial":[]}]}]]},
	// 	{"nodeManagementDetailedDiscoveryData":[
	// 		{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:19667_PorscheEVSE-00016544"}]}]}]},
	// 		{"entityInformation":[[
	// 			{"description":[
	// 				{"entityAddress":[{"entity":[1,1]}]},{"entityType":"EV"},
	// 				{"lastStateChange":"added"},
	// 				{"description":"Electric Vehicle"}]}
	// 		]]},
	// 		{"featureInformation":[
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":1}]},{"featureType":"LoadControl"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"loadControlLimitDescriptionListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"loadControlLimitListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]
	// 				]},
	// 				{"description":"Load Control"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":2}]},{"featureType":"ElectricalConnection"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"electricalConnectionParameterDescriptionListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"electricalConnectionDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionPermittedValueSetListData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Electrical Connection"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":3}]},{"featureType":"Measurement"},{"specificUsage":["Electrical"]},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"measurementListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"measurementDescriptionListData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Measurements"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":5}]},{"featureType":"DeviceConfiguration"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"deviceConfigurationKeyValueDescriptionListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"deviceConfigurationKeyValueListData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Device Configuration EV"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":6}]},{"featureType":"DeviceClassification"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Device Classification for EV"}]
	// 			}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":7}]},{"featureType":"TimeSeries"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"timeSeriesConstraintsListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"timeSeriesDescriptionListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"timeSeriesListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]
	// 				]},
	// 				{"description":"Time Series"}]
	// 			}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":8}]},{"featureType":"IncentiveTable"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"incentiveTableConstraintsData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"incentiveTableData"},{"possibleOperations":[{"read":[]},{"write":[]}]}],
	// 					[{"function":"incentiveTableDescriptionData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]
	// 				]},
	// 				{"description":"Incentive Table"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":9}]},{"featureType":"DeviceDiagnosis"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"deviceDiagnosisStateData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Device Diagnosis EV"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":10}]},{"featureType":"Identification"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"identificationListData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Identification for EV"}
	// 			]}]
	// 		]}
	// 	]}
	// ]]}

	// is this addition?
	if lastStateChange == model.NetworkManagementStateChangeTypeAdded {
		if err := f.addEntityAndFeatures(ctrl, remoteDevice, data); err != nil {
			return err
		}

		ctrl.UpdateDevice(lastStateChange)
	}

	// removal example:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[0]},{"feature":0}]},{"msgCounter":4835},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"function":"nodeManagementDetailedDiscoveryData"},{"filter":[[{"cmdControl":[{"partial":[]}]}]]},{"nodeManagementDetailedDiscoveryData":[{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:19667_PorscheEVSE-00016544"}]}]}]},{"entityInformation":[[{"description":[{"entityAddress":[{"entity":[1,1]}]},{"lastStateChange":"removed"}]}]]}]}]]}]}]}}]}
	// {
	// 	"cmd": [[
	// 			{"function": "nodeManagementDetailedDiscoveryData"},
	// 			{"filter": [[{"cmdControl": [{"partial": []}]}]]},
	// 			{"nodeManagementDetailedDiscoveryData": [
	// 					{"deviceInformation": [{"description": [{"deviceAddress": [{"device": "d:_i:19667_PorscheEVSE-00016544"}]}]}]},
	// 					{"entityInformation": [[
	// 							{
	// 								"description": [
	// 									{"entityAddress": [{"entity": [1,1]}]},
	// 									{"lastStateChange": "removed"}
	// 								]
	// 							}
	// 						]]
	// 					}
	// 				]
	// 			}
	// 	]]
	// }

	// is this removal?
	if lastStateChange == model.NetworkManagementStateChangeTypeRemoved {
		for _, ei := range data.EntityInformation {
			if err := f.checkEntityInformation(ei); err != nil {
				return err
			}

			entityAddress := ei.Description.EntityAddress.Entity

			remoteDevice.RemoveByAddress(entityAddress)
		}

		ctrl.UpdateDevice(lastStateChange)
	}

	return nil
}

func (f *NodeManagement) announceFeatureDiscovery(ctrl spine.Context, e spine.Entity) error {
	entity := f.GetEntity()
	if entity == nil {
		return errors.New("announceFeatureDiscovery: entity not found")
	}
	device := entity.GetDevice()
	if device == nil {
		return errors.New("announceFeatureDiscovery: device not found")
	}
	entities := device.GetEntities()
	if entities == nil {
		return errors.New("announceFeatureDiscovery: entities not found")
	}

	for _, le := range entities {
		for _, lf := range le.GetFeatures() {

			// connect client to server features
			for _, rf := range e.GetFeatures() {
				lr := lf.GetRole()
				rr := rf.GetRole()
				rolesValid := (lr == model.RoleTypeSpecial && rr == model.RoleTypeSpecial) || (lr == model.RoleTypeClient && rr == model.RoleTypeServer)
				if lf.GetType() == rf.GetType() && rolesValid {
					if cf, ok := lf.(spine.ClientFeature); ok {
						if err := cf.ServerFound(ctrl, rf); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func (f *NodeManagement) handleDetailedDiscoveryData(ctrl spine.Context, op model.CmdClassifierType, data *model.NodeManagementDetailedDiscoveryDataType, isPartialForCmd bool) error {
	switch op {
	case model.CmdClassifierTypeRead:
		return f.readDetailedDiscoveryData(ctrl, *data)

	case model.CmdClassifierTypeReply:
		return f.replyDetailedDiscoveryData(ctrl, *data)

	case model.CmdClassifierTypeNotify:
		return f.notifyDetailedDiscoveryData(ctrl, *data, isPartialForCmd)

	default:
		return fmt.Errorf("nodemanagement.handleDetailedDiscoveryData: NodeManagementDetailedDiscoveryData CmdClassifierType not implemented: %s", op)
	}
}
