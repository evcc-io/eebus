package entity

import (
	"github.com/evcc-io/eebus/device/feature"
	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

//  Entities:
//   e[0] type=DeviceInformation
//  Features:
//   e[0] f-0 special.NodeManagement
//    {RO} nodeManagementDetailedDiscoveryData
//    {--} nodeManagementSubscriptionRequestCall
//    {--} nodeManagementBindingRequestCall
//    {--} nodeManagementSubscriptionDeleteCall
//    {--} nodeManagementBindingDeleteCall
//    {RO} nodeManagementSubscriptionData
//    {RO} nodeManagementBindingData
//    {RO} nodeManagementUseCaseData
//   e[0] f-1 server.DeviceClassification
//    {RO} deviceClassificationManufacturerData
func DeviceInformation() spine.Entity {
	var entityType model.EntityTypeType = model.EntityTypeType(model.EntityTypeEnumTypeDeviceInformation)
	entity := &spine.EntityImpl{
		Type: entityType,
	}

	fid := FeatureNumerator(0)

	{
		f := feature.NewNodeManagement()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewDeviceClassificationServer()
		f.SetID(fid())
		entity.Add(f)
	}

	return entity
}
