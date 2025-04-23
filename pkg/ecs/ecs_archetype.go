package ecs

type Archetype struct {
	Signature Signature
	// Components holds all components indexed by ComponentType
	// an entity id uses the same index for the different components.
	Components map[ComponentType]*Buffer[Component]
	// Entities are all entities matching this Archetype.
	Entities *Buffer[EntityID]
	// EntityToComponentIndex contains the index in the Components array for this EntityID.
	EntityToComponentIndex map[EntityID]int
}

func (a *Archetype) Add(eId uint64, cs ...Component) {
	eIdx := a.Entities.Size()

	a.Entities.Add(eId)
	a.EntityToComponentIndex[eId] = eIdx

	for _, component := range cs {
		a.putComponent(eIdx, component)
	}
}

func (a *Archetype) Remove(eId EntityID) {
	cIdx := a.EntityToComponentIndex[eId]

	lastEId := a.Entities.Last()

	a.EntityToComponentIndex[lastEId] = cIdx

	for _, cs := range a.Components {
		// this shifts last to cIdx
		cs.RemoveIndex(cIdx)
	}

	a.Entities.Remove(eId)

	delete(a.EntityToComponentIndex, eId)
}

func (a *Archetype) putComponent(index int, component Component) {
	cType := component.Type()

	if _, ok := component.(*TagComponent); ok {
		component = nil
	}

	if _, ok := a.Components[cType]; !ok {
		a.Components[cType] = NewBuffer[Component](256, 256)
	}

	a.Components[cType].Put(index, component)
}

func (m *EntityManager) GetOrCreateArchetype(sig Signature) *Archetype {
	if archetype, ok := m.signatureToArchetype[sig]; ok {
		return archetype
	}

	archetype := &Archetype{
		Signature:              sig,
		Components:             make(map[ComponentType]*Buffer[Component]),
		Entities:               NewBuffer[EntityID](1024, 256),
		EntityToComponentIndex: make(map[EntityID]int),
	}

	m.signatureToArchetype[sig] = archetype

	// todo: keep queries in the ecs?

	return archetype
}
