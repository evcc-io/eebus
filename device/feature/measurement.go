package feature

import (
	"fmt"
	"time"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type MeasurementDatasetDefinitionsType struct {
	MeasurementId   uint
	MeasurementType model.MeasurementTypeEnumType
	ScopeType       model.ScopeTypeEnumType
}

type MeasurementConstraintsDefinitionsType struct {
	MeasurementId uint
	MinValue      float64
	MaxValue      float64
	StepSize      float64
}

type MeasurementDatasetDataType struct {
	Timestamp     time.Time
	MeasurementId uint
	Value         float64
}

type MeasurementData interface {
	GetMeasurementDescription() []MeasurementDatasetDefinitionsType
	GetMeasurementData() []MeasurementDatasetDataType
}

type MeasurementDelegate interface {
	UpdateMeasurementData(*Measurement)
}

type Measurement struct {
	*spine.FeatureImpl
	Delegate               MeasurementDelegate
	datasetDefinitions     []MeasurementDatasetDefinitionsType
	constraintsDefinitions []MeasurementConstraintsDefinitionsType
	datasetData            []MeasurementDatasetDataType
}

func NewMeasurementClient() spine.Feature {
	f := &Measurement{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeMeasurement,
			Role: model.RoleTypeClient,
		},
	}

	return f
}

func (f *Measurement) GetMeasurementDescription() []MeasurementDatasetDefinitionsType {
	return f.datasetDefinitions
}

func (f *Measurement) GetMeasurementData() []MeasurementDatasetDataType {
	return f.datasetData
}

func (f *Measurement) requestDescriptionListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		MeasurementDescriptionListData: &model.MeasurementDescriptionListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *Measurement) replyDescriptionListData(ctrl spine.Context, data model.MeasurementDescriptionListDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":3}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":3}]},{"msgCounter":6977},{"msgCounterReference":15},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"measurementDescriptionListData":[{"measurementDescriptionData":[[{"measurementId":1},{"measurementType":"current"},{"commodityType":"electricity"},{"unit":"A"},{"scopeType":"acCurrent"}],[{"measurementId":4},{"measurementType":"power"},{"commodityType":"electricity"},{"unit":"W"},{"scopeType":"acPower"}],[{"measurementId":7},{"measurementType":"energy"},{"commodityType":"electricity"},{"unit":"Wh"},{"scopeType":"charge"}]]}]}]]}]}]}}]}

	f.datasetDefinitions = nil
	for _, item := range data.MeasurementDescriptionData {
		newItem := MeasurementDatasetDefinitionsType{
			MeasurementId:   uint(*item.MeasurementId),
			MeasurementType: model.MeasurementTypeEnumType(*item.MeasurementType),
			ScopeType:       model.ScopeTypeEnumType(*item.ScopeType),
		}
		f.datasetDefinitions = append(f.datasetDefinitions, newItem)
	}

	return nil
}

func (f *Measurement) requestConstraintsListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		MeasurementConstraintsListData: &model.MeasurementConstraintsListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *Measurement) replyConstraintsListData(ctrl spine.Context, data model.MeasurementConstraintsListDataType) error {
	// example data:

	f.datasetDefinitions = nil
	for _, item := range data.MeasurementConstraintsData {
		newItem := MeasurementConstraintsDefinitionsType{
			MeasurementId: uint(*item.MeasurementId),
			MinValue:      item.ValueRangeMin.GetValue(),
			MaxValue:      item.ValueRangeMax.GetValue(),
			StepSize:      item.ValueStepSize.GetValue(),
		}
		f.constraintsDefinitions = append(f.constraintsDefinitions, newItem)
	}

	return nil
}

func (f *Measurement) requestListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		MeasurementListData: &model.MeasurementListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *Measurement) replyListData(ctrl spine.Context, data model.MeasurementListDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":3}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":3}]},{"msgCounter":15971},{"msgCounterReference":33},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"measurementListData":[{"measurementData":[[{"measurementId":1},{"valueType":"value"},{"timestamp":"2021-04-23T12:39:19.037Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":4},{"valueType":"value"},{"timestamp":"2021-04-23T12:39:19.037Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":2},{"valueType":"value"},{"timestamp":"2021-04-23T12:39:19.037Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":5},{"valueType":"value"},{"timestamp":"2021-04-23T12:39:19.037Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":3},{"valueType":"value"},{"timestamp":"2021-04-23T12:39:19.037Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":6},{"valueType":"value"},{"timestamp":"2021-04-23T12:39:19.037Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}],[{"measurementId":7},{"valueType":"value"},{"timestamp":"2021-04-23T12:39:19.037Z"},{"value":[{"number":0},{"scale":0}]},{"valueSource":"measuredValue"}]]}]}]]}]}]}}]}

	f.datasetData = nil
	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil || item.Value == nil {
			continue
		}

		timestamp := time.Time{}
		if item.Timestamp != nil {
			stamp, err := time.Parse(time.RFC3339, *item.Timestamp)
			if err == nil {
				timestamp = stamp
			}
		}

		newItem := MeasurementDatasetDataType{
			MeasurementId: uint(*item.MeasurementId),
			Timestamp:     timestamp,
			Value:         item.Value.GetValue(),
		}
		f.datasetData = append(f.datasetData, newItem)
	}

	if f.Delegate != nil {
		f.Delegate.UpdateMeasurementData(f)
	}

	return nil
}

func (f *Measurement) notifyListData(ctrl spine.Context, data model.MeasurementListDataType, isPartialForCmd bool) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.3.0"},{"addressSource":[{"device":"d:_i:47859_Elli-Wallbox-2019A0OV8H"},{"entity":[1,1]},{"feature":11}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":3}]},{"msgCounter":811},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"function":"measurementListData"},{"filter":[[{"cmdControl":[{"partial":[]}]}]]},{"measurementListData":[{"measurementData":[[{"measurementId":0},{"valueType":"value"},{"value":[{"number":608},{"scale":-2}]}],[{"measurementId":1},{"valueType":"value"},{"value":[{"number":587},{"scale":-2}]}],[{"measurementId":2},{"valueType":"value"},{"value":[{"number":604},{"scale":-2}]}]]}]}]]}]}]}}]}

	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil || item.Value == nil {
			continue
		}

		timestamp := time.Now()
		if item.Timestamp != nil {
			stamp, err := time.Parse(time.RFC3339, *item.Timestamp)
			if err == nil {
				timestamp = stamp
			}
		}

		var updatedDataSetItem MeasurementDatasetDataType
		var updatedDataSetIndex int
		itemFound := false
		for index, datasetItem := range f.datasetData {
			if datasetItem.MeasurementId == uint(*item.MeasurementId) {
				updatedDataSetItem = datasetItem
				updatedDataSetIndex = index
				itemFound = true
				break
			}
		}

		if !itemFound {
			newItem := MeasurementDatasetDataType{
				MeasurementId: uint(*item.MeasurementId),
				Timestamp:     timestamp,
				Value:         item.Value.GetValue(),
			}
			f.datasetData = append(f.datasetData, newItem)
		} else {
			updatedDataSetItem.Timestamp = timestamp
			updatedDataSetItem.Value = item.Value.GetValue()
			f.datasetData[updatedDataSetIndex] = updatedDataSetItem
		}
	}

	if f.Delegate != nil {
		f.Delegate.UpdateMeasurementData(f)
	}

	return nil
}

func (f *Measurement) HandleRequest(ctrl spine.Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf spine.Feature) (*model.MsgCounterType, error) {
	switch fct {
	case model.FunctionEnumTypeMeasurementDescriptionListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestDescriptionListData(ctrl, rf)
		}
		return nil, fmt.Errorf("measurement.handleRequest: FunctionEnumTypeMeasurementDescriptionListData op not implemented: %s", op)

	case model.FunctionEnumTypeMeasurementConstraintsListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestConstraintsListData(ctrl, rf)
		}
		return nil, fmt.Errorf("measurement.handleRequest: FunctionEnumTypeMeasurementConstraintsListData op not implemented: %s", op)

	case model.FunctionEnumTypeMeasurementListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestListData(ctrl, rf)
		}
		return nil, fmt.Errorf("measurement.handleRequest: FunctionEnumTypeMeasurementListData op not implemented: %s", op)
	}

	return nil, fmt.Errorf("measurement.handleRequest: FunctionEnumType not implemented: %s", fct)
}

func (f *Measurement) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.MeasurementDescriptionListData != nil:
		data := cmd.MeasurementDescriptionListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyDescriptionListData(ctrl, *data)

		case model.CmdClassifierTypeNotify:
			return f.replyDescriptionListData(ctrl, *data)

		default:
			return fmt.Errorf("measurement.handle: MeasurementDescriptionListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.MeasurementConstraintsListData != nil:
		data := cmd.MeasurementConstraintsListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyConstraintsListData(ctrl, *data)

		case model.CmdClassifierTypeNotify:
			return f.replyConstraintsListData(ctrl, *data)

		default:
			return fmt.Errorf("measurement.handle: MeasurementConstraintsListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.MeasurementListData != nil:
		data := cmd.MeasurementListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyListData(ctrl, *data)

		case model.CmdClassifierTypeNotify:
			return f.notifyListData(ctrl, *data, isPartialForCmd)

		default:
			return fmt.Errorf("measurement.handle: MeasurementListData CmdClassifierType not implemented: %s", op)
		}

	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("measurement.handle: CmdType not implemented: %s", populatedFields(cmd))
	}
}

func (f *Measurement) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	return ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
}
