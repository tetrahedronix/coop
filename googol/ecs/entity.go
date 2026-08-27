package ecs

import "fmt"

// Entities have IDs
type Guid interface {
	Id() uint64
}

// Each Entity is nothing more than a Globally Unique Identifier (GUID)
// with components attached to it.
type Entity struct {
	// Use Sonyflake to get unique GUID
	guid             uint64
	componentsPast   []Component
	componentsFuture []Component
}

// Add a component to the entity
// Metodo pubblico per aggiungere al past (solo inizializzazione)
func (e *Entity) AddComponent(c Component, data ...interface{}) {

	e.componentsPast = append(e.componentsPast, c)

	for _, d := range data {
		e.componentsPast[len(e.componentsPast)-1].Add(d)
	}
}

func (e *Entity) AddFutureComponent(c Component, data ...any) {

	e.componentsFuture = append(e.componentsFuture, c)

	for _, d := range data {
		e.componentsFuture[len(e.componentsFuture)-1].Add((d))
	}
}

func (e *Entity) LenComponent() int {
	return len(e.componentsPast)
}

func (e *Entity) LenFutureComponent() int {
	return len(e.componentsFuture)
}

// Metodo pubblico per accedere al past (lettura)
func (e *Entity) GetComponent(i int) Component {

	return e.componentsPast[i]
}

// Metodo pubblico per accedere al future (scrittura)
func (e *Entity) GetFutureComponent(i int) Component {
	return e.componentsFuture[i]
}

// GetWritable restituisce un componente modificabile dal future per l'indice i.
// Se il future non ha un componente all'indice i, e il past ce l'ha, lo clona

func (e *Entity) GetWritable(i int) Component {
	if i < 0 {
		return nil
	}

	if len(e.componentsFuture) <= i {
		// e.componentsFuture = append(e.componentsFuture, nil)
		// Crea un array di lunghezza i+1, preservando gli esistenti
		newFuture := make([]Component, i+1)
		copy(newFuture, e.componentsFuture)
		e.componentsFuture = newFuture
	}

	if e.componentsFuture[i] != nil {
		return e.componentsFuture[i]
	}

	if i < len(e.componentsPast) && e.componentsPast[i] != nil {
		// Usa la funzione helper di copia esterna
		e.componentsFuture[i] = CopyComponent(e.componentsPast[i])

		return e.componentsFuture[i]
	}

	return nil

}

func (e *Entity) String() string {
	return fmt.Sprintf("Entity[%x] (%d past, %d future)",
		e.guid, len(e.componentsPast), len(e.componentsFuture))
}

func (e *Entity) Id() uint64 {

	return e.guid
}

// NewEntity crea un nuovo Entity con l'ID specificato
func NewEntity(id uint64) *Entity {
	return &Entity{
		guid: id,
	}
}
