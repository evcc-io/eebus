package app

import (
	"github.com/evcc-io/eebus/communication"
	"github.com/evcc-io/eebus/device/entity"
	"github.com/evcc-io/eebus/spine"
	"github.com/evcc-io/eebus/spine/model"
)

func HEMS(details communication.ManufacturerDetails) spine.Device {
	localDeviceName := model.DeviceClassificationStringType(details.DeviceName)
	localDeviceCode := model.DeviceClassificationStringType(details.DeviceCode)
	localBrandName := model.DeviceClassificationStringType(details.BrandName)
	localDeviceAddress := model.AddressDeviceType(details.DeviceAddress)

	manufacturerData := model.DeviceClassificationManufacturerDataType{
		DeviceName: &localDeviceName,
		DeviceCode: &localDeviceCode,
		BrandName:  &localBrandName,
		VendorName: &localBrandName,
	}

	operationState := model.DeviceDiagnosisOperatingStateType(model.DeviceDiagnosisOperatingStateEnumTypeNormalOperation)

	dev := &spine.DeviceImpl{
		Address: localDeviceAddress,
		Type:    model.DeviceTypeType(model.DeviceTypeEnumTypeEnergyManagementSystem),
	}

	eid := entity.Numerator([]uint{0})

	{
		e := entity.DeviceInformation()
		e.SetAddress(eid())
		dev.Add(e)
	}
	{
		e := entity.CEM()
		e.SetAddress(eid())
		e.SetManufacturerData(manufacturerData)
		e.SetOperationState(operationState)
		dev.Add(e)
	}

	return dev
}
