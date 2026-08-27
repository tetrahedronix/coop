package googol

import (
	"github.com/sony/sonyflake"
	"github.com/tetrahedronix/coop/googol/ecs"

	"log"
)

func init() {

	// Create a new Sonyflake instance 'flake' configured with
	// the given argument. https://github.com/sony/Sonyflake */
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})

	if flake == nil {
		panic("sonyflake not created")
	}
}

var flake *sonyflake.Sonyflake

type World struct {
	// Whether the world tick should execute
	enabled bool
}

func NewWorld() *World {

	return &World{true}
}

// createEntity creates a new entity struct
func (w *World) CreateEntity() *ecs.Entity {

	// Gets a new unique ID with Sonyflake package.
	id, err := flake.NextID()

	if err != nil {
		// DEBUG
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}

	log.Printf("a new Entity has been initialized with EGUID: [%x]\n", id)

	return ecs.NewEntity(id)

}

func (w *World) CreateSystem() *ecs.System {
	return &ecs.System{}
}
