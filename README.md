# Coop‑Googol Engine

**Coop‑Googol** is a Go‑based Entity‑Component‑System (ECS) engine that implements the *Tharsis* double‑buffering pattern. It powers a 2‑D tile‑based demo inspired by 'Road to Moscow' game (1984), built with Ebiten and optional Box2D for physics.

The included examples demonstrate rendering, input handling, path‑finding, and optional physics integration. The roadmap is split into staged plans, key characteristics, and a priority‑ordered implementation list.


---

## Quick‑Start

1. **Prerequisites** – Go 1.20+ (module `go.mod`), `git`.
2. Clone the repository and initialise the module:
   ```bash
   git clone https://github.com/tetrahedronix/coop.git
   cd coop
   go mod tidy
   ```
3. Build and run the demo:
   ```bash
   go run ./src/examples/createEntity/main.go
   ```
   This example initializes a Googol world, creates an entity, attaches Shape and Position components, adds a dummy system, then prints the entity ID and each component's data.

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

## Disclaimer

This software was developed with the assistance of AI under continuous human supervision throughout the entire development process. AI-generated code is never accepted as-is: every contribution is reviewed, verified, tested, and consciously approved before being integrated into the project. AI is used as an assistant to make the development process faster, more efficient, and aligned with current technologies, while all final decisions remain under human responsibility.

<p align="center">
  <img src="share/ai-assisted.png" alt="AI Assisted disclaimer"/>
</p>



## License

This project is licensed under the **Apache License 2.0**.
