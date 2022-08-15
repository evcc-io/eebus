package spine

import (
	"errors"
	"fmt"
	"io"

	"github.com/evcc-io/eebus/spine/model"
)

type ClientFeature interface {
	ServerFound(ctrl Context, e Feature) error
}

type Feature interface {
	GetAddress() *model.FeatureAddressType

	GetID() uint
	SetID(id uint)

	GetEntity() Entity
	SetEntity(e Entity)

	GetType() model.FeatureTypeEnumType
	GetRole() model.RoleType

	Add(fun model.FunctionEnumType, r, w bool)

	SupportForFunctionAvailable(fun model.FunctionEnumType) bool

	Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType
	Dump(w io.Writer)

	EVDisconnect()

	HandleRequest(ctrl Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf Feature) (*model.MsgCounterType, error)
	Handle(ctrl Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error
	HandleResultData(ctrl Context, op model.CmdClassifierType) error
}

var _ Feature = (*FeatureImpl)(nil)

type FeatureImpl struct {
	Entity        Entity
	ID            uint
	Type          model.FeatureTypeEnumType
	Description   model.DescriptionType
	Role          model.RoleType
	Functions     map[model.FunctionEnumType]RW
	Subscriptions []model.SubscriptionManagementEntryDataType
}

func (f *FeatureImpl) GetAddress() *model.FeatureAddressType {
	return FeatureAddressType(f)
}

func (f *FeatureImpl) GetID() uint {
	return f.ID
}

func (f *FeatureImpl) SetID(id uint) {
	f.ID = id
}

func (f *FeatureImpl) GetEntity() Entity {
	return f.Entity
}

func (f *FeatureImpl) SetEntity(e Entity) {
	f.Entity = e
}

func (f *FeatureImpl) GetType() model.FeatureTypeEnumType {
	return f.Type
}

func (f *FeatureImpl) GetRole() model.RoleType {
	return f.Role
}

func (f *FeatureImpl) Add(fun model.FunctionEnumType, r, w bool) {
	if f.Functions == nil {
		f.Functions = make(map[model.FunctionEnumType]RW)
	}

	if f.Role == model.RoleTypeClient {
		panic("cannot add functions to client role")
	}

	f.Functions[fun] = RW{r, w}
}

func (f *FeatureImpl) SupportForFunctionAvailable(fun model.FunctionEnumType) bool {
	_, found := f.Functions[fun]
	return found
}

func (f *FeatureImpl) Information() *model.NodeManagementDetailedDiscoveryFeatureInformationType {
	var funs []model.FunctionPropertyType
	for fun, rw := range f.Functions {
		var functionType model.FunctionType = model.FunctionType(fun)
		sf := model.FunctionPropertyType{
			Function:           &functionType,
			PossibleOperations: rw.Information(),
		}

		funs = append(funs, sf)
	}

	var featureType model.FeatureTypeType = model.FeatureTypeType(f.Type)
	var featureRole model.RoleType = model.RoleType(f.Role)

	res := model.NodeManagementDetailedDiscoveryFeatureInformationType{
		Description: &model.NetworkManagementFeatureDescriptionDataType{
			FeatureAddress:    FeatureAddressType(f),
			FeatureType:       &featureType,
			Role:              &featureRole,
			SupportedFunction: funs,
		},
	}

	return &res
}

func (f *FeatureImpl) Dump(w io.Writer) {
	for fun, ops := range f.Functions {
		fmt.Fprintf(w, "      {%s} %s\n", ops, fun)
	}
}

func (f *FeatureImpl) EVDisconnect() {}

func (f *FeatureImpl) HandleRequest(ctrl Context, fct model.FunctionEnumType, op model.CmdClassifierType, rf Feature) (*model.MsgCounterType, error) {
	return nil, errors.New("HandleRequest() not implemented")
}

func (f *FeatureImpl) Handle(ctrl Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	return errors.New("Handle() not implemented")
}

func (f *FeatureImpl) HandleResultData(ctrl Context, op model.CmdClassifierType) error {
	switch op {
	case model.CmdClassifierTypeResult:
		// TODO process the return result data for the message sent with the ID in msgCounterReference
		// error numbers explained in Resource Spec 3.11
		return nil

	default:
		return fmt.Errorf("ResultData CmdClassifierType %s not implemented", op)
	}
}

func UnmarshalFeature(
	// entityID uint,
	featureData model.NodeManagementDetailedDiscoveryFeatureInformationType,
) *FeatureImpl {
	var f *FeatureImpl

	if fid := featureData.Description; fid != nil {
		f = &FeatureImpl{
			Type: model.FeatureTypeEnumType(*fid.FeatureType),
		}

		if fid.Description != nil {
			f.Description = *fid.Description
		}

		if fid.Role != nil {
			f.Role = *fid.Role
		}

		if addr := fid.FeatureAddress; addr != nil {
			f.ID = uint(*addr.Feature)
		}

		for _, sf := range fid.SupportedFunction {
			f.Add(model.FunctionEnumType(*sf.Function), sf.PossibleOperations.Read != nil, sf.PossibleOperations.Write != nil)
		}
	}

	return f
}

func FeatureAddressType(f Feature) *model.FeatureAddressType {
	e := f.GetEntity()
	featureId := model.AddressFeatureType(f.GetID())

	res := model.FeatureAddressType{
		Entity:  e.GetAddress(),
		Feature: &featureId,
	}

	device := e.GetDevice()
	if device != nil {
		addr := device.GetAddress()
		res.Device = &addr
	}

	return &res
}
