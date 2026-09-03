package ecs

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

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

// componentRegistry assegna ComponentID crescenti e densi ai tipi di
// componente registrati a runtime (tipicamente in init()). Non usa Sonyflake
// o altri generatori di ID sparsi: qui serve un contatore piccolo e denso,
// adatto a essere usato come indice di slice.
type componentRegistry struct {
	mu     sync.RWMutex
	nextID ComponentID
	byID   map[ComponentID]reflect.Type
	byType map[reflect.Type]ComponentID
}

func newComponentRegistry() *componentRegistry {
	return &componentRegistry{
		byID:   make(map[ComponentID]reflect.Type),
		byType: make(map[reflect.Type]ComponentID),
	}
}

// registry è l'istanza globale del componentRegistry del motore ECS
var registry = newComponentRegistry()

// RegisterComponent registra un tipo di componente T presso il registry
// globale e ne restituisce il ComponentID. Se T implementa Typed, la
// registrazione viene rifiutata: quei componenti usano il fast-path bitmask
// e non devono passare dal registry dinamico. Chiamate multiple con lo stesso
// tipo T restituiscono dsempre lo stesso ID (idempotente).
func RegisterComponent[T Component]() (ComponentID, error) {
	var zero T

	if _, ok := any(zero).(TypedComponent); ok {
		return 0, fmt.Errorf("Tenet: %T implementa TypedComponent, non va registrato nel registro dinamico", zero)
	}

	t := reflect.TypeOf(zero)

	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t == nil {
		return 0, errors.New("Tenet: impossibile determinare il tipo del componente")
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if id, ok := registry.byType[t]; ok {
		return id, nil
	}

	registry.nextID++
	id := registry.nextID
	registry.byID[id] = t
	registry.byType[t] = id

	return id, nil
}
