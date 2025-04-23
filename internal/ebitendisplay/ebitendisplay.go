package ebitendisplay

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nitwhiz/reactor/pkg/ecs"
	"github.com/nitwhiz/reactor/pkg/geometry"
	"github.com/nitwhiz/reactor/pkg/sim"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
)

type Reactor struct {
	world *ecs.EntityManager

	offscreenImage *ebiten.Image
	imageDrawOp    *ebiten.DrawImageOptions
	textDrawOp     *text.DrawOptions

	neutronCount int

	textFaceSource *text.GoTextFaceSource

	renderQuery *ecs.Query

	screenWidth  int
	screenHeight int
}

func NewReactor(world *ecs.EntityManager, screenWidth int, screenHeight int) *Reactor {
	renderQuery := world.Query(sim.ComponentTypeBody | sim.ComponentTypeParticle | sim.ComponentTypeRender)
	renderQuery.RegisterHook(sim.SortHook)

	ff, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))

	if err != nil {
		log.Fatal(err)
	}

	textDrawOp := text.DrawOptions{}

	textDrawOp.GeoM.Translate(float64(screenWidth)/2.0, 60)
	textDrawOp.LayoutOptions.PrimaryAlign = text.AlignCenter
	textDrawOp.ColorScale.ScaleWithColor(color.Black)

	return &Reactor{
		world:          world,
		offscreenImage: ebiten.NewImage(screenWidth, screenHeight),
		imageDrawOp:    &ebiten.DrawImageOptions{},
		textDrawOp:     &textDrawOp,
		renderQuery:    renderQuery,
		textFaceSource: ff,
		screenWidth:    screenWidth,
		screenHeight:   screenHeight,
	}
}

func (r *Reactor) Render() {
	r.offscreenImage.Fill(colornames.White)

	for _, eId := range r.renderQuery.Ids() {
		bodyComponent := r.world.GetComponent(sim.ComponentTypeBody, eId).(*sim.BodyComponent)
		particleTypeComponent := r.world.GetComponent(sim.ComponentTypeParticle, eId).(*sim.ParticleComponent)

		body := bodyComponent.Body
		bodyLoc := body.Location()

		switch b := body.(type) {
		case *geometry.Rectangle:
			bodyLoc := b.TopLeft()

			vector.DrawFilledRect(
				r.offscreenImage,
				bodyLoc.X,
				bodyLoc.Y,
				b.Width,
				b.Height,
				particleTypeComponent.ParticleType.Color,
				false,
			)
			break
		case *geometry.Circle:
			vector.DrawFilledCircle(
				r.offscreenImage,
				bodyLoc.X,
				bodyLoc.Y,
				b.Radius,
				particleTypeComponent.ParticleType.Color,
				false,
			)
			break
		default:
			log.Println("unknown body type")
		}
	}

	text.Draw(
		r.offscreenImage,
		fmt.Sprintf("%d", r.neutronCount),
		&text.GoTextFace{Source: r.textFaceSource, Size: 24},
		r.textDrawOp,
	)
}

func (r *Reactor) Update() error {
	r.world.Update()

	r.neutronCount = r.world.CountComponent(sim.TagThermalNeutron) + r.world.CountComponent(sim.TagFastNeutron)

	return nil
}

func (r *Reactor) Draw(screen *ebiten.Image) {
	r.Render()

	r.imageDrawOp.GeoM.Reset()
	screen.DrawImage(r.offscreenImage, r.imageDrawOp)

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"FPS: %.1f\nTPS: %.1f\nE: %d",
			ebiten.ActualFPS(),
			ebiten.ActualTPS(),
			r.world.EntityCount(),
		),
	)
}

func (r *Reactor) Layout(_, _ int) (screenWidth, screenHeight int) {
	return r.screenWidth, r.screenHeight
}

func (r *Reactor) Start() error {
	ebiten.SetWindowSize(r.screenWidth, r.screenHeight)
	ebiten.SetWindowTitle("reactor")
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(60)

	return ebiten.RunGame(r)
}
