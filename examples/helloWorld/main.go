package main

import (
	"fmt"
	"os"

	"github.com/tetrahedronix/coop/googol"
)

func main() {
	gw := googol.NewWorld()

	entity, err := gw.CreateEntity()

	if err != nil {
		os.Exit(1)
	}

	entity.AddComponent(googol.NewShape(), "Hello World")

	fmt.Println(entity, entity.GetComponent(0))
}
