package ebitendisplay

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nitwhiz/reactor/pkg/geometry"
	"github.com/nitwhiz/reactor/pkg/sim"
	"golang.org/x/image/colornames"
	"log"
)

type Reactor struct {
	world *sim.EntityManager

	offscreenImage *ebiten.Image
	drawOp         *ebiten.DrawImageOptions

	renderQuery *sim.Query

	screenWidth  int
	screenHeight int
}

func NewReactor(em *sim.EntityManager, screenWidth int, screenHeight int) *Reactor {
	renderQuery := em.Query(sim.ComponentTypeBody | sim.ComponentTypeParticleType | sim.ComponentTypeZIndex)
	renderQuery.RegisterHook(sim.SortHook)

	return &Reactor{
		world:          em,
		offscreenImage: ebiten.NewImage(screenWidth, screenHeight),
		drawOp:         &ebiten.DrawImageOptions{},
		renderQuery:    renderQuery,
		screenWidth:    screenWidth,
		screenHeight:   screenHeight,
	}
}

func (r *Reactor) Render() {
	r.offscreenImage.Fill(colornames.White)

	for _, eId := range r.renderQuery.Ids() {
		bodyComponent := r.world.GetComponent(sim.ComponentTypeBody, eId).(*sim.BodyComponent)
		particleTypeComponent := r.world.GetComponent(sim.ComponentTypeParticleType, eId).(*sim.ParticleTypeComponent)

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
}

func (r *Reactor) Update() error {
	r.world.Update()
	r.Render()

	return nil
}

func (r *Reactor) Draw(screen *ebiten.Image) {
	r.drawOp.GeoM.Reset()
	screen.DrawImage(r.offscreenImage, r.drawOp)

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
