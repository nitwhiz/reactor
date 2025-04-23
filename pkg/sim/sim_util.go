package sim

import "github.com/nitwhiz/reactor/pkg/ecs"

func eachCollision(em *ecs.EntityManager, eId ecs.EntityID, signature ecs.Signature, callback func(eId ecs.EntityID, body *BodyComponent) bool) {
	b := em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

	em.EachEntity(signature|ComponentTypeBody, func(q *ecs.Query, otherEntityId ecs.EntityID) bool {
		otherBody := em.GetComponent(ComponentTypeBody, otherEntityId).(*BodyComponent)

		if eId != otherEntityId && b.Body.Overlaps(otherBody.Body) {
			if !callback(otherEntityId, b) {
				return false
			}
		}

		return true
	})
}

func eachBodyCollision(em *ecs.EntityManager, colliderSignature ecs.Signature, collideWith *BodyComponent, callback func() bool) {
	ecs.EachComponent[*BodyComponent](em, colliderSignature|ComponentTypeBody, ComponentTypeBody, func(c *BodyComponent) bool {
		if c.Body != collideWith.Body && collideWith.Body.Overlaps(c.Body) {
			return callback()
		}

		return true
	})
}
