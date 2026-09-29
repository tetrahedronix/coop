package ecs

type TypedComponentID uint64

const (
	ComponentTypePosition TypedComponentID = 1 << iota
	ComponentTypeSelectable
	ComponentTypeShape
	ComponentTypeVelocity
	ComponentTypeTile
)

// Component è il contratto minimio di ogni componente dell'engine. Un
// componente che necessita di copia profonda implementa in aggiunta
// l'interfaccia DeepCopier: non è più richiesto implementare
// Add/Copy/Get/Reset per essere un componente valido.
type Component interface {
	ECSComponent()
}

// TypedComponent è un'interfaccia satellite opzionale: i componenti che la
// implementano ottengono un fast-path O(1) basato su bitmask fissa a
// compile-time. I componenti che non la implementano vengono comunque
// identificati tramite il registry (via relect.Tye), senza bisogno di alcun
// metodo stub.
type TypedComponent interface {
	TypedID() TypedComponentID
}

type ComponentID uint64

// TypedStore[T] è il contenitore dati per un singolo tipo di componente
// bitmask (T deve implementare TypedComponent). Layout: Structure-of-Arrays.
// Campi: index mappa EntityID -> indice di riga; past e future sono slice
// dense parallele, stesso indice di riga per entrambe. A differenza del
// modello OOP "entity-as-container", un'entità che non possiede T
// semplicemente non ha una entry in index, di conseguenza nessuno slot viene
// sprecato.
type TypedStore[T TypedComponent] struct {
	index  map[EntityID]int
	past   []T
	future []T
}

// NewTypedStore crea un TypedStore vuoto, pronto per ricevere componenti.
func NewTypedStore[T TypedComponent]() *TypedStore[T] {
	return &TypedStore[T]{
		index: make(map[EntityID]int),
	}
}
