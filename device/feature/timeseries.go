package feature

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type TimeSeriesDescriptionListDatasetType struct {
	TimeSeriesId     uint
	TimeSeriesType   model.TimeSeriesTypeEnumType
	IsSeriesWritable bool
	Unit             string
}

type TimeSeriesDatasetType struct {
	TimeSeriesId    uint
	TimePeriod      *model.TimePeriodType
	TimeSeriesSlots []model.TimeSeriesSlotType
}

type TimeSeriesData interface {
	GetTimeSeriesDescriptionData() []TimeSeriesDescriptionListDatasetType
	GetTimeSeriesData() []TimeSeriesDatasetType
	WriteTimeSeriesPlanData(ctrl spine.Context, rf spine.Feature, seriesData []model.TimeSeriesDataType) error
}

type TimeSeriesDelegate interface {
	UpdateTimeSeriesData(*TimeSeries)
}

type TimeSeries struct {
	*spine.FeatureImpl
	Delegate                  TimeSeriesDelegate
	timeSeriesDescriptionData []TimeSeriesDescriptionListDatasetType
	timeSeriesData            []TimeSeriesDatasetType
}

func NewTimeSeriesClient() spine.Feature {
	f := &TimeSeries{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeTimeSeries,
			Role: model.RoleTypeClient,
		},
	}

	return f
}

func (f *TimeSeries) GetTimeSeriesDescriptionData() []TimeSeriesDescriptionListDatasetType {
	return f.timeSeriesDescriptionData
}

func (f *TimeSeries) GetTimeSeriesData() []TimeSeriesDatasetType {
	return f.timeSeriesData
}

func (f *TimeSeries) replyConstraintsData(ctrl spine.Context, data model.TimeSeriesConstraintsDataType, isPartialCmd bool) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":1226},{"cmdClassifier":"notify"}]},{"payload":[
	// {"cmd":[[
	//   {"timeSeriesConstraintsListData":[
	// 	   {"timeSeriesConstraintsData":[
	// 		   [{"timeSeriesId":1},{"slotCountMax":29}]
	// 	   ]}
	//   ]}
	// ]]}]}]}}]}

	// TODO: implement processing

	return nil
}

func (f *TimeSeries) requestDescriptionListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		TimeSeriesDescriptionListData: &model.TimeSeriesDescriptionListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *TimeSeries) replyDescriptionListData(ctrl spine.Context, data model.TimeSeriesDescriptionListDataType, isPartialForCmd bool) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":1593590},{"msgCounterReference":25},{"cmdClassifier":"reply"}]},{"payload":[
	// {"cmd":[[
	// 	{"timeSeriesDescriptionListData":[
	// 		{"timeSeriesDescriptionData":[
	// 			[{"timeSeriesId":1},{"timeSeriesType":"constraints"},{"timeSeriesWriteable":false},{"updateRequired":false},{"unit":"W"}],
	// 			[{"timeSeriesId":2},{"timeSeriesType":"plan"},{"timeSeriesWriteable":false},{"unit":"W"}],
	// 			[{"timeSeriesId":3},{"timeSeriesType":"singleDemand"},{"timeSeriesWriteable":false},{"unit":"Wh"}]
	// 		]}
	// 	]}
	// ]]}

	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":1598941},{"cmdClassifier":"notify"}]},{"payload":[
	// {"cmd":[[
	// 	{"timeSeriesDescriptionListData":[
	// 		{"timeSeriesDescriptionData":[
	// 			[{"timeSeriesId":1},{"timeSeriesType":"constraints"},{"timeSeriesWriteable":true},{"updateRequired":false},{"unit":"W"}],
	// 			[{"timeSeriesId":2},{"timeSeriesType":"plan"},{"timeSeriesWriteable":false},{"unit":"W"}],
	// 			[{"timeSeriesId":3},{"timeSeriesType":"singleDemand"},{"timeSeriesWriteable":false},{"unit":"Wh"}]
	// 		]}
	// 	]}
	// ]]}

	if !isPartialForCmd {
		f.timeSeriesDescriptionData = nil
	}
	for _, item := range data.TimeSeriesDescriptionData {
		newItem := TimeSeriesDescriptionListDatasetType{
			TimeSeriesId:     uint(*item.TimeSeriesId),
			TimeSeriesType:   model.TimeSeriesTypeEnumType(*item.TimeSeriesType),
			IsSeriesWritable: *item.TimeSeriesWriteable,
			Unit:             string(*item.Unit),
		}
		f.timeSeriesDescriptionData = append(f.timeSeriesDescriptionData, newItem)
		if !isPartialForCmd {
			f.timeSeriesDescriptionData = append(f.timeSeriesDescriptionData, newItem)
		} else {
			replaceIndex := -1
			for index, element := range f.timeSeriesDescriptionData {
				if element.TimeSeriesId == newItem.TimeSeriesId {
					replaceIndex = index
				}
			}
			if replaceIndex != -1 {
				f.timeSeriesDescriptionData[replaceIndex] = newItem
			} else {
				f.timeSeriesDescriptionData = append(f.timeSeriesDescriptionData, newItem)
			}
		}
	}

	if f.Delegate != nil {
		f.Delegate.UpdateTimeSeriesData(f)
	}

	return nil
}

func (f *TimeSeries) requestListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *TimeSeries) replyListData(ctrl spine.Context, data model.TimeSeriesListDataType, isPartialForCmd bool) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":1593609},{"msgCounterReference":43},{"cmdClassifier":"reply"}]},{"payload":[
	// {"cmd":[[
	// 	{"timeSeriesListData":[
	// 		{"timeSeriesData":[
	// 			[
	// 				{"timeSeriesId":1}
	// 			],
	//			[
	//				{"timeSeriesId":2},
	//				{"timePeriod":[{"startTime":"PT0S"}]},
	//				{"timeSeriesSlot":[[{"timeSeriesSlotId":0},{"duration":"P1D"},{"maxValue":[{"number":0},{"scale":0}]}]]}
	//			],
	//			[
	//				{"timeSeriesId":3},
	//				{"timePeriod":[{"startTime":"PT0S"}]},
	//				{"timeSeriesSlot":[[{"timeSeriesSlotId":1},{"value":[{"number":0},{"scale":0}]}]]}
	// 			]
	// 		]}
	// 	]}
	// ]]}

	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":1599878},{"cmdClassifier":"notify"}]},{"payload":[
	// {"cmd":[[
	// 	{"function":"timeSeriesListData"},
	//  {"filter":[[{"cmdControl":[{"partial":[]}]}]]},
	// 	{"timeSeriesListData":[
	// 		{"timeSeriesData":[
	// 			[
	// 				{"timeSeriesId":2},
	// 				{"timePeriod":[{"startTime":"PT0S"}]},
	// 				{"timeSeriesSlot":[
	// 				  [{"timeSeriesSlotId":0},{"duration":"PT1H22M37S"},{"maxValue":[{"number":11052},{"scale":0}]}],
	// 				  [{"timeSeriesSlotId":1},{"duration":"PT19M43S"},{"maxValue":[{"number":6526},{"scale":0}]}],
	// 				  [{"timeSeriesSlotId":2},{"duration":"P1D"},{"maxValue":[{"number":0},{"scale":0}]}]
	// 				]}
	// 			]
	// 		]}
	// 	]}
	// ]]}

	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":1703680},{"cmdClassifier":"notify"},{"ackRequest":true}]},{"payload":[
	// {"cmd":[[
	// 	{"function":"timeSeriesListData"},
	// 	{"filter":[[{"cmdControl":[{"partial":[]}]}]]},
	// 	{"timeSeriesListData":[
	// 		{"timeSeriesData":[
	// 			[{"timeSeriesId":3},{"timePeriod":[{"startTime":"PT0S"}]},{"timeSeriesSlot":[[{"timeSeriesSlotId":1},{"value":[{"number":16200},{"scale":0}]}]]}]
	// 		]}
	// 	]}
	// ]]}

	if !isPartialForCmd {
		f.timeSeriesData = nil
	}
	for _, item := range data.TimeSeriesData {
		newItem := TimeSeriesDatasetType{
			TimeSeriesId: uint(*item.TimeSeriesId),
		}
		if item.TimePeriod != nil {
			newItem.TimePeriod = item.TimePeriod
		}
		if item.TimeSeriesSlot != nil {
			newItem.TimeSeriesSlots = item.TimeSeriesSlot
		}

		if !isPartialForCmd {
			f.timeSeriesData = append(f.timeSeriesData, newItem)
		} else {
			replaceIndex := -1
			for index, element := range f.timeSeriesData {
				if element.TimeSeriesId == newItem.TimeSeriesId {
					replaceIndex = index
				}
			}
			if replaceIndex != -1 {
				f.timeSeriesData[replaceIndex] = newItem
			} else {
				f.timeSeriesData = append(f.timeSeriesData, newItem)
			}
		}
	}

	if f.Delegate != nil {
		f.Delegate.UpdateTimeSeriesData(f)
	}

	return nil
}

func (f *TimeSeries) WriteTimeSeriesPlanData(ctrl spine.Context, rf spine.Feature, seriesData []model.TimeSeriesDataType) error {
	fct := model.FunctionType(model.FunctionEnumTypeTimeSeriesListData)
	partial := model.ElementTagType{}
	cmdControl := model.CmdControlType{Partial: &partial}
	filter := []model.FilterType{{CmdControl: &cmdControl}}

	res := []model.CmdType{{
		Function: &fct,
		Filter:   filter,
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: seriesData,
		},
	}}

	return ctrl.Write(spine.FeatureAddressType(f), spine.FeatureAddressType(rf), res)
}

func (f *TimeSeries) HandleRequest(ctrl spine.Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf spine.Feature) (*model.MsgCounterType, error) {
	switch fct {
	case model.FunctionEnumTypeTimeSeriesDescriptionListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestDescriptionListData(ctrl, rf)
		}
		return nil, fmt.Errorf("timeseries.handleRequest: FunctionEnumTypeTimeSeriesDescriptionListData op not implemented: %s", op)

	case model.FunctionEnumTypeTimeSeriesListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestListData(ctrl, rf)
		}
		return nil, fmt.Errorf("timeseries.handleRequest: FunctionEnumTypeTimeSeriesListData op not implemented: %s", op)
	}

	return nil, fmt.Errorf("timeseries.handleRequest: FunctionEnumType not implemented: %s", fct)
}

func (f *TimeSeries) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.TimeSeriesConstraintsData != nil:
		data := cmd.TimeSeriesConstraintsData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyConstraintsData(ctrl, *data, isPartialForCmd)
		case model.CmdClassifierTypeNotify:
			return f.replyConstraintsData(ctrl, *data, isPartialForCmd)
		default:
			return fmt.Errorf("timeseries.handle: TimeSeriesConstraintsData CmdClassifierType not implemented: %s", op)
		}
	case cmd.TimeSeriesDescriptionListData != nil:
		data := cmd.TimeSeriesDescriptionListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyDescriptionListData(ctrl, *data, isPartialForCmd)
		case model.CmdClassifierTypeNotify:
			return f.replyDescriptionListData(ctrl, *data, isPartialForCmd)
		default:
			return fmt.Errorf("timeseries.handle: TimeSeriesDescriptionListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.TimeSeriesListData != nil:
		data := cmd.TimeSeriesListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyListData(ctrl, *data, isPartialForCmd)
		case model.CmdClassifierTypeNotify:
			return f.replyListData(ctrl, *data, isPartialForCmd)
		default:
			return fmt.Errorf("timeseries.handle: TimeSeriesListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("timeseries.Handle: CmdType not implemented: %s", populatedFields(cmd))
	}
}

func (f *TimeSeries) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	return ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
}
