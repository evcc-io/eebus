package feature

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type ElectricalConnectionParameterDescriptionDataType struct {
	ElectricalConnectionId uint
	ParameterId            uint
	MeasurementId          uint
	Phase                  uint
	ScopeType              model.ScopeTypeEnumType
}

type ElectricalConnectionDatasetDataType struct {
	ElectricalConnectionId uint
	ConnectedPhases        uint
}

type ElectricalConnectionPermittedDataType struct {
	ElectricalConnectionId uint
	ParameterId            uint
	Value                  float64
	MinValue               float64
	MaxValue               float64
}

type ElectricalConnectionDescriptionData interface {
	GetElectricalConnectionDescription() []ElectricalConnectionParameterDescriptionDataType
	GetElectricalConnectionData() []ElectricalConnectionDatasetDataType
	GetElectricalConnectionPermittedData() []ElectricalConnectionPermittedDataType
}

type ElectricalConnectionDelegate interface {
	UpdateElectricalConnectionData(*ElectricalConnection)
}

type ElectricalConnection struct {
	*spine.FeatureImpl
	Delegate                 ElectricalConnectionDelegate
	parameterDescriptionData []ElectricalConnectionParameterDescriptionDataType
	descriptionData          []ElectricalConnectionDatasetDataType
	permittedData            []ElectricalConnectionPermittedDataType
}

func NewElectricalConnectionClient() spine.Feature {
	f := &ElectricalConnection{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeElectricalConnection,
			Role: model.RoleTypeClient,
		},
	}

	return f
}

func (f *ElectricalConnection) GetElectricalConnectionDescription() []ElectricalConnectionParameterDescriptionDataType {
	return f.parameterDescriptionData
}

func (f *ElectricalConnection) GetElectricalConnectionData() []ElectricalConnectionDatasetDataType {
	return f.descriptionData
}

func (f *ElectricalConnection) GetElectricalConnectionPermittedData() []ElectricalConnectionPermittedDataType {
	return f.permittedData
}

func (f *ElectricalConnection) requestParameterDescriptionListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		ElectricalConnectionParameterDescriptionListData: &model.ElectricalConnectionParameterDescriptionListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *ElectricalConnection) replyParameterDescriptionListData(ctrl spine.Context, data model.ElectricalConnectionParameterDescriptionListDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":2}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":8}]},{"msgCounter":15976},{"msgCounterReference":34},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"electricalConnectionParameterDescriptionListData":[{"electricalConnectionParameterDescriptionData":[[{"electricalConnectionId":0},{"parameterId":1},{"measurementId":1},{"voltageType":"ac"},{"acMeasuredPhases":"a"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":2},{"measurementId":4},{"voltageType":"ac"},{"acMeasuredPhases":"a"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":3},{"measurementId":2},{"voltageType":"ac"},{"acMeasuredPhases":"b"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":4},{"measurementId":5},{"voltageType":"ac"},{"acMeasuredPhases":"b"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":5},{"measurementId":3},{"voltageType":"ac"},{"acMeasuredPhases":"c"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":6},{"measurementId":6},{"voltageType":"ac"},{"acMeasuredPhases":"c"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":7},{"measurementId":7},{"voltageType":"ac"},{"acMeasuredPhases":"abc"},{"acMeasuredInReferenceTo":"neutral"},{"acMeasurementType":"real"},{"acMeasurementVariant":"rms"}],[{"electricalConnectionId":0},{"parameterId":8},{"acMeasuredPhases":"abc"},{"scopeType":"acPowerTotal"}]]}]}]]}]}]}}]}

	// TODO make this work with any kind of data, not only currents on single phases
	var phases = map[string]uint{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	f.parameterDescriptionData = nil
	for _, item := range data.ElectricalConnectionParameterDescriptionData {
		if item.ElectricalConnectionId == nil || item.ParameterId == nil || item.AcMeasuredPhases == nil {
			continue
		}

		phasesValue := string(*item.AcMeasuredPhases)
		if phasesValue == "a" || phasesValue == "b" || phasesValue == "c" || phasesValue == "abc" {
			newItem := ElectricalConnectionParameterDescriptionDataType{
				ElectricalConnectionId: uint(*item.ElectricalConnectionId),
				ParameterId:            uint(*item.ParameterId),
			}
			if item.MeasurementId != nil {
				newItem.MeasurementId = uint(*item.MeasurementId)
			}
			if item.ScopeType != nil {
				newItem.ScopeType = model.ScopeTypeEnumType(*item.ScopeType)
			}
			if phaseValue, ok := phases[phasesValue]; ok {
				newItem.Phase = phaseValue
			}
			f.parameterDescriptionData = append(f.parameterDescriptionData, newItem)
		}
	}

	return nil
}

func (f *ElectricalConnection) requestDescriptionListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		ElectricalConnectionDescriptionListData: &model.ElectricalConnectionDescriptionListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *ElectricalConnection) replyDescriptionListData(ctrl spine.Context, data model.ElectricalConnectionDescriptionListDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":2}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":8}]},{"msgCounter":15981},{"msgCounterReference":35},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"electricalConnectionDescriptionListData":[{"electricalConnectionDescriptionData":[[{"electricalConnectionId":0},{"powerSupplyType":"ac"},{"acConnectedPhases":3},{"positiveEnergyDirection":"consume"}]]}]}]]}]}]}}]}

	f.descriptionData = nil
	for _, item := range data.ElectricalConnectionDescriptionData {
		if item.ElectricalConnectionId == nil {
			continue
		}

		newItem := ElectricalConnectionDatasetDataType{
			ElectricalConnectionId: uint(*item.ElectricalConnectionId),
		}

		if item.AcConnectedPhases != nil {
			newItem.ConnectedPhases = *item.AcConnectedPhases
		} else {
			// Assume this
			newItem.ConnectedPhases = 3
		}
		f.descriptionData = append(f.descriptionData, newItem)
	}

	if f.Delegate != nil {
		f.Delegate.UpdateElectricalConnectionData(f)
	}

	return nil
}

func (f *ElectricalConnection) requestPermittedValueSetData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		ElectricalConnectionPermittedValueSetListData: &model.ElectricalConnectionPermittedValueSetListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *ElectricalConnection) replyPermittedValueSetData(ctrl spine.Context, data model.ElectricalConnectionPermittedValueSetListDataType, isPartialForCmd bool, filter []model.FilterType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":2}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":8}]},{"msgCounter":1793},{"msgCounterReference":35},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"electricalConnectionPermittedValueSetListData":[{"electricalConnectionPermittedValueSetData":[[{"electricalConnectionId":0},{"parameterId":1},{"permittedValueSet":[[{"value":[[{"number":100},{"scale":-3}]]},{"range":[[{"min":[{"number":2},{"scale":0}]},{"max":[{"number":16},{"scale":0}]}]]}]]}],[{"electricalConnectionId":0},{"parameterId":8},{"permittedValueSet":[[{"value":[[{"number":100},{"scale":-3}]]},{"range":[[{"min":[{"number":490},{"scale":0}]},{"max":[{"number":3920},{"scale":0}]}]]}]]}]]}]}]]}]}]}}]}
	// {"cmd":[[
	// 	{"electricalConnectionPermittedValueSetListData":[
	// 		{"electricalConnectionPermittedValueSetData":[
	// 			[
	// 				{"electricalConnectionId":0},{"parameterId":1},
	// 				{"permittedValueSet":[[
	// 					{"value":[[
	// 						{"number":100},
	// 						{"scale":-3}
	// 					]]},
	// 					{"range":[[
	// 						{"min":[{"number":2},{"scale":0}]},
	// 						{"max":[{"number":16},{"scale":0}]}
	// 					]]}
	// 				]]}
	// 			],
	// 			[
	// 				{"electricalConnectionId":0},{"parameterId":8},
	// 				{"permittedValueSet":[[
	// 					{"value":[[
	// 						{"number":100},
	// 						{"scale":-3}
	// 					]]},
	// 					{"range":[[
	// 						{"min":[{"number":490},{"scale":0}]},
	// 						{"max":[{"number":3920},{"scale":0}]}
	// 					]]}
	// 				]]}
	// 			]
	// 		]}
	// 	]}
	// ]]}
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.3.0"},{"addressSource":[{"device":"d:_i:47859_Elli-Wallbox-2019A0OV8H"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":8}]},{"msgCounter":768},{"cmdClassifier":"notify"}]},{"payload":[
	// {"cmd":[[
	// 	{"function":"electricalConnectionPermittedValueSetListData"},
	// 	{"filter":[[{"cmdControl":[{"partial":[]}]}]]},
	// 	{"electricalConnectionPermittedValueSetListData":[
	// 		{"electricalConnectionPermittedValueSetData":[
	// 			[{"electricalConnectionId":0},{"parameterId":1},{"permittedValueSet":[[{"range":[[{"min":[{"number":6},{"scale":0}]},{"max":[{"number":6},{"scale":0}]}]]}]]}],
	// 			[{"electricalConnectionId":0},{"parameterId":2},{"permittedValueSet":[[{"range":[[{"min":[{"number":6},{"scale":0}]},{"max":[{"number":6},{"scale":0}]}]]}]]}],
	// 			[{"electricalConnectionId":0},{"parameterId":3},{"permittedValueSet":[[{"range":[[{"min":[{"number":6},{"scale":0}]},{"max":[{"number":6},{"scale":0}]}]]}]]}]
	// ]}]}]]}]}]}}]}
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.3.0"},{"addressSource":[{"device":"d:_i:47859_Elli-Wallbox-2019A0OV8H"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":8}]},{"msgCounter":148},{"cmdClassifier":"notify"}]},{"payload":[
	// {"cmd":[[
	// 	{"function":"electricalConnectionPermittedValueSetListData"},
	// 	{"filter":[[{"cmdControl":[{"partial":[]}]}]]},
	// 	{"electricalConnectionPermittedValueSetListData":[
	// 		{"electricalConnectionPermittedValueSetData":[
	// 			[{"electricalConnectionId":0},{"parameterId":1},{"permittedValueSet":[[{"range":[[{"min":[{"number":6},{"scale":0}]},{"max":[{"number":16},{"scale":0}]}]]}]]}],
	// 			[{"electricalConnectionId":0},{"parameterId":2},{"permittedValueSet":[[{"range":[[{"min":[{"number":6},{"scale":0}]},{"max":[{"number":16},{"scale":0}]}]]}]]}],
	// 			[{"electricalConnectionId":0},{"parameterId":3},{"permittedValueSet":[[{"range":[[{"min":[{"number":6},{"scale":0}]},{"max":[{"number":16},{"scale":0}]}]]}]]}]
	// ]}]}]]}]}]}}]}
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.3.0"},{"addressSource":[{"device":"d:_i:47859_Elli-Wallbox-2019A0OV8H"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":8}]},{"msgCounter":266},{"cmdClassifier":"notify"}]},{"payload": [
	// {"cmd":[[
	// 	{"function":"electricalConnectionPermittedValueSetListData"},
	// 	{"filter": [
	// 		[{"cmdControl":[{"delete":[]}]},{"electricalConnectionPermittedValueSetListDataSelectors":[{"electricalConnectionId":0},{"parameterId":0}]}],
	// 		[{"cmdControl":[{"partial":[]}]}]]
	// 	},
	// 	{"electricalConnectionPermittedValueSetListData":[
	// 		{"electricalConnectionPermittedValueSetData":[[
	// 			{"electricalConnectionId":0},
	// 			{"parameterId":0},
	// 			{"permittedValueSet":[[{"range":[[{"min":[{"number":1},{"scale":0}]}]]}]]}
	// ]]}]}]]}]}]}}]}

	if !isPartialForCmd {
		f.permittedData = nil
	}

	for _, filterItem := range filter {
		if filterItem.CmdControl != nil && filterItem.CmdControl.Delete != nil {
			selector := filterItem.ElectricalConnectionPermittedValueSetListDataSelectors
			if selector != nil && selector.ElectricalConnectionId != nil && selector.ParameterId != nil {

				connectionID := uint(*selector.ElectricalConnectionId)
				parameterID := uint(*selector.ParameterId)

				newValueSetData := make([]ElectricalConnectionPermittedDataType, 0)
				for _, item := range f.permittedData {
					if connectionID != uint(item.ElectricalConnectionId) || parameterID != uint(item.ParameterId) {
						newValueSetData = append(newValueSetData, item)
					}
				}
				f.permittedData = newValueSetData
			}
		}
	}

	for _, item := range data.ElectricalConnectionPermittedValueSetData {
		if item.ElectricalConnectionId == nil || item.ParameterId == nil {
			continue
		}

		dataSetItem := ElectricalConnectionPermittedDataType{
			ElectricalConnectionId: uint(*item.ElectricalConnectionId),
			ParameterId:            uint(*item.ParameterId),
		}

		replaceIndex := -1
		if isPartialForCmd {
			for index, datasetItem := range f.permittedData {
				if datasetItem.ElectricalConnectionId == uint(*item.ElectricalConnectionId) && datasetItem.ParameterId == uint(*item.ParameterId) {
					replaceIndex = index
					break
				}
			}
		}

		if len(item.PermittedValueSet) > 0 {
			valueData := item.PermittedValueSet[0].Value
			if len(valueData) > 0 {
				valueItem := valueData[0]
				dataSetItem.Value = valueItem.GetValue()
			}
			rangeData := item.PermittedValueSet[0].Range
			if len(rangeData) > 0 {
				rangeValue := rangeData[0]
				if rangeValue.Min != nil {
					dataSetItem.MinValue = rangeValue.Min.GetValue()
				}
				if rangeValue.Max != nil {
					dataSetItem.MaxValue = rangeValue.Max.GetValue()
				}
			}
		}
		if replaceIndex != -1 {
			f.permittedData[replaceIndex] = dataSetItem
		} else {
			f.permittedData = append(f.permittedData, dataSetItem)
		}
	}

	if f.Delegate != nil {
		f.Delegate.UpdateElectricalConnectionData(f)
	}

	return nil
}

func (f *ElectricalConnection) HandleRequest(ctrl spine.Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf spine.Feature) (*model.MsgCounterType, error) {
	switch fct {
	case model.FunctionEnumTypeElectricalConnectionParameterDescriptionListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestParameterDescriptionListData(ctrl, rf)
		}
		return nil, fmt.Errorf("electricalconnection.handleRequest: FunctionEnumTypeElectricalConnectionParameterDescriptionListData op not implemented: %s", op)

	case model.FunctionEnumTypeElectricalConnectionDescriptionListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestDescriptionListData(ctrl, rf)
		}
		return nil, fmt.Errorf("electricalconnection.handleRequest: FunctionEnumTypeElectricalConnectionDescriptionListData op not implemented: %s", op)

	case model.FunctionEnumTypeElectricalConnectionPermittedValueSetListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestPermittedValueSetData(ctrl, rf)
		}
		return nil, fmt.Errorf("electricalconnection.handleRequest: FunctionEnumTypeElectricalConnectionPermittedValueSetListData op not implemented: %s", op)

	}

	return nil, fmt.Errorf("electricalconnection.handleRequest: FunctionEnumType not implemented: %s", fct)
}

func (f *ElectricalConnection) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.ElectricalConnectionParameterDescriptionListData != nil:
		data := cmd.ElectricalConnectionParameterDescriptionListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyParameterDescriptionListData(ctrl, *data)

		case model.CmdClassifierTypeNotify:
			return f.replyParameterDescriptionListData(ctrl, *data)

		default:
			return fmt.Errorf("electricalconnection.Handle: ElectricalConnectionParameterDescriptionListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.ElectricalConnectionDescriptionListData != nil:
		data := cmd.ElectricalConnectionDescriptionListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyDescriptionListData(ctrl, *data)

		case model.CmdClassifierTypeNotify:
			return f.replyDescriptionListData(ctrl, *data)

		default:
			return fmt.Errorf("electricalconnection.Handle: ElectricalConnectionDescriptionListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.ElectricalConnectionPermittedValueSetListData != nil:
		data := cmd.ElectricalConnectionPermittedValueSetListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyPermittedValueSetData(ctrl, *data, isPartialForCmd, cmd.Filter)

		case model.CmdClassifierTypeNotify:
			return f.replyPermittedValueSetData(ctrl, *data, isPartialForCmd, cmd.Filter)

		default:
			return fmt.Errorf("electricalconnection.Handle: ElectricalConnectionPermittedValueSetListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("electricalconnection.Handle: CmdType not implemented: %s", populatedFields(cmd))
	}
}

func (f *ElectricalConnection) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	return ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
}
