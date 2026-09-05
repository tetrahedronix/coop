package ecs

// Entities have IDs
type Guid interface {
	Id() uint64
}

// Each EntityID is nothing more than a Globally Unique Identifier (GUID)
// with components attached to it.
type EntityID uint64

// Id resituisce l'identificatore numerico dell'entità. È l'unica operazione
// che ha senso su un EntityID puro.
func (id EntityID) Id() uint64 {
	return uint64(id)
}
