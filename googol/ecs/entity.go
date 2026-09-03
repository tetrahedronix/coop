package ecs

import "fmt"

// Entities have IDs
type Guid interface {
	Id() uint64
}

// Each Entity is nothing more than a Globally Unique Identifier (GUID)
// with components attached to it.
// ----------------------------------------------------------------------------
// EXTENSIBILITY NOTE (identità ibrida: bitmask + registry)
//
// L'entità traccia i componenti con due meccanismi distinti, perché
// rispondono a domande diverse e non sono unificabili in un solo campo
// senza rischiare collisioni tra ID di natura diversa:
//
//   - signature (TypedComponentID, bitmask a 64 bit): fast-path O(1) per i
//     componenti che implementano TypedComponent. Ogni tipo occupa un bit
//     riservato
//     a compile-time (1 << iota). Limite: max 64 tipi TypedComponent.
//
//   - dynamicSignature (set di ComponentID): componenti registrati a
//     runtime tramite RegisterComponent, per i tipi che NON implementano
//     TypedComponent. Nessun limite di cardinalità, ma lookup non O(1) come il
//     bitmask (costo di un accesso a mappa).
//
// Percorso di refactoring futuro (da valutare SOLO a POC completo e
// funzionante, per non introdurre complessità non necessaria ora):
// sostituire questo schema a due campi con un bitset a lunghezza variabile
// unico, in cui i bit 0-63 restano riservati ai componenti TypedComponent e
// i bit 64+ vengono assegnati dinamicamente agli ID di registry man mano che
// vengono registrati. Questo richiede un tipo bitset custom (non un
// singolo uint64) e la riscrittura di HasComponent/AddComponent per
// operare su di esso; è un cambiamento strutturale, non incrementale,
// quindi va pianificato come refactoring dedicato a parte.
// ------------------------------------------------------------------
type Entity struct {
	// Use Sonyflake to get unique GUID
	guid uint64
	// fast-path: solo componenti Typed
	signature TypedComponentID
	// componenti registrati dinamicamente
	dynamicSignature map[ComponentID]struct{}
	componentsPast   []Component
	componentsFuture []Component
}

// Add a component to the entity
// Metodo pubblico per aggiungere al past (solo inizializzazione)
func (e *Entity) AddComponent(c Component, data ...interface{}) {

	e.componentsPast = append(e.componentsPast, c)
	e.signature |= c.TypeID()

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

// Metodo pubblico per accedere al past (lettura)
func (e *Entity) GetComponent(i int) Component {

	return e.componentsPast[i]
}

// GetComponents returns a slice of the entity's past components.
// The returned slice is read‑only; mutating it may break the double‑buffer invariants.
func (e *Entity) GetComponents() []Component {
	return e.componentsPast
}

// Metodo pubblico per accedere al future (scrittura)
func (e *Entity) GetFutureComponent(i int) Component {
	return e.componentsFuture[i]
}

// GetFutureComponents returns a slice of the entity's future components.
// The returned slice is read‑only; mutating it may break the double‑buffer invariants.
func (e *Entity) GetFutureComponents() []Component {
	return e.componentsFuture
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

func (e *Entity) HasComponent(tid ComponentTypeID) bool {
	return (e.signature & tid) != 0
}

func (e *Entity) LenComponent() int {
	return len(e.componentsPast)
}

func (e *Entity) LenFutureComponent() int {
	return len(e.componentsFuture)
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

func Normalize(e *Entity) {
	// Allarga la slice componentsFuture affinché abbia la stessa lunghezza del past.
	if e.LenFutureComponent() < e.LenComponent() {
		newFuture := make([]Component, e.LenComponent())
		copy(newFuture, e.componentsFuture)
		e.componentsFuture = newFuture
	}

	// Copia (clona) i componenti dal past al future dove il future è nil.
	for i := 0; i < e.LenComponent(); i++ {
		if e.componentsFuture[i] == nil && i < len(e.componentsPast) && e.componentsPast[i] != nil {
			e.componentsFuture[i] = CopyComponent(e.componentsPast[i])
		}
	}
}

// SwapBuffers scambia i due buffer di componenti dell'entità.
// Restituisce false se l'entità è nil o se i buffer hanno lunghezze diverse.
func SwapBuffers(e *Entity) bool {
	if e == nil {
		return false
	}

	if e.LenComponent() != e.LenFutureComponent() {
		return false
	}

	e.componentsPast, e.componentsFuture = e.componentsFuture, e.componentsPast

	return true
}
