package scenes

import (
	"github.com/atercattus/golangconf19_examples/handmade/render"
)

type (
	SceneManager struct {
		scenes   map[string]Scener
		curScene Scener

		renderer *render.WebGLRender
	}
)

func NewSceneManager(renderer *render.WebGLRender) *SceneManager {
	sm := &SceneManager{
		scenes: make(map[string]Scener),

		renderer: renderer,
	}
	return sm
}

func (sm *SceneManager) Register(name string, scene Scener) {
	if _, ok := sm.scenes[name]; ok {
		println(`Scene "` + name + `" is already registered`)
	} else {
		scene.Setup(sm.renderer, sm)
		sm.scenes[name] = scene
	}
}

func (sm *SceneManager) Goto(name string) {
	if scene, ok := sm.scenes[name]; !ok {
		println(`Unknown scene name: `, name)
		return
	} else if scene == sm.curScene {
		return
	} else {
		if sm.curScene != nil {
			sm.curScene.Hide()
		}
		println(`Goto scene ` + name)
		sm.curScene = scene
		sm.curScene.Show()
	}
}

func (sm *SceneManager) Current() Scener {
	return sm.curScene
}
