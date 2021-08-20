package feature

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type DeviceDiagnosisDataType struct {
	OperationState model.DeviceDiagnosisOperatingStateEnumType
}

type DeviceDiagnosisData interface {
	GetDeviceDiagnosisDataData() DeviceDiagnosisDataType
}

type DeviceDiagnosisDelegate interface {
	UpdateDeviceDiagnosisData(*DeviceDiagnosis, model.FeatureAddressType, DeviceDiagnosisDataType)
}

type DeviceDiagnosis struct {
	*spine.FeatureImpl
	Delegate DeviceDiagnosisDelegate
	data     DeviceDiagnosisDataType
}

func NewDeviceDiagnosisServer() spine.Feature {
	f := &DeviceDiagnosis{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeDeviceDiagnosis,
			Role: model.RoleTypeServer,
		},
	}

	f.Add(model.FunctionEnumTypeDeviceDiagnosisStateData, true, false)
	f.Add(model.FunctionEnumTypeDeviceDiagnosisHeartbeatData, true, false)

	return f
}

func NewDeviceDiagnosisClient() spine.Feature {
	f := &DeviceDiagnosis{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeDeviceDiagnosis,
			Role: model.RoleTypeClient,
		},
	}

	return f
}

func (f *DeviceDiagnosis) GetDeviceDiagnosisDataData() DeviceDiagnosisDataType {
	return f.data
}

func (f *DeviceDiagnosis) readHeartbeatData(ctrl spine.Context, data model.DeviceDiagnosisHeartbeatDataType) error {
	// TODO is this all we need here?

	var heartBeatTimeout string = "PT4S"

	res := model.CmdType{
		DeviceDiagnosisHeartbeatData: &model.DeviceDiagnosisHeartbeatDataType{
			HeartbeatCounter: ctrl.HeartbeatCounter(),
			HeartbeatTimeout: &heartBeatTimeout,
		},
	}

	err := ctrl.Reply(model.CmdClassifierTypeReply, res)

	return err
}

func (f *DeviceDiagnosis) requestStateData(ctrl spine.Context, rf spine.Feature) error {
	res := []model.CmdType{{
		DeviceDiagnosisStateData: &model.DeviceDiagnosisStateDataType{},
	}}

	_, err := ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
	return err
}

func (f *DeviceDiagnosis) readStateData(ctrl spine.Context, data model.DeviceDiagnosisStateDataType) error {
	// TODO is this all we need here?

	operationState := f.Entity.GetOperationState()

	res := model.CmdType{
		DeviceDiagnosisStateData: &model.DeviceDiagnosisStateDataType{
			OperatingState: &operationState,
		},
	}

	err := ctrl.Reply(model.CmdClassifierTypeReply, res)

	return err
}

func (f *DeviceDiagnosis) replyStateData(ctrl spine.Context, rf model.FeatureAddressType, data model.DeviceDiagnosisStateDataType) error {
	if f.Delegate != nil {
		f.Delegate.UpdateDeviceDiagnosisData(f, rf, DeviceDiagnosisDataType{OperationState: model.DeviceDiagnosisOperatingStateEnumType(*data.OperatingState)})
	}

	return nil
}

func (f *DeviceDiagnosis) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.DeviceDiagnosisHeartbeatData != nil:
		data := cmd.DeviceDiagnosisHeartbeatData
		switch op {
		case model.CmdClassifierTypeRead:
			return f.readHeartbeatData(ctrl, *data)

		default:
			return fmt.Errorf("devicediagnosis.Handle: DeviceDiagnosisHeartbeatData CmdClassifierType not implemented. %s", op)
		}

	case cmd.DeviceDiagnosisStateData != nil:
		data := cmd.DeviceDiagnosisStateData
		switch op {
		case model.CmdClassifierTypeRead:
			return f.readStateData(ctrl, *data)

		case model.CmdClassifierTypeReply:
			return f.replyStateData(ctrl, rf, *data)

		case model.CmdClassifierTypeNotify:
			return f.replyStateData(ctrl, rf, *data)

		default:
			return fmt.Errorf("devicediagnosis.Handle: DeviceDiagnosisStateData CmdClassifierType not implemented. %s", op)
		}

	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("devicediagnosis.Handle: CmdType not implemented. %s", populatedFields(cmd))
	}
}

func (f *DeviceDiagnosis) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	err := ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
	if err != nil {
		return err
	}
	return f.requestStateData(ctrl, rf)
}
