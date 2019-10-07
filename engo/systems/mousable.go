package systems

import (
	"github.com/EngoEngine/ecs"
)

type (
	MouseEvent     func(dt float32)
	MousableSystem struct {
		Callback MouseEvent
	}
)

var (
	_ ecs.System = &MousableSystem{}
)

func (*MousableSystem) Remove(ecs.BasicEntity) {}

func (s *MousableSystem) Update(dt float32) {
	if s.Callback != nil {
		s.Callback(dt)
	}
}
