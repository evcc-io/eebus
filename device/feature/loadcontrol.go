package feature

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type LoadControlLimitDescriptionDataType struct {
	LimitId       uint
	LimitType     model.LoadControlLimitTypeEnumType
	MeasurementId uint
	ScopeType     model.ScopeTypeEnumType
}

type LoadControlLimitDatasetType struct {
	LimitId           uint
	IsLimitChangeable bool
	IsLimitActive     bool
	Value             float64
}

type LoadControlData interface {
	GetLoadControlLimitDescriptionData() []LoadControlLimitDescriptionDataType
	GetLoadControlLimitData() []LoadControlLimitDatasetType
	WriteLoadControlLimitListData(ctrl spine.Context, rf spine.Feature, limits []LoadControlLimitDatasetType) error
}

type LoadControlDelegate interface {
	UpdateLoadControlLimitData(*LoadControl)
}

type LoadControl struct {
	*spine.FeatureImpl
	Delegate             LoadControlDelegate
	limitDescriptionData []LoadControlLimitDescriptionDataType
	limitData            []LoadControlLimitDatasetType
}

func NewLoadControlClient() spine.Feature {
	f := &LoadControl{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeLoadControl,
			Role: model.RoleTypeClient,
		},
	}

	return f
}

func (f *LoadControl) GetLoadControlLimitDescriptionData() []LoadControlLimitDescriptionDataType {
	return f.limitDescriptionData
}

func (f *LoadControl) GetLoadControlLimitData() []LoadControlLimitDatasetType {
	return f.limitData
}

func (f *LoadControl) requestLimitDescriptionListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		LoadControlLimitDescriptionListData: &model.LoadControlLimitDescriptionListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *LoadControl) replyLimitDescriptionListData(ctrl spine.Context, data model.LoadControlLimitDescriptionListDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":1}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":6}]},{"msgCounter":6898},{"msgCounterReference":18},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"loadControlLimitDescriptionListData":[{"loadControlLimitDescriptionData":[[{"limitId":1},{"limitType":"maxValueLimit"},{"limitCategory":"obligation"},{"limitDirection":"consume"},{"measurementId":1},{"unit":"A"},{"scopeType":"overloadProtection"}],[{"limitId":2},{"limitType":"maxValueLimit"},{"limitCategory":"recommendation"},{"limitDirection":"consume"},{"measurementId":1},{"unit":"A"},{"scopeType":"selfConsumption"}]]}]}]]}]}]}}]}
	// {"cmd":[[
	// 	{"loadControlLimitDescriptionListData":[
	// 		{"loadControlLimitDescriptionData":[
	// 			[
	// 				{"limitId":1},
	// 				{"limitType":"maxValueLimit"},
	// 				{"limitCategory":"obligation"},
	// 				{"limitDirection":"consume"},
	// 				{"measurementId":1},
	// 				{"unit":"A"},
	// 				{"scopeType":"overloadProtection"}
	// 			],
	// 			[
	// 				{"limitId":2},
	// 				{"limitType":"maxValueLimit"},
	// 				{"limitCategory":"recommendation"},
	// 				{"limitDirection":"consume"},
	// 				{"measurementId":1},
	// 				{"unit":"A"},
	// 				{"scopeType":"selfConsumption"}
	// 			]
	// 		]}
	// 	]}
	// ]]}

	f.limitDescriptionData = nil
	for _, item := range data.LoadControlLimitDescriptionData {
		newItem := LoadControlLimitDescriptionDataType{
			LimitId:       uint(*item.LimitId),
			LimitType:     model.LoadControlLimitTypeEnumType(*item.LimitType),
			MeasurementId: uint(*item.MeasurementId),
			ScopeType:     model.ScopeTypeEnumType(*item.ScopeType),
		}
		f.limitDescriptionData = append(f.limitDescriptionData, newItem)
	}

	if f.Delegate != nil {
		f.Delegate.UpdateLoadControlLimitData(f)
	}

	return nil
}

func (f *LoadControl) requestLimitListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		LoadControlLimitListData: &model.LoadControlLimitListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *LoadControl) replyLimitListData(ctrl spine.Context, data model.LoadControlLimitListDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":1}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":6}]},{"msgCounter":6928},{"msgCounterReference":34},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"loadControlLimitListData":[{"loadControlLimitData":[[{"limitId":1},{"isLimitChangeable":true},{"isLimitActive":false},{"value":[{"number":0},{"scale":0}]}],[{"limitId":2},{"isLimitChangeable":true},{"isLimitActive":false},{"value":[{"number":0},{"scale":0}]}]]}]}]]}]}]}}]}
	// {"cmd":[[
	// 	{"loadControlLimitListData":[
	// 		{"loadControlLimitData":[
	// 			[
	// 				{"limitId":1},
	// 				{"isLimitChangeable":true},
	// 				{"isLimitActive":false},
	// 				{"value":[{"number":0},{"scale":0}]}
	// 			],
	// 			[
	// 				{"limitId":2},
	// 				{"isLimitChangeable":true},
	// 				{"isLimitActive":false},
	// 				{"value":[{"number":0},{"scale":0}]}
	// 			]
	// 		]}
	// 	]}
	// ]]}

	f.limitData = nil
	for _, item := range data.LoadControlLimitData {
		if item.Value == nil || item.LimitId == nil || item.IsLimitActive == nil {
			continue
		}
		newItem := LoadControlLimitDatasetType{
			LimitId:       uint(*item.LimitId),
			IsLimitActive: *item.IsLimitActive,
			Value:         item.Value.GetValue(),
		}
		if item.IsLimitChangeable != nil {
			newItem.IsLimitChangeable = *item.IsLimitChangeable
		}
		f.limitData = append(f.limitData, newItem)
	}

	if f.Delegate != nil {
		f.Delegate.UpdateLoadControlLimitData(f)
	}

	return nil
}

func (f *LoadControl) WriteLoadControlLimitListData(ctrl spine.Context, rf spine.Feature, limits []LoadControlLimitDatasetType) error {
	var data []model.LoadControlLimitDataType

	for _, item := range limits {
		itemId := model.LoadControlLimitIdType(item.LimitId)
		// itemIsChangeable := true
		itemIsActive := true

		itemValue := model.NewScaledNumberType(item.Value)

		newItem := model.LoadControlLimitDataType{
			LimitId: &itemId,
			// IsLimitChangeable: &itemIsChangeable,
			IsLimitActive: &itemIsActive,
			Value:         itemValue,
		}
		data = append(data, newItem)
	}

	res := []model.CmdType{{
		LoadControlLimitListData: &model.LoadControlLimitListDataType{
			LoadControlLimitData: data,
		},
	}}

	return ctrl.Write(spine.FeatureAddressType(f), spine.FeatureAddressType(rf), res)
}

func (f *LoadControl) HandleRequest(ctrl spine.Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf spine.Feature) (*model.MsgCounterType, error) {
	switch fct {
	case model.FunctionEnumTypeLoadControlLimitDescriptionListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestLimitDescriptionListData(ctrl, rf)
		}
		return nil, fmt.Errorf("loadcontrol.handleRequest: FunctionEnumTypeLoadControlLimitDescriptionListData op not implemented: %s", op)

	case model.FunctionEnumTypeLoadControlLimitListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestLimitListData(ctrl, rf)
		}
		return nil, fmt.Errorf("loadcontrol.handleRequest: FunctionEnumTypeLoadControlLimitListData op not implemented: %s", op)
	}

	return nil, fmt.Errorf("loadcontrol.handleRequest: FunctionEnumType not implemented: %s", fct)
}

func (f *LoadControl) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.LoadControlLimitDescriptionListData != nil:
		data := cmd.LoadControlLimitDescriptionListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyLimitDescriptionListData(ctrl, *data)
		case model.CmdClassifierTypeNotify:
			return f.replyLimitDescriptionListData(ctrl, *data)
		default:
			return fmt.Errorf("loadcontrol.handle: LoadControlLimitDescriptionListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.LoadControlLimitListData != nil:
		data := cmd.LoadControlLimitListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyLimitListData(ctrl, *data)
		case model.CmdClassifierTypeNotify:
			return f.replyLimitListData(ctrl, *data)
		default:
			return fmt.Errorf("loadcontrol.handle: LoadControlLimitListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("loadcontrol.Handle: CmdType not implemented: %s", populatedFields(cmd))
	}
}

func (f *LoadControl) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	return ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
}
