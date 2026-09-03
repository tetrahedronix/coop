package ecs

import (
	"fmt"
	"reflect"
)

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

// AddComponent aggiunge un componente al buffer past (solo inizializzazione).
// Se c implementa TypedComponent, la sua identità viene registrata nella
// bitmask signature (fast-path). Altrimenti, l'identità viene cercata nel
// registry globale e registrata in dynamicSignature (slow-path). Il
// componente deve essere già stato registrato con RegisterComponent se non è
// TypedComponent, altrimenti AddComponent restituisce un errore.
func (e *Entity) AddComponent(c Component) error {

	if tc, ok := c.(TypedComponent); ok {
		e.signature |= tc.TypedID()
		e.componentsPast = append(e.componentsPast, c)

		return nil
	}

	t := reflect.TypeOf(c)

	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	registry.mu.RLock()
	id, ok := registry.byType[t]
	registry.mu.RUnlock()

	if !ok {
		return fmt.Errorf("Tenet: componente %T non registrato; chiamare RegisterComponent prima di AddComponent", c)
	}

	if e.dynamicSignature == nil {
		e.dynamicSignature = make(map[ComponentID]struct{})
	}

	e.dynamicSignature[id] = struct{}{}
	e.componentsPast = append(e.componentsPast, c)

	return nil
}

// HasComponent verifica se l'entità possiede un componente registrato
// dinamicamente (slow-path, via registry). Per i componente tipizzati
// (bitmask) usare HasTypedComponent.
func (e *Entity) HasComponent(id ComponentID) bool {
	if e.dynamicSignature == nil {
		return false
	}

	_, ok := e.dynamicSignature[id]

	return ok
}

// HasTypedComponent verifica se l'entità possiede un componente Typed
// (fast-path O(1), bia bitmask). Per i componenti registrati dinamicamente
// usare invece HasComponent.
func (e *Entity) HasTypedComponent(tid TypedComponentID) bool {
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
