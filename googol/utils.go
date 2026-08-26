package googol

import "gitlab.com/tetrahedronix/coop/googol/ecs"

func EqualPosition(a, b *ecs.Position) bool {

	//return a.coordinate == b.coordinate

	return a.Coordinate == b.Coordinate

}

// Helper per sistemi che vogliono leggere dal past
func GetPastComponent(e *ecs.Entity, i int) ecs.Component {

	if i < 0 || i >= e.LenComponent() {
		return nil
	}

	return e.GetComponent(i)
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

func NewShape() *Shape {
	return ecs.NewShape().(*Shape)
}
