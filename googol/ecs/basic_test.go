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
// The Add method should set the coordinate field from a Coordinate value passed as interface{}.
func TestPositionAdd(t *testing.T) {
	// Create a new Position instance
	pos := &Position{}

	// Test adding a valid Coordinate
	testCoord := Coordinate{100.0, 200.0}
	pos.Add(testCoord)

	// Verify the coordinate was set correctly using Get()
	result := pos.Get().(Coordinate)
	if result[0] != testCoord[0] || result[1] != testCoord[1] {
		t.Errorf("Position.Get() = [%f, %f], want [%f, %f]",
			result[0], result[1], testCoord[0], testCoord[1])
	}

	// Test adding zero coordinate
	pos.Add(Coordinate{0, 0})
	result = pos.Get().(Coordinate)
	if result[0] != 0 || result[1] != 0 {
		t.Errorf("Position.Get() = [%f, %f], want [0, 0] after adding zero", result[0], result[1])
	}

	// Test adding negative coordinates
	negativeCoord := Coordinate{-50.5, -75.3}
	pos.Add(negativeCoord)
	result = pos.Get().(Coordinate)
	if result[0] != negativeCoord[0] || result[1] != negativeCoord[1] {
		t.Errorf("Position.Get() = [%f, %f], want [%f, %f]",
			result[0], result[1], negativeCoord[0], negativeCoord[1])
	}
}

// TestPositionCopy tests the Copy method of the Position struct.
// The Copy method should copy the coordinate from another Position instance.
func TestPositionCopy(t *testing.T) {
	// Create source and destination positions
	srcPos := &Position{}
	srcPos.Add(Coordinate{123.45, 678.90})
	dstPos := &Position{}

	// Copy from source to destination
	dstPos.Copy(srcPos)

	// Verify the coordinate was copied correctly
	srcCoord := srcPos.Get().(Coordinate)
	dstCoord := dstPos.Get().(Coordinate)
	if dstCoord[0] != srcCoord[0] || dstCoord[1] != srcCoord[1] {
		t.Errorf("DstPos.Get() = [%f, %f], want [%f, %f] (SrcPos.Get())",
			dstCoord[0], dstCoord[1], srcCoord[0], srcCoord[1])
	}

	// Verify that modifying source doesn't affect destination (deep copy because Coordinate is a value type)
	srcPos.Add(Coordinate{999.99, 888.88})
	srcCoord = srcPos.Get().(Coordinate)
	dstCoord = dstPos.Get().(Coordinate)
	if dstCoord[0] == srcCoord[0] || dstCoord[1] == srcCoord[1] {
		t.Error("DstPos coordinate was modified when SrcPos changed, copy should be independent")
	}

	// Test copying zero coordinate
	srcPos.Add(Coordinate{0, 0})
	dstPos.Copy(srcPos)
	dstCoord = dstPos.Get().(Coordinate)
	if dstCoord[0] != 0 || dstCoord[1] != 0 {
		t.Errorf("DstPos.Get() = [%f, %f], want [0, 0] after copying zero", dstCoord[0], dstCoord[1])
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
	pos.Add(testCoord)
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
// The Reset method should set the coordinate field to [0, 0].
func TestPositionReset(t *testing.T) {
	pos := &Position{}
	pos.Add(Coordinate{111.11, 222.22})

	// Reset the position
	pos.Reset()

	// Verify coordinate is now [0, 0]
	result := pos.Get().(Coordinate)
	if result[0] != 0 || result[1] != 0 {
		t.Errorf("Position.Get() = [%f, %f] after Reset(), want [0, 0]", result[0], result[1])
	}

	// Verify Reset on already zero coordinate doesn't cause issues
	pos.Reset()
	result = pos.Get().(Coordinate)
	if result[0] != 0 || result[1] != 0 {
		t.Errorf("Position.Get() = [%f, %f] after second Reset(), want [0, 0]", result[0], result[1])
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

	// Verify the coordinate is initialized to zero (default value)
	result := pos.Get().(Coordinate)
	if result[0] != 0 || result[1] != 0 {
		t.Errorf("NewPosition().Get() = [%f, %f], want [0, 0] (default initialization)", result[0], result[1])
	}

	// Verify it implements the Component interface by checking required methods exist
	_ = component.Add
	_ = component.Copy
	_ = component.Get
	_ = component.Reset
}

// TestSelectableAdd tests the Add method of the Selectable struct.
// The Add method should set the selected field from a bool value passed as interface{}.
func TestSelectableAdd(t *testing.T) {
	// Create a new Selectable instance
	sel := &Selectable{}

	// Test adding true value
	sel.Add(true)
	result := sel.Get().(bool)
	if result != true {
		t.Errorf("Selectable.Get() = %v, want true", result)
	}

	// Test adding false value
	sel.Add(false)
	result = sel.Get().(bool)
	if result != false {
		t.Errorf("Selectable.Get() = %v, want false", result)
	}

	// Test toggling multiple times
	sel.Add(true)
	sel.Add(false)
	sel.Add(true)
	result = sel.Get().(bool)
	if result != true {
		t.Errorf("Selectable.Get() = %v after toggling, want true", result)
	}
}

// TestSelectableCopy tests the Copy method of the Selectable struct.
// The Copy method should copy the selected field from another Selectable instance.
func TestSelectableCopy(t *testing.T) {
	// Create source and destination selectables
	srcSel := &Selectable{}
	srcSel.Add(true)
	dstSel := &Selectable{}

	// Copy from source to destination
	dstSel.Copy(srcSel)

	// Verify the selected field was copied correctly
	srcResult := srcSel.Get().(bool)
	dstResult := dstSel.Get().(bool)
	if dstResult != srcResult {
		t.Errorf("DstSel.Get() = %v, want %v (SrcSel.Get())", dstResult, srcResult)
	}

	// Verify that modifying source doesn't affect destination
	srcSel.Add(false)
	srcResult = srcSel.Get().(bool)
	dstResult = dstSel.Get().(bool)
	if dstResult == srcResult {
		t.Error("DstSel was modified when SrcSel changed, copy should be independent")
	}

	// Test copying false value
	srcSel.Add(false)
	dstSel.Copy(srcSel)
	dstResult = dstSel.Get().(bool)
	if dstResult != false {
		t.Errorf("DstSel.Get() = %v, want false after copying false", dstResult)
	}
}

// TestSelectableGet tests the Get method of the Selectable struct.
// The Get method should return the selected field as an interface{}.
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
	sel.Add(true)
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
// The Reset method should set the selected field to false.
func TestSelectableReset(t *testing.T) {
	sel := &Selectable{}
	sel.Add(true)

	// Reset the selectable
	sel.Reset()

	// Verify selected is now false
	result := sel.Get().(bool)
	if result != false {
		t.Errorf("Selectable.Get() = %v after Reset(), want false", result)
	}

	// Verify Reset on already false doesn't cause issues
	sel.Reset()
	result = sel.Get().(bool)
	if result != false {
		t.Errorf("Selectable.Get() = %v after second Reset(), want false", result)
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

	// Verify the selected field is initialized to false (default value)
	result := sel.Get().(bool)
	if result != false {
		t.Errorf("NewSelectable().Get() = %v, want false (default initialization)", result)
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

	// Verify the primitive was set correctly using Get()
	result := shape.Get().(string)
	if result != testPrimitive {
		t.Errorf("Shape.Get() = %q, want %q", result, testPrimitive)
	}

	// Test adding empty string
	shape.Add("")
	result = shape.Get().(string)
	if result != "" {
		t.Errorf("Shape.Get() = %q, want \"\" after adding empty string", result)
	}

	// Test adding different shape types
	shape.Add("box")
	result = shape.Get().(string)
	if result != "box" {
		t.Errorf("Shape.Get() = %q, want \"box\"", result)
	}

	shape.Add("polygon")
	result = shape.Get().(string)
	if result != "polygon" {
		t.Errorf("Shape.Get() = %q, want \"polygon\"", result)
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
	srcResult := srcShape.Get().(string)
	dstResult := dstShape.Get().(string)
	if dstResult != srcResult {
		t.Errorf("DstShape.Get() = %q, want %q (SrcShape.Get())", dstResult, srcResult)
	}

	// Verify that modifying source doesn't affect destination (strings are immutable in Go, but we verify independence)
	srcShape.Add("box")
	srcResult = srcShape.Get().(string)
	dstResult = dstShape.Get().(string)
	if dstResult == srcResult {
		t.Error("DstShape was modified when SrcShape changed, copy should be independent")
	}

	// Test copying empty string
	srcShape.Add("")
	dstShape.Copy(srcShape)
	dstResult = dstShape.Get().(string)
	if dstResult != "" {
		t.Errorf("DstShape.Get() = %q, want \"\" after copying empty string", dstResult)
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
	result := shape.Get().(string)
	if result != "" {
		t.Errorf("Shape.Get() = %q after Reset(), want \"\"", result)
	}

	// Verify Reset on already empty primitive doesn't cause issues
	shape.Reset()
	result = shape.Get().(string)
	if result != "" {
		t.Errorf("Shape.Get() = %q after second Reset(), want \"\"", result)
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
	result := shape.Get().(string)
	if result != "" {
		t.Errorf("NewShape().Get() = %q, want \"\" (default initialization)", result)
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
	testVel := Velocity{speed: 10.5, direction: 45.0}
	vel.Add(testVel)

	// Verify the fields were set correctly using Get()
	result := vel.Get().(Velocity)
	if result.speed != testVel.speed || result.direction != testVel.direction {
		t.Errorf("Velocity.Get() = {Speed: %f, direction: %f}, want {Speed: %f, direction: %f}",
			result.speed, result.direction, testVel.speed, testVel.direction)
	}

	// Test adding zero velocity
	vel.Add(Velocity{speed: 0, direction: 0})
	result = vel.Get().(Velocity)
	if result.speed != 0 || result.direction != 0 {
		t.Errorf("Velocity.Get() = {Speed: %f, direction: %f}, want {Speed: 0, direction: 0} after adding zero", result.speed, result.direction)
	}

	// Test adding negative direction (e.g., moving backwards or in opposite direction)
	negativeVel := Velocity{speed: 5.0, direction: -90.0}
	vel.Add(negativeVel)
	result = vel.Get().(Velocity)
	if result.speed != negativeVel.speed || result.direction != negativeVel.direction {
		t.Errorf("Velocity.Get() = {Speed: %f, direction: %f}, want {Speed: %f, direction: %f}",
			result.speed, result.direction, negativeVel.speed, negativeVel.direction)
	}
}

// TestVelocityCopy tests the Copy method of the Velocity struct.
// The Copy method should copy Speed and Direction from another Velocity instance.
func TestVelocityCopy(t *testing.T) {
	// Create source and destination velocities
	srcVel := &Velocity{}
	srcVel.Add(Velocity{speed: 100.0, direction: 180.0})
	dstVel := &Velocity{}

	// Copy from source to destination
	dstVel.Copy(srcVel)

	// Verify the fields were copied correctly
	srcResult := srcVel.Get().(Velocity)
	dstResult := dstVel.Get().(Velocity)
	if dstResult.speed != srcResult.speed || dstResult.direction != srcResult.direction {
		t.Errorf("DstVel.Get() = {Speed: %f, direction: %f}, want {Speed: %f, direction: %f} (SrcVel.Get())",
			dstResult.speed, dstResult.direction, srcResult.speed, srcResult.direction)
	}

	// Verify that modifying source doesn't affect destination
	srcVel.Add(Velocity{speed: 999.0, direction: 270.0})
	srcResult = srcVel.Get().(Velocity)
	dstResult = dstVel.Get().(Velocity)
	if dstResult.speed == srcResult.speed || dstResult.direction == srcResult.direction {
		t.Error("DstVel was modified when SrcVel changed, copy should be independent")
	}

	// Test copying zero velocity
	srcVel.Add(Velocity{speed: 0, direction: 0})
	dstVel.Copy(srcVel)
	dstResult = dstVel.Get().(Velocity)
	if dstResult.speed != 0 || dstResult.direction != 0 {
		t.Errorf("DstVel.Get() = {Speed: %f, direction: %f}, want {Speed: 0, direction: 0} after copying zero", dstResult.speed, dstResult.direction)
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
	if returnedVel.speed != 0 || returnedVel.direction != 0 {
		t.Errorf("Velocity.Get() = {Speed: %f, direction: %f}, want {Speed: 0, direction: 0} (initial value)",
			returnedVel.speed, returnedVel.direction)
	}

	// Test Get after setting specific values
	testVel := Velocity{speed: 25.5, direction: 135.0}
	vel.Add(testVel)
	result = vel.Get()

	// Type assert the result back to Velocity
	returnedVel, ok = result.(Velocity)
	if !ok {
		t.Fatalf("Velocity.Get() returned non-Velocity type: %T", result)
	}

	if returnedVel.speed != testVel.speed || returnedVel.direction != testVel.direction {
		t.Errorf("Velocity.Get() = {Speed: %f, direction: %f}, want {Speed: %f, direction: %f}",
			returnedVel.speed, returnedVel.direction, testVel.speed, testVel.direction)
	}
}

// TestVelocityReset tests the Reset method of the Velocity struct.
// The Reset method should set both Speed and Direction to zero.
func TestVelocityReset(t *testing.T) {
	vel := &Velocity{}
	vel.Add(Velocity{speed: 50.0, direction: 90.0})

	// Reset the velocity
	vel.Reset()

	// Verify both fields are now zero
	result := vel.Get().(Velocity)
	if result.speed != 0 || result.direction != 0 {
		t.Errorf("Velocity.Get() = {Speed: %f, direction: %f} after Reset(), want {Speed: 0, direction: 0}", result.speed, result.direction)
	}

	// Verify Reset on already zero values doesn't cause issues
	vel.Reset()
	result = vel.Get().(Velocity)
	if result.speed != 0 || result.direction != 0 {
		t.Errorf("Velocity.Get() = {Speed: %f, direction: %f} after second Reset(), want {Speed: 0, direction: 0}", result.speed, result.direction)
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
	result := vel.Get().(Velocity)
	if result.speed != 0 || result.direction != 0 {
		t.Errorf("NewVelocity().Get() = {Speed: %f, direction: %f}, want {Speed: 0, direction: 0} (default initialization)",
			result.speed, result.direction)
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
	result := pos.Get().(Coordinate)
	if result[0] != largeCoord[0] || result[1] != largeCoord[1] {
		t.Errorf("Position failed with large values: got [%f, %f], want [%f, %f]",
			result[0], result[1], largeCoord[0], largeCoord[1])
	}

	// Test Velocity with very small floating point values
	vel := &Velocity{}
	smallVel := Velocity{speed: 0.0001, direction: 0.0001}
	vel.Add(smallVel)
	velResult := vel.Get().(Velocity)
	if velResult.speed != smallVel.speed || velResult.direction != smallVel.direction {
		t.Errorf("Velocity failed with small values: got {Speed: %f, direction: %f}, want {Speed: %f, direction: %f}",
			velResult.speed, velResult.direction, smallVel.speed, smallVel.direction)
	}

	// Test Shape with special characters in primitive name
	shape := &Shape{}
	specialName := "shape-with_special.chars123"
	shape.Add(specialName)
	shapeResult := shape.Get().(string)
	if shapeResult != specialName {
		t.Errorf("Shape failed with special characters: got %q, want %q", shapeResult, specialName)
	}

	// Test self-copy for Position
	pos.Add(Coordinate{1.0, 2.0})
	pos.Copy(pos)
	result = pos.Get().(Coordinate)
	if result[0] != 1.0 || result[1] != 2.0 {
		t.Errorf("Position self-copy failed: got [%f, %f], want [1.0, 2.0]", result[0], result[1])
	}

	// Test self-copy for Velocity
	vel.Add(Velocity{speed: 3.0, direction: 4.0})
	vel.Copy(vel)
	velResult = vel.Get().(Velocity)
	if velResult.speed != 3.0 || velResult.direction != 4.0 {
		t.Errorf("Velocity self-copy failed: got {Speed: %f, direction: %f}, want {Speed: 3.0, direction: 4.0}", velResult.speed, velResult.direction)
	}
}
