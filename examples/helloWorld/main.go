package main

import (
	"fmt"

	"github.com/tetrahedronix/coop/googol"
)

func main() {
	gw := googol.NewWorld()

	e := gw.CreateEntity()

	e.AddComponent(googol.NewShape(), "Hello World")

	fmt.Println(e)
}
