package ecs

import "log"

type Query struct {
	m *EntityManager

	signature          Signature
	matchingArchetypes *Buffer[*Archetype]

	entities *Buffer[EntityID]

	hooks []func(*Query)
}

func (q *Query) Sort(cmp func(a, b EntityID) int) {
	q.entities.Sort(cmp)
}

func (q *Query) EntityManager() *EntityManager {
	return q.m
}

func (q *Query) Ids() []EntityID {
	return q.entities.Elements()
}

func (q *Query) removeResult(eId EntityID) {
	q.entities.Remove(eId)
}

func (q *Query) addResult(eId EntityID) {
	q.entities.Add(eId)
}

func (q *Query) RegisterHook(hook func(*Query)) {
	q.hooks = append(q.hooks, hook)
}

func (q *Query) RunHooks() {
	for _, hook := range q.hooks {
		hook(q)
	}
}

func (m *EntityManager) Query(signature Signature) *Query {
	q, ok := m.queryCache[signature]

	if !ok {
		log.Println("create new query")

		q = &Query{
			m:                  m,
			signature:          signature,
			matchingArchetypes: NewBuffer[*Archetype](8, 8),
			entities:           NewBuffer[EntityID](1024, 256),
			hooks:              make([]func(*Query), 0),
		}

		m.queryCache[signature] = q
	}

	if q.matchingArchetypes.Size() == 0 {
		for sig, archetype := range m.signatureToArchetype {
			if sig&signature == signature {
				q.matchingArchetypes.Add(archetype)
			}
		}
	}

	if len(q.entities.buf) == 0 {
		for _, archetype := range q.matchingArchetypes.Elements() {
			q.entities.Append(archetype.Entities.Elements())
		}
	}

	return q
}
