package sim

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log/slog"
	"slices"
	"sync"
	"time"
)

type EnvSettings struct {
	RoomTemperature            float32
	WaterVaporizeTemperature   float32
	WaterTemperatureChangeRate float32
	WaterNeutronAbsorbRate     float32
	NeutronWaterHeating        float32

	UpdateWaterTemperature bool
}

type Env struct {
	objects  map[int64]Object
	layers   map[int]map[int64]Object
	zIndexes []int

	lastUpdated time.Time

	mu *sync.Mutex

	settings *EnvSettings
}

func NewEnv(settings *EnvSettings) *Env {
	return &Env{
		objects:     map[int64]Object{},
		layers:      map[int]map[int64]Object{},
		zIndexes:    []int{},
		lastUpdated: time.Now(),
		mu:          &sync.Mutex{},
		settings:    settings,
	}
}

func (e *Env) Update() {
	e.mu.Lock()
	defer e.mu.Unlock()

	d := time.Since(e.lastUpdated)
	e.lastUpdated = time.Now()

	for _, o := range e.objects {
		o.Update(d)

		if o.Location().X < 0 || o.Location().Y < 0 || o.Location().X > 1280 || o.Location().Y > 720 {
			e.Remove(o)
		}
	}

	for _, o1 := range e.objects {
		for _, o2 := range e.objects {
			if o1 != o2 && o1.Intersects(o2) {
				o1.React(o2)
			}
		}
	}
}

func (e *Env) Draw(screen *ebiten.Image) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, z := range e.zIndexes {
		for _, o := range e.layers[z] {
			o.Draw(screen)
		}
	}
}

func (e *Env) Add(o Object) {
	id := o.Id()

	slog.Info("adding body", "id", id)

	_, ok := e.objects[id]

	if ok {
		panic(fmt.Sprintf("body already exists: %d", id))
	}

	e.objects[id] = o

	z := o.ZIndex()

	l, ok := e.layers[z]

	if !ok {
		l = map[int64]Object{}
		e.layers[z] = l
	}

	l[id] = o

	if !slices.Contains(e.zIndexes, z) {
		e.zIndexes = append(e.zIndexes, z)
		slices.Sort(e.zIndexes)
	}
}

func (e *Env) Remove(o Object) {
	id := o.Id()

	slog.Info("removing body", "id", id)

	delete(e.objects, id)
	delete(e.layers[o.ZIndex()], id)
}
