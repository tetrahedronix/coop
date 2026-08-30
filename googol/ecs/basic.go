package ecs

import "fmt"

// Per il momento solo mondi 2D
type Coordinate [2]float64

// The position of the entity on the screen
type Position struct {
	componentType ComponentTypeID
	coordinate    Coordinate
}

// Metodo Add: accetta esplicitamente Coordinate
func (p *Position) Add(data any) {
	// Type assertion su Coordinate (non su [2]float64)
	p.coordinate = data.(Coordinate)
}

// Copy performs a deep copy of coordinates from src (must be *Position) to p.
// It panics if src is not a *Position. The copy is independent because
// coordinates are a value array [2]float64.
func (p *Position) Copy(src any) {
	p.coordinate = src.(*Position).coordinate
}

// Metodo Get: restituisce Coordinate
func (p *Position) Get() any {
	return p.coordinate
}

func (p *Position) Reset() {
	p.coordinate = Coordinate{0, 0}
}

func (p *Position) TypeID() ComponentTypeID {
	return p.componentType
}

func NewPosition() Component {

	return &Position{}
}

type Selectable struct {
	componentType ComponentTypeID
	selected      bool
}

func (s *Selectable) Add(data any) {
	s.selected = data.(bool)
}

func (s *Selectable) Copy(src any) {
	s.selected = src.(*Selectable).selected
}

func (s *Selectable) Get() any {
	return s.selected
}

func (s *Selectable) Reset() {
	s.selected = false
}

func (s *Selectable) TypeID() ComponentTypeID {
	return s.componentType
}

func NewSelectable() Component {
	return &Selectable{}
}

// The type of shape: for example circle or box
type Shape struct {
	componentType ComponentTypeID
	primitive     string
}

func (s *Shape) Add(data any) {

	s.primitive = data.(string)
}

func (s *Shape) Copy(src any) {
	s.primitive = src.(*Shape).primitive
}

func (s *Shape) Get() any {
	return s.primitive
}

func (s *Shape) Reset() {
	s.primitive = ""
}

func (s *Shape) TypeID() ComponentTypeID {
	return s.componentType
}

func (s *Shape) String() string {
	return fmt.Sprintf("%s", s.primitive)
}

func NewShape() Component {

	return &Shape{}
}

type Speed float64

type Direction float64

// The speed and direction in which the entity moves
type Velocity struct {
	componentType ComponentTypeID
	speed         Speed
	direction     Direction
}

func (v *Velocity) Add(data any) {
	// Si aspetta un valore di tipo Velocity (non puntatore)
	d := data.(Velocity)
	v.speed = d.speed
	v.direction = d.direction
}

func (v *Velocity) Copy(src any) {
	s := src.(*Velocity)
	v.speed = s.speed
	v.direction = s.direction
}

func (v *Velocity) Get() any {
	return Velocity{speed: v.speed, direction: v.direction}
}

func (v *Velocity) Reset() {
	v.speed, v.direction = 0, 0
}

func (v *Velocity) TypeID() ComponentTypeID {
	return v.componentType
}

func NewVelocity() Component {
	return &Velocity{}
}
