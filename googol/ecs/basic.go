package ecs

type Primitive string

// The type of shape: for example circle or box
type Shape struct {
	componentType ComponentTypeID
	Primitive     Primitive
}

func NewShape() Component {

	return &Shape{}
}

func (s *Shape) Add(data any) {

	s.Primitive = data.(Primitive)
}

func (s *Shape) Copy(src any) {
	s.Primitive = src.(*Shape).Primitive
}

func (s *Shape) Get() any {
	return s.Primitive
}

func (s *Shape) Reset() {

}

type Coordinate [2]float64

// The position of the entity on the screen
type Position struct {
	ComponentType ComponentTypeID
	Coordinate    Coordinate
}

func NewPosition() Component {

	return &Position{}
}

// Metodo Add: accetta esplicitamente Coordinate
func (p *Position) Add(data any) {
	// Type assertion su Coordinate (non su [2]float64)
	p.Coordinate = data.(Coordinate)
}

// Copy performs a deep copy of coordinates from src (must be *Position) to p.
// It panics if src is not a *Position. The copy is independent because
// coordinates are a value array [2]float64.
func (p *Position) Copy(src any) {
	p.Coordinate = src.(*Position).Coordinate
}

// Metodo Get: restituisce Coordinate
func (p *Position) Get() any {
	return p.Coordinate
}

func (p *Position) Reset() {

}

// The speed and direction in which the entity moves
type Velocity struct {
	v, d float64
}
