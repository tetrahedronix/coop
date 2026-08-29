package ecs

import (
"testing"
)

// TestCoordinateType verifies that Coordinate is correctly defined as a 2D array.
func TestCoordinateType(t *testing.T) {
var coord Coordinate
if len(coord) != 2 {
t.Errorf("Coordinate length = %d, want 2", len(coord))
}
if coord[0] != 0.0 || coord[1] != 0.0 {
t.Errorf("Coordinate default values = [%f, %f], want [0.0, 0.0]", coord[0], coord[1])
}
coord = Coordinate{10.5, 20.7}
if coord[0] != 10.5 || coord[1] != 20.7 {
t.Errorf("Coordinate values = [%f, %f], want [10.5, 20.7]", coord[0], coord[1])
}
}

// TestPositionAdd tests the Add method of the Position struct.
func TestPositionAdd(t *testing.T) {
pos := &Position{}
testCoord := Coordinate{100.0, 200.0}
pos.Add(testCoord)
result := pos.Get().(Coordinate)
if result[0] != testCoord[0] || result[1] != testCoord[1] {
t.Errorf("Position.Get() = [%f, %f], want [%f, %f]", result[0], result[1], testCoord[0], testCoord[1])
}
pos.Add(Coordinate{0, 0})
result = pos.Get().(Coordinate)
if result[0] != 0 || result[1] != 0 {
t.Errorf("Position.Get() = [%f, %f], want [0, 0] after adding zero", result[0], result[1])
}
negativeCoord := Coordinate{-50.5, -75.3}
pos.Add(negativeCoord)
result = pos.Get().(Coordinate)
if result[0] != negativeCoord[0] || result[1] != negativeCoord[1] {
t.Errorf("Position.Get() = [%f, %f], want [%f, %f]", result[0], result[1], negativeCoord[0], negativeCoord[1])
}
}

// TestPositionCopy tests the Copy method of the Position struct.
func TestPositionCopy(t *testing.T) {
srcPos := &Position{}
srcPos.Add(Coordinate{123.45, 678.90})
dstPos := &Position{}
dstPos.Copy(srcPos)
srcResult := srcPos.Get().(Coordinate)
dstResult := dstPos.Get().(Coordinate)
if dstResult[0] != srcResult[0] || dstResult[1] != srcResult[1] {
t.Errorf("DstPos.Get() = [%f, %f], want [%f, %f]", dstResult[0], dstResult[1], srcResult[0], srcResult[1])
}
srcPos.Add(Coordinate{999.99, 888.88})
srcResult = srcPos.Get().(Coordinate)
dstResult = dstPos.Get().(Coordinate)
if dstResult[0] == srcResult[0] || dstResult[1] == srcResult[1] {
t.Error("DstPos coordinate was modified when SrcPos coordinate changed")
}
srcPos.Add(Coordinate{0, 0})
dstPos.Copy(srcPos)
dstResult = dstPos.Get().(Coordinate)
if dstResult[0] != 0 || dstResult[1] != 0 {
t.Errorf("DstPos.Get() = [%f, %f], want [0, 0]", dstResult[0], dstResult[1])
}
}

// TestPositionGet tests the Get method of the Position struct.
func TestPositionGet(t *testing.T) {
pos := &Position{}
result := pos.Get()
coord, ok := result.(Coordinate)
if !ok {
t.Fatalf("Position.Get() returned non-Coordinate type: %T", result)
}
if coord[0] != 0.0 || coord[1] != 0.0 {
t.Errorf("Position.Get() = [%f, %f], want [0.0, 0.0]", coord[0], coord[1])
}
testCoord := Coordinate{42.0, 84.0}
pos.Add(testCoord)
result = pos.Get()
coord, ok = result.(Coordinate)
if !ok {
t.Fatalf("Position.Get() returned non-Coordinate type: %T", result)
}
if coord[0] != testCoord[0] || coord[1] != testCoord[1] {
t.Errorf("Position.Get() = [%f, %f], want [%f, %f]", coord[0], coord[1], testCoord[0], testCoord[1])
}
}

// TestPositionReset tests the Reset method of the Position struct.
func TestPositionReset(t *testing.T) {
pos := &Position{}
pos.Add(Coordinate{111.11, 222.22})
pos.Reset()
result := pos.Get().(Coordinate)
if result[0] != 0 || result[1] != 0 {
t.Errorf("Position.Get() = [%f, %f] after Reset(), want [0, 0]", result[0], result[1])
}
pos.Reset()
result = pos.Get().(Coordinate)
if result[0] != 0 || result[1] != 0 {
t.Errorf("Position.Get() = [%f, %f] after second Reset()", result[0], result[1])
}
}

// TestNewPosition tests the NewPosition constructor function.
func TestNewPosition(t *testing.T) {
component := NewPosition()
if component == nil {
t.Fatal("NewPosition() returned nil")
}
pos, ok := component.(*Position)
if !ok {
t.Fatalf("NewPosition() returned non-*Position type: %T", component)
}
result := pos.Get().(Coordinate)
if result[0] != 0 || result[1] != 0 {
t.Errorf("NewPosition().Get() = [%f, %f], want [0, 0]", result[0], result[1])
}
_ = component.Add
_ = component.Copy
_ = component.Get
_ = component.Reset
}

// TestSelectableAdd tests the Add method of the Selectable struct.
func TestSelectableAdd(t *testing.T) {
sel := &Selectable{}
sel.Add(true)
result := sel.Get().(bool)
if result != true {
t.Errorf("Selectable.Get() = %v, want true", result)
}
sel.Add(false)
result = sel.Get().(bool)
if result != false {
t.Errorf("Selectable.Get() = %v, want false", result)
}
sel.Add(true)
sel.Add(false)
sel.Add(true)
result = sel.Get().(bool)
if result != true {
t.Errorf("Selectable.Get() = %v after toggling, want true", result)
}
}

// TestSelectableCopy tests the Copy method of the Selectable struct.
func TestSelectableCopy(t *testing.T) {
srcSel := &Selectable{}
srcSel.Add(true)
dstSel := &Selectable{}
dstSel.Copy(srcSel)
srcResult := srcSel.Get().(bool)
dstResult := dstSel.Get().(bool)
if dstResult != srcResult {
t.Errorf("DstSel.Get() = %v, want %v", dstResult, srcResult)
}
srcSel.Add(false)
srcResult = srcSel.Get().(bool)
dstResult = dstSel.Get().(bool)
if dstResult == srcResult {
t.Error("DstSel.selected was modified when SrcSel.selected changed")
}
srcSel.Add(false)
dstSel.Copy(srcSel)
dstResult = dstSel.Get().(bool)
if dstResult != false {
t.Errorf("DstSel.Get() = %v, want false", dstResult)
}
}

// TestSelectableGet tests the Get method of the Selectable struct.
func TestSelectableGet(t *testing.T) {
sel := &Selectable{}
result := sel.Get()
selected, ok := result.(bool)
if !ok {
t.Fatalf("Selectable.Get() returned non-bool type: %T", result)
}
if selected != false {
t.Errorf("Selectable.Get() = %v, want false", selected)
}
sel.Add(true)
result = sel.Get()
selected, ok = result.(bool)
if !ok {
t.Fatalf("Selectable.Get() returned non-bool type: %T", result)
}
if selected != true {
t.Errorf("Selectable.Get() = %v, want true", selected)
}
}

// TestSelectableReset tests the Reset method of the Selectable struct.
func TestSelectableReset(t *testing.T) {
sel := &Selectable{}
sel.Add(true)
sel.Reset()
result := sel.Get().(bool)
if result != false {
t.Errorf("Selectable.Get() = %v after Reset(), want false", result)
}
sel.Reset()
result = sel.Get().(bool)
if result != false {
t.Errorf("Selectable.Get() = %v after second Reset()", result)
}
}

// TestNewSelectable tests the NewSelectable constructor function.
func TestNewSelectable(t *testing.T) {
component := NewSelectable()
if component == nil {
t.Fatal("NewSelectable() returned nil")
}
sel, ok := component.(*Selectable)
if !ok {
t.Fatalf("NewSelectable() returned non-*Selectable type: %T", component)
}
result := sel.Get().(bool)
if result != false {
t.Errorf("NewSelectable().Get() = %v, want false", result)
}
_ = component.Add
_ = component.Copy
_ = component.Get
_ = component.Reset
}

// TestShapeAdd tests the Add method of the Shape struct.
func TestShapeAdd(t *testing.T) {
shape := &Shape{}
testPrimitive := "circle"
shape.Add(testPrimitive)
result := shape.Get().(string)
if result != testPrimitive {
t.Errorf("Shape.Get() = %q, want %q", result, testPrimitive)
}
shape.Add("")
result = shape.Get().(string)
if result != "" {
t.Errorf("Shape.Get() = %q, want \"\"", result)
}
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
func TestShapeCopy(t *testing.T) {
srcShape := &Shape{}
srcShape.Add("circle")
dstShape := &Shape{}
dstShape.Copy(srcShape)
srcResult := srcShape.Get().(string)
dstResult := dstShape.Get().(string)
if dstResult != srcResult {
t.Errorf("DstShape.Get() = %q, want %q", dstResult, srcResult)
}
srcShape.Add("box")
srcResult = srcShape.Get().(string)
dstResult = dstShape.Get().(string)
if dstResult == srcResult {
t.Error("DstShape.primitive was modified when SrcShape.primitive changed")
}
srcShape.Add("")
dstShape.Copy(srcShape)
dstResult = dstShape.Get().(string)
if dstResult != "" {
t.Errorf("DstShape.Get() = %q, want \"\"", dstResult)
}
}

// TestShapeGet tests the Get method of the Shape struct.
func TestShapeGet(t *testing.T) {
shape := &Shape{}
result := shape.Get()
primitive, ok := result.(string)
if !ok {
t.Fatalf("Shape.Get() returned non-string type: %T", result)
}
if primitive != "" {
t.Errorf("Shape.Get() = %q, want \"\"", primitive)
}
testPrimitive := "triangle"
shape.Add(testPrimitive)
result = shape.Get()
primitive, ok = result.(string)
if !ok {
t.Fatalf("Shape.Get() returned non-string type: %T", result)
}
if primitive != testPrimitive {
t.Errorf("Shape.Get() = %q, want %q", primitive, testPrimitive)
}
}

// TestShapeReset tests the Reset method of the Shape struct.
func TestShapeReset(t *testing.T) {
shape := &Shape{}
shape.Add("hexagon")
shape.Reset()
result := shape.Get().(string)
if result != "" {
t.Errorf("Shape.Get() = %q after Reset(), want \"\"", result)
}
shape.Reset()
result = shape.Get().(string)
if result != "" {
t.Errorf("Shape.Get() = %q after second Reset()", result)
}
}

// TestNewShape tests the NewShape constructor function.
func TestNewShape(t *testing.T) {
component := NewShape()
if component == nil {
t.Fatal("NewShape() returned nil")
}
shape, ok := component.(*Shape)
if !ok {
t.Fatalf("NewShape() returned non-*Shape type: %T", component)
}
result := shape.Get().(string)
if result != "" {
t.Errorf("NewShape().Get() = %q, want \"\"", result)
}
_ = component.Add
_ = component.Copy
_ = component.Get
_ = component.Reset
}

// TestVelocityAdd tests the Add method of the Velocity struct.
func TestVelocityAdd(t *testing.T) {
vel := &Velocity{}
testVel := Velocity{0, Speed(10.5), Direction(45.0)}
vel.Add(testVel)
result := vel.Get().(Velocity)
if result.speed != testVel.speed || result.direction != testVel.direction {
t.Errorf("Velocity.Get() = {%f, %f}, want {%f, %f}", result.speed, result.direction, testVel.speed, testVel.direction)
}
vel.Add(Velocity{0, 0, 0})
result = vel.Get().(Velocity)
if result.speed != 0 || result.direction != 0 {
t.Errorf("Velocity.Get() = {%f, %f}, want {0, 0}", result.speed, result.direction)
}
negativeVel := Velocity{0, Speed(5.0), Direction(-90.0)}
vel.Add(negativeVel)
result = vel.Get().(Velocity)
if result.speed != negativeVel.speed || result.direction != negativeVel.direction {
t.Errorf("Velocity.Get() = {%f, %f}, want {%f, %f}", result.speed, result.direction, negativeVel.speed, negativeVel.direction)
}
}

// TestVelocityCopy tests the Copy method of the Velocity struct.
func TestVelocityCopy(t *testing.T) {
srcVel := &Velocity{}
srcVel.Add(Velocity{0, Speed(100.0), Direction(180.0)})
dstVel := &Velocity{}
dstVel.Copy(srcVel)
srcResult := srcVel.Get().(Velocity)
dstResult := dstVel.Get().(Velocity)
if dstResult.speed != srcResult.speed || dstResult.direction != srcResult.direction {
t.Errorf("DstVel.Get() = {%f, %f}, want {%f, %f}", dstResult.speed, dstResult.direction, srcResult.speed, srcResult.direction)
}
srcVel.Add(Velocity{0, Speed(999.0), Direction(270.0)})
srcResult = srcVel.Get().(Velocity)
dstResult = dstVel.Get().(Velocity)
if dstResult.speed == srcResult.speed || dstResult.direction == srcResult.direction {
t.Error("DstVel was modified when SrcVel changed")
}
srcVel.Add(Velocity{0, 0, 0})
dstVel.Copy(srcVel)
dstResult = dstVel.Get().(Velocity)
if dstResult.speed != 0 || dstResult.direction != 0 {
t.Errorf("DstVel.Get() = {%f, %f}, want {0, 0}", dstResult.speed, dstResult.direction)
}
}

// TestVelocityGet tests the Get method of the Velocity struct.
func TestVelocityGet(t *testing.T) {
vel := &Velocity{}
result := vel.Get()
returnedVel, ok := result.(Velocity)
if !ok {
t.Fatalf("Velocity.Get() returned non-Velocity type: %T", result)
}
if returnedVel.speed != 0 || returnedVel.direction != 0 {
t.Errorf("Velocity.Get() = {%f, %f}, want {0, 0}", returnedVel.speed, returnedVel.direction)
}
testVel := Velocity{0, Speed(25.5), Direction(135.0)}
vel.Add(testVel)
result = vel.Get()
returnedVel, ok = result.(Velocity)
if !ok {
t.Fatalf("Velocity.Get() returned non-Velocity type: %T", result)
}
if returnedVel.speed != testVel.speed || returnedVel.direction != testVel.direction {
t.Errorf("Velocity.Get() = {%f, %f}, want {%f, %f}", returnedVel.speed, returnedVel.direction, testVel.speed, testVel.direction)
}
}

// TestVelocityReset tests the Reset method of the Velocity struct.
func TestVelocityReset(t *testing.T) {
vel := &Velocity{}
vel.Add(Velocity{0, Speed(50.0), Direction(90.0)})
vel.Reset()
result := vel.Get().(Velocity)
if result.speed != 0 || result.direction != 0 {
t.Errorf("Velocity.Get() = {%f, %f} after Reset(), want {0, 0}", result.speed, result.direction)
}
vel.Reset()
result = vel.Get().(Velocity)
if result.speed != 0 || result.direction != 0 {
t.Errorf("Velocity.Get() = {%f, %f} after second Reset()", result.speed, result.direction)
}
}

// TestNewVelocity tests the NewVelocity constructor function.
func TestNewVelocity(t *testing.T) {
component := NewVelocity()
if component == nil {
t.Fatal("NewVelocity() returned nil")
}
vel, ok := component.(*Velocity)
if !ok {
t.Fatalf("NewVelocity() returned non-*Velocity type: %T", component)
}
result := vel.Get().(Velocity)
if result.speed != 0 || result.direction != 0 {
t.Errorf("NewVelocity().Get() = {%f, %f}, want {0, 0}", result.speed, result.direction)
}
_ = component.Add
_ = component.Copy
_ = component.Get
_ = component.Reset
}

// TestVelocityComponentInterface verifies that Velocity implements Component.
func TestVelocityComponentInterface(t *testing.T) {
var _ Component = &Velocity{}
var _ Component = NewVelocity()
}

// TestBasicComponentsEdgeCases tests edge cases for basic components.
func TestBasicComponentsEdgeCases(t *testing.T) {
pos := &Position{}
largeCoord := Coordinate{1e10, -1e10}
pos.Add(largeCoord)
result := pos.Get().(Coordinate)
if result[0] != largeCoord[0] || result[1] != largeCoord[1] {
t.Errorf("Position failed with large values: got [%f, %f]", result[0], result[1])
}
vel := &Velocity{}
smallVel := Velocity{0, Speed(0.0001), Direction(0.0001)}
vel.Add(smallVel)
velResult := vel.Get().(Velocity)
if velResult.speed != smallVel.speed || velResult.direction != smallVel.direction {
t.Errorf("Velocity failed with small values: got {%f, %f}", velResult.speed, velResult.direction)
}
shape := &Shape{}
specialName := "shape-with_special.chars123"
shape.Add(specialName)
shapeResult := shape.Get().(string)
if shapeResult != specialName {
t.Errorf("Shape failed with special characters: got %q", shapeResult)
}
pos.Add(Coordinate{1.0, 2.0})
pos.Copy(pos)
result = pos.Get().(Coordinate)
if result[0] != 1.0 || result[1] != 2.0 {
t.Errorf("Position self-copy failed: got [%f, %f]", result[0], result[1])
}
vel.Add(Velocity{0, Speed(3.0), Direction(4.0)})
vel.Copy(vel)
velResult = vel.Get().(Velocity)
if velResult.speed != 3.0 || velResult.direction != 4.0 {
t.Errorf("Velocity self-copy failed: got {%f, %f}", velResult.speed, velResult.direction)
}
}
