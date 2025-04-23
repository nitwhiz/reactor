package sim

import (
	"github.com/nitwhiz/reactor/pkg/ecs"
)

var SortHook = func(q *ecs.Query) {
	q.Sort(func(a, b ecs.EntityID) int {
		zIndexComponentA := q.EntityManager().GetComponent(ComponentTypeRender, a).(*RenderComponent)
		zIndexComponentB := q.EntityManager().GetComponent(ComponentTypeRender, b).(*RenderComponent)

		return zIndexComponentA.ZIndex - zIndexComponentB.ZIndex
	})
}
