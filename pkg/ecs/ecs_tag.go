package ecs

type TagComponent struct {
	typ uint64
}

func NewTagComponent(typ uint64) *TagComponent {
	return &TagComponent{
		typ: typ,
	}
}

func (c *TagComponent) Type() uint64 {
	return c.typ
}
