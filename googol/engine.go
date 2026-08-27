package googol

import "github.com/tetrahedronix/coop/googol/ecs"

type Loop struct {
}

func NewLoop() {

}

func (l Loop) Update() {}

func (l Loop) MainLoop() {}

func (l Loop) Add(sys ecs.System) {}
