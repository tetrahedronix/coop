package tenet

import "github.com/tetrahedronix/coop/tenet/ecs"

// Alias per tipi di componenti
type (
	Coordinate = ecs.Coordinate
	Position   = ecs.Position
	Shape      = ecs.Shape
	Tile       = ecs.Tile
	Velocity   = ecs.Velocity
	Selectable = ecs.Selectable
)

// Alias per l'entità
type Entity = ecs.Entity

// Alias per il sistema
type System = ecs.System

// Aliasi per l'interfaccia Component
type Component = ecs.Component
