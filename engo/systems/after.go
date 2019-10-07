package systems

import (
	"github.com/EngoEngine/ecs"
	"time"
)

type AfterSystem struct {
	world    *ecs.World
	delay    float32
	callback func()
}

var (
	_ ecs.System = &AfterSystem{}
)

func After(world *ecs.World, delay time.Duration, callback func()) {
	world.AddSystem(&AfterSystem{
		world:    world,
		delay:    float32(float64(delay) / float64(time.Second)),
		callback: callback,
	})
}

func (*AfterSystem) Remove(ecs.BasicEntity) {
}

func (s *AfterSystem) Update(dt float32) {
	if s.delay <= 0 {
		return
	}

	if s.delay -= dt; s.delay <= 0 {
		// ToDo: удалять событие (завести минхип срабатываний)
		s.callback()
	}
}
