# Engine Specifications (Coop) — Tharsis Paradigm

This document is the **single source of truth** for the engine implementation. All paths, dependencies, and non‑obvious behaviors are defined here.

---

## 1. Core Architecture – Tharsis ECS

The engine follows the **Tharsis double‑buffering** paradigm: each entity maintains two component slices, `componentsPast` and `componentsFuture`. Systems **read** from the past (current state) and **write** to the future (next state). At the end of each frame, the future becomes the new past via a buffer swap.

### 1.1 Package Layout (final)
```
/googol/                    # high‑level engine (World, Loop, Systems)
├─ engine.go                # Loop implementation (thread‑safe, ticker, system execution)
├─ world.go                 # World, entity registry, SwapBuffers, GUID handling (sonyflake)
├─ entity.go                # Entity struct, component slices (past/future), AddComponent, GetComponent
├─ system.go                # System interface (Process), base system implementations (Render, Input, Pathfinding, Save, Load)
└─ ...

/googol/ecs/                # core ECS components (Component, Shape, Position, Velocity, Tile, Selectable)
├─ component.go             # Component interface definition (Add, Copy, Reset)
├─ shape.go                 # Shape component (primitive string)
├─ position.go              # Position component (x, y)
├─ velocity.go              # Velocity component (v, d)
├─ tile.go                  # Tile component (climate, terrain, obstacle, height, walkable)
└─ selectable.go            # Selectable component (selected flag, highlight color)
```

### 1.2 Entity

```go
type Entity struct {
    guid       uint64                 // Sonyflake ID
    componentsPast   []Component      // current state (read‑only for systems)
    componentsFuture []Component      // next state (writable by systems)
    // optional: componentTypeIndex map[ComponentTypeID]int for fast lookup
}
```

- `AddComponent(c Component, data ...interface{})`:
  - Appends the component to `componentsFuture` (and also `componentsPast` if not present) and calls `c.Add(data...)`.
  - The component is placed in both buffers to ensure consistency.
- `GetPastComponent(index int) Component` returns the component from the past slice.
- `GetWritableComponent(index int) Component` returns the component from the future slice, **extending the future slice with `nil`** if necessary to match the past length. This method is used by systems to obtain a reference to the future component for modification.
- `SwapBuffers()` swaps `componentsPast` and `componentsFuture` (simple pointer swap). After swap, the future slice becomes the new past, and the old past becomes the new future (which may be cleared or reused).

### 1.3 World

```go
type World struct {
    mu        sync.RWMutex
    entities  []*Entity          // registry of all entities
    // … other fields (map, systems, etc.)
}
```

- `CreateEntity() *Entity` creates a new entity, generates a GUID (Sonyflake), and appends it to `World.entities`.
- `GetEntities() []*Entity` returns the current list of entities (or a copy if needed for safe iteration).
- `SwapBuffers()` iterates over all entities and calls `SwapBuffers()` on each.

### 1.4 Component Index Management (optional but recommended)

To avoid magic numbers, a **component type registry** can map a `ComponentTypeID` (e.g., a string or integer) to an index in the component slices. This can be implemented as a global map or per‑world registry. The specification does not mandate a specific approach, but it is **strongly recommended** to use a central registry to simplify system development.

---

## 2. Loop (main game loop)

The `Loop` is thread‑safe and runs at a fixed tick rate (default 60 Hz).

```go
type Loop struct {
    mu      sync.RWMutex
    world   *World
    systems []System
    ticker  *time.Ticker
}
```

- `NewLoop(w *World) *Loop` initialises the loop with the given world.
- `AddSystem(s System)` and `RemoveSystem(s System)` are thread‑safe (they acquire the write lock).
- `Run()`:
  ```go
  func (l *Loop) Run() {
      for range l.ticker.C {
          l.mu.RLock()
          // 1. Execute all systems concurrently (see §3)
          // 2. After all systems finish, call l.world.SwapBuffers()
          l.mu.RUnlock()
      }
  }
  ```
- Systems are executed **in registration order**, but they may run concurrently (see §3.1) if parallelism is enabled. The lock ensures that system registration is not modified mid‑frame.

---

## 3. Systems

### 3.1 System Interface

```go
type System interface {
    // Process reads the past state of the given entity and writes to its future state.
    // It is called for every entity in the world.
    Process(e *Entity)
}
```

- Systems **must not** read from the future buffer or write to the past buffer.
- The `Loop` is responsible for calling `Process` on every entity, for each system, before swapping buffers.
- Systems may be executed **sequentially** or **concurrently** (goroutines) — the specification allows parallelism as a medium‑priority feature. When concurrent, the `Loop` must wait for all systems to finish using a `sync.WaitGroup` before swapping.

### 3.2 System Implementations

| System | Responsibilities | Key Methods |
|--------|------------------|-------------|
| `RenderSystem` | Draws the map, tiles, entities, and selection highlights using Ebiten. | `Process(e *Entity)` (reads `Position`, `Shape`, `Selectable`) |
| `InputSystem` | Handles mouse clicks, CTRL+click multi‑selection, keys A/C/D for actions, S/L for save/load. | `Process(e *Entity)` (updates `Selectable` flags, queues movement commands) |
| `PathFindingSystem` | Computes A* paths when an entity receives a move command. | `Process(e *Entity)` (reads pending move command, writes new `Position` or `Velocity` to future) |
| `SaveSystem` | Serialises the world state to JSON when the user presses **S**. | `Process(e *Entity)` (triggers save) |
| `LoadSystem` | Deserialises `save.json` and rebuilds the world when **L** is pressed. | `Process(e *Entity)` (triggers load) |

**Note:** For `SaveSystem` and `LoadSystem`, they may not process every entity individually; they can act as singletons that check a global flag. The specification allows systems to have internal state (e.g., `SaveSystem` could maintain a `saveRequested` boolean). The `Process` method is still called for each entity, but it can ignore entities and act on a global trigger.

---

## 4. Map Handling

- **Format**: Tiled **JSON** export.
- **Parser**: `github.com/lafriks/go-tiled` or manual JSON unmarshalling.
- **Data model**:
  ```go
  type Tile struct {
      Climate   string `json:"climate"`
      Terrain   string `json:"terrain"`
      Obstacle  bool   `json:"obstacle"`
      Height    int    `json:"height"`
      Walkable  bool   `json:"walkable"`
  }
  ```
- The map is stored as `[][]*Tile` in `World`.

---

## 5. Path‑finding (A*)

- Input: `start Point`, `goal Point`, reference to `World` map.
- Heuristic: Manhattan distance.
- Returns a slice of `Point` representing the tile path.
- Collision detection uses `Tile.Walkable` (obstacle tiles are non‑walkable). Box2D integration is optional for finer obstacles.

---

## 6. Persistence

- **File**: `save.json` at the repository root (or configurable path).
- **Structure**: `WorldState` containing:
  - Map tiles (array of tiles).
  - Entities: a slice of `EntityState` (GUID + component data serialised as JSON).
- **Serialisation**: `encoding/json` – human‑readable and easy to debug.

---

## 7. Demo – `demos/lastduel/main.go`

- **Entry point**:
  - Creates a `World`.
  - Loads a Tiled JSON map.
  - Registers `RenderSystem`, `InputSystem`, `PathFindingSystem`, `SaveSystem`, `LoadSystem` with the `Loop`.
  - Starts the loop.
- **Workflow**:
  1. Mouse click selects an entity (highlight color changes). CTRL+click toggles multi‑selection.
  2. Press **A** → if there is a selection and a destination tile (clicked after), A* path is calculated and the entity moves.
  3. Press **C** → confirms attack (placeholder; logs to output).
  4. Press **D** → disbands the entity (removes from `World.entities`).
  5. **S/L** → saves/loads the state.
- **Rendering**: Ebiten draws the map grid, entities as simple icons, and selected entities with a brighter color.

---

## 8. Testing & Benchmarking

- **Unit tests** for all components (`Add/Copy/Reset`), `Entity`, and `World` (GUID uniqueness).
- **Integration test** (headless) that:
  - Starts a `Loop` with all systems.
  - Loads a small JSON map.
  - Simulates clicks and key presses.
  - Verifies that the entity ends up on the target tile.
- **Benchmarks**:
  - `BenchmarkAStar-8` on 50×50, 100×100, and 200×200 maps.
  - `BenchmarkLoop-8` for 60 fps with 100 entities.

---

## 10. Tharsis Paradigm – Additional Specifications (from TODO §3)

The following items are **critical** for a correct Tharsis implementation and are integrated into the above specification.

### 10.1 Double‑Buffering (SwapBuffers)
- **Critical**: Every entity must have `SwapBuffers()` that swaps its `componentsPast` and `componentsFuture`.
- `World.SwapBuffers()` calls `SwapBuffers()` on all entities.
- The `Loop` calls `world.SwapBuffers()` **after** all systems have finished processing and before the next tick.

### 10.2 Entity Registry
- `World` holds a slice `entities []*Entity`.
- `CreateEntity()` appends the new entity to this slice.
- `World.GetEntities()` returns the slice (or a copy for safe iteration).
- Systems iterate over `World.GetEntities()` and call `Process(e)` on each.

### 10.3 System Processing
- The `System` interface defines `Process(*Entity)`.
- The `Loop` calls `Process(e)` for every entity, for each system, in order (or concurrently if parallelism is enabled).
- Systems **must** use `GetPastComponent` for reading and `GetWritableComponent` for writing to the future.

### 10.4 Concurrency (optional, medium priority)
- When enabled, each system runs in its own goroutine.
- A `sync.WaitGroup` ensures all systems finish before swapping buffers.
- Systems must be stateless with respect to shared mutable data (or use appropriate synchronisation).

### 10.5 Component Index Management (recommended)
- To avoid hard‑coded indices, a component type registry (e.g., `map[reflect.Type]int` or a global `ComponentID`) can map component types to slice indices.
- This is not mandatory but strongly advised for maintainability.

### 10.6 Memory Management (medium priority)
- Slices should be pre‑allocated with sufficient capacity to avoid frequent reallocations.
- `GetWritableComponent` should ensure that the future slice has the same length as the past slice (extending with `nil` as needed).

### 10.7 Future Extensions (low priority)
- **Query/archetype systems** – entities grouped by component signature for cache‑friendly iteration.
- **Serialisation/deserialisation** – `Marshal/Unmarshal` methods for components.
- **Events** – an event bus for cross‑system communication.
- **Profiling** – metrics and debug panels.

---

*This document is the authoritative reference. All implementations must adhere to these specifications.*