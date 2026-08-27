package main

import (
	"fmt"

	"github.com/tetrahedronix/coop/googol"
)

func main() {
	// Binding to Googol engine: Crea un mondo
	//
	gw := googol.NewWorld()

	// Crea una nuova entità nel mondo
	e := gw.CreateEntity()

	// Aggiunge un componente Shape (con dato "circle")
	e.AddComponent(googol.NewShape(), "circle")

	// Aggiunge un componente Position (con coordinate)
	e.AddComponent(googol.NewPosition(), googol.Coordinate{10.0, 20.0})

	// Crea un sistema (demo, non fa nulla)
	gw.CreateSystem()

	// Stampa le informazioni dell'entità
	fmt.Printf("Entity ID: %x\n", e.Id())
	fmt.Printf("Numero di componenti: %d\n", e.LenComponent())

	// Itera e stampa ogni componente
	for i := 0; i < e.LenComponent(); i++ {
		comp := e.GetComponent(i)
		fmt.Printf("Componente %d: %T %+v\n", i, comp, comp)
		// Mostra anche i dati interni tramite il metodo Get()
		if comp != nil {
			fmt.Printf("	Dati: %v\n", comp.Get())
		}
	}

}
