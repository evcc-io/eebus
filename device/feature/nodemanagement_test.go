package feature

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type mockContext struct {
	dev spine.Device
}

func (c *mockContext) AddressSource() *model.FeatureAddressType {
	address := model.FeatureAddressType{}
	return &address
}

func (c *mockContext) HeartbeatCounter() *uint64 {
	var i uint64 = 0
	return &i
}

func (c *mockContext) LocalDeviceFeature(featureType model.FeatureTypeEnumType, role model.RoleType) (*spine.Feature, error) {
	return nil, nil
}

func (c *mockContext) Subscribe(localFeature spine.Feature, remoteFeature spine.Feature, serverFeatureType model.FeatureTypeType) error {
	return nil
}

func (c *mockContext) ProcessSequenceFlowRequest(featureType model.FeatureTypeEnumType, functionType model.FunctionEnumType, cmdClassifier model.CmdClassifierType) (*model.MsgCounterType, error) {
	return nil, nil
}

func (c *mockContext) Request(cmdClassifier model.CmdClassifierType, senderAddress model.FeatureAddressType, destinationAddress model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error) {
	return nil, nil
}

func (c *mockContext) Reply(cmdClassifier model.CmdClassifierType, cmd model.CmdType) error {
	return nil
}

func (c *mockContext) Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error {
	return nil
}

func (c *mockContext) Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error {
	return nil
}

func (c *mockContext) SetDevice(device spine.Device) {
	c.dev = device
}

func (c *mockContext) GetDevice() spine.Device {
	return c.dev
}

func (c *mockContext) UpdateDevice(stateChange model.NetworkManagementStateChangeType) {}

func (c *mockContext) AddSubscription(data model.SubscriptionManagementRequestCallType) error {
	return nil
}

func (c *mockContext) RemoveSubscription(data model.SubscriptionManagementDeleteCallType) error {
	return nil
}

func (c *mockContext) Subscriptions() []model.SubscriptionManagementEntryDataType {
	return nil
}

func TestReplyDetailedDiscoveryDataWithoutSubEntities(t *testing.T) {
	nodeManagementFeature := &NodeManagement{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeNodeManagement,
			Role: model.RoleTypeSpecial,
		},
	}

	cmdJson := `[{"specificationVersionList":[{"specificationVersion":["1.1.1"]}]},{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:EVSE"}]},{"deviceType":"ChargingStation"},{"networkFeatureSet":"smart"}]}]},{"entityInformation":[[{"description":[{"entityAddress":[{"entity":[0]}]},{"entityType":"DeviceInformation"}]}],[{"description":[{"entityAddress":[{"entity":[1]}]},{"entityType":"EVSE"},{"description":"Electric Vehicle Supply Equipment"}]}]]},{"featureInformation":[[{"description":[{"featureAddress":[{"entity":[0]},{"feature":0}]},{"featureType":"NodeManagement"},{"role":"special"},{"supportedFunction":[[{"function":"nodeManagementDetailedDiscoveryData"},{"possibleOperations":[{"read":[]}]}],[{"function":"nodeManagementSubscriptionRequestCall"},{"possibleOperations":[]}],[{"function":"nodeManagementBindingRequestCall"},{"possibleOperations":[]}],[{"function":"nodeManagementSubscriptionDeleteCall"},{"possibleOperations":[]}],[{"function":"nodeManagementBindingDeleteCall"},{"possibleOperations":[]}],[{"function":"nodeManagementSubscriptionData"},{"possibleOperations":[{"read":[]}]}],[{"function":"nodeManagementBindingData"},{"possibleOperations":[{"read":[]}]}],[{"function":"nodeManagementUseCaseData"},{"possibleOperations":[{"read":[]}]}]]}]}],[{"description":[{"featureAddress":[{"entity":[0]},{"feature":1}]},{"featureType":"DeviceClassification"},{"role":"server"},{"supportedFunction":[[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]]}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":1}]},{"featureType":"DeviceClassification"},{"role":"server"},{"supportedFunction":[[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Classification for EVSE"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":2}]},{"featureType":"DeviceDiagnosis"},{"role":"server"},{"supportedFunction":[[{"function":"deviceDiagnosisStateData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Diagnosis EVSE"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":3}]},{"featureType":"DeviceClassification"},{"role":"client"},{"description":"EMS Device Classification"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":5}]},{"featureType":"DeviceDiagnosis"},{"role":"client"},{"description":"Device Diagnosis Heartbeat"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":6}]},{"featureType":"Bill"},{"role":"server"},{"supportedFunction":[[{"function":"billDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"billConstraintsListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"billListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Bill Feature for EVSE"}]}]]}]`
	data := model.NodeManagementDetailedDiscoveryDataType{}

	err := json.Unmarshal(json.RawMessage(cmdJson), &data)

	if err != nil {
		t.Errorf("Unmarshal failed error = %v", err)
	}

	context := new(mockContext)
	err = nodeManagementFeature.replyDetailedDiscoveryData(context, data)

	if err != nil {
		t.Errorf("replyDetailedDiscoveryData failed error = %v", err)
	}

	if context.dev != nil {
		context.dev.Dump(os.Stdout)
	}
}

func TestReplyDetailedDiscoveryDataWithSubEntities(t *testing.T) {
	nodeManagementFeature := &NodeManagement{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeNodeManagement,
			Role: model.RoleTypeSpecial,
		},
	}

	cmdJson := `[{"specificationVersionList":[{"specificationVersion":["1.1.1"]}]},{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:EVSE"}]},{"deviceType":"ChargingStation"},{"networkFeatureSet":"smart"}]}]},{"entityInformation":[[{"description":[{"entityAddress":[{"entity":[0]}]},{"entityType":"DeviceInformation"}]}],[{"description":[{"entityAddress":[{"entity":[1]}]},{"entityType":"EVSE"},{"description":"Electric Vehicle Supply Equipment"}]}],[{"description":[{"entityAddress":[{"entity":[1,1]}]},{"entityType":"EV"},{"description":"Electric Vehicle"}]}]]},{"featureInformation":[[{"description":[{"featureAddress":[{"entity":[0]},{"feature":0}]},{"featureType":"NodeManagement"},{"role":"special"},{"supportedFunction":[[{"function":"nodeManagementDetailedDiscoveryData"},{"possibleOperations":[{"read":[]}]}],[{"function":"nodeManagementSubscriptionRequestCall"},{"possibleOperations":[]}],[{"function":"nodeManagementBindingRequestCall"},{"possibleOperations":[]}],[{"function":"nodeManagementSubscriptionDeleteCall"},{"possibleOperations":[]}],[{"function":"nodeManagementBindingDeleteCall"},{"possibleOperations":[]}],[{"function":"nodeManagementSubscriptionData"},{"possibleOperations":[{"read":[]}]}],[{"function":"nodeManagementBindingData"},{"possibleOperations":[{"read":[]}]}],[{"function":"nodeManagementUseCaseData"},{"possibleOperations":[{"read":[]}]}]]}]}],[{"description":[{"featureAddress":[{"entity":[0]},{"feature":1}]},{"featureType":"DeviceClassification"},{"role":"server"},{"supportedFunction":[[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]]}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":1}]},{"featureType":"DeviceClassification"},{"role":"server"},{"supportedFunction":[[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Classification for EVSE"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":2}]},{"featureType":"DeviceDiagnosis"},{"role":"server"},{"supportedFunction":[[{"function":"deviceDiagnosisStateData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Diagnosis EVSE"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":3}]},{"featureType":"DeviceClassification"},{"role":"client"},{"description":"EMS Device Classification"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":5}]},{"featureType":"DeviceDiagnosis"},{"role":"client"},{"description":"Device Diagnosis Heartbeat"}]}],[{"description":[{"featureAddress":[{"entity":[1]},{"feature":6}]},{"featureType":"Bill"},{"role":"server"},{"supportedFunction":[[{"function":"billDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"billConstraintsListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"billListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Bill Feature for EVSE"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":1}]},{"featureType":"LoadControl"},{"role":"server"},{"supportedFunction":[[{"function":"loadControlLimitDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"loadControlLimitListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Load Control"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":2}]},{"featureType":"ElectricalConnection"},{"role":"server"},{"supportedFunction":[[{"function":"electricalConnectionParameterDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionPermittedValueSetListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Electrical Connection"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":3}]},{"featureType":"Measurement"},{"specificUsage":["Electrical"]},{"role":"server"},{"supportedFunction":[[{"function":"measurementListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"measurementDescriptionListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Measurements"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":5}]},{"featureType":"DeviceConfiguration"},{"role":"server"},{"supportedFunction":[[{"function":"deviceConfigurationKeyValueDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"deviceConfigurationKeyValueListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Configuration EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":6}]},{"featureType":"DeviceClassification"},{"role":"server"},{"supportedFunction":[[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Classification for EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":7}]},{"featureType":"TimeSeries"},{"role":"server"},{"supportedFunction":[[{"function":"timeSeriesConstraintsListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"timeSeriesDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"timeSeriesListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Time Series"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":8}]},{"featureType":"IncentiveTable"},{"role":"server"},{"supportedFunction":[[{"function":"incentiveTableConstraintsData"},{"possibleOperations":[{"read":[]}]}],[{"function":"incentiveTableData"},{"possibleOperations":[{"read":[]},{"write":[]}]}],[{"function":"incentiveTableDescriptionData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Incentive Table"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":9}]},{"featureType":"DeviceDiagnosis"},{"role":"server"},{"supportedFunction":[[{"function":"deviceDiagnosisStateData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Diagnosis EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":10}]},{"featureType":"Identification"},{"role":"server"},{"supportedFunction":[[{"function":"identificationListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Identification for EV"}]}]]}]`
	data := model.NodeManagementDetailedDiscoveryDataType{}

	err := json.Unmarshal(json.RawMessage(cmdJson), &data)

	if err != nil {
		t.Errorf("Unmarshal failed error = %v", err)
	}

	context := new(mockContext)
	err = nodeManagementFeature.replyDetailedDiscoveryData(context, data)

	if err != nil {
		t.Errorf("replyDetailedDiscoveryData failed error = %v", err)
	}

	if context.dev != nil {
		context.dev.Dump(os.Stdout)
	}
}
