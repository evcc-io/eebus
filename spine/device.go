package spine

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/evcc-io/eebus/spine/model"
)

type Device interface {
	GetAddress() model.AddressDeviceType
	GetEntities() []Entity
	GetType() model.DeviceTypeType

	Add(e Entity)
	RemoveByAddress(addr []model.AddressEntityType)
	Entity(addr []model.AddressEntityType) Entity
	EntityByType(typ model.EntityTypeType) Entity

	SetUseCaseActor(actorName string, useCases []model.UseCaseSupportType)
	ResetUseCaseActors()
	GetUseCaseActors() []string
	GetUseCaseForActor(actorName string) ([]model.UseCaseSupportType, error)

	Information() *model.NodeManagementDetailedDiscoveryDeviceInformationType
	Dump(w io.Writer)
}

var _ Device = (*DeviceImpl)(nil)

type DeviceImpl struct {
	Address       model.AddressDeviceType
	Type          model.DeviceTypeType
	Entities      []Entity
	ActorUseCases map[string][]model.UseCaseSupportType
}

func (d *DeviceImpl) SetUseCaseActor(actorName string, useCases []model.UseCaseSupportType) {
	if d.ActorUseCases == nil {
		d.ActorUseCases = make(map[string][]model.UseCaseSupportType)
	}
	d.ActorUseCases[actorName] = useCases
}

func (d *DeviceImpl) ResetUseCaseActors() {
	d.ActorUseCases = nil
}

func (d *DeviceImpl) GetUseCaseActors() []string {
	var actors []string
	for actorName := range d.ActorUseCases {
		actors = append(actors, actorName)
	}
	return actors
}

func (d *DeviceImpl) GetUseCaseForActor(actorName string) ([]model.UseCaseSupportType, error) {
	usecases, found := d.ActorUseCases[actorName]
	if found {
		return usecases, nil
	}
	return nil, errors.New("actor not found")
}

func (d *DeviceImpl) GetAddress() model.AddressDeviceType {
	return d.Address
}

func (d *DeviceImpl) GetEntities() []Entity {
	return d.Entities
}

func (d *DeviceImpl) GetType() model.DeviceTypeType {
	return d.Type
}

func (d *DeviceImpl) Add(e Entity) {
	e.SetDevice(d)
	d.Entities = append(d.Entities, e)
}

func (d *DeviceImpl) RemoveByAddress(addr []model.AddressEntityType) {
	entityForRemoval := d.Entity(addr)

	var newEntities []Entity
	for _, item := range d.Entities {
		if !reflect.DeepEqual(item, entityForRemoval) {
			newEntities = append(newEntities, item)
		}
	}

	d.Entities = newEntities
}

func (d *DeviceImpl) Entity(id []model.AddressEntityType) Entity {
	for _, e := range d.Entities {
		if reflect.DeepEqual(id, e.GetAddress()) {
			return e
		}
	}
	return nil
}

func (d *DeviceImpl) EntityByType(typ model.EntityTypeType) Entity {
	if d != nil {
		for _, e := range d.Entities {
			if e.GetType() == typ {
				return e
			}
		}
	}
	return nil
}

func (d *DeviceImpl) Information() *model.NodeManagementDetailedDiscoveryDeviceInformationType {
	res := model.NodeManagementDetailedDiscoveryDeviceInformationType{
		Description: &model.NetworkManagementDeviceDescriptionDataType{
			DeviceAddress: &model.DeviceAddressType{
				Device: &d.Address,
			},
			DeviceType: &d.Type,
			// TODO NetworkFeatureSet
			// NetworkFeatureSet: &smart,
		},
	}
	return &res
}

func (d *DeviceImpl) Dump(w io.Writer) {
	fmt.Fprintf(w, "Details: device=%s, type=%s\n", d.Address, d.Type)

	fmt.Fprintln(w, "  Entities:")
	for _, e := range d.Entities {
		addr := EntityAddressString(e)
		fmt.Fprintf(w, "    e[%s] type=%s\n", addr, e.GetType())
	}

	fmt.Fprintln(w, "  Features:")
	for _, e := range d.Entities {
		e.Dump(w)
	}
}

func UnmarshalDevice(data model.NodeManagementDetailedDiscoveryDataType) *DeviceImpl {
	var dev *DeviceImpl

	if di := data.DeviceInformation; di != nil && di.Description != nil {
		dev = new(DeviceImpl)
		did := di.Description

		if did.DeviceType != nil {
			dev.Type = *did.DeviceType
		}

		if did.DeviceAddress != nil && did.DeviceAddress.Device != nil {
			dev.Address = *did.DeviceAddress.Device
		}
	}

	return dev
}
