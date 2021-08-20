package feature

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type IdentificationDatasetDataType struct {
	IdentificationId    uint
	IdentificationType  model.IdentificationTypeEnumType
	IdentificationValue string
}

type IdentificationDelegate interface {
	UpdateIdentificationData(*Identification, []IdentificationDatasetDataType)
}

type Identification struct {
	*spine.FeatureImpl
	Delegate    IdentificationDelegate
	datasetData []IdentificationDatasetDataType
}

func NewIdentificationClient() spine.Feature {
	f := &Identification{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeIdentification,
			Role: model.RoleTypeClient,
		},
	}

	return f
}

func (f *Identification) requestListData(ctrl spine.Context, rf spine.Feature) (*model.MsgCounterType, error) {
	res := []model.CmdType{{
		IdentificationListData: &model.IdentificationListDataType{},
	}}

	return ctrl.Request(model.CmdClassifierTypeRead, *spine.FeatureAddressType(f), *spine.FeatureAddressType(rf), true, res)
}

func (f *Identification) replyListData(ctrl spine.Context, data model.IdentificationListDataType) error {
	// example data:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[1,1]},{"feature":10}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[1]},{"feature":7}]},{"msgCounter":21495},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"identificationListData":[{"identificationData":[[{"identificationId":0},{"identificationType":"eui48"},{"identificationValue":"F0:7F:0C:07:9B:C7"}]]}]}]]}]}]}}]}

	f.datasetData = nil
	for _, item := range data.IdentificationData {
		newItem := IdentificationDatasetDataType{
			IdentificationId:    uint(*item.IdentificationId),
			IdentificationType:  model.IdentificationTypeEnumType(*item.IdentificationType),
			IdentificationValue: string(*item.IdentificationValue),
		}
		f.datasetData = append(f.datasetData, newItem)
	}

	if f.Delegate != nil {
		f.Delegate.UpdateIdentificationData(f, f.datasetData)
	}

	return nil
}

func (f *Identification) HandleRequest(ctrl spine.Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf spine.Feature) (*model.MsgCounterType, error) {
	switch fct {
	case model.FunctionEnumTypeIdentificationListData:
		if op == model.CmdClassifierTypeRead {
			return f.requestListData(ctrl, rf)
		}
		return nil, fmt.Errorf("identification.handleRequest: FunctionEnumTypeIdentificationListData op not implemented: %s", op)
	}

	return nil, fmt.Errorf("identification.handleRequest: FunctionEnumType not implemented: %s", fct)
}

func (f *Identification) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.IdentificationListData != nil:
		data := cmd.IdentificationListData
		switch op {
		case model.CmdClassifierTypeReply:
			return f.replyListData(ctrl, *data)

		case model.CmdClassifierTypeNotify:
			return f.replyListData(ctrl, *data)

		default:
			return fmt.Errorf("identification.Handle: IdentificationListData CmdClassifierType not implemented: %s", op)
		}
	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("identification.Handle: CmdType not implemented: %s", populatedFields(cmd))
	}
}

func (f *Identification) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	return ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
}
