package googol

import (
	"github.com/tetrahedronix/coop/googol/ecs"
)

func EqualPosition(a, b *ecs.Position) bool {

	//return a.coordinate == b.coordinate

	return a.Get() == b.Get()

}

// Helper per i flag di flip
func FlippedHorizontally(t *ecs.Tile) bool {

	return (t.GID & ecs.TileFlipHorizFlag) != 0
}

func FlippedVertically(t *ecs.Tile) bool {

	return (t.GID & ecs.TileFlipVertFlag) != 0
}

func FlippedDiagonally(t *ecs.Tile) bool {

	return (t.GID & ecs.TileFlipDiagFlag) != 0
}

// Helper per sistemi che vogliono leggere dal past
func GetPastComponent(e *ecs.Entity, i int) ecs.Component {

	if i < 0 || i >= e.LenComponent() {
		return nil
	}

	return e.GetComponent(i)
}

// GetPastComponentByType restituisce il componente del buffer "past"
// dell'entità che ha il ComponentTypeID specificato. Se non presente,
// restituisce nil.
func GetPastComponentByType(e *ecs.Entity, tid ecs.ComponentTypeID) ecs.Component {
	if e == nil {
		return nil
	}

	// Controllo rapido tramite signature (metodo pubblico)
	if !e.HasComponent(tid) {
		return nil
	}

	// Scorre la slice appropriata (past o future) linearmente e restituire la
	// prima componente che ha TypeID() uguale a tid: semplice, lento O(n) ma
	// con pochi componenti di base è OK.
	for i := 0; i < e.LenComponent(); i++ {
		c := e.GetComponent(i)

		if c != nil && c.TypeID() == tid {
			return c
		}
	}

	return nil
}

// GetFutureCompoentByType restituisce il componente del buffer "future"
// dell'entità che ha il ComponentTypeId specificato. Se non presente,
// restituisce nil.
func GetFutureCompoentByType(e *ecs.Entity, tid ecs.ComponentTypeID) ecs.Component {
	if e == nil {
		return nil
	}

	// Nota: la signature riflette i componenti nel past, non nel future.
	// Se un componente è stato aggiunto solo al future (durante l'elaborazione
	// ), la signature potrebbe non averlo ancora. Quindi per il future non
	// è possibile usare HasComponent(). Occorre scandire direttamente.
	for i := 0; i < e.LenComponent(); i++ {
		c := e.GetFutureComponent(i)

		if c != nil && c.TypeID() == tid {
			return c
		}

	}

	return nil
}

// GetComponentFuture restituisce una copia del componente dall'indice
// specificato nella slice past e lo inserisce nella slice future, pronto per
// essere modificato.
// Se l'indice è fuori range, restituisce nil.
// Helper per sistemi che vogliono scrivere sul future
func GetFutureComponent(e *ecs.Entity, i int) ecs.Component {

	if i < 0 || i >= e.LenFutureComponent() {
		return nil
	}

	return e.GetFutureComponent(i)
}

// Helper per ottenere un componente un componente modificabile (copia dal past
// al future)
func GetWritableComponent(e *ecs.Entity, i int) ecs.Component {
	if i < 0 {
		return nil
	}

	return e.GetWritable(i) // chiama il metodo pubblico

}

// Funzioni helper per creare componenti

func NewPosition() *Position {
	return ecs.NewPosition().(*Position)
}

func NewSelectable() *Selectable {
	return ecs.NewSelectable().(*Selectable)
}

func NewShape() *Shape {
	return ecs.NewShape().(*Shape)
}

func NewTile() *Tile {
	return ecs.NewTile().(*Tile)
}

func NewVelocity() *Velocity {
	return ecs.NewVelocity().(*Velocity)
}

// TileID estrae il tile ID puro
func TileID(t *Tile) uint32 {
	return t.GID & ecs.TileIDMask
}
