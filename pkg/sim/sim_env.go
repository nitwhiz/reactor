package sim

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log/slog"
	"sync"
	"time"
)

type Env struct {
	objects     map[int64]Object
	bg          map[int64]Object
	obj0        map[int64]Object
	lastUpdated time.Time

	mu *sync.Mutex
}

func NewEnv() *Env {
	return &Env{
		objects:     map[int64]Object{},
		bg:          map[int64]Object{},
		obj0:        map[int64]Object{},
		lastUpdated: time.Now(),
		mu:          &sync.Mutex{},
	}
}

func (e *Env) Update() {
	e.mu.Lock()
	defer e.mu.Unlock()

	d := time.Since(e.lastUpdated)
	e.lastUpdated = time.Now()

	for _, o := range e.objects {
		o.Update(d)
	}

	for _, o1 := range e.objects {
		for _, o2 := range e.objects {
			if o1 != o2 && o1.Intersects(o2) {
				o1.React(e, o2)
			}
		}
	}
}

func (e *Env) Draw(screen *ebiten.Image) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, o := range e.bg {
		o.Draw(screen)
	}

	for _, o := range e.obj0 {
		o.Draw(screen)
	}
}

func (e *Env) Add(o Object) {
	id := o.Id()

	slog.Info("adding object", "id", id)

	_, ok := e.objects[id]

	if ok {
		panic(fmt.Sprintf("object already exists: %d", id))
	}

	e.objects[id] = o

	switch o.(type) {
	case *Particle:
		e.obj0[id] = o
		break
	case *Body:
		e.bg[id] = o
		break
	}
}

func (e *Env) Remove(o Object) {
	id := o.Id()

	slog.Info("removing object", "id", id)

	delete(e.objects, id)
	delete(e.bg, id)
	delete(e.obj0, id)
}
