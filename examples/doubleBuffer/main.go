// Esempio più complesso che dimostra il funzionamento del doppio buffer
// (past/future) e l'uso di GetWritable per modificare i componenti senza
// alterare il past.
package main

import (
	"fmt"

	"github.com/tetrahedronix/coop/googol"
)

// main esegue questi task:
// * Crea due entità  con posizioni iniziali.
// * Mostra le posizioni dal past (invariate)
// * Modifica le posizioni nel future usando GetWritable
// * Mostra il past è rimasto invariato mentre il future ha i nuovi valori

func main() {
	world := googol.NewWorld()

	// Creazione entità #e1
	e1 := world.CreateEntity()
	e1.AddComponent(googol.NewShape(), "circle")
	e1.AddComponent(googol.NewPosition(), googol.Coordinate{0.0, 0.0})

	// Creazione entità #e2
	e2 := world.CreateEntity()
	e2.AddComponent(googol.NewShape(), "square")
	e2.AddComponent(googol.NewPosition(), googol.Coordinate{5.0, 5.0})

	// Stampa posizioni iniziali (dal past)
	fmt.Println("== POSIZIONI INIZIALI (PAST) ===")
	printEntity("Entità 1", e1, 1)
	printEntity("Entità 2", e2, 1)

	// Modifica posizioni nel future
	fmt.Println("\n=== MODIFICA POSIZIONI NEL FUTURE ===")
	modifyPositionFuture(e1, 1, googol.Coordinate{1.0, 1.0})
	modifyPositionFuture(e2, 1, googol.Coordinate{6.0, 6.0})

	// Stampa past (invariato) vs future (modificato)
	fmt.Println("\n=== CONFRONTO PAST vs FUTURE ===")
	compareEntity("Entità 1", e1, 1)
	compareEntity("Entità 2", e2, 1)

	// Aggiunta di un nuovo componente solo nel future
	fmt.Println("\n=== AGGIUNTA COMPONENTE NEL FUTURE===")
	e1.AddFutureComponent(googol.NewShape(), "triangle") // aggiunge shape al future
	fmt.Printf("Entità 1: Future ha ora %d componenti (past ne ha %d)\n",
		e1.LenFutureComponent(), e1.LenComponent())
}

// Funzione helper per stampare un componente Position dal past
func printEntity(name string, e *googol.Entity, idx int) {
	comp := googol.GetPastComponent(e, idx)

	if comp != nil {
		pos := comp.(*googol.Position)
		fmt.Printf("%s: Position = %v\n", name, pos.Coordinate)
	} else {
		fmt.Printf("%s: nessun Position all'indice %d\n", name, idx)
	}
}

// Modifica un Position nel future (usando GetWritable)
func modifyPositionFuture(e *googol.Entity, idx int, newCoord googol.Coordinate) {
	w := googol.GetWritableComponent(e, idx)

	if w != nil {
		pos := w.(*googol.Position)
		pos.Coordinate = newCoord
		fmt.Printf(" Nuova posizione scritta nel future: %v\n", newCoord)
	} else {
		fmt.Printf("	Impossibile modificare: indice %d non valido\n", idx)
	}
}

// Confronta il componente all'indice idx: past vs future
func compareEntity(name string, e *googol.Entity, idx int) {
	pastComp := googol.GetPastComponent(e, idx)
	futureComp := googol.GetFutureComponent(e, idx)

	pastVal, futureVal := "nil", "nil"

	if pastComp != nil {
		pastVal = fmt.Sprintf("%v", pastComp.(*googol.Position).Coordinate)
	}

	if futureComp != nil {
		futureVal = fmt.Sprintf("%v", futureComp.(*googol.Position).Coordinate)
	}

	fmt.Printf("%s:\n", name)
	fmt.Printf("	Past 	%s\n", pastVal)
	fmt.Printf("	Future: %s\n", futureVal)
}
