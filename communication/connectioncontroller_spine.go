package communication

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

// handle transmission success acknowledgment messages. Protocol 5.2.4 & 5.2.5
var cmdClassifierRequiresTransmissionAck = map[model.CmdClassifierType]bool{
	model.CmdClassifierTypeRead:   false,
	model.CmdClassifierTypeWrite:  false,
	model.CmdClassifierTypeCall:   false,
	model.CmdClassifierTypeReply:  true,
	model.CmdClassifierTypeNotify: true,
	model.CmdClassifierTypeResult: false,
}

// TODO check if msgCounter value should also be global on CEM level and not on connection level
func (c *ConnectionController) msgCounter() *model.MsgCounterType {
	i := model.MsgCounterType(atomic.AddUint64(&c.msgNum, 1))
	return &i
}

func (c *ConnectionController) sendSpineMessage(datagram model.DatagramType) error {
	data := &model.CmiDatagramType{
		Datagram: datagram,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	cmdClassifier := datagram.Header.CmdClassifier
	cmd := datagram.Payload.Cmd[0]

	destinationAddress := datagram.Header.AddressDestination
	if c.remoteDevice != nil {
		remoteEntity := c.remoteDevice.Entity(destinationAddress.Entity)
		if remoteEntity == nil {
			return errors.New("sendSpineMessage: invalid remote entity address")
		}
		remoteFeature := remoteEntity.Feature(uint(*destinationAddress.Feature))
		c.log.Printf("send: %s %s:%s %s", *cmdClassifier, remoteEntity.GetType(), remoteFeature.GetType(), c.cmdDetails(cmd))
	} else {
		c.log.Printf("send: %s %s", *cmdClassifier, c.cmdDetails(cmd))

	}

	err = c.conn.Write(json.RawMessage(payload))

	// sending a proper spine message succeeded
	if err == nil {
		c.spineMsgMux.Lock()
		c.lastSpineMsg = time.Now()
		c.spineMsgMux.Unlock()
	}

	return err
}

func (c *ConnectionController) sendAcknowledgementMessage(err error, featureSource, featureDestination *model.FeatureAddressType, msgCounterReference *model.MsgCounterType) error {
	// send result message, see protocol spec 5.2.5
	// "0" in case of success, any other value in case of an error

	cmdClassifier := model.CmdClassifierTypeResult
	var resultSuccess model.ErrorNumberType
	var resultDescription model.DescriptionType
	var resultData model.ResultDataType

	if err == nil {
		resultSuccess = 0
		resultData = model.ResultDataType{
			ErrorNumber: &resultSuccess,
		}
	} else {
		resultSuccess = 1
		resultDescription = model.DescriptionType(err.Error())
		resultData = model.ResultDataType{
			ErrorNumber: &resultSuccess,
			Description: &resultDescription,
		}
	}

	responseDatagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &c.specificationVersion,
			AddressSource:        featureSource,
			AddressDestination:   featureDestination,
			MsgCounter:           c.msgCounter(),
			MsgCounterReference:  msgCounterReference,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{
				ResultData: &resultData,
			}},
		},
	}

	return c.sendSpineMessage(responseDatagram)
}

func (c *ConnectionController) processDatagram(datagram model.DatagramType) error {
	// TODO validate datagram.Header.AddressSource

	destAddr := datagram.Header.AddressDestination
	if destAddr == nil {
		return errors.New("processDatagram: invalid datagram.Header.AddressDestination")
	}

	// if destDevice is empty assume it's us
	destDevice := destAddr.Device
	if destDevice != nil {
		if c.localDevice.GetAddress() != *destDevice {
			return nil // not us
		}
	}

	if destAddr.Entity == nil {
		return errors.New("processDatagram: invalid datagram.Header.AddressDestination.Entity")
	}

	entity := c.localDevice.Entity(destAddr.Entity)
	if entity == nil {
		return errors.New("processDatagram: invalid entity address")
	}

	featureID := destAddr.Feature
	if featureID == nil {
		return errors.New("processDatagram: invalid datagram.Header.AddressDestination.Feature")
	}

	feature := entity.Feature(uint(*featureID))
	if feature == nil {
		return errors.New("processDatagram: invalid feature address")
	}

	err := c.processCmd(datagram, entity, feature)

	// handle processing success acknowledgement message. Protocol 5.2.4 & 5.2.5
	cmdClassifier := datagram.Header.CmdClassifier
	ackRequest := datagram.Header.AckRequest != nil && *datagram.Header.AckRequest
	if !cmdClassifierRequiresTransmissionAck[*cmdClassifier] && ackRequest {
		featureSource := datagram.Header.AddressDestination
		featureDestination := datagram.Header.AddressSource
		msgCounter := datagram.Header.MsgCounter
		ackErr := c.sendAcknowledgementMessage(err, featureSource, featureDestination, msgCounter)
		if ackErr != nil {
			return ackErr
		}
	}

	if err != nil {
		return err
	}

	// we received a proper SPINE message
	if err == nil {
		c.spineMsgMux.Lock()
		c.lastSpineMsg = time.Now()
		c.spineMsgMux.Unlock()
	}

	// TODO we need to process resultData responses also!
	//   these can also be a response to a sequence and those error results
	//   need to be considered!
	if datagram.Header.MsgCounterReference != nil {
		ctx := c.context(nil)
		return c.sequencesController.ProcessResponseInSequences(ctx, datagram.Header.MsgCounterReference)
	}

	return nil
}

func (c *ConnectionController) sendRequestToEVForFeatureAndFunction(featureType model.FeatureTypeEnumType, functionType model.FunctionEnumType, cmdClassifier model.CmdClassifierType) (*model.MsgCounterType, error) {
	le := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM))
	if le == nil {
		return nil, fmt.Errorf("local entity not found: %s", model.EntityTypeEnumTypeCEM)
	}

	re := c.remoteDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeEV))
	if re == nil {
		return nil, fmt.Errorf("remote entity not found: %s", model.EntityTypeEnumTypeEV)
	}

	lf := le.FeatureByProps(featureType, model.RoleTypeClient)
	if lf == nil {
		return nil, fmt.Errorf("local entity client feature not found: %s", featureType)
	}

	rf := re.FeatureByProps(featureType, model.RoleTypeServer)
	if rf == nil {
		return nil, fmt.Errorf("remote entity server feature not found: %s", featureType)
	}

	ctx := c.context(nil)
	return lf.HandleRequest(ctx, functionType, cmdClassifier, rf)
}

func (c *ConnectionController) cmdDetails(cmd model.CmdType) string {
	switch {
	case cmd.DeviceClassificationManufacturerData != nil:
		return "DeviceClassificationManufacturerData"
	case cmd.DeviceConfigurationKeyValueDescriptionListData != nil:
		return "DeviceConfigurationKeyValueDescriptionListData"
	case cmd.DeviceConfigurationKeyValueListData != nil:
		return "DeviceConfigurationKeyValueListData"
	case cmd.DeviceDiagnosisHeartbeatData != nil:
		return "DeviceDiagnosisHeartbeatData"
	case cmd.DeviceDiagnosisStateData != nil:
		return "DeviceDiagnosisStateData"
	case cmd.ElectricalConnectionDescriptionListData != nil:
		return "ElectricalConnectionDescriptionListData"
	case cmd.ElectricalConnectionParameterDescriptionListData != nil:
		return "ElectricalConnectionParameterDescriptionListData"
	case cmd.ElectricalConnectionPermittedValueSetListData != nil:
		return "ElectricalConnectionPermittedValueSetListData"
	case cmd.IdentificationListData != nil:
		return "IdentificationListData"
	case cmd.IncentiveTableDescriptionData != nil:
		return "IncentiveTableDescriptionData"
	case cmd.IncentiveTableConstraintsData != nil:
		return "IncentiveTableConstraintsData"
	case cmd.IncentiveTableData != nil:
		return "IncentiveTableData"
	case cmd.LoadControlLimitDescriptionListData != nil:
		return "LoadControlLimitDescriptionListData"
	case cmd.LoadControlLimitListData != nil:
		return "LoadControlLimitListData"
	case cmd.NodeManagementBindingRequestCall != nil:
		return "NodeManagementBindingRequestCall"
	case cmd.NodeManagementDestinationListData != nil:
		return "NodeManagementDestinationListData"
	case cmd.NodeManagementDetailedDiscoveryData != nil:
		return "NodeManagementDetailedDiscoveryData"
	case cmd.NodeManagementSubscriptionData != nil:
		return "NodeManagementSubscriptionData"
	case cmd.NodeManagementSubscriptionRequestCall != nil:
		return "NodeManagementSubscriptionRequestCall"
	case cmd.NodeManagementSubscriptionDeleteCall != nil:
		return "NodeManagementSubscriptionDeleteCall"
	case cmd.NodeManagementUseCaseData != nil:
		return "NodeManagementUseCaseData"
	case cmd.MeasurementConstraintsListData != nil:
		return "MeasurementConstraintsListData"
	case cmd.MeasurementDescriptionListData != nil:
		return "MeasurementDescriptionListData"
	case cmd.MeasurementListData != nil:
		return "MeasurementListData"
	case cmd.TimeSeriesConstraintsData != nil:
		return "TimeSeriesConstraintsListData"
	case cmd.TimeSeriesDescriptionListData != nil:
		return "TimeSeriesDescriptionListData"
	case cmd.TimeSeriesListData != nil:
		return "TimeSeriesListData"
	case cmd.ResultData != nil:
		return "ResultData"
	}

	return "unknown"
}

func (c *ConnectionController) processCmd(datagram model.DatagramType, localEntity spine.Entity, localFeature spine.Feature) error {
	cmdClassifier := datagram.Header.CmdClassifier
	if cmdClassifier == nil {
		return errors.New("processCmd: invalid datagram.Header.CmdClassifier")
	}

	if len(datagram.Payload.Cmd) != 1 {
		return errors.New("processCmd: invalid datagram.Payload.Cmd")
	}
	cmd := datagram.Payload.Cmd[0]

	c.log.Printf("recv: %s %s:%s %s", *cmdClassifier, localEntity.GetType(), localFeature.GetType(), c.cmdDetails(cmd))

	// isPartial
	isPartial := false
	functionCmd := cmd.Function
	filterCmd := cmd.Filter

	if functionCmd != nil && filterCmd != nil {
		// TODO check if the function is the same as the provided cmd value
		if len(filterCmd) > 0 {
			for _, filter := range filterCmd {
				if filter.CmdControl.Partial != nil {
					isPartial = true
					break
				}
			}
		}
	}

	return localFeature.Handle(c.context(&datagram), *datagram.Header.AddressSource, *cmdClassifier, cmd, isPartial)
}
