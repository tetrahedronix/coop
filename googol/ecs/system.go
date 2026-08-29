package ecs

type System interface {
	// Match restituisce true se il sistema deve processare questa entità
	Match(e *Entity) bool
	// Process elabora l'entità (lettura dal past, scrittura nel future)
	Process(e *Entity)
}

// Systems have logic
// func NewSystem() {
// Create a new ring of size 5
// r := ring.New(5)
//
// Get the length of the ring
// n := r.Len()
//
// Initialize the ring with some integer values
// for i := 0; i < n; i++ {
// 	r.Value = "hello"
// 	r = r.Next()
// }
//
// Iterate through the ring and print its contents
// r.Do(func(p any) {
// 	fmt.Println(p.(string))
// })
// }
