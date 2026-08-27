package googol

import "github.com/tetrahedronix/coop/googol/ecs"

// Alias per tipi di componenti
type (
	Coordinate = ecs.Coordinate
	Position   = ecs.Position
	Shape      = ecs.Shape
	Velocity   = ecs.Velocity
)

// Alias per l'entità
type Entity = ecs.Entity

// Alias per il sistema
type System = ecs.System

// Aliasi per l'interfaccia Component
type Component = ecs.Component
