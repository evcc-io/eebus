package feature

import (
	"fmt"
	"time"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type IncentiveConstraintsDataType struct {
	TariffID             uint
	MaxTiersPerTariff    uint
	MaxBoundariesPerTier uint
	MaxIncentivesPerTier uint
	SlotCountMax         uint
}

type IncentiveChargingSlot struct {
	Duration time.Duration
	Pricing  float64
}

type IncentiveChargingPlan struct {
	Duration time.Duration
	Slots    []IncentiveChargingSlot
}

type IncentiveDelegate interface {
	UpdateIncentiveConstraintsData(*IncentiveTable)
	UpdateIncentiveData(*IncentiveTable)
}

type IncentiveTableData interface {
	GetIncentiveTableConstraintsDataType() IncentiveConstraintsDataType
	WriteIncentiveTablePlanData(ctrl spine.Context, rf spine.Feature, chargingPlan IncentiveChargingPlan) error
}

type IncentiveTable struct {
	*spine.FeatureImpl
	Delegate        IncentiveDelegate
	constraintsData IncentiveConstraintsDataType
}

func NewIncentiveTableClient() spine.Feature {
	f := &IncentiveTable{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeIncentiveTable,
			Role: model.RoleTypeClient,
		},
	}

	return f
}

func (f *IncentiveTable) GetIncentiveTableConstraintsDataType() IncentiveConstraintsDataType {
	return f.constraintsData
}

func (f *IncentiveTable) requestDescriptionData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		IncentiveTableDescriptionData: &model.IncentiveTableDescriptionDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *IncentiveTable) replyDescriptionData(ctrl spine.Context, data model.IncentiveTableDescriptionDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":8}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":10}]},{"msgCounter":1593868},{"msgCounterReference":47},{"cmdClassifier":"reply"}]},{"payload":[
	// {"cmd":[[
	// 	{"incentiveTableDescriptionData":[
	// 		{"incentiveTableDescription":[
	// 			[{"tariffDescription":[{"tariffId":1},{"tariffWriteable":false},{"updateRequired":false},{"scopeType":"simpleIncentiveTable"}]}]
	// 		]}
	// 	]}
	// ]]}

	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":8}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":10}]},{"msgCounter":1598939},{"cmdClassifier":"notify"}]},{"payload":[
	// {"cmd":[[
	// 	{"incentiveTableDescriptionData":[
	// 		{"incentiveTableDescription":[
	// 			[{"tariffDescription":[{"tariffId":1},{"tariffWriteable":true},{"updateRequired":false},{"scopeType":"simpleIncentiveTable"}]}]
	// 		]}
	// 	]}
	// ]]}

	// f.descriptionData = nil
	// for _, item := range data.TimeSeriesDescriptionListData {
	// 	newItem := TimeSeriesDescriptionDataType{
	// 		LimitId:       uint(*item.LimitId),
	// 		LimitType:     model.LoadControlLimitTypeEnumType(*item.LimitType),
	// 		MeasurementId: uint(*item.MeasurementId),
	// 		ScopeType:     model.ScopeTypeEnumType(*item.ScopeType),
	// 	}
	// 	f.limitDescriptionData = append(f.limitDescriptionData, newItem)
	// }

	// if f.Delegate != nil {
	// 	f.Delegate.UpdateLoadControlLimitData(f)
	// }

	return nil
}

func (f *IncentiveTable) WriteDescriptionData(ctrl spine.Context, rf spine.Feature) error {
	// example data:
	// {"datagram":[{"header":[{"specificationVersion":"1.1.0"},{"addressSource":[{"device":"EVCC_HEMS"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":8}]},{"msgCounter":3016},{"cmdClassifier":"write"},{"ackRequest":true}]},{"payload":[
	// {"cmd":[[
	// 	{"incentiveTableDescriptionData":[
	// 		{"incentiveTableDescription":[[
	// 			{"tariffDescription":[{"tariffId":1}]},
	// 			{"tier":[
	// 				[
	// 		  		{"tierDescription":[{"tierId":1},{"tierType":"dynamicCost"}]},
	// 			  	{"boundaryDescription":[[{"boundaryId":1},{"boundaryType":"powerBoundary"},{"boundaryUnit":"W"}]]},
	// 		  		{"incentiveDescription":[[{"incentiveId":1},{"incentiveType":"absoluteCost"},{"currency":"EUR"}]]}
	// 				]
	// 			]}
	// 		]]}
	// 	]}
	// ]]}]}]}

	fct := model.FunctionType(model.FunctionEnumTypeIncentiveTableDescriptionData)
	partial := model.ElementTagType{}
	cmdControl := model.CmdControlType{Partial: &partial}
	filter := []model.FilterType{{CmdControl: &cmdControl}}

	// we limit this to one tariff, one tier, one boundary and one incentive
	tarrifId := model.TariffIdType(f.constraintsData.TariffID)
	tierId := model.TierIdType(1)
	tierType := model.TierTypeType(model.TierTypeEnumTypeDynamiccost)
	boundaryId := model.TierBoundaryIdType(1)
	boundaryType := model.TierBoundaryTypeType(model.TierBoundaryTypeEnumTypePowerBoundary)
	boundaryUnit := model.UnitOfMeasurementType(model.UnitOfMeasurementEnumTypeW)
	incentiveId := model.IncentiveIdType(1)
	incentiveType := model.IncentiveTypeType(model.IncentiveTypeEnumTypeAbsolutecost)
	currency := model.CurrencyType(model.CurrencyEnumTypeEur)

	res := []model.CmdType{{
		Function: &fct,
		Filter:   filter,
		IncentiveTableDescriptionData: &model.IncentiveTableDescriptionDataType{
			IncentiveTableDescription: []model.IncentiveTableDescriptionType{
				{
					TariffDescription: &model.TariffDescriptionDataType{
						TariffId: &tarrifId,
					},
					Tier: []model.IncentiveTableDescriptionTierType{
						{
							TierDescription: &model.TierDescriptionDataType{
								TierId:   &tierId,
								TierType: &tierType,
							},
							BoundaryDescription: []model.TierBoundaryDescriptionDataType{
								{
									BoundaryId:   &boundaryId,
									BoundaryType: &boundaryType,
									BoundaryUnit: &boundaryUnit,
								},
							},
							IncentiveDescription: []model.IncentiveDescriptionDataType{
								{
									IncentiveId:   &incentiveId,
									IncentiveType: &incentiveType,
									Currency:      &currency,
								},
							},
						},
					},
				},
			},
		},
	}}

	return ctrl.Write(spine.FeatureAddressType(f), spine.FeatureAddressType(rf), res)
}

func (f *IncentiveTable) requestConstraintsData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		IncentiveTableConstraintsData: &model.IncentiveTableConstraintsDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *IncentiveTable) replyConstraintsData(ctrl spine.Context, data model.IncentiveTableConstraintsDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":8}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":10}]},{"msgCounter":1593875},{"msgCounterReference":49},{"cmdClassifier":"reply"}]},{"payload":[
	// {"cmd":[[
	// 	{"incentiveTableConstraintsData":[
	// 		{"incentiveTableConstraints":[
	// 			[
	// 				{"tariff":[{"tariffId":1}]},
	// 				{"tariffConstraints":[{"maxTiersPerTariff":3},{"maxBoundariesPerTier":1},{"maxIncentivesPerTier":3}]},
	// 				{"incentiveSlotConstraints":[{"slotCountMax":29}]}
	// 			]
	// 		]}
	// 	]}
	// ]]}

	for _, constraints := range data.IncentiveTableConstraints {
		f.constraintsData = IncentiveConstraintsDataType{
			TariffID:             uint(*constraints.Tariff.TariffId),
			MaxTiersPerTariff:    uint(*constraints.TariffConstraints.MaxTiersPerTariff),
			MaxBoundariesPerTier: uint(*constraints.TariffConstraints.MaxBoundariesPerTier),
			MaxIncentivesPerTier: uint(*constraints.TariffConstraints.MaxIncentivesPerTier),
			SlotCountMax:         uint(*constraints.IncentiveSlotConstraints.SlotCountMax),
		}
	}

	if f.Delegate != nil {
		f.Delegate.UpdateIncentiveConstraintsData(f)
	}

	return nil
}

func (f *IncentiveTable) requestData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		IncentiveTableData: &model.IncentiveTableDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *IncentiveTable) replyData(ctrl spine.Context, data model.IncentiveTableDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":8}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":10}]},{"msgCounter":1593885},{"msgCounterReference":51},{"cmdClassifier":"reply"}]},{"payload":[
	// {"cmd":[[
	// 	{"incentiveTableData":[
	// 		{"incentiveTable":[]}
	// 	]}
	// ]]}

	// f.limitData = nil
	// for _, item := range data.LoadControlLimitData {
	// 	newItem := LoadControlLimitDatasetType{
	// 		LimitId:           uint(*item.LimitId),
	// 		IsLimitChangeable: *item.IsLimitChangeable,
	// 		IsLimitActive:     *item.IsLimitActive,
	// 		Value:             item.Value.GetValue(),
	// 	}
	// 	f.limitData = append(f.limitData, newItem)
	// }

	// if f.Delegate != nil {
	// 	f.Delegate.UpdateLoadControlLimitData(f)
	// }

	return nil
}

func (f *IncentiveTable) WriteIncentiveTablePlanData(ctrl spine.Context, rf spine.Feature, chargingPlan IncentiveChargingPlan) error {
	tariffId := model.TariffIdType(1)
	var incentiveSlots []model.IncentiveTableIncentiveSlotType

	tierId := model.TierIdType(1)
	boundaryId := model.TierBoundaryIdType(1)
	incentiveId := model.IncentiveIdType(1)
	var totalDuration time.Duration

	for index, item := range chargingPlan.Slots {
		timeInterval := model.TimeTableDataType{
			StartTime: &model.AbsoluteOrRecurringTimeType{
				Relative: model.NewISO8601Duration(totalDuration),
			},
		}
		totalDuration += item.Duration

		// last item needs an endTime
		if index == len(chargingPlan.Slots)-1 {
			timeInterval.EndTime = &model.AbsoluteOrRecurringTimeType{
				Relative: model.NewISO8601Duration(totalDuration),
			}
		}

		incentiveSlot := model.IncentiveTableIncentiveSlotType{
			TimeInterval: &timeInterval,
			Tier: []model.IncentiveTableTierType{
				{
					Tier: &model.TierDataType{
						TierId: &tierId,
					},
					Boundary: []model.TierBoundaryDataType{
						{
							BoundaryId:         &boundaryId,
							LowerBoundaryValue: model.NewScaledNumberType(0),
						},
					},
					Incentive: []model.IncentiveDataType{
						{
							IncentiveId: &incentiveId,
							Value:       model.NewScaledNumberType(item.Pricing),
						},
					},
				},
			},
		}
		incentiveSlots = append(incentiveSlots, incentiveSlot)
	}

	res := []model.CmdType{{
		IncentiveTableData: &model.IncentiveTableDataType{
			IncentiveTable: []model.IncentiveTableType{
				{
					Tariff:        &model.TariffDataType{TariffId: &tariffId},
					IncentiveSlot: incentiveSlots,
				},
			},
		},
	}}

	return ctrl.Write(spine.FeatureAddressType(f), spine.FeatureAddressType(rf), res)
}

func (f *IncentiveTable) HandleRequest(ctrl spine.Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf spine.Feature) (*model.MsgCounterType, error) {
	switch fct {
	case model.FunctionEnumTypeIncentiveTableDescriptionData:
		if op == model.CmdClassifierTypeRead {
			return f.requestDescriptionData(ctrl, rf)
		}
		return nil, fmt.Errorf("timeseries.handleRequest: FunctionEnumTypeIncentiveTableDescriptionData op not implemented: %s", op)

	case model.FunctionEnumTypeIncentiveTableConstraintsData:
		if op == model.CmdClassifierTypeRead {
			return f.requestConstraintsData(ctrl, rf)
		}
		return nil, fmt.Errorf("incentivetable.handleRequest: FunctionEnumTypeIncentiveTableConstraintsData op not implemented: %s", op)

	case model.FunctionEnumTypeIncentiveTableData:
		if op == model.CmdClassifierTypeRead {
			return f.requestData(ctrl, rf)
		}
		return nil, fmt.Errorf("incentivetable.handleRequest: FunctionEnumTypeIncentiveTableData op not implemented: %s", op)
	}

	return nil, fmt.Errorf("incentivetable.handleRequest: FunctionEnumType not implemented: %s", fct)
}

func (f *IncentiveTable) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.IncentiveTableDescriptionData != nil:
		data := cmd.IncentiveTableDescriptionData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyDescriptionData(ctrl, *data)
		case model.CmdClassifierTypeNotify:
			return f.replyDescriptionData(ctrl, *data)
		default:
			return fmt.Errorf("incentivetable.handle: IncentiveTableDescriptionData CmdClassifierType not implemented: %s", op)
		}

	case cmd.IncentiveTableConstraintsData != nil:
		data := cmd.IncentiveTableConstraintsData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyConstraintsData(ctrl, *data)
		case model.CmdClassifierTypeNotify:
			return f.replyConstraintsData(ctrl, *data)
		default:
			return fmt.Errorf("incentivetable.handle: IncentiveTableConstraintsData CmdClassifierType not implemented: %s", op)
		}

	case cmd.IncentiveTableData != nil:
		data := cmd.IncentiveTableData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyData(ctrl, *data)
		case model.CmdClassifierTypeNotify:
			return f.replyData(ctrl, *data)
		default:
			return fmt.Errorf("incentivetable.handle: IncentiveTableData CmdClassifierType not implemented: %s", op)
		}

	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("incentivetable.Handle: CmdType not implemented: %s", populatedFields(cmd))
	}
}

func (f *IncentiveTable) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	return ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
}
