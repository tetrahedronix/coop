# TODO – Road to Moscow demo (Coop-Googol engine)

## 1️⃣ Staged Plan (numbered)

### 1.1 ECS Engine Implementation (core)

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| 1 | **Environment setup** | Verify Go 1.20+ toolchain, add dependencies (`ebiten/v2`, `box2d`, `go‑tiled`). Create `dev` branch. | Updated `go.mod`, local build OK | ✅ |
| 2 | **Components** | Implement `Add/Copy/Reset` methods for `Shape`, `Position`, `Velocity`. Add `Tile` and `Selectable` components. | Components with unit tests | ⬜ |
| 3 | **Core systems** | Define the `System` interface. Implement `RenderSystem` (ebiten) and `InputSystem` (mouse click, CTRL+click, A/C/D keys). | Systems registrable via `Loop.AddSystem` | ⬜ |
| 4 | **Thread‑safe loop** | Implement `Loop` with `sync.RWMutex`, methods `AddSystem`, `RemoveSystem`, `Run` (60 fps ticker). | Loop ready for dynamically adding systems | ⬜ |

### 1.2 Ebiten Integration

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| 5 | **Map loader & A\*** | JSON parser for Tiled maps (`climate`, `terrain`, `obstacle`, `height`, `walkable` properties). Implement A* algorithm on a grid. | Loadable maps, working path‑finding | ⬜ |
| 6 | **Persistence** | Save/load map state (entities, positions) to JSON (`save.json`). S/L keys for actions. | Single‑player save system | ⬜ |
| 8 | **Integration testing & benchmarking** | Headless tests (`ebiten` in headless mode) to verify the full flow. Benchmark A* on 50‑200×200 maps and the 60 fps Loop. | ≥ 80% coverage + benchmark | ⬜ |

### 1.3 Box2D Integration

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| *(no steps currently defined)* |

### 1.4 "Road to Moscow" Demo Implementation

| # | Step | Description | Expected Output | Status |
|---|------|-------------|------------------|--------|
| 7 | **Final demo** | Move the demo to `demos/lastduel/main.go`. Update imports to `github.com/tetrahedronix/coop/googol` and `github.com/tetrahedronix/coop/googol/ecs`. Display the map, unit selection (color change), movement with the **A** key, attack **C**, disband **D**. | Playable demo, runnable with `go run …` | ⬜ |
| 9 | **CI & Docker** | GitHub Actions (build, test, lint, benchmark). Multistage Docker: builder → minimal runtime. | Green CI, ready Docker image | ⬜ |
|10 | **Documentation** | Update `README` (quick‑start, screenshots). Update `AGENTS.md` with any new commands. | Complete docs | ⬜ |

## 2️⃣ Key Characteristics & Required Software

| Category | Detail | Software / Library |
|----------|--------|---------------------|
| Rendering | 2D, tile‑based, selection highlighting | `github.com/hajimehoshi/ebiten/v2` |
| Maps | Tiled JSON export (climate, terrain, obstacles, elevation, walkability) | `github.com/lafriks/go-tiled` (or manual parsing) |
| Path‑finding | A* on a grid, uniform weight, non‑walkable obstacles | In‑house implementation (Go) |
| Collisions (opt.) | Box2D for basic collisions (only if needed) | `github.com/ByteArena/box2d` |
| Persistence | Save map state to JSON | Standard library `encoding/json` |
| Concurrency | Thread‑safe loop, dynamic add/remove of systems | `sync.RWMutex` (std) |
| Testing | Unit, headless integration, benchmarking | `testing`, `go test -race`, `ebiten` headless mode |
| CI/Deploy | GitHub Actions, multistage Docker | GitHub Actions, `docker`, `alpine` |

---

## 3️⃣ Tharsis Paradigm Roadmap (ECS core – double buffering)

Priority-ordered implementation roadmap for the Googol engine, based on the current codebase and the goals of the Tharsis paradigm. These items map onto and refine the core ECS work in **section 1.1** above (components, systems, loop), and should be treated as the detailed technical breakdown of that section.

### 🔴 CRITICAL PRIORITY (Do immediately)

| # | Step | Why | What to do | Where | Status |
|---|------|-----|-------------|-------|--------|
| 11 | **Buffer swapping (SwapBuffers)** | Without this, the future buffer is never promoted to past, so any changes made during the frame are lost. Tharsis requires that at the end of the frame the future becomes the new past. | Add `SwapBuffers()` to `ecs.Entity`; add an `entities` slice to `World` and a `SwapBuffers()` method that iterates over all entities; in `Loop.Update()`, call `world.SwapBuffers()`. | `entities.go` (`SwapBuffers` method), `game.go` (`entities` field on `World`, `SwapBuffers` method), `engine.go` (call inside `Update`) | ⬜ |
| 12 | **Entity registry in the World** | The engine needs to iterate over all entities to apply systems and swap buffers. | Add `entities []*Entity` to `ecs.World` (or `googol.World`); `CreateEntity()` must add the entity to this slice; add `GetEntities()` to return the slice (or an iterator). | `game.go` (`entities` field, `CreateEntity` modification) | ⬜ |

### 🟠 HIGH PRIORITY (Next step)

| # | Step | Why | What to do | Where | Status |
|---|------|-----|-------------|-------|--------|
| 13 | **Systems that process all entities** | Systems are currently empty. For a functioning ECS, systems must iterate over all entities and process their components. | Define a `System` interface with a `Process(*Entity)` method; create a `World.AddSystem(System)` method; in `Loop.Update()`, before `SwapBuffers`, call `system.Process(e)` for every entity; systems should read from the past (via `GetPastComponent`) and write to the future (via `GetWritableComponent`). | `system.go` (interface and methods), `game.go` (`systems` slice on `World`), `engine.go` (execution loop) | ⬜ |
| 14 | **Component index management** | Using magic numbers (`e1.AddComponent(..., 1)`) is fragile and hard to read. A system to map component types to indices is needed. | Define a global or per‑entity registry mapping `ComponentTypeID` → index; or use a `map[ComponentTypeID]int` on the entity; alternatively, use a `map[ComponentTypeID]Component` instead of a slice (more flexible). | `components.go` (registry or map), `entities.go` (modify internal structure) | ⬜ |

### 🟡 MEDIUM PRIORITY (Improvements)

| # | Step | Why | What to do | Where | Status |
|---|------|-----|-------------|-------|--------|
| 15 | **Parallelism (goroutines)** | Tharsis allows concurrent execution of systems thanks to double buffering. This is the real advantage of the paradigm. | Run each system in a separate goroutine; use `sync.WaitGroup` to wait for all of them to finish before swapping; ensure systems only read from the past and only write to the future (already guaranteed if using `GetPastComponent` and `GetWritableComponent`). | `engine.go` (modify `Update` to launch goroutines) | ⬜ |
| 16 | **Resize / memory management** | Currently `GetWritable` extends the future buffer with `nil` using `make` and `copy`. More efficient memory management can improve performance. | Preallocate slices with sufficient capacity (e.g. `make([]Component, 0, expectedSize)`); when `AddComponent` is called, ensure future and past have the same length (if needed). | `entities.go` (`AddComponent`, `GetWritable` methods) | ⬜ |

### 🟢 LOW PRIORITY (Optional / future improvements)

| # | Step | Why | What to do | Where | Status |
|---|------|-----|-------------|-------|--------|
| 17 | **Query/archetype systems** | For advanced ECS engines, it's useful to group entities by component type (archetype) for faster, cache‑friendly iteration. | Each entity has a "signature" (bitmask) of its present components; systems specify which components they require; the engine tracks entities by signature. | `entities.go` (`signature` field), `world.go` (`signature → []Entity` map) | ⬜ |
| 18 | **Serialization / Deserialization** | For saving/loading game state or networking. | Add `Marshal/Unmarshal` methods to `Component`; add `Serialize/Deserialize` methods to `Entity` and `World`. | `components.go` (extended `Component` interface), `entities.go` (serialization methods) | ⬜ |
| 19 | **Events / messaging between systems** | Sometimes a system needs to notify other systems of events (e.g. collision, entity death). | Create a global event bus; systems can publish events during `Process()`; other systems can listen for them. | New `events` package or dedicated file | ⬜ |
| 20 | **Profiling and debugging** | To optimize performance and understand what's happening. | Add metrics (system execution times, number of entities processed, etc.); add a debug panel (if graphical). | `engine.go` (metrics), new `debug` package | ⬜ |

### 🗺️ Priority Summary (Tharsis roadmap)

| Priority | Task |
|----------|------|
| 🔴 **Critical** | 11. SwapBuffers, 12. Entity registry |
| 🟠 **High** | 13. Systems processing all entities, 14. Index management |
| 🟡 **Medium** | 15. Parallelism, 16. Memory management |
| 🟢 **Low** | 17‑20. Query/archetype, Serialization, Events, Profiling |

### ✅ Recommended next step

**Start with SwapBuffers**, since without it the engine isn't a true Tharsis implementation. Then move to the entity registry in the World, followed by systems that process all entities.

---

*The tables above highlight all the non‑obvious requirements agents should be aware of. Section 3 items (11‑20) are the detailed ECS/Tharsis breakdown that underlies the core engine work described in section 1.1, and should be tackled together with — or before — steps 2‑4.*

## Analysis Summary (English)

### 🔴 Critico

| # | Issue | What to do | Where |
|---|-------|-----------|-------|
| 1 | **Limited ComponentTypeID** – Only 5 component types defined via `iota` (`component.go:6-11`). Adding new components requires modifying `iota` and all `CopyComponent` cases. | Extend `ComponentTypeID` to support more types (e.g., use `map[ComponentTypeID]Component` or increase bit-width). | `component.go` |
| 2 | **Incomplete `CopyComponent`** – Only handles `*Position` and `*Shape`; panics with `unsupported component type` for `Velocity`, `Tile` (`component.go:30-45`). | Add cases for `*Velocity` and `*Tile`; consider a generic copy strategy. | `component.go` |
| 3 | **`SwapBuffers` length check** – Returns `false` if `LenComponent() != LenFutureComponent()` (`entity.go:163-165`). `Normalize` must be called manually before swap; not automatic. | Make `SwapBuffers` internally call `Normalize` or auto‑sync lengths; or document that `Normalize` must precede swap. | `entity.go` (`SwapBuffers`, `Normalize`) |
| 4 | **No component removal API** – No method to remove a single component from an entity; only `PurgeRemoved` removes entire entities (`game.go:74-94`). | Add `RemoveComponent(tid ComponentTypeID) bool` to `Entity`. | `entity.go` |

### 🔶 Urgente

| # | Issue | What to do | Where |
|---|-------|-----------|-------|
| 5 | **Entity iteration without archetype/query** – `Loop.Update()` fetches all entities and iterates O(n*m) per system (`engine.go:28-37`). No signature‑based filtering. | Implement a query/system that filters entities by required component signatures; optionally introduce archetype grouping. | `engine.go`, `entity.go` (signature) |
| 6 | **`LenComponent()` counts only past** – Returns `len(e.componentsPast)` (`entity.go:115-117`). `SwapBuffers` checks this against `LenFutureComponent()`, causing confusion about total component count. | Document that `LenComponent()` = past only; add `TotalComponents()` if needed; ensure `Normalize` is called before swap. | `entity.go` |
| 7 | **`GetFutureComponentByType` depends on past signature** – Comment notes signature reflects past, not future (`utils.go:76-79`). Scanning future linearly is O(n) and error‑prone. | Add a `futureSignature` field to `Entity` updated during `AddFutureComponent`/`SwapBuffers`; use it for O(1) lookup. | `entity.go`, `utils.go` |
| 8 | **Unclear `AddComponent`/`AddFutureComponent` API** – `AddComponent` only for past (initialization); `AddFutureComponent` adds to future. Systems may bypass double‑buffer invariants by directly mutating future. | Document the intended workflow: initialize with `AddComponent`, modify during frame with `GetWritableComponent`/`AddFutureComponent`; consider deprecating direct `AddFutureComponent` in favor of `GetWritable`. | `entity.go` |
| 9 | **No public `Loop.Stop()`** – `stopCh` channel exists but not exposed in a clean public API (`engine.go:59-82`). | Add `Stop()` method that closes `stopCh` and waits for goroutine termination; expose in `Loop` struct. | `engine.go` |

### ⚪ Necessario

| # | Issue | What to do | Where |
|---|-------|-----------|-------|
| 10 | **Missing component serialization** – No `Marshal/Unmarshal` methods on `Component` or `Entity`/`World`. Required for save/load or networking. | Add `Component` interface methods `Marshal() ([]byte, error)` and `Unmarshal(data []byte) error`; implement for each component type. | `component.go`, `entity.go` |
| 11 | **No events/messaging between systems** – No event bus; systems cannot notify each other of collisions, entity death, etc. | Create a simple global event bus (publish/subscribe) or integrate with existing `World` logger channel. | New `events` package or `world.go` |
| 12 | **No parallel goroutine execution of systems** – `Update` runs systems sequentially (`engine.go:31-37`). Tharsis advantage (concurrent systems) not realized. | Implement goroutine-per-system with `sync.WaitGroup` in `Update`; ensure systems only read past and write future (already enforced by API). | `engine.go` (`Update`) |
| 13 | **Type assertions in `Add` without error handling** – `Position.Add`, `Selectable.Add`, `Velocity.Add`, `Tile.Add` all do `data.(Type)` panics on wrong type (`basic.go:15-18,121-126`). | Add type switch with fallback or return error; or document strict type requirements. | `basic.go` |
| 14 | **No memory preallocation optimization** – `GetWritable` uses `make`+`copy` each time (`entity.go:88-93`). No capacity preallocation based on expected entity count. | Pre‑allocate `componentsFuture` with capacity during `Normalize` or `World.CreateEntity`; reuse slices when possible. | `entity.go` (`Normalize`, `GetWritable`) |

> **Recommended next step**: Start with **🔴 Critico** item 3 (fix `SwapBuffers`/`Normalize`), then address 🔶 Urgente item 5 (archetype/query system) and item 7 (future signature). These three steps give the engine a functional double‑buffer ECS core.