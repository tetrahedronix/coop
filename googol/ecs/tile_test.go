package ecs

import (
	"testing"
)

// TestTileConstants verifies that the tile flag constants are correctly defined.
// These constants are used for bit manipulation of the GID (Global Tile ID).
// TileFlipHorizFlag: horizontal flip flag (bit 31)
// TileFlipVertFlag: vertical flip flag (bit 30)
// TileFlipDiagFlag: diagonal flip flag (bit 29)
// TileIDMask: mask to extract the actual tile ID (lower 29 bits)
func TestTileConstants(t *testing.T) {
	// Verify horizontal flip flag value (0x80000000 = 2147483648 in decimal)
	if TileFlipHorizFlag != 0x80000000 {
		t.Errorf("TileFlipHorizFlag = %x, want %x", TileFlipHorizFlag, 0x80000000)
	}

	// Verify vertical flip flag value (0x40000000 = 1073741824 in decimal)
	if TileFlipVertFlag != 0x40000000 {
		t.Errorf("TileFlipVertFlag = %x, want %x", TileFlipVertFlag, 0x40000000)
	}

	// Verify diagonal flip flag value (0x20000000 = 536870912 in decimal)
	if TileFlipDiagFlag != 0x20000000 {
		t.Errorf("TileFlipDiagFlag = %x, want %x", TileFlipDiagFlag, 0x20000000)
	}

	// Verify tile ID mask (0x1FFFFFFF = 536870911 in decimal, lower 29 bits)
	if TileIDMask != 0x1FFFFFFF {
		t.Errorf("TileIDMask = %x, want %x", TileIDMask, 0x1FFFFFFF)
	}
}

// TestTileAdd tests the Add method of the Tile struct.
// The Add method should set the GID field from a uint32 value passed as interface{}.
func TestTileAdd(t *testing.T) {
	// Create a new Tile instance
	tile := &Tile{}

	// Test adding a valid uint32 GID
	testGID := uint32(12345)
	tile.Add(testGID)

	// Verify the GID was set correctly
	if tile.GID != testGID {
		t.Errorf("Tile.GID = %d, want %d", tile.GID, testGID)
	}

	// Test adding zero GID
	tile.Add(uint32(0))
	if tile.GID != 0 {
		t.Errorf("Tile.GID = %d, want 0 after adding zero", tile.GID)
	}

	// Test adding maximum uint32 value
	maxGID := uint32(0xFFFFFFFF)
	tile.Add(maxGID)
	if tile.GID != maxGID {
		t.Errorf("Tile.GID = %d, want %d (max uint32)", tile.GID, maxGID)
	}
}

// TestTileCopy tests the Copy method of the Tile struct.
// The Copy method should copy the GID from another Tile instance.
func TestTileCopy(t *testing.T) {
	// Create source and destination tiles
	srcTile := &Tile{GID: 54321}
	dstTile := &Tile{}

	// Copy from source to destination
	dstTile.Copy(srcTile)

	// Verify the GID was copied correctly
	if dstTile.GID != srcTile.GID {
		t.Errorf("DstTile.GID = %d, want %d (SrcTile.GID)", dstTile.GID, srcTile.GID)
	}

	// Verify that modifying source doesn't affect destination (deep copy for primitive type)
	srcTile.GID = 99999
	if dstTile.GID == srcTile.GID {
		t.Error("DstTile.GID was modified when SrcTile.GID changed, copy should be independent")
	}

	// Test copying zero GID
	srcTile.GID = 0
	dstTile.Copy(srcTile)
	if dstTile.GID != 0 {
		t.Errorf("DstTile.GID = %d, want 0 after copying zero GID", dstTile.GID)
	}
}

// TestTileGet tests the Get method of the Tile struct.
// The Get method should return the raw GID as an interface{}.
func TestTileGet(t *testing.T) {
	tile := &Tile{}

	// Test Get with initial zero value
	result := tile.Get()
	if result != uint32(0) {
		t.Errorf("Tile.Get() = %v, want 0 (initial value)", result)
	}

	// Test Get after setting a specific GID
	testGID := uint32(67890)
	tile.GID = testGID
	result = tile.Get()

	// Type assert the result back to uint32
	gid, ok := result.(uint32)
	if !ok {
		t.Fatalf("Tile.Get() returned non-uint32 type: %T", result)
	}

	if gid != testGID {
		t.Errorf("Tile.Get() = %d, want %d", gid, testGID)
	}
}

// TestTileReset tests the Reset method of the Tile struct.
// The Reset method should set the GID field to zero.
func TestTileReset(t *testing.T) {
	tile := &Tile{GID: 11111}

	// Reset the tile
	tile.Reset()

	// Verify GID is now zero
	if tile.GID != 0 {
		t.Errorf("Tile.GID = %d after Reset(), want 0", tile.GID)
	}

	// Verify Reset on already zero GID doesn't cause issues
	tile.Reset()
	if tile.GID != 0 {
		t.Errorf("Tile.GID = %d after second Reset(), want 0", tile.GID)
	}
}

// TestNewTile tests the NewTile constructor function.
// NewTile should return a new Tile instance implementing the Component interface.
func TestNewTile(t *testing.T) {
	// Call the constructor
	component := NewTile()

	// Verify the returned value is not nil
	if component == nil {
		t.Fatal("NewTile() returned nil")
	}

	// Type assert to *Tile
	tile, ok := component.(*Tile)
	if !ok {
		t.Fatalf("NewTile() returned non-*Tile type: %T", component)
	}

	// Verify the GID is initialized to zero (default value)
	if tile.GID != 0 {
		t.Errorf("NewTile().GID = %d, want 0 (default initialization)", tile.GID)
	}

	// Verify it implements the Component interface by checking required methods exist
	_ = component.Add
	_ = component.Copy
	_ = component.Get
	_ = component.Reset
}

// TestTileWithFlags tests Tile GID manipulation with flip flags.
// This test demonstrates how to combine tile IDs with transformation flags.
func TestTileWithFlags(t *testing.T) {
	tile := NewTile().(*Tile)

	// Base tile ID (without any flags)
	baseTileID := uint32(100)

	// Test horizontal flip
	horizFlipped := baseTileID | TileFlipHorizFlag
	tile.Add(horizFlipped)
	if tile.GID != horizFlipped {
		t.Errorf("Horizontal flip GID = %x, want %x", tile.GID, horizFlipped)
	}

	// Extract the base ID using the mask to verify it's correct
	extractedID := tile.GID & TileIDMask
	if extractedID != baseTileID {
		t.Errorf("Extracted base ID = %d, want %d", extractedID, baseTileID)
	}

	// Test vertical flip
	tile.Reset()
	vertFlipped := baseTileID | TileFlipVertFlag
	tile.Add(vertFlipped)
	if (tile.GID & TileFlipVertFlag) == 0 {
		t.Error("Vertical flip flag not set")
	}

	// Test diagonal flip
	tile.Reset()
	diagFlipped := baseTileID | TileFlipDiagFlag
	tile.Add(diagFlipped)
	if (tile.GID & TileFlipDiagFlag) == 0 {
		t.Error("Diagonal flip flag not set")
	}

	// Test multiple flags combined
	tile.Reset()
	combinedFlags := baseTileID | TileFlipHorizFlag | TileFlipVertFlag | TileFlipDiagFlag
	tile.Add(combinedFlags)

	// Verify all flags are present
	if (tile.GID & TileFlipHorizFlag) == 0 {
		t.Error("Horizontal flip flag missing in combined flags")
	}
	if (tile.GID & TileFlipVertFlag) == 0 {
		t.Error("Vertical flip flag missing in combined flags")
	}
	if (tile.GID & TileFlipDiagFlag) == 0 {
		t.Error("Diagonal flip flag missing in combined flags")
	}

	// Verify base ID is still extractable
	extractedID = tile.GID & TileIDMask
	if extractedID != baseTileID {
		t.Errorf("Extracted base ID with combined flags = %d, want %d", extractedID, baseTileID)
	}
}

// TestTileComponentInterface verifies that Tile properly implements the Component interface.
// This ensures compile-time compliance with the interface contract.
func TestTileComponentInterface(t *testing.T) {
	// This will fail to compile if Tile doesn't implement Component
	var _ Component = &Tile{}

	// Also verify NewTile returns a Component
	var _ Component = NewTile()
}

// TestTileEdgeCases tests edge cases and boundary conditions for Tile operations.
func TestTileEdgeCases(t *testing.T) {
	tile := &Tile{}

	// Test with maximum tile ID (all 29 bits set, no flags)
	maxBaseID := uint32(TileIDMask)
	tile.Add(maxBaseID)
	if tile.GID != maxBaseID {
		t.Errorf("Max base ID failed: got %x, want %x", tile.GID, maxBaseID)
	}

	// Test that flags don't interfere with ID extraction
	tile.Reset()
	testID := uint32(0x12345678) // A random ID within the mask range
	fullGID := testID | TileFlipHorizFlag
	tile.Add(fullGID)

	extracted := tile.GID & TileIDMask
	if extracted != testID {
		t.Errorf("ID extraction failed: got %x, want %x", extracted, testID)
	}

	// Test Copy with same instance (self-copy)
	tile.GID = 42
	tile.Copy(tile)
	if tile.GID != 42 {
		t.Errorf("Self-copy failed: got %d, want 42", tile.GID)
	}
}
