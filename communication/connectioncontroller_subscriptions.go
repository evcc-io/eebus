package communication

import (
	"errors"
	"reflect"
	"sync/atomic"

	"github.com/evcc-io/eebus/spine/model"
)

func (c *ConnectionController) subscriptionId() *model.SubscriptionIdType {
	i := model.SubscriptionIdType(atomic.AddUint64(&c.subscriptionNum, 1))
	return &i
}

func (c *ConnectionController) addSubscription(data model.SubscriptionManagementRequestCallType) error {
	var requestAllowed bool

	localAddress, localErr := c.featureAddressForTypeAndRole(
		c.localDevice,
		"local",
		model.EntityTypeEnumTypeDeviceInformation,
		model.FeatureTypeEnumTypeNodeManagement,
		model.RoleTypeSpecial,
	)

	remoteAddress, remoteErr := c.featureAddressForTypeAndRole(
		c.remoteDevice,
		"remote",
		model.EntityTypeEnumTypeDeviceInformation,
		model.FeatureTypeEnumTypeNodeManagement,
		model.RoleTypeSpecial,
	)

	// check if this is a subscription from nodemgmt to nodemgmt
	if localErr == nil && remoteErr == nil {
		if reflect.DeepEqual(data.ServerAddress, localAddress) && reflect.DeepEqual(data.ClientAddress, remoteAddress) {
			requestAllowed = true
		}
	}

	// check if this is a subscription from a client to a server of the same feature type
	if !requestAllowed {
		localAddress, localErr = c.featureAddressForTypeAndRole(
			c.localDevice,
			"local",
			model.EntityTypeEnumTypeCEM,
			model.FeatureTypeEnumType(*data.ServerFeatureType),
			model.RoleTypeServer,
		)

		remoteAddress, remoteErr = c.featureAddressForTypeAndRole(
			c.remoteDevice,
			"remote",
			model.EntityTypeEnumType(c.remoteDevice.Entity(data.ClientAddress.Entity).GetType()),
			model.FeatureTypeEnumType(*data.ServerFeatureType),
			model.RoleTypeClient,
		)

		// quick hack for ID. Charger which sends a subscription from featureType "Generic" to featureType "DeviceDiagnosis"
		altRemoteAddress, altRemoteErr := c.featureAddressForTypeAndRole(
			c.remoteDevice,
			"remote",
			model.EntityTypeEnumType(c.remoteDevice.Entity(data.ClientAddress.Entity).GetType()),
			model.FeatureTypeEnumTypeGeneric,
			model.RoleTypeClient,
		)

		if localErr == nil && (remoteErr == nil || altRemoteErr == nil) {
			if reflect.DeepEqual(data.ServerAddress, localAddress) && (reflect.DeepEqual(data.ClientAddress, remoteAddress) || reflect.DeepEqual(data.ClientAddress, altRemoteAddress)) {
				requestAllowed = true
			}
		}
	}

	if !requestAllowed {
		msg := "subscription request not conforming a request from a client to a server of the same type"
		c.log.Println(msg)

		// if this is an Elli device allow invalid requests, otherwise don't allow it
		if c.clientData.EVSEData.Manufacturer.BrandName != "Elli" {
			return errors.New(msg)
		}
	}

	subscriptionEntry := model.SubscriptionManagementEntryDataType{
		SubscriptionId: c.subscriptionId(),
		ClientAddress:  data.ClientAddress,
		ServerAddress:  data.ServerAddress,
	}

	c.subscriptionEntries = append(c.subscriptionEntries, subscriptionEntry)

	if model.FeatureTypeEnumType(*data.ServerFeatureType) == model.FeatureTypeEnumTypeDeviceDiagnosis {
		c.startHeartBeatSend()
	}

	return nil
}

func (c *ConnectionController) removeSubscription(data model.SubscriptionManagementDeleteCallType) error {
	// TODO: test this!!!

	var newSubscriptionEntries []model.SubscriptionManagementEntryDataType

	// according to the spec 7.4.4
	// a. The absence of "subscriptionDelete. clientAddress. device" SHALL be treated as if it was
	//    present and set to the sender's "device" address part.
	// b. The absence of "subscriptionDelete. serverAddress. device" SHALL be treated as if it was
	//    present and set to the recipient's "device" address part.

	clientAddress := data.ClientAddress
	if data.ClientAddress.Device == nil {
		clientAddress.Device = c.remoteDevice.Information().Description.DeviceAddress.Device
	}

	serverAddress := data.ServerAddress
	if data.ServerAddress.Device == nil {
		serverAddress.Device = c.localDevice.Information().Description.DeviceAddress.Device
	}

	for _, item := range c.subscriptionEntries {

		if reflect.DeepEqual(item.ClientAddress, clientAddress) {
			newSubscriptionEntries = append(newSubscriptionEntries, item)
		}
	}

	if len(newSubscriptionEntries) == len(c.subscriptionEntries) {
		return errors.New("could not find requested SubscriptionId to be removed")
	}

	c.subscriptionEntries = newSubscriptionEntries

	return nil
}

func (c *ConnectionController) subscriptions() []model.SubscriptionManagementEntryDataType {
	return c.subscriptionEntries
}
