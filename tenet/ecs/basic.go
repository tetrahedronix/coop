package ecs

// Per il momento solo mondi 2D
type Coordinate [2]float64

// The position of the entity on the screen.
// Il campo Coordinate è accessibile e modificabile direttamente da Store[t] e
// dai System. La copia è una semplice copia di struct (value type pure),
// quindi Position non necessita di implementare DeepCopier.
type Position struct {
	Coordinate Coordinate
}

func (p *Position) TypeID() TypedComponentID {
	return ComponentTypePosition
}

// Selectable indica se l'entità è correntemente selezionata (es da un sistema
// di input).
type Selectable struct {
	Selected bool
}

func (s *Selectable) TypedID() TypedComponentID {
	return ComponentTypeSelectable
}

// Shape descrive la primitiva grafica da disegnare (es. "circle", "box", ecc.)
type Shape struct {
	Primitive string
}

func (s *Shape) TypedID() TypedComponentID {
	return ComponentTypeShape
}

type Speed float64

type Direction float64

// Velocity descrive la velocità e direzione di movimento dell'entità
type Velocity struct {
	Speed     Speed
	Direction Direction
}

func (v *Velocity) TypedID() TypedComponentID {
	return ComponentTypeVelocity
}

const (
	TileFlipHorizFlag = 0x80000000
	TileFlipVertFlag  = 0x40000000
	TileFlipDiagFlag  = 0x20000000
	TileIDMask        = 0x1FFFFFFF
)

// Tile rappresenta un tile della mappa. GID incorpora sia l'ID grezzo del
// tile sia i flag di flip (orizzontale/verticale/diagonale), estraibili
// tramite le costanti sopra.
type Tile struct {
	GID uint32 // Global Tile ID (inclide flag di flip)
}

func (t *Tile) TypedID() TypedComponentID {
	return ComponentTypeTile
}
