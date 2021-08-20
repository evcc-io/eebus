package feature

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type DeviceClassificationDelegate interface {
	UpdateDeviceClassificationData(*DeviceClassification, model.FeatureAddressType, model.DeviceClassificationManufacturerDataType)
}

type DeviceClassification struct {
	*spine.FeatureImpl
	Delegate DeviceClassificationDelegate
}

func NewDeviceClassificationServer() spine.Feature {
	f := &DeviceClassification{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeDeviceClassification,
			Role: model.RoleTypeServer,
		},
	}

	f.Add(model.FunctionEnumTypeDeviceClassificationManufacturerData, true, false)

	return f
}

func NewDeviceClassificationClient() spine.Feature {
	f := &DeviceClassification{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeDeviceClassification,
			Role: model.RoleTypeClient,
		},
	}

	return f
}

func (f *DeviceClassification) requestManufacturerData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *DeviceClassification) readManufacturerData(ctrl spine.Context, data model.DeviceClassificationManufacturerDataType) error {
	manufacturerData := f.Entity.GetManufacturerData()

	res := model.CmdType{
		DeviceClassificationManufacturerData: &manufacturerData,
	}

	err := ctrl.Reply(model.CmdClassifierTypeReply, res)

	return err
}

func (f *DeviceClassification) replyManufacturerData(ctrl spine.Context, rf model.FeatureAddressType, data model.DeviceClassificationManufacturerDataType) error {
	if f.Delegate != nil {
		f.Delegate.UpdateDeviceClassificationData(f, rf, data)
	}

	return nil
}

func (f *DeviceClassification) HandleRequest(ctrl spine.Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf spine.Feature) (*model.MsgCounterType, error) {
	switch fct {
	case model.FunctionEnumTypeDeviceClassificationManufacturerData:
		if op == model.CmdClassifierTypeRead {
			return f.requestManufacturerData(ctrl, rf)
		}
		return nil, fmt.Errorf("deviceclassification.handleRequest: FunctionEnumTypeDeviceClassificationManufacturerData op not implemented: %s", op)
	}

	return nil, fmt.Errorf("deviceclassification.handleRequest: FunctionEnumType not implemented: %s", fct)
}

func (f *DeviceClassification) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.DeviceClassificationManufacturerData != nil:
		data := cmd.DeviceClassificationManufacturerData
		switch op {
		case model.CmdClassifierTypeRead:
			return f.readManufacturerData(ctrl, *data)

		case model.CmdClassifierTypeReply:
			return f.replyManufacturerData(ctrl, rf, *data)

		default:
			return fmt.Errorf("deviceclassification.Handle: DeviceClassificationManufacturerData CmdClassifierType not implemented: %s", op)
		}

	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("deviceclassification.Handle: CmdType not implemented: %s", populatedFields(cmd))
	}
}

func (f *DeviceClassification) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	err := ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
	if err != nil {
		return err
	}
	// this is a workaround so it will work for EVSE also right now
	_, err = f.requestManufacturerData(ctrl, rf)
	return err
}
