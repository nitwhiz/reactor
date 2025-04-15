package sim

func eachCollision(em *EntityManager, eId EntityID, signature Signature, callback func(eId EntityID, body *BodyComponent) bool) {
	b := em.GetComponent(ComponentTypeBody, eId).(*BodyComponent)

	em.EachEntity(signature|ComponentTypeBody, func(q *Query, otherEntityId EntityID) bool {
		otherBody := em.GetComponent(ComponentTypeBody, otherEntityId).(*BodyComponent)

		if eId != otherEntityId && b.Body.Overlaps(otherBody.Body) {
			if !callback(otherEntityId, b) {
				return false
			}
		}

		return true
	})
}
