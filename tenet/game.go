package tenet

import (
	"fmt"
	"io"
	"sync"

	"github.com/sony/sonyflake"
	"github.com/tetrahedronix/coop/tenet/ecs"

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
	// Proteggere a slice di entità con un mutex
	mu           sync.RWMutex
	entities     []*ecs.Entity
	deleteEntity map[uint64]bool
	Logger       *log.Logger
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

	e := ecs.NewEntity(id)

	w.mu.Lock()
	w.entities = append(w.entities, e)
	w.mu.Unlock()
	w.Logger.Printf("a new Entity has been initialized with EGUID: [%x]\n", id)

	return e, nil
}

func (w *World) GetEntities() []*ecs.Entity {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return append([]*ecs.Entity{}, w.entities...)
}

func (w *World) MarkForRemoval(id uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.deleteEntity == nil {
		w.deleteEntity = make(map[uint64]bool)
	}

	w.deleteEntity[id] = true
}

func (w *World) PurgeRemoved() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.deleteEntity) == 0 {
		return
	}

	keep := w.entities[:0]

	for _, e := range w.entities {
		if !w.deleteEntity[e.Id()] {
			keep = append(keep, e)
		}
	}

	w.entities = keep

	// Svuota la mpaa
	w.deleteEntity = make(map[uint64]bool)
}

// RemoveEntity rimuove l'entità con l'ID specificato dal mondo.
// Restituisce true se l'entità è stata trovata e rimossa, false altrimenti
func (w *World) RemoveEntity(id uint64) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	for i, e := range w.entities {
		if e.Id() == id {
			// Rimuove l'elemento in posizione i
			w.entities = append(w.entities[:i], w.entities[i+1:]...)
			return true
		}
	}

	return false
}

func (w *World) SwapBuffers() {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, e := range w.entities {
		if !ecs.SwapBuffers(e) {
			w.Logger.Printf("swap failed for entity %x", e.Id())
		}
	}
}

func NewWorld() *World {

	return &World{
		Logger:  log.New(io.Discard, "", 0),
		enabled: true,
	}
}
