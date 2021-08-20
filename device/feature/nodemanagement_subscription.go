package feature

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

// route subscription request calls to the appropriate feature implementation and add the subscription to the current list
func (f *NodeManagement) replySubscriptionData(ctrl spine.Context) error {
	res := model.CmdType{
		NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{
			SubscriptionEntry: ctrl.Subscriptions(),
		},
	}

	return ctrl.Reply(model.CmdClassifierTypeReply, res)
}

func (f *NodeManagement) handleSubscriptionData(ctrl spine.Context, op model.CmdClassifierType, data *model.NodeManagementSubscriptionDataType, isPartialForCmd bool) error {
	switch op {
	case model.CmdClassifierTypeCall:
		return f.replySubscriptionData(ctrl)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", op)
	}
}

func (f *NodeManagement) handleSubscriptionRequestCall(ctrl spine.Context, op model.CmdClassifierType, data *model.NodeManagementSubscriptionRequestCallType, isPartialForCmd bool) error {
	switch op {
	case model.CmdClassifierTypeCall:
		if err := ctrl.AddSubscription(*data.SubscriptionRequest); err != nil {
			// subscription failed, send an resulterror reply?
		}
		return nil

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionRequestCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", op)
	}
}

func (f *NodeManagement) handleSubscriptionDeleteCall(ctrl spine.Context, op model.CmdClassifierType, data *model.NodeManagementSubscriptionDeleteCallType, isPartialForCmd bool) error {
	switch op {
	case model.CmdClassifierTypeCall:
		return ctrl.RemoveSubscription(*data.SubscriptionDelete)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", op)
	}
}
