package feature

import (
	"fmt"
	"time"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type TimeSeriesDescriptionListDatasetType struct {
	TimeSeriesId     uint
	TimeSeriesType   model.TimeSeriesTypeEnumType
	IsSeriesWritable bool
	UpdateRequired   bool
	Unit             string
}

type TimeSeriesDatasetType struct {
	TimeSeriesId    uint
	TimePeriod      *model.TimePeriodType
	TimeSeriesSlots []model.TimeSeriesSlotType
}

type TimeSeriesChargingSlot struct {
	Duration time.Duration
	MaxValue float64 // Watts
}

type TimeSeriesChargingPlan struct {
	Duration time.Duration
	Slots    []TimeSeriesChargingSlot
}

type TimeSeriesData interface {
	GetTimeSeriesDescriptionData() []TimeSeriesDescriptionListDatasetType
	GetTimeSeriesData() []TimeSeriesDatasetType
	GetTimeSeriesTypeForId(id uint) model.TimeSeriesTypeEnumType
	WriteTimeSeriesPlanData(ctrl spine.Context, rf spine.Feature, chargingPlan TimeSeriesChargingPlan) error
}

type TimeSeriesDelegate interface {
	UpdateTimeSeriesDescriptionData(*TimeSeries)
	UpdateTimeSeriesData(*TimeSeries, TimeSeriesDatasetType)
}

type TimeSeries struct {
	*spine.FeatureImpl
	Delegate                  TimeSeriesDelegate
	timeSeriesDescriptionData []TimeSeriesDescriptionListDatasetType
	timeSeriesData            []TimeSeriesDatasetType
	timeSeriesMaxCount        uint
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

func (f *TimeSeries) GetTimeSeriesTypeForId(id uint) (model.TimeSeriesTypeEnumType, error) {
	for _, d := range f.timeSeriesDescriptionData {
		if d.TimeSeriesId == id {
			return d.TimeSeriesType, nil
		}
	}

	return "", fmt.Errorf("timeseries.GetTimeSeriesTypeForId: id not found: %d", id)
}

func (f *TimeSeries) getTimeSeriesIdForType(typeEnum model.TimeSeriesTypeEnumType) uint {
	for _, d := range f.timeSeriesDescriptionData {
		if d.TimeSeriesType == typeEnum {
			return d.TimeSeriesId
		}
	}

	return 0
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

	f.timeSeriesMaxCount = uint(*data.SlotCountMax)

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

	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00009463"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":7497510},{"cmdClassifier":"notify"}]},{"payload":[
	// {"cmd":[[
	// 	{"timeSeriesDescriptionListData":[
	// 		{"timeSeriesDescriptionData":[
	// 			[{"timeSeriesId":1},{"timeSeriesType":"constraints"},{"timeSeriesWriteable":true},{"updateRequired":true},{"unit":"W"}],
	// 			[{"timeSeriesId":2},{"timeSeriesType":"plan"},{"timeSeriesWriteable":false},{"unit":"W"}],
	// 			[{"timeSeriesId":3},{"timeSeriesType":"singleDemand"},{"timeSeriesWriteable":false},{"unit":"Wh"}]
	// 		]}
	// 	]}
	// ]]}]}]}}]}

	if !isPartialForCmd {
		f.timeSeriesDescriptionData = nil
	}
	for _, item := range data.TimeSeriesDescriptionData {
		if item.TimeSeriesId == nil || item.TimeSeriesType == nil || item.Unit == nil {
			continue
		}

		isWriteable := false
		if item.TimeSeriesWriteable != nil {
			isWriteable = *item.TimeSeriesWriteable
		}

		newItem := TimeSeriesDescriptionListDatasetType{
			TimeSeriesId:     uint(*item.TimeSeriesId),
			TimeSeriesType:   model.TimeSeriesTypeEnumType(*item.TimeSeriesType),
			IsSeriesWritable: isWriteable,
			Unit:             string(*item.Unit),
		}
		if item.UpdateRequired != nil {
			newItem.UpdateRequired = *item.UpdateRequired
		} else {
			newItem.UpdateRequired = false
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
		f.Delegate.UpdateTimeSeriesDescriptionData(f)
	}

	return nil
}

func (f *TimeSeries) WriteTimeSeriesDescriptionData(ctrl spine.Context, rf spine.Feature) error {
	id1 := model.TimeSeriesIdType(1)
	type1 := model.TimeSeriesTypeType(model.TimeSeriesTypeEnumTypeConstraints)
	writable1 := false
	unit1 := model.UnitOfMeasurementType(model.UnitOfMeasurementEnumTypeW)

	id2 := model.TimeSeriesIdType(2)
	type2 := model.TimeSeriesTypeType(model.TimeSeriesTypeEnumTypePlan)
	writable2 := false
	unit2 := model.UnitOfMeasurementType(model.UnitOfMeasurementEnumTypeW)

	id3 := model.TimeSeriesIdType(3)
	type3 := model.TimeSeriesTypeType(model.TimeSeriesTypeEnumTypeSingleDemand)
	writable3 := false
	unit3 := model.UnitOfMeasurementType(model.UnitOfMeasurementEnumTypeWh)

	res := []model.CmdType{{
		TimeSeriesDescriptionListData: &model.TimeSeriesDescriptionListDataType{
			TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
				{
					TimeSeriesId:        &id1,
					TimeSeriesType:      &type1,
					TimeSeriesWriteable: &writable1,
					Unit:                &unit1,
				},
				{
					TimeSeriesId:        &id2,
					TimeSeriesType:      &type2,
					TimeSeriesWriteable: &writable2,
					Unit:                &unit2,
				},
				{
					TimeSeriesId:        &id3,
					TimeSeriesType:      &type3,
					TimeSeriesWriteable: &writable3,
					Unit:                &unit3,
				},
			},
		},
	}}

	return ctrl.Write(spine.FeatureAddressType(f), spine.FeatureAddressType(rf), res)
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
	// 		 [
	//			 {"timeSeriesId":3},
	//			 {"timePeriod":[{"startTime":"PT0S"}]},
	//			 {"timeSeriesSlot":[
	//				 [{"timeSeriesSlotId":1},{"value":[{"number":16200},{"scale":0}]}]
	//				]}
	//			]
	// 		]}
	// 	]}
	// ]]}

	// returning the sent time series constraints
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":7497520},{"cmdClassifier":"notify"}]},{"payload":[
	// {"cmd":[[
	// 	{"function":"timeSeriesListData"},
	// 	{"filter":[[{"cmdControl":[{"partial":[]}]}]]},
	// 	{"timeSeriesListData":[
	// 		{"timeSeriesData":[
	// 			[
	// 				{"timeSeriesId":1},
	// 				{"timePeriod":[{"startTime":"PT0S"},{"endTime":"P1D"}]},
	// 				{"timeSeriesSlot":[
	// 					[{"timeSeriesSlotId":0},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":1},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":2},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":3},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":4},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":5},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":6},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":7},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":8},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":9},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":10},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":11},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":12},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":13},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":14},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":15},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":16},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":17},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":18},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":19},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":20},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":21},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":22},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":23},{"duration":"PT1H"},{"maxValue":[{"number":11040},{"scale":0}]}]
	// 				]}
	// 			]
	// 		]}
	// 	]}
	// ]]}
	// ]}]}}]}

	// actual charging plan report, direct charging
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":7}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":9}]},{"msgCounter":7497547},{"cmdClassifier":"notify"}]},{"payload":[
	// {"cmd":[[
	// 	{"function":"timeSeriesListData"},
	// 	{"filter":[[{"cmdControl":[{"partial":[]}]}]]},
	// 	{"timeSeriesListData":[
	// 		{"timeSeriesData":[
	// 			[
	// 				{"timeSeriesId":2},
	// 				{"timePeriod":[{"startTime":"PT0S"}]},
	// 				{"timeSeriesSlot":[
	// 					[{"timeSeriesSlotId":0},{"duration":"PT10H41M48S"},{"maxValue":[{"number":4547},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":1},{"duration":"PT45M1S"},{"maxValue":[{"number":2526},{"scale":0}]}],
	// 					[{"timeSeriesSlotId":2},{"duration":"P1D"},{"maxValue":[{"number":0},{"scale":0}]}]
	// 				]}
	// 			]
	// 		]}
	// 	]}
	// ]]}
	// ]}]}}]}

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
		if f.Delegate != nil {
			f.Delegate.UpdateTimeSeriesData(f, newItem)
		}
	}

	return nil
}

// sends a charging plan to the EVSE
// duration is the duration of the charging plan in seconds
func (f *TimeSeries) WriteTimeSeriesPlanData(ctrl spine.Context, rf spine.Feature, chargingPlan TimeSeriesChargingPlan) error {
	seriesId := model.TimeSeriesIdType(f.getTimeSeriesIdForType(model.TimeSeriesTypeEnumTypeConstraints))
	startTime := model.NewISO8601Duration(time.Duration(0) * time.Second)
	endTime := model.NewISO8601Duration(chargingPlan.Duration)
	seriesPeriod := model.TimePeriodType{StartTime: startTime, EndTime: endTime}

	timeSeriesData := model.TimeSeriesDataType{
		TimeSeriesId: &seriesId,
		TimePeriod:   &seriesPeriod,
	}

	var timeSeriesSlots []model.TimeSeriesSlotType
	for index, item := range chargingPlan.Slots {
		slotId := model.TimeSeriesSlotIdType(index)
		duration := model.NewISO8601Duration(item.Duration)

		timeSeriesSlots = append(timeSeriesSlots, model.TimeSeriesSlotType{
			TimeSeriesSlotId: &slotId,
			Duration:         duration,
			MaxValue:         model.NewScaledNumberType(item.MaxValue),
		})
	}
	timeSeriesData.TimeSeriesSlot = timeSeriesSlots

	res := []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: []model.TimeSeriesDataType{timeSeriesData},
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
