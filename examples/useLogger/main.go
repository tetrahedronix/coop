package main

import (
	"io"
	"log"
	"os"

	"github.com/tetrahedronix/coop/googol"
)

func main() {
	// Creazione del mondo
	world := googol.NewWorld()

	// Configurazione del logger: scrive su stdout con timestamp e prefisso
	world.Logger = log.New(os.Stdout, "GOOGOL: ", log.LstdFlags)

	// Creazione di un'entità (i log interni di CreateEntity verranno stampati)
	entity, err := world.CreateEntity()

	if err != nil {
		log.Fatal("Impossibile creare l'entità: flake.NextID ha fallito")
	}

	// Aggiunta di un componente Shape con dati
	entity.AddComponent(googol.NewShape(), "Cerchio")

	// Stampa dell'entità per vedere lo stato tramite String
	world.Logger.Println(entity)

	// Esempio di lettura di un componente (past)
	if entity.LenComponent() > 0 {
		comp := entity.GetComponent(0)
		if s, ok := comp.(*googol.Shape); ok {
			val := s.Get().(string)
			world.Logger.Printf("Il componente Shape dell'entità contiene: %s", val)
		}
	}

	// Disattiva il logger
	world.Logger = log.New(io.Discard, "", 0)
	// Non stampa nulla
	world.Logger.Println(world)
}
