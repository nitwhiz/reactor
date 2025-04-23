package ecs

type EntityManager struct {
	signatureToArchetype map[Signature]*Archetype
	entityArchetype      map[EntityID]*Archetype
	entityAlive          map[EntityID]bool

	removeEntities *Buffer[EntityID]
	futureEntities *Buffer[*FutureEntity]

	runHooks map[Signature]*Query

	systems []System

	queryCache map[Signature]*Query
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		signatureToArchetype: make(map[Signature]*Archetype),
		entityArchetype:      make(map[EntityID]*Archetype),
		entityAlive:          make(map[EntityID]bool),

		removeEntities: NewBuffer[EntityID](128, 32),
		futureEntities: NewBuffer[*FutureEntity](128, 32),

		runHooks: make(map[Signature]*Query),

		systems: []System{},

		queryCache: make(map[Signature]*Query),
	}
}

func (m *EntityManager) EntityCount() int {
	return len(m.entityArchetype)
}

func (m *EntityManager) GetComponent(cType ComponentType, eId EntityID) Component {
	archetype := m.entityArchetype[eId]
	eIdx := archetype.EntityToComponentIndex[eId]
	cs := archetype.Components[cType]

	return cs.At(eIdx)
}

func (m *EntityManager) EachEntity(signature uint64, callback func(q *Query, eId EntityID) bool) {
	q := m.Query(signature)

	for _, eId := range q.Ids() {
		if alive, ok := m.entityAlive[eId]; ok && alive {
			if !callback(q, eId) {
				return
			}
		}
	}
}

func EachComponent[T Component](em *EntityManager, signature Signature, componentType ComponentType, callback func(component T) bool) {
	for s, a := range em.signatureToArchetype {
		if s&signature == signature {
			cs := a.Components[componentType]

			for _, c := range cs.Elements() {
				if !callback(c.(T)) {
					return
				}
			}
		}
	}
}

func (m *EntityManager) CountComponent(componentType ComponentType) int {
	res := 0

	for s, a := range m.signatureToArchetype {
		if s&componentType == componentType {
			cs := a.Components[componentType]

			res += cs.Size()
		}
	}

	return res
}

func (m *EntityManager) AddSystem(s System) {
	m.systems = append(m.systems, s)
}

func (m *EntityManager) AddEntity(cs ...Component) {
	m.futureEntities.Add(&FutureEntity{
		id:         newId(),
		components: cs,
	})

	return
}

func (m *EntityManager) RemoveEntity(eId uint64) {
	m.removeEntities.Add(eId)
	m.entityAlive[eId] = false
}

func (m *EntityManager) updateEntities() {
	for _, eId := range m.removeEntities.Elements() {
		archetype, ok := m.entityArchetype[eId]

		if !ok {
			continue
		}

		archetype.Remove(eId)

		for _, q := range m.queryCache {
			if archetype.Signature&q.signature == q.signature {
				q.removeResult(eId)
				m.runHooks[q.signature] = q
			}
		}

		delete(m.entityArchetype, eId)
		delete(m.entityAlive, eId)
	}

	for _, f := range m.futureEntities.Elements() {
		sig := getSignature(f.components)

		archetype := m.GetOrCreateArchetype(sig)

		archetype.Add(f.id, f.components...)

		for _, q := range m.queryCache {
			if archetype.Signature&q.signature == q.signature {
				q.addResult(f.id)
				m.runHooks[q.signature] = q
			}
		}

		m.entityArchetype[f.id] = archetype

		m.entityAlive[f.id] = true
	}

	m.removeEntities.Clear()
	m.futureEntities.Clear()

	for _, q := range m.runHooks {
		q.RunHooks()
	}

	clear(m.runHooks)
}

func (m *EntityManager) Update() {
	for _, s := range m.systems {
		s.Update()
	}

	m.updateEntities()
}
