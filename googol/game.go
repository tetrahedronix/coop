package googol

import (
	"fmt"
	"io"

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
	Logger *log.Logger
	// Whether the world tick should execute
	enabled bool
}

// createEntity creates a new entity struct
func (w *World) CreateEntity() (*ecs.Entity, error) {

	// Gets a new unique ID with Sonyflake package.
	id, err := flake.NextID()

	if err != nil {
		return nil, fmt.Errorf("flake.NextID() failed with %w\n", err)
	}

	w.Logger.Printf("a new Entity has been initialized with EGUID: [%x]\n", id)

	return ecs.NewEntity(id), nil

}

func (w *World) CreateSystem() *ecs.System {
	return &ecs.System{}
}

func NewWorld() *World {

	return &World{
		Logger:  log.New(io.Discard, "", 0),
		enabled: true,
	}
}
