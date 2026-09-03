package ecs

const (
	TileFlipHorizFlag = 0x80000000
	TileFlipVertFlag  = 0x40000000
	TileFlipDiagFlag  = 0x20000000
	TileIDMask        = 0x1FFFFFFF
)

type Tile struct {
	ComponentType ComponentTypeID
	GID           uint32 // Global Tile ID (inclide flag di flip)
}

// Add imposta il GID. data deve essere uint32
func (t *Tile) Add(data any) {
	t.GID = data.(uint32)
}

// Copy copia il GID da un altro Tile
func (t *Tile) Copy(src any) {
	t.GID = src.(*Tile).GID
}

// Get restituisce il GID grezzo
func (t *Tile) Get() any {
	return t.GID
}

// Reset azzera il GID
func (t *Tile) Reset() {
	t.GID = 0
}

func (t *Tile) TypeID() ComponentTypeID {
	return t.ComponentType
}

// NewTile restituisce un nuovo componente Tile
func NewTile() Component {
	return &Tile{}
}
