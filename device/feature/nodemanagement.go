package feature

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

type NodeManagementDelegate interface {
	UpdateUseCaseSupportData(*NodeManagement, model.UseCaseNameType, bool)
}

type NodeManagement struct {
	*spine.FeatureImpl
	Delegate NodeManagementDelegate
}

func NewNodeManagement() spine.Feature {
	f := &NodeManagement{
		FeatureImpl: &spine.FeatureImpl{
			Type: model.FeatureTypeEnumTypeNodeManagement,
			Role: model.RoleTypeSpecial,
		},
	}

	f.Add(model.FunctionEnumTypeNodeManagementDetailedDiscoveryData, true, false)
	f.Add(model.FunctionEnumTypeNodeManagementSubscriptionRequestCall, false, false)
	f.Add(model.FunctionEnumTypeNodeManagementBindingRequestCall, false, false)
	f.Add(model.FunctionEnumTypeNodeManagementSubscriptionDeleteCall, false, false)
	f.Add(model.FunctionEnumTypeNodeManagementBindingDeleteCall, false, false)
	f.Add(model.FunctionEnumTypeNodeManagementSubscriptionData, true, false)
	f.Add(model.FunctionEnumTypeNodeManagementBindingData, true, false)
	f.Add(model.FunctionEnumTypeNodeManagementUseCaseData, true, false)

	return f
}

func (f *NodeManagement) Handle(ctrl spine.Context, rf model.FeatureAddressType, op model.CmdClassifierType, cmd model.CmdType, isPartialForCmd bool) error {
	switch {
	case cmd.NodeManagementDestinationListData != nil:
		return f.handleNodeManagementDestinationListData(ctrl, op, cmd.NodeManagementDestinationListData, isPartialForCmd)

	case cmd.NodeManagementDetailedDiscoveryData != nil:
		return f.handleDetailedDiscoveryData(ctrl, op, cmd.NodeManagementDetailedDiscoveryData, isPartialForCmd)

	case cmd.NodeManagementSubscriptionRequestCall != nil:
		return f.handleSubscriptionRequestCall(ctrl, op, cmd.NodeManagementSubscriptionRequestCall, isPartialForCmd)

	case cmd.NodeManagementSubscriptionDeleteCall != nil:
		return f.handleSubscriptionDeleteCall(ctrl, op, cmd.NodeManagementSubscriptionDeleteCall, isPartialForCmd)

	case cmd.NodeManagementSubscriptionData != nil:
		return f.handleSubscriptionData(ctrl, op, cmd.NodeManagementSubscriptionData, isPartialForCmd)

	case cmd.NodeManagementUseCaseData != nil:
		return f.handleUseCaseData(ctrl, op, cmd.NodeManagementUseCaseData, isPartialForCmd)

	case cmd.ResultData != nil:
		return f.HandleResultData(ctrl, op)

	default:
		return fmt.Errorf("nodemanagement.Handle: CmdType not implemented: %s", populatedFields(cmd))
	}
}

func (f *NodeManagement) ServerFound(ctrl spine.Context, rf spine.Feature) error {
	return ctrl.Subscribe(f, rf, model.FeatureTypeType(f.Type))
}
