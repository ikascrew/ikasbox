# AGENTS.md

This file provides guidance to coding agents when working with code in this repository.

## Task guides — read before starting

| When the task involves… | Read first |
|---|---|
| adding/changing a DB column or table, editing `db/*_gen.go` or `migrate()`, "schema", "column", "スキーマ", "テーブル", "列" | [`_docs/db-schema-change.md`](_docs/db-schema-change.md) |
| adding/changing an endpoint under `handler/api/`, `apiMap`, `AddEndpoint`, a new `API.post/patch/delete` call, "API", "endpoint", "エンドポイント" | [`_docs/api-endpoint.md`](_docs/api-endpoint.md) |

## Project Overview

ikasbox is a Go server + React SPA (Material UI) for managing video/image content (groups, projects, thumbnails) used by the ikascrew VJ ecosystem. It stores metadata in SQLite (`ikasbox.db`) and uses GoCV (OpenCV) for thumbnail generation, so a working OpenCV install is required to build the Go side. Documentation and comments are largely in Japanese. The old server-rendered Go template UI (MDL) has been removed; the Go server now only serves the JSON API plus the embedded React build.

## Commands

### Go server

```
cd cmd
go run main.go init            # create SQLite database (fails if ikasbox.db exists)
go run main.go start           # start the server on port 5555
go run main.go group ...       # group subcommands: register, import, check, list, remove
go run main.go content ...     # content register <group-id> <name> <type> [params-json] — generated (plugin) contents
go run main.go project ...     # project subcommands: register, list, add, remove
```

- Flags (defined in `cmd/main.go`): `-db <file>` (default `ikasbox.db`), `-ext <patterns>` (import extensions, default `*.mp4,*.mpeg,*.png,*.jpg,*.jpeg`). A DB filename can also be embedded at build time via `-ldflags "-X main.embedDB=xxxx.db"`.
- `cmd/develop.sh` runs the server with `skewer` for live reload during development.
- Tests: `go test ./...` (tests live in `db/`, `handler/api/`, `contentimport/`, `config/`). Run a single test with e.g. `go test ./db/ -run TestName`.
- Frontend tests: `npm test` (in `frontend/`, runs `vitest run`). See the React frontend section below.
- **`npm run build` (in `frontend/`) must be run at least once before `go build`/`go run` on a fresh checkout.** The React build output is `go:embed`-ed directly into the `handler/internal` package (see below) — the embed directive fails to compile if that directory has no files in it. A placeholder `.gitkeep` is committed so a pristine clone still compiles; after the first `npm run build` the real build output satisfies the embed.

### React frontend (`frontend/`)

```
npm run dev      # Vite dev server on http://localhost:3000, proxies /api, /thumb, /content/media to :5555
npm run build    # production build, output straight into handler/internal/_assets/spa (embedded by the Go server)
npm run preview  # preview the production build
npm test         # vitest run — unit/component tests, jsdom environment
```

Development flow (the maintainer's own session): run the Go server (port 5555) and `npm run dev` together; the dev server proxies API calls. In production, `go run main.go start` alone serves everything: the built React SPA (embedded from `handler/internal/_assets/spa`) plus the JSON API — no separate static file server or `npm run dev` needed.

### Verifying changes

- Verify with static checks and tests only: `go build ./...`, `go vet ./...`, `go test ./...`, and in `frontend/` `npm run build` (confirms the Vite/embed build compiles) + `npm test`. Add colocated vitest tests for new pure logic or presentational components where practical.
- Do **not** start or stop the Go server (`go run main.go start`) or the Vite dev server (`npm run dev`), and never kill whatever is listening on port 5555/3000 to "free" it — the maintainer runs their own dev session and checks the UI in their own browser. Don't drive a browser to click through the UI either. If a live check is genuinely needed, ask the maintainer to run it.

## Architecture

Details and background for each package: [`_docs/architecture.md`](_docs/architecture.md).

- Go: `cmd/main.go` (CLI) → `ikasbox.go` (subcommands, multicast + HTTP start) / `contentimport/` (scan, register, thumbnails; `RegisterGenerated` for plugin contents) / `config/` / `db/` / `handler/` (media, thumb, SPA embed) / `handler/api/` (JSON API).
- React: `frontend/src/` — class components + MUI + react-router v7, `API.js` fetch wrapper, theme in `src/theme.js`. JSX goes in `.jsx` files.
- `db/*_gen.go` are argen-generated but argen no longer runs — edit them by hand following [`_docs/db-schema-change.md`](_docs/db-schema-change.md).
- **`/project/content/list/{id}` and `handler.ProjectResponse` are used directly by the sibling `server` and `client` repos.** Do not remove or reshape them without updating both repos in lockstep.