package communication

import (
	"fmt"

	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

func (c *ConnectionController) context(datagram *model.DatagramType) spine.Context {
	ctrl := &contextImpl{
		ConnectionController: c,
	}
	if datagram != nil {
		ctrl.datagram = *datagram
	}
	return ctrl
}

func (c *ConnectionController) GetDevice() spine.Device {
	return c.remoteDevice
}

func (c *ConnectionController) SetDevice(device spine.Device) {
	c.remoteDevice = device
}

func (c *ConnectionController) isEVConnected() bool {
	if c.remoteDevice == nil {
		return false
	}

	entity := c.remoteDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeEV))

	return entity != nil
}

func (c *ConnectionController) UpdateDevice(stateChange model.NetworkManagementStateChangeType) {
	isEVConnected := c.isEVConnected()

	if stateChange == model.NetworkManagementStateChangeTypeAdded && isEVConnected {
		c.log.Println("detected ev connection")

		c.clientData.EVData.ChargeState = EVChargeStateEnumTypeActive

		err := c.requestNodeManagementUseCaseData()
		if err != nil {
			c.log.Println("Sending UseCaseData read request failed!")
		}

		// TODO these actions should be usecase support specific!
		ctx := c.context(nil)
		err = c.sequencesController.StartSequenceFlow(ctx, SequenceEnumTypeEV)
		if err != nil {
			c.log.Println("error processing EV sequence")
		}

		le := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM))
		lf := le.FeatureByProps(model.FeatureTypeEnumTypeLoadControl, model.RoleTypeClient)

		re := c.remoteDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeEV))
		rf := re.FeatureByProps(model.FeatureTypeEnumTypeLoadControl, model.RoleTypeServer)

		serverFeatureType := model.FeatureTypeType(model.FeatureTypeEnumTypeLoadControl)

		msgCounter, err := c.callNodeManagementBindingRequest(lf, rf, serverFeatureType)
		if err != nil {
			c.log.Println(msgCounter, err)
		}
		c.callDataUpdateHandler(EVDataElementUpdateEVConnectionState)
	} else if !isEVConnected && stateChange == model.NetworkManagementStateChangeTypeRemoved {
		c.log.Println("detected ev disconnection")
		c.clientData.EVData.ChargeState = EVChargeStateEnumTypeUnplugged
		c.remoteDevice.ResetUseCaseActors()
		c.callDataUpdateHandler(EVDataElementUpdateEVConnectionState)
	}
}

// TODO move this into NodeManagement feature implementation
func (c *ConnectionController) requestNodeManagementDetailedDiscoveryData() error {
	cmdClassifier := model.CmdClassifierTypeRead
	nodeMgmtF, featureDestination := c.remoteNodeManagementFeature()

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &c.specificationVersion,
			AddressSource:        spine.FeatureAddressType(nodeMgmtF),
			AddressDestination:   &featureDestination,
			MsgCounter:           c.msgCounter(),
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{
				NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
			}},
		},
	}

	return c.sendSpineMessage(datagram)
}

func (c *ConnectionController) requestNodeManagementUseCaseData() error {
	cmdClassifier := model.CmdClassifierTypeRead
	nodeMgmtF, featureDestination := c.remoteNodeManagementFeature()

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &c.specificationVersion,
			AddressSource:        spine.FeatureAddressType(nodeMgmtF),
			AddressDestination:   &featureDestination,
			MsgCounter:           c.msgCounter(),
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{{
				NodeManagementUseCaseData: &model.NodeManagementUseCaseDataType{},
			}},
		},
	}

	return c.sendSpineMessage(datagram)
}

func (c *ConnectionController) remoteNodeManagementFeature() (spine.Feature, model.FeatureAddressType) {
	deviceInfoE := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeDeviceInformation))
	nodeMgmtF := deviceInfoE.FeatureByProps(model.FeatureTypeEnumTypeNodeManagement, model.RoleTypeSpecial)

	var feature0 model.AddressFeatureType = 0
	featureDestination := model.FeatureAddressType{
		Entity:  []model.AddressEntityType{0},
		Feature: &feature0,
	}

	if c.remoteDevice != nil {
		featureDestination.Device = c.remoteDevice.Information().Description.DeviceAddress.Device
	}

	return nodeMgmtF, featureDestination
}

func (c *ConnectionController) callNodeManagementBindingRequest(lf, rf spine.Feature, featureType model.FeatureTypeType) (*model.MsgCounterType, error) {

	localDeviceInfoE := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeDeviceInformation))
	localNodeMgmtF := localDeviceInfoE.FeatureByProps(model.FeatureTypeEnumTypeNodeManagement, model.RoleTypeSpecial)

	remoteDeviceInfoE := c.remoteDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeDeviceInformation))
	remoteNodeMgmtF := remoteDeviceInfoE.FeatureByProps(model.FeatureTypeEnumTypeNodeManagement, model.RoleTypeSpecial)

	res := []model.CmdType{{
		NodeManagementBindingRequestCall: &model.NodeManagementBindingRequestCallType{
			BindingRequest: &model.BindingManagementRequestCallType{
				ClientAddress:     spine.FeatureAddressType(lf),
				ServerAddress:     spine.FeatureAddressType(rf),
				ServerFeatureType: &featureType,
			},
		},
	}}

	ctrl := c.context(nil)

	return ctrl.Request(model.CmdClassifierTypeCall, *spine.FeatureAddressType(localNodeMgmtF), *spine.FeatureAddressType(remoteNodeMgmtF), true, res)
}

func (c *ConnectionController) featureAddressForTypeAndRole(device spine.Device, deviceLocation string, entityType model.EntityTypeEnumType, featureType model.FeatureTypeEnumType, role model.RoleType) (*model.FeatureAddressType, error) {
	entity := device.EntityByType(model.EntityTypeType(entityType))
	if entity == nil {
		err := fmt.Errorf("couldn't find device with entity %s on %s device %v", entityType, deviceLocation, device)
		c.log.Println(err)
		return nil, err
	}

	feature := entity.FeatureByProps(featureType, role)
	if feature == nil {
		err := fmt.Errorf("couldn't find feature %s with role %s in entity %s on %s device %v", featureType, role, entityType, deviceLocation, device)
		return nil, err
	}

	address := spine.FeatureAddressType(feature)
	return address, nil
}
