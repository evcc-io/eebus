package entity

import (
	"github.com/evcc-io/eebus/device/feature"
	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

//  Entities:
//   e[1] type=CEM, CEM Energy Guard
//  Features:
//   e[1] f-1 client.DeviceClassification - Device Classification
//   e[1] f-2 client.DeviceDiagnosis - Device Diagnosis
//   e[1] f-3 client.Measurement - Measurement for client
//   e[1] f-4 client.DeviceConfiguration - Device Configuration
//   e[1] f-5 server.DeviceDiagnosis - DeviceDiag
//    {RO} deviceDiagnosisStateData
//    {RO} deviceDiagnosisHeartbeatData
//   e[1] f-7 client.LoadControl - LoadControl client for CEM
//   e[1] f-8 client.Identification - EV identification
//   e[1] f-9 client.ElectricalConnection - Electrical Connection
func CEM() spine.Entity {
	var entityType model.EntityTypeType = model.EntityTypeType(model.EntityTypeEnumTypeCEM)
	entity := &spine.EntityImpl{
		Type: entityType,
	}

	fid := FeatureNumerator(1)
	{
		f := feature.NewDeviceClassificationClient()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewDeviceDiagnosisClient()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewMeasurementClient()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewDeviceConfigurationClient()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewDeviceDiagnosisServer()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewLoadControlClient()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewIdentificationClient()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewElectricalConnectionClient()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewTimeSeriesClient()
		f.SetID(fid())
		entity.Add(f)
	}
	{
		f := feature.NewIncentiveTableClient()
		f.SetID(fid())
		entity.Add(f)
	}

	return entity
}
