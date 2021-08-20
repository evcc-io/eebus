package spine

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/evcc-io/eebus/spine/model"
	"github.com/thoas/go-funk"
)

type Entity interface {
	GetAddress() []model.AddressEntityType
	SetAddress([]model.AddressEntityType)

	GetDevice() Device
	SetDevice(d Device)

	GetFeatures() []Feature
	GetType() model.EntityTypeType

	GetManufacturerData() model.DeviceClassificationManufacturerDataType
	SetManufacturerData(model.DeviceClassificationManufacturerDataType)

	GetOperationState() model.DeviceDiagnosisOperatingStateType
	SetOperationState(model.DeviceDiagnosisOperatingStateType)

	Add(f Feature)
	Feature(id uint) Feature
	FeatureByProps(typ model.FeatureTypeEnumType, role model.RoleType) Feature

	Information() *model.NodeManagementDetailedDiscoveryEntityInformationType
	Dump(w io.Writer)
}

var _ Entity = (*EntityImpl)(nil)

type EntityImpl struct {
	Device           Device
	Address          []model.AddressEntityType
	Type             model.EntityTypeType
	Description      model.DescriptionType
	Parent           Entity
	Entities         []Entity
	Features         []Feature
	ManufacturerData model.DeviceClassificationManufacturerDataType
	OperationState   model.DeviceDiagnosisOperatingStateType
}

func (e *EntityImpl) GetAddress() []model.AddressEntityType {
	return e.Address
}

func (e *EntityImpl) SetAddress(addr []model.AddressEntityType) {
	e.Address = addr
}

func (e *EntityImpl) GetDevice() Device {
	return e.Device
}

func (e *EntityImpl) SetDevice(d Device) {
	e.Device = d
}

func (e *EntityImpl) GetFeatures() []Feature {
	return e.Features
}

func (e *EntityImpl) GetType() model.EntityTypeType {
	return e.Type
}

func (e *EntityImpl) GetManufacturerData() model.DeviceClassificationManufacturerDataType {
	return e.ManufacturerData
}

func (e *EntityImpl) SetManufacturerData(data model.DeviceClassificationManufacturerDataType) {
	e.ManufacturerData = data
}

func (e *EntityImpl) GetOperationState() model.DeviceDiagnosisOperatingStateType {
	return e.OperationState
}

func (e *EntityImpl) SetOperationState(data model.DeviceDiagnosisOperatingStateType) {
	e.OperationState = data
}

// TODO maintain feature id when adding
func (e *EntityImpl) Add(f Feature) {
	f.SetEntity(e)
	e.Features = append(e.Features, f)
}

func (e *EntityImpl) Feature(id uint) Feature {
	if e != nil {
		for _, f := range e.Features {
			if f.GetID() == id {
				return f
			}
		}
	}
	return nil
}

func (e *EntityImpl) FeatureByProps(typ model.FeatureTypeEnumType, role model.RoleType) Feature {
	if e != nil {
		for _, f := range e.Features {
			if f.GetType() == typ && f.GetRole() == role {
				return f
			}
		}
	}
	return nil
}

func (e *EntityImpl) Information() *model.NodeManagementDetailedDiscoveryEntityInformationType {
	res := model.NodeManagementDetailedDiscoveryEntityInformationType{
		Description: &model.NetworkManagementEntityDescriptionDataType{
			EntityAddress: EntityAddressType(e),
			EntityType:    &e.Type,
		},
	}

	return &res
}

func (e *EntityImpl) Dump(w io.Writer) {
	for _, f := range e.Features {
		addr := EntityAddressString(e)
		fmt.Fprintf(w, "    e[%s] f-%d type=%s.%s\n", addr, f.GetID(), f.GetRole(), f.GetType())
		f.Dump(w)
	}
	for _, child := range e.Entities {
		child.Dump(w)
	}
}

func UnmarshalEntity(
	deviceAddress model.AddressDeviceType,
	entityData model.NodeManagementDetailedDiscoveryEntityInformationType,
) *EntityImpl {
	var e *EntityImpl

	if eid := entityData.Description; eid != nil {
		e = new(EntityImpl)

		if eid.EntityType != nil {
			e.Type = *eid.EntityType
		}

		if eid.Description != nil {
			e.Description = *eid.Description
		}

		if ea := eid.EntityAddress; ea != nil {
			e.Address = ea.Entity

			if ea.Device != nil && *ea.Device != deviceAddress {
				return nil
			}
		}
	}

	return e
}

func EntityAddressString(e Entity) string {
	return strings.Join(funk.Map(e.GetAddress(), func(id model.AddressEntityType) string {
		return strconv.FormatUint(uint64(id), 10)
	}).([]string), ",")
}

func EntityAddressType(e Entity) *model.EntityAddressType {
	addr := e.GetDevice().GetAddress()

	res := model.EntityAddressType{
		Device: &addr,
		Entity: e.GetAddress(),
	}

	return &res
}
