package sim

import (
	"github.com/nitwhiz/reactor/pkg/geometry"
	"image/color"
)

const (
	ComponentTypeBody = uint64(1) << uint64(iota)
	ComponentTypeVelocity
	ComponentTypeParticle
	ComponentTypeTemperature
	ComponentTypeRender
	TagEmitNeutrons
	TagFission
	TagControlRod
	TagControlRodSet1
	TagControlRodSet2
	TagThermalNeutron
	TagFastNeutron
	TagWater
	TagNonFissile
	TagXenon
	TagModerator
)

type BodyComponent struct {
	Body geometry.Body
}

func NewBodyComponent(body geometry.Body) *BodyComponent {
	return &BodyComponent{
		body,
	}
}

func (b *BodyComponent) Type() uint64 {
	return ComponentTypeBody
}

type VelocityComponent struct {
	velocity geometry.Vec2
}

func NewVelocityComponent(vx, vy float32) *VelocityComponent {
	return &VelocityComponent{
		geometry.Vec2{
			X: vx,
			Y: vy,
		},
	}
}

func (v *VelocityComponent) Type() uint64 {
	return ComponentTypeVelocity
}

type ParticleComponent struct {
	ParticleType *ParticleType
}

func NewParticleTypeComponent(typ *ParticleType) *ParticleComponent {
	return &ParticleComponent{
		typ,
	}
}

func (p *ParticleComponent) Type() uint64 {
	return ComponentTypeParticle
}

type TemperatureComponent struct {
	Temperature float32
	BaseColor   color.Color
}

func NewTemperatureComponent(baseColor color.Color) *TemperatureComponent {
	return &TemperatureComponent{
		Temperature: 20.0,
		BaseColor:   baseColor,
	}
}

func (h *TemperatureComponent) Type() uint64 {
	return ComponentTypeTemperature
}

func (h *TemperatureComponent) CurrentColor() color.Color {
	// todo: lookup table for this

	if h.Temperature >= 100 {
		return color.Transparent
	}

	r, g, b, _ := h.BaseColor.RGBA()

	r8 := float32(r >> 8)
	g8 := float32(g >> 8)
	b8 := float32(b >> 8)

	t := max(0.0, min(1.0, (h.Temperature-20.0)/80.0))

	newR := uint8(r8 + t*(255.0-r8))
	newG := uint8(g8 + t*(0-g8))
	newB := uint8(b8 + t*(0-b8))
	// todo: alpha?

	return color.RGBA{
		R: newR,
		G: newG,
		B: newB,
		A: 255,
	}
}

type RenderComponent struct {
	ZIndex int
}

func (z *RenderComponent) Type() uint64 {
	return ComponentTypeRender
}

func NewRenderComponent(z int) *RenderComponent {
	return &RenderComponent{
		ZIndex: z,
	}
}
