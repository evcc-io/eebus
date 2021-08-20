package communication

import (
	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

var _ spine.Context = (*contextImpl)(nil)

// contextImpl provides request context to features for processing
type contextImpl struct {
	*ConnectionController
	datagram model.DatagramType
}

func (c *contextImpl) AddressSource() *model.FeatureAddressType {
	return c.datagram.Header.AddressSource
}

func (c *contextImpl) AddSubscription(data model.SubscriptionManagementRequestCallType) error {
	return c.addSubscription(data)
}

func (c *contextImpl) RemoveSubscription(data model.SubscriptionManagementDeleteCallType) error {
	return c.removeSubscription(data)
}

func (c *contextImpl) HeartbeatCounter() *uint64 {
	return c.heartBeatCounter()
}

func (c *contextImpl) Subscriptions() []model.SubscriptionManagementEntryDataType {
	return c.subscriptions()
}

// Send a subscription request to a remove server feature
func (c *contextImpl) Subscribe(localFeature, remoteFeature spine.Feature, serverFeatureType model.FeatureTypeType) error {
	cmd := model.CmdType{
		NodeManagementSubscriptionRequestCall: &model.NodeManagementSubscriptionRequestCallType{
			SubscriptionRequest: &model.SubscriptionManagementRequestCallType{
				ClientAddress:     spine.FeatureAddressType(localFeature),
				ServerAddress:     spine.FeatureAddressType(remoteFeature),
				ServerFeatureType: &serverFeatureType,
			},
		},
	}

	// we always send it to the remode NodeManagment feature, which always is at entity:[0],feature:0
	var feature0 model.AddressFeatureType = 0
	remoteAddress := model.FeatureAddressType{
		Entity:  []model.AddressEntityType{0},
		Feature: &feature0,
	}
	remoteEntity := remoteFeature.GetEntity()
	if remoteEntity != nil {
		remoteDevice := remoteEntity.GetDevice()
		if remoteDevice != nil {
			deviceAddress := remoteDevice.GetAddress()
			remoteAddress.Device = &deviceAddress
		}
	}

	cmdClassifier := model.CmdClassifierTypeCall
	ackRequired := true

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &c.specificationVersion,
			AddressSource:        spine.FeatureAddressType(localFeature),
			AddressDestination:   &remoteAddress,
			MsgCounter:           c.msgCounter(),
			CmdClassifier:        &cmdClassifier,
			AckRequest:           &ackRequired,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return c.sendSpineMessage(datagram)
}

func (c *contextImpl) ProcessSequenceFlowRequest(featureType model.FeatureTypeEnumType, functionType model.FunctionEnumType, cmdClassifier model.CmdClassifierType) (*model.MsgCounterType, error) {
	return c.sendRequestToEVForFeatureAndFunction(featureType, functionType, cmdClassifier)
}

// Sends read request
func (c *contextImpl) Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error) {
	msgCounter := c.msgCounter()

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &c.specificationVersion,
			AddressSource:        &senderAddress,
			AddressDestination:   &destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: cmd,
		},
	}

	if ackRequest {
		datagram.Header.AckRequest = &ackRequest
	}

	return msgCounter, c.sendSpineMessage(datagram)
}

// Reply sends reply to original sender
func (c *contextImpl) Reply(cmdClassifier model.CmdClassifierType, cmd model.CmdType) error {
	// TODO where ack handling?

	// if ackRequest {
	// 	_ = c.sendAcknowledgementMessage(nil, featureSource, featureDestination, msgCounterReference)
	// }

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &c.specificationVersion,
			AddressSource:        c.datagram.Header.AddressDestination,
			AddressDestination:   c.datagram.Header.AddressSource,
			MsgCounter:           c.msgCounter(),
			MsgCounterReference:  c.datagram.Header.MsgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return c.sendSpineMessage(datagram)
}

// Notify sends notification to destination
func (c *contextImpl) Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error {
	cmdClassifier := model.CmdClassifierTypeNotify

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &c.specificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           c.msgCounter(),
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: cmd,
		},
	}

	return c.sendSpineMessage(datagram)
}

// Write sends notification to destination
func (c *contextImpl) Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error {
	cmdClassifier := model.CmdClassifierTypeWrite
	ackRequest := true

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &c.specificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           c.msgCounter(),
			CmdClassifier:        &cmdClassifier,
			AckRequest:           &ackRequest,
		},
		Payload: model.PayloadType{
			Cmd: cmd,
		},
	}

	return c.sendSpineMessage(datagram)
}
