package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/atercattus/golangconf19_examples/engo/common"
)

type (
	RoadSystem struct {
		lines      []*common.Sprite
		Speed      float32
		MinX, MaxX float32
	}
)

var (
	_ ecs.System = &RoadSystem{}
)

func (r *RoadSystem) Add(line *common.Sprite) {
	r.lines = append(r.lines, line)
}

func (*RoadSystem) Remove(ecs.BasicEntity) {}

func (r *RoadSystem) Update(dt float32) {
	if r.Speed == 0 {
		return
	}

	diff := dt * r.Speed
	for _, line := range r.lines {
		line.Position.X -= diff
		if line.Position.X <= r.MinX {
			line.Position.X += r.MaxX
		}
	}
}
