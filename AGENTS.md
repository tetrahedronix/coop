**AGENTS.md – High‑Signal Guidance for OpenCode Agents**

- **Go toolchain** – The repo does not include a `go.mod`. Before running any `go` command initialise a module:
  ```bash
  go mod init github.com/tetrahedronix/coop/googol && go mod tidy
  ```
  (Agents often miss this step and get `cannot find module` errors.)

- **Running examples** – Two runnable demos are provided:
  - `go run ./examples/createEntity/main.go` – Shows basic world/entity creation and component addition.
  - `go run ./examples/doubleBuffer/main.go` – Demonstrates the past/future double‑buffer model and the `GetWritableComponent` helper.

- **Core entry points** (agents frequently look for a `main` package and miss these):
  - `googol.NewWorld()` – creates a `World` with an enabled tick flag.
  - `World.CreateEntity()` – returns a new `*ecs.Entity` with a Sonyflake GUID.
  - `World.CreateSystem()` – returns a fresh `*ecs.System` (currently a stub).
  - Component helpers:
    - `googol.NewPosition()` / `googol.NewShape()` return concrete `Component` implementations.
    - `e.AddComponent(comp, data…)` adds to the *past* slice (initialisation).
    - `e.AddFutureComponent(comp, data…)` adds to the *future* slice.

- **Double‑buffer helpers (non‑obvious API)** – located in `googol/utils.go`:
  - `GetPastComponent(e, idx)` – read‑only access to the *past* component.
  - `GetFutureComponent(e, idx)` – direct access to the *future* component (may be nil).
  - `GetWritableComponent(e, idx)` – returns a mutable component from the *future*; if the future slot is empty it clones the past component via `ecs.CopyComponent`.
  - These helpers are essential for the Tharsis “double buffering” pattern; agents that bypass them will break the intended immutability.

- **Docker workflow** – `Dockerfile` builds a container with:
  - `golang:1.26-bookworm`
  - Node LTS via NVM, placed on `PATH` (`/usr/local/nvm-default/bin`).
  - Global `opencode-ai` installation (`npm install -g --allow-scripts=opencode-ai`).
  - Workdir set to `/workspace`.
  - Typical usage:
    ```bash
    docker build -t googol .
    docker run --rm -v "$PWD":/workspace -w /workspace googol go run ./examples/createEntity/main.go
    ```
  - Agents often forget to mount the source directory; without `-v` the container sees an empty workspace.

- **Missing implementation warnings** – `googol/engine.go` and `googol/engine.go` contain only stubs (`NewLoop`, `Update`, `MainLoop`, `Add`). The current library does not provide a running game loop; attempts to call these will be no‑ops.

- **Road‑map reference** – `TODO.md` lists concrete steps (e.g., implementing `SwapBuffers`, entity registry, system processing). Treat it as the single source of truth for pending features; agents should not infer missing code from comments.

- **Common pitfalls**
  1. Forgetting to import the correct path (`github.com/tetrahedronix/coop/googol`). The repo's import statements are in Italian comments but the actual path is required.
  2. Assuming `Entity.LenComponent()` returns the total of past + future components – it only counts past components.
  3. Using `e.componentsFuture` directly; always use the provided helpers to preserve double‑buffer invariants.
  4. Running `go test` – there are no test files; the project relies on example programs for manual verification.

- **Formatting / linting** – No explicit linter config is present; default `go vet` and `golint` apply.

---
*This file is deliberately concise, focusing on facts an autonomous agent would otherwise miss. It should be kept in sync with the repository’s evolving architecture.*