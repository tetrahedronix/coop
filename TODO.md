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

> ⚠️ **ARCHITECTURAL PIVOT (see item 22)** — Items 11, 12 and 13 below were written under an **entity-as-container** model (`Entity` owning `componentsPast`/`componentsFuture` directly). This model was identified as a deviation toward OOP-style composition rather than a data-oriented ECS. The engine is pivoting to a **Store-per-type (SoA)** model: `Entity` becomes a plain ID, and each component type owns its own double-buffered `Store[T]`. Items 11-13 are kept below for historical context (they explain *why* the pivot happened) but are **superseded by item 22** and should not be implemented as originally written.

### 🔴 CRITICAL PRIORITY (Do immediately)

| # | Step | Why | What to do | Where | Status |
|---|------|-----|-------------|-------|--------|
| 11 | ~~**Buffer swapping (SwapBuffers)**~~ **[SUPERSEDED — see #22]** | Without this, the future buffer is never promoted to past, so any changes made during the frame are lost. Tharsis requires that at the end of the frame the future becomes the new past. | ~~Add `SwapBuffers()` to `ecs.Entity`~~ → SwapBuffers now belongs to `Store[T]`, one per component type, not to `Entity`. | `store.go` (`SwapBuffers` method per store), `game.go` (`World` iterates registered stores) | 🔄 superseded |
| 12 | **Entity registry in the World** | The engine needs to iterate over all entities to apply systems and swap buffers. | Still valid in spirit: `World` needs a registry of live entity IDs. Changes: entities are now plain `EntityID` values (no embedded component data); `CreateEntity()` returns an `EntityID`; `GetEntities()` returns `[]EntityID`. | `game.go` (`entities []EntityID` field on `World`) | ⬜ (re-scoped) |

### 🟠 HIGH PRIORITY (Next step)

| # | Step | Why | What to do | Where | Status |
|---|------|-----|-------------|-------|--------|
| 13 | ~~**Systems that process all entities**~~ **[SUPERSEDED — see #22]** | Systems are currently empty. For a functioning ECS, systems must iterate over all entities and process their components. | ~~`Process(*Entity)`~~ → systems now receive an `EntityID` and read/write component data via the relevant `Store[T]`, not via methods on `Entity`. | `system.go` (interface reworked), `store.go` | 🔄 superseded |
| 14 | **Component index management** | Using magic numbers (`e1.AddComponent(..., 1)`) is fragile and hard to read. A system to map component types to indices is needed. | **Done and unaffected by the pivot.** Hybrid identity scheme in place: `ComponentTypeID` (bitmask, `TypedComponent` fast-path) + `ComponentID` (dense registry, `RegisterComponent[T]`), each with distinct methods (`HasTypedComponent` / `HasComponent`). This scheme survives the AoS→SoA pivot unchanged — it will simply be consulted by `World`/`Store` instead of by `Entity`. | `component.go` (registry, `TypedComponent`), `entity.go` (identity checks, to be relocated to `World` in item 22) | ✅ |

### 🟡 MEDIUM PRIORITY (Improvements)

| # | Step | Why | What to do | Where | Status |
|---|------|-----|-------------|-------|--------|
| 15 | **Parallelism (goroutines)** | Tharsis allows concurrent execution of systems thanks to double buffering. This is the real advantage of the paradigm. | Run each system in a separate goroutine; use `sync.WaitGroup` to wait for all of them to finish before swapping; ensure systems only read from the past and only write to the future. Becomes more natural under the Store-per-type model: each store can be swapped/processed independently. | `engine.go` (modify `Update` to launch goroutines) | ⬜ |
| 16 | ~~**Resize / memory management**~~ **[SUPERSEDED — see #22]** | Currently `GetWritable` extends the future buffer with `nil` using `make` and `copy`. More efficient memory management can improve performance. | This concern moves to `Store[T]`: preallocate `values`/`entities` slices with sufficient capacity; the dense, per-type layout of a Store makes this simpler than the old per-entity heterogeneous slice. | `store.go` | 🔄 superseded |

### 🟢 LOW PRIORITY (Optional / future improvements)

| # | Step | Why | What to do | Where | Status |
|---|------|-----|-------------|-------|--------|
| 17 | **Query/archetype systems** | For advanced ECS engines, it's useful to group entities by component type (archetype) for faster, cache‑friendly iteration. | Under the Store-per-type model this becomes more natural: `World` can intersect the entity sets of the stores a system requires, instead of scanning a per-entity bitmask. | `world.go` (`signature → []EntityID` map, or per-store entity index) | ⬜ |
| 18 | **Serialization / Deserialization** | For saving/loading game state or networking. | Add `Marshal/Unmarshal` methods to `Component`; add `Serialize/Deserialize` to `Store[T]` and `World`. Consider adding `byName` to the registry at this point, for stable on-disk identity (deliberately deferred until now — see design discussion). | `component.go` (extended interface), `store.go`, `registry.go` | ⬜ |
| 19 | **Events / messaging between systems** | Sometimes a system needs to notify other systems of events (e.g. collision, entity death). | Create a global event bus; systems can publish events during `Process()`; other systems can listen for them. | New `events` package or dedicated file | ⬜ |
| 20 | **Profiling and debugging** | To optimize performance and understand what's happening. | Add metrics (system execution times, number of entities processed, etc.); add a debug panel (if graphical). | `engine.go` (metrics), new `debug` package | ⬜ |
| 21 | **Unified variable bitset** | The dual-field identity scheme (`signature` bitmask + `dynamicSignature` map) works, but it requires two separate checks during every lookup. A unified variable bitset would combine the fast-path and registry into a single mechanism. | Replace `signature`/`dynamicSignature` with a variable-length bitset; bits 0–63 reserved for `TypedComponent`, 64+ dynamically assigned by the registry. Rewrite `HasComponent`/`HasTypedComponent`/`AddComponent` accordingly. **To be started only once the POC is complete and working**, and only after item 22 (the AoS→SoA pivot) has settled, since the bitset will most likely live on `World`/`Store`, not on `Entity`. | `world.go` or `store.go`, `component.go` | ⬜ |
| 22 | **Pivot AoS→SoA (store‑per‑type architecture)** | `Entity` owning `componentsPast`/`componentsFuture` directly is an Array‑of‑Structures model — closer to OOP‑style composition than to a data‑oriented Tharsis ECS. In a classic Tharsis/ECS design, component data lives in per‑type stores (Structure‑of‑Arrays), and `Entity` is nothing more than an ID. This item formalizes the decision to correct that deviation before proceeding further. | Redefine `Entity` as a plain ID type (e.g. `type EntityID uint64`, generated via Sonyflake as today). Introduce a generic `Store[T Component]` with its own `past`/`future` slices and its own `SwapBuffers()` — one store instance per component type, analogous to `ComponentStore[T]` already prototyped in the MLOps variant (`mlops.go`). `World` holds the registry of live `EntityID`s plus a lightweight per‑entity signature (reusing the hybrid `ComponentTypeID`/`ComponentID` scheme from item 14) so that `System.Match` can query "which entities have component X" without scanning every store. Systems read from `Store[T].past` and write to `Store[T].future` for the entities they match, instead of calling methods on `Entity`. | New `entity.go` (ID-only), new `store.go` (generic double-buffered store), `game.go` (`World` holds stores + entity registry + signatures), `system.go` (interface reworked to operate via `EntityID` + stores) | 🔄 in corso |

### 🗺️ Priority Summary (Tharsis roadmap)

| Priority | Task |
|----------|------|
| 🔴 **Critical** | 12. Entity registry (re-scoped to plain IDs), 22. AoS→SoA pivot |
| 🟠 **High** | 14. Index management (done, unaffected by pivot) |
| 🟡 **Medium** | 15. Parallelism, 16. Memory management (superseded, moves to Store) |
| 🟢 **Low** | 17‑21. Query/archetype, Serialization, Events, Profiling, Unified bitset |

### ✅ Recommended next step

**Complete item 22** (the AoS→SoA pivot): define `EntityID`, then `Store[T]`, then rework `World` to hold stores + a lightweight entity signature reusing the hybrid identity scheme from item 14. Only once `Store[T]` exists does it make sense to revisit parallelism (15) and archetype queries (17), since both become substantially simpler under a per-type store layout.

---

*The tables above highlight all the non‑obvious requirements agents should be aware of. Section 3 items (11‑22) are the detailed ECS/Tharsis breakdown that underlies the core engine work described in section 1.1, and should be tackled together with — or before — steps 2‑4. Items marked SUPERSEDED reflect design decisions later corrected by item 22; they are kept for historical traceability, not as active work items.*

## Analysis Summary (English)

### 🔴 Critico

| # | Issue | What to do | Where |  Status |
|---|-------|-----------|-------|----------|
| 1 | ~~**Limited ComponentTypeID**~~ – Only 5 component types defined via `iota` (`component.go:6-11`). Adding new components requires modifying `iota` and all `CopyComponent` cases. | **Resolved.** `ComponentTypeID` (bitmask, `TypedComponent` fast-path) now coexists with `ComponentID` (dense registry via `RegisterComponent[T]`). The 64-type limit applies only to the opt-in `TypedComponent` fast-path; components registered dynamically have no cardinality limit. | `component.go` | ✅ |
| 2 | **Incomplete `CopyComponent`** – Only handles `*Position` and `*Shape`; panics with `unsupported component type` for `Velocity`, `Tile` (`component.go:30-45`). | Not yet addressed. Will likely be replaced entirely by a `DeepCopier` satellite interface (mirroring `mlops.go`'s `DeepCopyComponent()`), consistent with the marker-interface direction taken for `Component`/`TypedComponent`. | `component.go` | ⬜ |
| 3 | ~~**`SwapBuffers` length check**~~ **[SUPERSEDED — see roadmap #22]** – Returns `false` if `LenComponent() != LenFutureComponent()` (`entity.go:163‑165`). | This class of problem is specific to the entity‑as‑container model and disappears once component data moves into per‑type `Store[T]`: each store manages its own past/future lengths independently, with no cross‑component length coupling to reconcile. | *(moves to `store.go`)* | 🔄 superseded |
| 4 | **No component removal API** – No method to remove a single component from an entity; only `PurgeRemoved` removes entire entities (`game.go:74-94`). | Still needed, re-scoped: `Store[T].Remove(EntityID)` under the new model, rather than `Entity.RemoveComponent`. | `store.go` (once created) | ⬜ (re-scoped) |

### 🔶 Urgente

| # | Issue | What to do | Where | Status |
|---|-------|-----------|-------|---------|
| 5 | **Entity iteration without archetype/query** – `Loop.Update()` fetches all entities and iterates O(n*m) per system (`engine.go:28-37`). No signature‑based filtering. | Still open. Becomes more tractable once entities carry a lightweight signature in `World` (item 14 scheme) and systems can pre-filter by it before touching any store. | `engine.go`, `world.go` (signature) | ⬜ |
| 6 | ~~**`LenComponent()` counts only past**~~ **[SUPERSEDED — see roadmap #22]** – Returns `len(e.componentsPast)` (`entity.go:115-117`). | This ambiguity is specific to the single-entity-owns-all-components model and disappears under Store-per-type: each `Store[T]` has its own unambiguous `Len()`. | *(moves to `store.go`)* | 🔄 superseded |
| 7 | ~~**`GetFutureComponentByType` depends on past signature**~~ **[SUPERSEDED — see roadmap #22]** – Comment notes signature reflects past, not future (`utils.go:76‑79`). | Under Store-per-type, "does entity X have a future value for component T" becomes a direct lookup in `Store[T].future` by `EntityID` — no signature ambiguity, since there's no shared cross-component signature to go stale. | *(moves to `store.go`)* | 🔄 superseded |
| 8 | ~~**Unclear `AddComponent`/`AddFutureComponent` API**~~ – `AddComponent` only for past (initialization); `AddFutureComponent` adds to future. Systems may bypass double‑buffer invariants by directly mutating future. | **Partially resolved, then superseded.** `Entity.AddComponent` was rewritten with a clear contract (dispatches to bitmask or registry identity, returns an error if a non-`TypedComponent` type isn't registered) — see engine refactoring log. However, this now moves to `Store[T].Add` / `Store[T].GetWritable` under item 22, since components no longer live on `Entity`. | `store.go` (once created) | 🔄 superseded |
| 9 | **No public `Loop.Stop()`** – `stopCh` channel exists but not exposed in a clean public API (`engine.go:59-82`). | Still open, unaffected by the pivot. Add `Stop()` method that closes `stopCh` and waits for goroutine termination. | `engine.go` | ⬜ |

### ⚪ Necessario

| # | Issue | What to do | Where |
|---|-------|-----------|-------|
| 10 | **Missing component serialization** – No `Marshal/Unmarshal` methods on `Component` or `Entity`/`World`. Required for save/load or networking. | Add `Component` interface methods, or a `DeepCopier`-style satellite interface for serialization; implement per component type. Natural point to also add `byName` to the registry for stable on-disk identity. | `component.go`, `store.go` |
| 11 | **No events/messaging between systems** – No event bus; systems cannot notify each other of collisions, entity death, etc. | Create a simple global event bus (publish/subscribe) or integrate with existing `World` logger channel. | New `events` package or `world.go` |
| 12 | **No parallel goroutine execution of systems** – `Update` runs systems sequentially (`engine.go:31-37`). Tharsis advantage (concurrent systems) not realized. | Implement goroutine-per-system with `sync.WaitGroup` in `Update`; more natural once systems operate on independent `Store[T]` instances. | `engine.go` (`Update`) |
| 13 | **Type assertions in `Add` without error handling** – **Largely moot**: the `Add(interface{})` method was removed from `Component` entirely (now a pure marker interface, `ECSComponent()`). Components are constructed fully-formed by the caller instead of being populated via `Add` after the fact. | No action needed on the old `Add` methods; verify no residual callers still depend on the old `Component.Add` contract when migrating `basic.go` types to the new interface. | `basic.go` | 🔄 in corso |
| 14 | **No memory preallocation optimization** – `GetWritable` uses `make`+`copy` each time (`entity.go:88-93`). No capacity preallocation based on expected entity count. | Re-scoped to `Store[T]`: preallocate `values`/`entities` with capacity during store construction; simpler than the old per-entity heterogeneous slice since each store is homogeneous. | `store.go` (once created) |

> **Recommended next step**: Complete the AoS→SoA pivot (roadmap item 22): define `EntityID`, then `Store[T]`, then rework `World`. Several items above (3, 6, 7, 8, 14) collapse or simplify automatically once component data lives in per-type stores rather than on `Entity` — resolving the pivot resolves them as a side effect.
