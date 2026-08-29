package ecs

import (
	"testing"
)

// TestCoordinateType verifies that Coordinate is correctly defined as a 2D array.
// Coordinate represents a 2D position with two float64 values [x, y].
func TestCoordinateType(t *testing.T) {
	// Create a zero-value coordinate
	var coord Coordinate

	// Verify it has exactly 2 elements
	if len(coord) != 2 {
		t.Errorf("Coordinate length = %d, want 2", len(coord))
	}

	// Verify default values are zero
	if coord[0] != 0.0 || coord[1] != 0.0 {
		t.Errorf("Coordinate default values = [%f, %f], want [0.0, 0.0]", coord[0], coord[1])
	}

	// Test creating a coordinate with specific values
	coord = Coordinate{10.5, 20.7}
	if coord[0] != 10.5 || coord[1] != 20.7 {
		t.Errorf("Coordinate values = [%f, %f], want [10.5, 20.7]", coord[0], coord[1])
	}
}

// TestPositionAdd tests the Add method of the Position struct.
// The Add method should set the Coordinate field from a Coordinate value passed as interface{}.
func TestPositionAdd(t *testing.T) {
	// Create a new Position instance
	pos := &Position{}

	// Test adding a valid Coordinate
	testCoord := Coordinate{100.0, 200.0}
	pos.Add(testCoord)

	// Verify the Coordinate was set correctly
	if pos.Coordinate[0] != testCoord[0] || pos.Coordinate[1] != testCoord[1] {
		t.Errorf("Position.Coordinate = [%f, %f], want [%f, %f]",
			pos.Coordinate[0], pos.Coordinate[1], testCoord[0], testCoord[1])
	}

	// Test adding zero coordinate
	pos.Add(Coordinate{0, 0})
	if pos.Coordinate[0] != 0 || pos.Coordinate[1] != 0 {
		t.Errorf("Position.Coordinate = [%f, %f], want [0, 0] after adding zero", pos.Coordinate[0], pos.Coordinate[1])
	}

	// Test adding negative coordinates
	negativeCoord := Coordinate{-50.5, -75.3}
	pos.Add(negativeCoord)
	if pos.Coordinate[0] != negativeCoord[0] || pos.Coordinate[1] != negativeCoord[1] {
		t.Errorf("Position.Coordinate = [%f, %f], want [%f, %f]",
			pos.Coordinate[0], pos.Coordinate[1], negativeCoord[0], negativeCoord[1])
	}
}

// TestPositionCopy tests the Copy method of the Position struct.
// The Copy method should copy the Coordinate from another Position instance.
func TestPositionCopy(t *testing.T) {
	// Create source and destination positions
	srcPos := &Position{Coordinate: Coordinate{123.45, 678.90}}
	dstPos := &Position{}

	// Copy from source to destination
	dstPos.Copy(srcPos)

	// Verify the Coordinate was copied correctly
	if dstPos.Coordinate[0] != srcPos.Coordinate[0] || dstPos.Coordinate[1] != srcPos.Coordinate[1] {
		t.Errorf("DstPos.Coordinate = [%f, %f], want [%f, %f] (SrcPos.Coordinate)",
			dstPos.Coordinate[0], dstPos.Coordinate[1], srcPos.Coordinate[0], srcPos.Coordinate[1])
	}

	// Verify that modifying source doesn't affect destination (deep copy because Coordinate is a value type)
	srcPos.Coordinate = Coordinate{999.99, 888.88}
	if dstPos.Coordinate[0] == srcPos.Coordinate[0] || dstPos.Coordinate[1] == srcPos.Coordinate[1] {
		t.Error("DstPos.Coordinate was modified when SrcPos.Coordinate changed, copy should be independent")
	}

	// Test copying zero coordinate
	srcPos.Coordinate = Coordinate{0, 0}
	dstPos.Copy(srcPos)
	if dstPos.Coordinate[0] != 0 || dstPos.Coordinate[1] != 0 {
		t.Errorf("DstPos.Coordinate = [%f, %f], want [0, 0] after copying zero", dstPos.Coordinate[0], dstPos.Coordinate[1])
	}
}

// TestPositionGet tests the Get method of the Position struct.
// The Get method should return the Coordinate as an interface{}.
func TestPositionGet(t *testing.T) {
	pos := &Position{}

	// Test Get with initial zero value
	result := pos.Get()
	coord, ok := result.(Coordinate)
	if !ok {
		t.Fatalf("Position.Get() returned non-Coordinate type: %T", result)
	}
	if coord[0] != 0.0 || coord[1] != 0.0 {
		t.Errorf("Position.Get() = [%f, %f], want [0.0, 0.0] (initial value)", coord[0], coord[1])
	}

	// Test Get after setting a specific coordinate
	testCoord := Coordinate{42.0, 84.0}
	pos.Coordinate = testCoord
	result = pos.Get()

	// Type assert the result back to Coordinate
	coord, ok = result.(Coordinate)
	if !ok {
		t.Fatalf("Position.Get() returned non-Coordinate type: %T", result)
	}

	if coord[0] != testCoord[0] || coord[1] != testCoord[1] {
		t.Errorf("Position.Get() = [%f, %f], want [%f, %f]", coord[0], coord[1], testCoord[0], testCoord[1])
	}
}

// TestPositionReset tests the Reset method of the Position struct.
// The Reset method should set the Coordinate field to [0, 0].
func TestPositionReset(t *testing.T) {
	pos := &Position{Coordinate: Coordinate{111.11, 222.22}}

	// Reset the position
	pos.Reset()

	// Verify Coordinate is now [0, 0]
	if pos.Coordinate[0] != 0 || pos.Coordinate[1] != 0 {
		t.Errorf("Position.Coordinate = [%f, %f] after Reset(), want [0, 0]", pos.Coordinate[0], pos.Coordinate[1])
	}

	// Verify Reset on already zero coordinate doesn't cause issues
	pos.Reset()
	if pos.Coordinate[0] != 0 || pos.Coordinate[1] != 0 {
		t.Errorf("Position.Coordinate = [%f, %f] after second Reset(), want [0, 0]", pos.Coordinate[0], pos.Coordinate[1])
	}
}

// TestNewPosition tests the NewPosition constructor function.
// NewPosition should return a new Position instance implementing the Component interface.
func TestNewPosition(t *testing.T) {
	// Call the constructor
	component := NewPosition()

	// Verify the returned value is not nil
	if component == nil {
		t.Fatal("NewPosition() returned nil")
	}

	// Type assert to *Position
	pos, ok := component.(*Position)
	if !ok {
		t.Fatalf("NewPosition() returned non-*Position type: %T", component)
	}

	// Verify the Coordinate is initialized to zero (default value)
	if pos.Coordinate[0] != 0 || pos.Coordinate[1] != 0 {
		t.Errorf("NewPosition().Coordinate = [%f, %f], want [0, 0] (default initialization)", pos.Coordinate[0], pos.Coordinate[1])
	}

	// Verify it implements the Component interface by checking required methods exist
	_ = component.Add
	_ = component.Copy
	_ = component.Get
	_ = component.Reset
}

// TestSelectableAdd tests the Add method of the Selectable struct.
// The Add method should set the Selected field from a bool value passed as interface{}.
func TestSelectableAdd(t *testing.T) {
	// Create a new Selectable instance
	sel := &Selectable{}

	// Test adding true value
	sel.Add(true)
	if sel.Selected != true {
		t.Errorf("Selectable.Selected = %v, want true", sel.Selected)
	}

	// Test adding false value
	sel.Add(false)
	if sel.Selected != false {
		t.Errorf("Selectable.Selected = %v, want false", sel.Selected)
	}

	// Test toggling multiple times
	sel.Add(true)
	sel.Add(false)
	sel.Add(true)
	if sel.Selected != true {
		t.Errorf("Selectable.Selected = %v after toggling, want true", sel.Selected)
	}
}

// TestSelectableCopy tests the Copy method of the Selectable struct.
// The Copy method should copy the Selected field from another Selectable instance.
func TestSelectableCopy(t *testing.T) {
	// Create source and destination selectables
	srcSel := &Selectable{Selected: true}
	dstSel := &Selectable{Selected: false}

	// Copy from source to destination
	dstSel.Copy(srcSel)

	// Verify the Selected field was copied correctly
	if dstSel.Selected != srcSel.Selected {
		t.Errorf("DstSel.Selected = %v, want %v (SrcSel.Selected)", dstSel.Selected, srcSel.Selected)
	}

	// Verify that modifying source doesn't affect destination
	srcSel.Selected = false
	if dstSel.Selected == srcSel.Selected {
		t.Error("DstSel.Selected was modified when SrcSel.Selected changed, copy should be independent")
	}

	// Test copying false value
	srcSel.Selected = false
	dstSel.Copy(srcSel)
	if dstSel.Selected != false {
		t.Errorf("DstSel.Selected = %v, want false after copying false", dstSel.Selected)
	}
}

// TestSelectableGet tests the Get method of the Selectable struct.
// The Get method should return the Selected field as an interface{}.
func TestSelectableGet(t *testing.T) {
	sel := &Selectable{}

	// Test Get with initial false value
	result := sel.Get()
	selected, ok := result.(bool)
	if !ok {
		t.Fatalf("Selectable.Get() returned non-bool type: %T", result)
	}
	if selected != false {
		t.Errorf("Selectable.Get() = %v, want false (initial value)", selected)
	}

	// Test Get after setting to true
	sel.Selected = true
	result = sel.Get()

	// Type assert the result back to bool
	selected, ok = result.(bool)
	if !ok {
		t.Fatalf("Selectable.Get() returned non-bool type: %T", result)
	}

	if selected != true {
		t.Errorf("Selectable.Get() = %v, want true", selected)
	}
}

// TestSelectableReset tests the Reset method of the Selectable struct.
// The Reset method should set the Selected field to false.
func TestSelectableReset(t *testing.T) {
	sel := &Selectable{Selected: true}

	// Reset the selectable
	sel.Reset()

	// Verify Selected is now false
	if sel.Selected != false {
		t.Errorf("Selectable.Selected = %v after Reset(), want false", sel.Selected)
	}

	// Verify Reset on already false doesn't cause issues
	sel.Reset()
	if sel.Selected != false {
		t.Errorf("Selectable.Selected = %v after second Reset(), want false", sel.Selected)
	}
}

// TestNewSelectable tests the NewSelectable constructor function.
// NewSelectable should return a new Selectable instance implementing the Component interface.
func TestNewSelectable(t *testing.T) {
	// Call the constructor
	component := NewSelectable()

	// Verify the returned value is not nil
	if component == nil {
		t.Fatal("NewSelectable() returned nil")
	}

	// Type assert to *Selectable
	sel, ok := component.(*Selectable)
	if !ok {
		t.Fatalf("NewSelectable() returned non-*Selectable type: %T", component)
	}

	// Verify the Selected field is initialized to false (default value)
	if sel.Selected != false {
		t.Errorf("NewSelectable().Selected = %v, want false (default initialization)", sel.Selected)
	}

	// Verify it implements the Component interface by checking required methods exist
	_ = component.Add
	_ = component.Copy
	_ = component.Get
	_ = component.Reset
}

// TestShapeAdd tests the Add method of the Shape struct.
// The Add method should set the primitive field from a string value passed as interface{}.
func TestShapeAdd(t *testing.T) {
	// Create a new Shape instance
	shape := &Shape{}

	// Test adding a valid string
	testPrimitive := "circle"
	shape.Add(testPrimitive)

	// Verify the primitive was set correctly
	if shape.primitive != testPrimitive {
		t.Errorf("Shape.primitive = %q, want %q", shape.primitive, testPrimitive)
	}

	// Test adding empty string
	shape.Add("")
	if shape.primitive != "" {
		t.Errorf("Shape.primitive = %q, want \"\" after adding empty string", shape.primitive)
	}

	// Test adding different shape types
	shape.Add("box")
	if shape.primitive != "box" {
		t.Errorf("Shape.primitive = %q, want \"box\"", shape.primitive)
	}

	shape.Add("polygon")
	if shape.primitive != "polygon" {
		t.Errorf("Shape.primitive = %q, want \"polygon\"", shape.primitive)
	}
}

// TestShapeCopy tests the Copy method of the Shape struct.
// The Copy method should copy the primitive field from another Shape instance.
func TestShapeCopy(t *testing.T) {
	// Create source and destination shapes
	srcShape := &Shape{}
	srcShape.Add("circle")
	dstShape := &Shape{}

	// Copy from source to destination
	dstShape.Copy(srcShape)

	// Verify the primitive was copied correctly
	if dstShape.primitive != srcShape.primitive {
		t.Errorf("DstShape.primitive = %q, want %q (SrcShape.primitive)", dstShape.primitive, srcShape.primitive)
	}

	// Verify that modifying source doesn't affect destination (strings are immutable in Go, but we verify independence)
	srcShape.Add("box")
	if dstShape.primitive == srcShape.primitive {
		t.Error("DstShape.primitive was modified when SrcShape.primitive changed, copy should be independent")
	}

	// Test copying empty string
	srcShape.Add("")
	dstShape.Copy(srcShape)
	if dstShape.primitive != "" {
		t.Errorf("DstShape.primitive = %q, want \"\" after copying empty string", dstShape.primitive)
	}
}

// TestShapeGet tests the Get method of the Shape struct.
// The Get method should return the primitive field as an interface{}.
func TestShapeGet(t *testing.T) {
	shape := &Shape{}

	// Test Get with initial empty value
	result := shape.Get()
	primitive, ok := result.(string)
	if !ok {
		t.Fatalf("Shape.Get() returned non-string type: %T", result)
	}
	if primitive != "" {
		t.Errorf("Shape.Get() = %q, want \"\" (initial value)", primitive)
	}

	// Test Get after setting a specific primitive
	testPrimitive := "triangle"
	shape.Add(testPrimitive)
	result = shape.Get()

	// Type assert the result back to string
	primitive, ok = result.(string)
	if !ok {
		t.Fatalf("Shape.Get() returned non-string type: %T", result)
	}

	if primitive != testPrimitive {
		t.Errorf("Shape.Get() = %q, want %q", primitive, testPrimitive)
	}
}

// TestShapeReset tests the Reset method of the Shape struct.
// The Reset method should set the primitive field to empty string.
func TestShapeReset(t *testing.T) {
	shape := &Shape{}
	shape.Add("hexagon")

	// Reset the shape
	shape.Reset()

	// Verify primitive is now empty string
	if shape.primitive != "" {
		t.Errorf("Shape.primitive = %q after Reset(), want \"\"", shape.primitive)
	}

	// Verify Reset on already empty primitive doesn't cause issues
	shape.Reset()
	if shape.primitive != "" {
		t.Errorf("Shape.primitive = %q after second Reset(), want \"\"", shape.primitive)
	}
}

// TestNewShape tests the NewShape constructor function.
// NewShape should return a new Shape instance implementing the Component interface.
func TestNewShape(t *testing.T) {
	// Call the constructor
	component := NewShape()

	// Verify the returned value is not nil
	if component == nil {
		t.Fatal("NewShape() returned nil")
	}

	// Type assert to *Shape
	shape, ok := component.(*Shape)
	if !ok {
		t.Fatalf("NewShape() returned non-*Shape type: %T", component)
	}

	// Verify the primitive field is initialized to empty string (default value)
	if shape.primitive != "" {
		t.Errorf("NewShape().primitive = %q, want \"\" (default initialization)", shape.primitive)
	}

	// Verify it implements the Component interface by checking required methods exist
	_ = component.Add
	_ = component.Copy
	_ = component.Get
	_ = component.Reset
}

// TestVelocityAdd tests the Add method of the Velocity struct.
// The Add method should set Speed and Direction fields from a Velocity value passed as interface{}.
func TestVelocityAdd(t *testing.T) {
	// Create a new Velocity instance
	vel := &Velocity{}

	// Test adding a valid Velocity
	testVel := Velocity{Speed: 10.5, Direction: 45.0}
	vel.Add(testVel)

	// Verify the fields were set correctly
	if vel.Speed != testVel.Speed || vel.Direction != testVel.Direction {
		t.Errorf("Velocity = {Speed: %f, Direction: %f}, want {Speed: %f, Direction: %f}",
			vel.Speed, vel.Direction, testVel.Speed, testVel.Direction)
	}

	// Test adding zero velocity
	vel.Add(Velocity{0, 0})
	if vel.Speed != 0 || vel.Direction != 0 {
		t.Errorf("Velocity = {Speed: %f, Direction: %f}, want {Speed: 0, Direction: 0} after adding zero", vel.Speed, vel.Direction)
	}

	// Test adding negative direction (e.g., moving backwards or in opposite direction)
	negativeVel := Velocity{Speed: 5.0, Direction: -90.0}
	vel.Add(negativeVel)
	if vel.Speed != negativeVel.Speed || vel.Direction != negativeVel.Direction {
		t.Errorf("Velocity = {Speed: %f, Direction: %f}, want {Speed: %f, Direction: %f}",
			vel.Speed, vel.Direction, negativeVel.Speed, negativeVel.Direction)
	}
}

// TestVelocityCopy tests the Copy method of the Velocity struct.
// The Copy method should copy Speed and Direction from another Velocity instance.
func TestVelocityCopy(t *testing.T) {
	// Create source and destination velocities
	srcVel := &Velocity{Speed: 100.0, Direction: 180.0}
	dstVel := &Velocity{}

	// Copy from source to destination
	dstVel.Copy(srcVel)

	// Verify the fields were copied correctly
	if dstVel.Speed != srcVel.Speed || dstVel.Direction != srcVel.Direction {
		t.Errorf("DstVel = {Speed: %f, Direction: %f}, want {Speed: %f, Direction: %f} (SrcVel)",
			dstVel.Speed, dstVel.Direction, srcVel.Speed, srcVel.Direction)
	}

	// Verify that modifying source doesn't affect destination
	srcVel.Speed = 999.0
	srcVel.Direction = 270.0
	if dstVel.Speed == srcVel.Speed || dstVel.Direction == srcVel.Direction {
		t.Error("DstVel was modified when SrcVel changed, copy should be independent")
	}

	// Test copying zero velocity
	srcVel.Speed = 0
	srcVel.Direction = 0
	dstVel.Copy(srcVel)
	if dstVel.Speed != 0 || dstVel.Direction != 0 {
		t.Errorf("DstVel = {Speed: %f, Direction: %f}, want {Speed: 0, Direction: 0} after copying zero", dstVel.Speed, dstVel.Direction)
	}
}

// TestVelocityGet tests the Get method of the Velocity struct.
// The Get method should return a Velocity value (not pointer) as an interface{}.
func TestVelocityGet(t *testing.T) {
	vel := &Velocity{}

	// Test Get with initial zero values
	result := vel.Get()
	returnedVel, ok := result.(Velocity)
	if !ok {
		t.Fatalf("Velocity.Get() returned non-Velocity type: %T", result)
	}
	if returnedVel.Speed != 0 || returnedVel.Direction != 0 {
		t.Errorf("Velocity.Get() = {Speed: %f, Direction: %f}, want {Speed: 0, Direction: 0} (initial value)",
			returnedVel.Speed, returnedVel.Direction)
	}

	// Test Get after setting specific values
	testVel := Velocity{Speed: 25.5, Direction: 135.0}
	vel.Speed = testVel.Speed
	vel.Direction = testVel.Direction
	result = vel.Get()

	// Type assert the result back to Velocity
	returnedVel, ok = result.(Velocity)
	if !ok {
		t.Fatalf("Velocity.Get() returned non-Velocity type: %T", result)
	}

	if returnedVel.Speed != testVel.Speed || returnedVel.Direction != testVel.Direction {
		t.Errorf("Velocity.Get() = {Speed: %f, Direction: %f}, want {Speed: %f, Direction: %f}",
			returnedVel.Speed, returnedVel.Direction, testVel.Speed, testVel.Direction)
	}
}

// TestVelocityReset tests the Reset method of the Velocity struct.
// The Reset method should set both Speed and Direction to zero.
func TestVelocityReset(t *testing.T) {
	vel := &Velocity{Speed: 50.0, Direction: 90.0}

	// Reset the velocity
	vel.Reset()

	// Verify both fields are now zero
	if vel.Speed != 0 || vel.Direction != 0 {
		t.Errorf("Velocity = {Speed: %f, Direction: %f} after Reset(), want {Speed: 0, Direction: 0}", vel.Speed, vel.Direction)
	}

	// Verify Reset on already zero values doesn't cause issues
	vel.Reset()
	if vel.Speed != 0 || vel.Direction != 0 {
		t.Errorf("Velocity = {Speed: %f, Direction: %f} after second Reset(), want {Speed: 0, Direction: 0}", vel.Speed, vel.Direction)
	}
}

// TestNewVelocity tests the NewVelocity constructor function.
// NewVelocity should return a new Velocity instance implementing the Component interface.
func TestNewVelocity(t *testing.T) {
	// Call the constructor
	component := NewVelocity()

	// Verify the returned value is not nil
	if component == nil {
		t.Fatal("NewVelocity() returned nil")
	}

	// Type assert to *Velocity
	vel, ok := component.(*Velocity)
	if !ok {
		t.Fatalf("NewVelocity() returned non-*Velocity type: %T", component)
	}

	// Verify the fields are initialized to zero (default value)
	if vel.Speed != 0 || vel.Direction != 0 {
		t.Errorf("NewVelocity() = {Speed: %f, Direction: %f}, want {Speed: 0, Direction: 0} (default initialization)",
			vel.Speed, vel.Direction)
	}

	// Verify it implements the Component interface by checking required methods exist
	_ = component.Add
	_ = component.Copy
	_ = component.Get
	_ = component.Reset
}

// TestPositionComponentInterface verifies that Position properly implements the Component interface.
// This ensures compile-time compliance with the interface contract.
func TestPositionComponentInterface(t *testing.T) {
	// This will fail to compile if Position doesn't implement Component
	var _ Component = &Position{}

	// Also verify NewPosition returns a Component
	var _ Component = NewPosition()
}

// TestSelectableComponentInterface verifies that Selectable properly implements the Component interface.
func TestSelectableComponentInterface(t *testing.T) {
	// This will fail to compile if Selectable doesn't implement Component
	var _ Component = &Selectable{}

	// Also verify NewSelectable returns a Component
	var _ Component = NewSelectable()
}

// TestShapeComponentInterface verifies that Shape properly implements the Component interface.
func TestShapeComponentInterface(t *testing.T) {
	// This will fail to compile if Shape doesn't implement Component
	var _ Component = &Shape{}

	// Also verify NewShape returns a Component
	var _ Component = NewShape()
}

// TestVelocityComponentInterface verifies that Velocity properly implements the Component interface.
func TestVelocityComponentInterface(t *testing.T) {
	// This will fail to compile if Velocity doesn't implement Component
	var _ Component = &Velocity{}

	// Also verify NewVelocity returns a Component
	var _ Component = NewVelocity()
}

// TestBasicComponentsEdgeCases tests edge cases and boundary conditions for basic components.
func TestBasicComponentsEdgeCases(t *testing.T) {
	// Test Position with very large coordinate values
	pos := &Position{}
	largeCoord := Coordinate{1e10, -1e10}
	pos.Add(largeCoord)
	if pos.Coordinate[0] != largeCoord[0] || pos.Coordinate[1] != largeCoord[1] {
		t.Errorf("Position failed with large values: got [%f, %f], want [%f, %f]",
			pos.Coordinate[0], pos.Coordinate[1], largeCoord[0], largeCoord[1])
	}

	// Test Velocity with very small floating point values
	vel := &Velocity{}
	smallVel := Velocity{Speed: 0.0001, Direction: 0.0001}
	vel.Add(smallVel)
	if vel.Speed != smallVel.Speed || vel.Direction != smallVel.Direction {
		t.Errorf("Velocity failed with small values: got {Speed: %f, Direction: %f}, want {Speed: %f, Direction: %f}",
			vel.Speed, vel.Direction, smallVel.Speed, smallVel.Direction)
	}

	// Test Shape with special characters in primitive name
	shape := &Shape{}
	specialName := "shape-with_special.chars123"
	shape.Add(specialName)
	if shape.primitive != specialName {
		t.Errorf("Shape failed with special characters: got %q, want %q", shape.primitive, specialName)
	}

	// Test self-copy for Position
	pos.Coordinate = Coordinate{1.0, 2.0}
	pos.Copy(pos)
	if pos.Coordinate[0] != 1.0 || pos.Coordinate[1] != 2.0 {
		t.Errorf("Position self-copy failed: got [%f, %f], want [1.0, 2.0]", pos.Coordinate[0], pos.Coordinate[1])
	}

	// Test self-copy for Velocity
	vel.Speed = 3.0
	vel.Direction = 4.0
	vel.Copy(vel)
	if vel.Speed != 3.0 || vel.Direction != 4.0 {
		t.Errorf("Velocity self-copy failed: got {Speed: %f, Direction: %f}, want {Speed: 3.0, Direction: 4.0}", vel.Speed, vel.Direction)
	}
}
