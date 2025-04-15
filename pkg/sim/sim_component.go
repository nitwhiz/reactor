package sim

import (
	"github.com/nitwhiz/reactor/pkg/geometry"
	"image/color"
)

const (
	ComponentTypeBody = uint64(1) << uint64(iota)
	ComponentTypeVelocity
	ComponentTypeParticleType
	ComponentTypeHeatTransfer
	ComponentTypeZIndex
	ComponentTypeFission
)

type BodyComponent struct {
	BaseComponent
	Body geometry.Body
}

func NewBodyComponent(body geometry.Body) *BodyComponent {
	return &BodyComponent{
		NewBaseComponent(),
		body,
	}
}

func (b *BodyComponent) Type() uint64 {
	return ComponentTypeBody
}

type VelocityComponent struct {
	BaseComponent
	velocity geometry.Vec2
}

func NewVelocityComponent(vx, vy float32) *VelocityComponent {
	return &VelocityComponent{
		NewBaseComponent(),
		geometry.Vec2{
			X: vx,
			Y: vy,
		},
	}
}

func (v *VelocityComponent) Type() uint64 {
	return ComponentTypeVelocity
}

type ParticleTypeComponent struct {
	BaseComponent
	ParticleType *ParticleType
}

func NewParticleTypeComponent(typ *ParticleType) *ParticleTypeComponent {
	return &ParticleTypeComponent{
		NewBaseComponent(),
		typ,
	}
}

func (p *ParticleTypeComponent) Type() uint64 {
	return ComponentTypeParticleType
}

type HeatTransferComponent struct {
	BaseComponent
	ParticleType uint64
	Temperature  float32
	BaseColor    color.Color
}

func NewHeatTransferComponent(baseColor color.Color, particleType uint64) *HeatTransferComponent {
	return &HeatTransferComponent{
		BaseComponent: NewBaseComponent(),
		ParticleType:  particleType,
		Temperature:   20.0,
		BaseColor:     baseColor,
	}
}

func (h *HeatTransferComponent) Type() uint64 {
	return ComponentTypeHeatTransfer
}

func (h *HeatTransferComponent) CurrentColor() color.Color {
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

type ZIndexComponent struct {
	BaseComponent
	ZIndex int
}

func (z *ZIndexComponent) Type() uint64 {
	return ComponentTypeZIndex
}

func NewZIndexComponent(z int) *ZIndexComponent {
	return &ZIndexComponent{
		BaseComponent: NewBaseComponent(),
		ZIndex:        z,
	}
}

type FissionComponent struct {
	BaseComponent
	InducingParticleType uint64
}

func NewFissionComponent(inducingParticleType uint64) *FissionComponent {
	return &FissionComponent{
		BaseComponent:        NewBaseComponent(),
		InducingParticleType: inducingParticleType,
	}
}

func (f *FissionComponent) Type() uint64 {
	return ComponentTypeFission
}
