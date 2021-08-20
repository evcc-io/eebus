package entity

import "github.com/evcc-io/eebus/spine/model"

func Numerator(ids []uint) func() []model.AddressEntityType {
	id := ids[len(ids)-1]
	start := make([]model.AddressEntityType, len(ids))
	for k, v := range ids {
		start[k] = model.AddressEntityType(v)
	}

	return func() []model.AddressEntityType {
		defer func() {
			id += 1
		}()

		addr := make([]model.AddressEntityType, len(start))
		_ = copy(addr, start)
		addr[len(addr)-1] = model.AddressEntityType(id)

		return addr
	}
}

func FeatureNumerator(id uint) func() uint {
	return func() uint {
		defer func() { id += 1 }()
		return id
	}
}
