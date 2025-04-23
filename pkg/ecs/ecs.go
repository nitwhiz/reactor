package ecs

type EntityID = uint64

type ComponentType = uint64

type Signature = uint64

var lastUid = uint64(0)

func newId() uint64 {
	lastUid++
	return lastUid
}

type Component interface {
	Type() ComponentType
}

type System interface {
	Update()
}

func getSignature(cs []Component) Signature {
	res := Signature(0)

	for _, c := range cs {
		res |= c.Type()
	}

	return res
}

type FutureEntity struct {
	id         EntityID
	components []Component
}
