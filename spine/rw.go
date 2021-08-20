package spine

import (
	"github.com/evcc-io/eebus/spine/model"
)

type RW struct {
	Read, Write bool
}

func (rw RW) String() string {
	switch {
	case rw.Read && !rw.Write:
		return "RO"
	case rw.Read && rw.Write:
		return "RW"
	default:
		return "--"
	}
}

func (rw RW) Information() *model.PossibleOperationsType {
	res := new(model.PossibleOperationsType)
	if rw.Read {
		res.Read = &model.PossibleOperationsReadType{}
	}
	if rw.Write {
		res.Write = &model.PossibleOperationsWriteType{}
	}

	return res
}
