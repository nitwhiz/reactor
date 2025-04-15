package sim

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

type BaseComponent struct{}

func NewBaseComponent() BaseComponent {
	return BaseComponent{}
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

type EntityManager struct {
	signatureToArchetype map[Signature]*Archetype
	entityArchetype      map[EntityID]*Archetype
	removeEntities       *Buffer[EntityID]
	futureEntities       *Buffer[*FutureEntity]
	runHooks             map[Signature]*Query

	systems []System

	queryCache map[Signature]*Query
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		signatureToArchetype: make(map[Signature]*Archetype),
		entityArchetype:      make(map[EntityID]*Archetype),
		removeEntities:       NewBuffer[EntityID](128, 32),
		futureEntities:       NewBuffer[*FutureEntity](128, 32),
		runHooks:             make(map[Signature]*Query),

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
		if !callback(q, eId) {
			return
		}
	}
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
}

func (m *EntityManager) updateEntities() {
	for _, eId := range m.removeEntities.Elements() {
		archetype := m.entityArchetype[eId]

		archetype.Remove(eId)

		for _, q := range m.queryCache {
			if archetype.Signature&q.signature == q.signature {
				q.removeResult(eId)
				m.runHooks[q.signature] = q
			}
		}

		delete(m.entityArchetype, eId)
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
