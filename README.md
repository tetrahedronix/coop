# Coop‑Googol Engine

**Coop‑Googol** is a Go‑based Entity‑Component‑System (ECS) engine that implements the *Tharsis* double‑buffering paradigm.  It powers a 2‑D tile‑based "Road to Moscow" demo built with **Ebiten** and optional **Box2D** physics.

**Coop‑Googol** is a Go‑based Entity‑Component‑System (ECS) engine that implements the *Tharsis* double‑buffering pattern. It powers a 2‑D tile‑based demo inspired by 'Road to Moscow', built with Ebiten and optional Box2D for physics.

---

## Quick‑Start

1. **Prerequisites** – Go 1.20+ (module `go.mod`), `git`.
2. Clone the repository and initialise the module:
   ```bash
   git clone <repo‑url>
   cd coop‑googol
   go mod tidy
   ```
3. Build and run the demo:
   ```bash
   go run ./src/demos/lastduel/main.go
   ```
   The demo loads a Tiled JSON map, displays entities, and supports selection, movement (A), attack (C) and disband (D).
4. **Persistence** – Press **S** to save the current world to `save.json`; press **L** to reload.

---

## Features

| Category | Detail | Implementation |
|----------|--------|----------------|
| Rendering | 2D tile‑based map, entity icons, selection highlight | `RenderSystem` (Ebiten) |
| Input | Mouse click, CTRL+click multi‑select, keys A/C/D, S/L for save/load | `InputSystem` |
| Path‑finding | A* on a grid, Manhattan heuristic, respecting `Tile.Walkable` | `PathFindingSystem` |
| Physics (optional) | Basic collision handling via Box2D | `box2d` integration (future) |
| Persistence | JSON serialisation of world state (`save.json`) | `SaveSystem` / `LoadSystem` |
| Concurrency | Thread‑safe `Loop`; systems can run in parallel goroutines | `Loop` with `sync.RWMutex` and `sync.WaitGroup` |
| CI / Docker | GitHub Actions (build, test, lint, benchmark) and multistage Docker image | `.github/workflows/*.yml`, `Dockerfile` |

## Documentation



## License

This project is licensed under the **Apache License 2.0**.

---

*For deeper implementation details see `doc/specs.md` and the source code under `googol/`.*

**Coop‑Googol** is a Go‑based ECS (Entity‑Component‑System) engine built around the *Tharsis* double‑buffering paradigm. The repository contains the core engine, example programs, and a detailed development roadmap.

The engine targets a 2‑D tile‑based game demo that demonstrates rendering, input handling, path‑finding, and optional physics integration. The roadmap is split into staged plans, key characteristics, and a priority‑ordered implementation list.
