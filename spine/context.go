package spine

import "github.com/evcc-io/eebus/spine/model"

type Context interface {
	CloseConnectionBecauseOfError(err error)
	SetDevice(Device)
	GetDevice() Device
	UpdateDevice(model.NetworkManagementStateChangeType)
	HeartbeatCounter() *uint64
	Subscribe(lf Feature, rf Feature, typ model.FeatureTypeType) error
	ProcessSequenceFlowRequest(featureType model.FeatureTypeEnumType, functionType model.FunctionEnumType, cmdClassifier model.CmdClassifierType) (*model.MsgCounterType, error)
	Request(model.CmdClassifierType, model.FeatureAddressType, model.FeatureAddressType, bool, []model.CmdType) (*model.MsgCounterType, error)
	Reply(model.CmdClassifierType, model.CmdType) error
	Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error
	Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error
	AddressSource() *model.FeatureAddressType
	AddSubscription(data model.SubscriptionManagementRequestCallType) error
	RemoveSubscription(data model.SubscriptionManagementDeleteCallType) error
	Subscriptions() []model.SubscriptionManagementEntryDataType
}
