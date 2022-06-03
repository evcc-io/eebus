package feature

import (
	"errors"
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type DeviceConfigurationDescriptionDataType struct {
	KeyId        uint
	KeyName      model.DeviceConfigurationKeyNameEnumType
	KeyValueType model.DeviceConfigurationKeyValueTypeType
}

type DeviceConfigurationDatasetDataType struct {
	KeyName         model.DeviceConfigurationKeyNameEnumType
	KeyValueType    model.DeviceConfigurationKeyValueTypeType
	KeyValueString  string
	KeyValueBoolean bool
}

type DeviceConfigurationDelegate interface {
	UpdateDeviceConfigurationData(*DeviceConfiguration, []DeviceConfigurationDatasetDataType)
}

type DeviceConfiguration struct {
	*spine.FeatureImpl
	Delegate        DeviceConfigurationDelegate
	descriptionData []DeviceConfigurationDescriptionDataType
	datasetData     []DeviceConfigurationDatasetDataType
}

func NewDeviceConfigurationClient() spine.Feature {
	f := &DeviceConfiguration{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeDeviceConfiguration,
			Role: model.RoleTypeClient,
		},
	}

	return f
}

func (f *DeviceConfiguration) requestKeyValueDescriptionListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		DeviceConfigurationKeyValueDescriptionListData: &model.DeviceConfigurationKeyValueDescriptionListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *DeviceConfiguration) replyKeyValueDescriptionListData(ctrl spine.Context, data model.DeviceConfigurationKeyValueDescriptionListDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":5}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":4}]},{"msgCounter":23313},{"msgCounterReference":19},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"deviceConfigurationKeyValueDescriptionListData":[{"deviceConfigurationKeyValueDescriptionData":[[{"keyId":1},{"keyName":"asymmetricChargingSupported"},{"valueType":"boolean"}],[{"keyId":2},{"keyName":"communicationsStandard"},{"valueType":"string"}]]}]}]]}]}]}}]}

	f.descriptionData = nil
	for _, item := range data.DeviceConfigurationKeyValueDescriptionData {
		newItem := DeviceConfigurationDescriptionDataType{
			KeyId:        uint(*item.KeyId),
			KeyName:      model.DeviceConfigurationKeyNameEnumType(*item.KeyName),
			KeyValueType: *item.ValueType,
		}
		f.descriptionData = append(f.descriptionData, newItem)
	}

	return nil
}

func (f *DeviceConfiguration) requestKeyValueListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		DeviceConfigurationKeyValueListData: &model.DeviceConfigurationKeyValueListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *DeviceConfiguration) replyKeyValueListData(ctrl spine.Context, data model.DeviceConfigurationKeyValueListDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":5}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":4}]},{"msgCounter":24307},{"msgCounterReference":34},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"deviceConfigurationKeyValueListData":[{"deviceConfigurationKeyValueData":[[{"keyId":1},{"value":[{"boolean":false}]}],[{"keyId":2},{"value":[{"string":"iso15118-2ed2"}]}]]}]}]]}]}]}}]}

	if f.descriptionData == nil {
		return errors.New("deviceconfiguration.replyKeyValueListData: descriptionData is not set, needs to be requested first")
	}

	f.datasetData = nil
	for _, item := range data.DeviceConfigurationKeyValueData {
		var keyId uint = uint(*item.KeyId)
		var valueTypeForKeyID model.DeviceConfigurationKeyValueTypeType
		var nameForKeyID model.DeviceConfigurationKeyNameEnumType
		found := false

		for _, descriptionItem := range f.descriptionData {
			if keyId == descriptionItem.KeyId {
				valueTypeForKeyID = descriptionItem.KeyValueType
				nameForKeyID = descriptionItem.KeyName
				found = true
				break
			}
		}

		if !found {
			continue
		}

		newItem := DeviceConfigurationDatasetDataType{
			KeyName:      nameForKeyID,
			KeyValueType: valueTypeForKeyID,
		}

		valid := true
		switch valueTypeForKeyID {
		case model.DeviceConfigurationKeyValueTypeTypeBoolean:
			newItem.KeyValueBoolean = *item.Value.Boolean
		case model.DeviceConfigurationKeyValueTypeTypeString:
			newItem.KeyValueString = string(*item.Value.String)
		default:
			valid = false
		}

		if valid {
			f.datasetData = append(f.datasetData, newItem)
		}
	}

	if f.Delegate != nil {
		f.Delegate.UpdateDeviceConfigurationData(f, f.datasetData)
	}

	return nil
}

func (f *DeviceConfiguration) notifyKeyValueListData(ctrl spine.Context, data model.DeviceConfigurationKeyValueListDataType, isPartialForCmd bool) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.3.0"},{"addressSource":[{"device":"d:_i:47859_Elli-Wallbox-2019A0OV8H"},{"entity":[1,1]},{"feature":24}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":4}]},{"msgCounter":767},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"function":"deviceConfigurationKeyValueListData"},{"filter":[[{"cmdControl":[{"partial":[]}]}]]},{"deviceConfigurationKeyValueListData":[{"deviceConfigurationKeyValueData":[[{"keyId":1},{"value":[{"string":"iec61851"}]}]]}]}]]}]}]}}]}

	if f.descriptionData == nil {
		return errors.New("deviceconfiguration.replyKeyValueListData: descriptionData is not set, needs to be requested first")
	}

	for _, item := range data.DeviceConfigurationKeyValueData {
		if item.KeyId == nil || item.Value == nil {
			continue
		}

		var keyId uint = uint(*item.KeyId)
		var valueTypeForKeyID model.DeviceConfigurationKeyValueTypeType
		var nameForKeyID model.DeviceConfigurationKeyNameEnumType
		found := false

		for _, descriptionItem := range f.descriptionData {
			if keyId == descriptionItem.KeyId {
				valueTypeForKeyID = descriptionItem.KeyValueType
				nameForKeyID = descriptionItem.KeyName
				found = true
				break
			}
		}

		if !found {
			continue
		}

		var updatedDataSetItem DeviceConfigurationDatasetDataType
		var updatedDataSetIndex int
		itemFound := false
		for index, datasetItem := range f.datasetData {
			if datasetItem.KeyName == nameForKeyID {
				updatedDataSetItem = datasetItem
				updatedDataSetIndex = index
				itemFound = true
				break
			}
		}

		var dataSetItem DeviceConfigurationDatasetDataType

		if !itemFound {
			dataSetItem = DeviceConfigurationDatasetDataType{
				KeyName:      nameForKeyID,
				KeyValueType: valueTypeForKeyID,
			}
		} else {
			dataSetItem = updatedDataSetItem
		}

		valid := true
		switch valueTypeForKeyID {
		case model.DeviceConfigurationKeyValueTypeTypeBoolean:
			dataSetItem.KeyValueBoolean = *item.Value.Boolean
		case model.DeviceConfigurationKeyValueTypeTypeString:
			dataSetItem.KeyValueString = string(*item.Value.String)
		default:
			valid = false
		}

		if valid {
			if !itemFound {
				f.datasetData = append(f.datasetData, dataSetItem)
			} else {
				f.datasetData[updatedDataSetIndex] = dataSetItem
			}
		}
	}

	if f.Delegate != nil {
		f.Delegate.UpdateDeviceConfigurationData(f, f.datasetData)
	}

	return nil
}

func (f *DeviceConfiguration) HandleRequest(ctrl spine.Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf spine.Feature) (*model.MsgCounterType, error) {
	switch fct {
	case model.FunctionEnumTypeDeviceConfigurationKeyValueDescriptionListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestKeyValueDescriptionListData(ctrl, rf)
		}
		return nil, fmt.Errorf("deviceconfiguration.handleRequest: FunctionEnumTypeDeviceConfigurationKeyValueDescriptionListData op not implemented: %s", op)

	case model.FunctionEnumTypeDeviceConfigurationKeyValueListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestKeyValueListData(ctrl, rf)
		}
		return nil, fmt.Errorf("deviceconfiguration.handleRequest: FunctionEnumTypeDeviceConfigurationKeyValueListData op not implemented: %s", op)
	}

	return nil, fmt.Errorf("deviceconfiguration.handleRequest: FunctionEnumType not implemented: %s", fct)
}

func (f *DeviceConfiguration) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.DeviceConfigurationKeyValueDescriptionListData != nil:
		data := cmd.DeviceConfigurationKeyValueDescriptionListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyKeyValueDescriptionListData(ctrl, *data)
		case model.CmdClassifierTypeNotify:
			return f.replyKeyValueDescriptionListData(ctrl, *data)

		default:
			return fmt.Errorf("deviceconfiguration.Handle: DeviceConfigurationKeyValueDescriptionListData CmdClassifierType not implemented: %s", op)
		}
	case cmd.DeviceConfigurationKeyValueListData != nil:
		data := cmd.DeviceConfigurationKeyValueListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyKeyValueListData(ctrl, *data)
		case model.CmdClassifierTypeNotify:
			return f.notifyKeyValueListData(ctrl, *data, isPartialForCmd)

		default:
			return fmt.Errorf("deviceconfiguration.Handle: DeviceConfigurationKeyValueListData CmdClassifierType not implemented: %s", op)
		}
	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("deviceconfiguration.Handle: CmdType not implemented: %s", populatedFields(cmd))
	}
}

func (f *DeviceConfiguration) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	return ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
}
