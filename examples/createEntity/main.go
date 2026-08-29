package main

import (
	"fmt"
	"os"

	"github.com/tetrahedronix/coop/googol"
)

func main() {
	// Binding to Googol engine: Crea un mondo
	//
	gw := googol.NewWorld()

	// Crea una nuova entità nel mondo
	entity, err := gw.CreateEntity()

	if err != nil {
		os.Exit(1)
	}

	// Aggiunge un componente Shape (con dato "circle")
	entity.AddComponent(googol.NewShape(), "circle")

	// Aggiunge un componente Position (con coordinate)
	entity.AddComponent(googol.NewPosition(), googol.Coordinate{10.0, 20.0})

	// Stampa le informazioni dell'entità
	fmt.Printf("Entity ID: %x\n", entity.Id())
	fmt.Printf("Numero di componenti: %d\n", entity.LenComponent())

	// Itera e stampa ogni componente
	for i := 0; i < entity.LenComponent(); i++ {
		comp := entity.GetComponent(i)
		fmt.Printf("Componente %d: %T %+v\n", i, comp, comp)
		// Mostra anche i dati interni tramite il metodo Get()
		if comp != nil {
			fmt.Printf("	Dati: %v\n", comp.Get())
		}
	}

}
