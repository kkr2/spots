# Spot Query

A small Go REST API for querying **GIS / location data**. Given a geographic
coordinate, it returns the nearby "spots" (restaurants, pubs, shops, hotels, …)
within a configurable range, ordered from nearest to furthest, with pagination.

Spatial lookups are powered by **PostgreSQL + PostGIS**; the HTTP layer is built
with the **Echo** framework. The codebase follows a clean, layered architecture
(domain → repository → service → delivery).

---

## Table of contents

- [What it does](#what-it-does)
- [Tech stack](#tech-stack)
- [Architecture](#architecture)
- [How the spatial query works](#how-the-spatial-query-works)
- [API reference](#api-reference)
- [Configuration](#configuration)
- [Getting started](#getting-started)
  - [Run with Docker Compose (recommended)](#run-with-docker-compose-recommended)
  - [Run locally](#run-locally)
- [Make targets](#make-targets)
- [Notes & known caveats](#notes--known-caveats)

---

## What it does

The service exposes a single domain endpoint that answers the question:

> *"What spots exist near this latitude/longitude, sorted by distance?"*

- Accepts a coordinate (`lat`, `lon`) and a search `range`.
- Finds all spots whose location falls inside the bounding box around that point.
- Orders results by true distance to the point (nearest first).
- Returns a paginated payload with `total_count`, `total_pages`, `page`, `size`,
  `has_more`, and the list of `spots`.

The database is seeded on first start with **~12,000 real-world spots** (pubs,
restaurants, hotels, shops across the UK, Europe and the US).

---

## Tech stack

| Concern            | Library / tool                                   |
| ------------------ | ------------------------------------------------ |
| Language           | Go 1.18                                           |
| HTTP framework     | [Echo v4](https://echo.labstack.com/)            |
| Database           | PostgreSQL 14 + [PostGIS](https://postgis.net/) 3.3 |
| DB access          | [sqlx](https://github.com/jmoiron/sqlx) + [pgx](https://github.com/jackc/pgx) driver |
| Migrations         | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Configuration      | [viper](https://github.com/spf13/viper) (YAML + env) |
| Logging            | [zap](https://github.com/uber-go/zap)            |
| Validation         | [validator/v10](https://github.com/go-playground/validator) |
| Input sanitizing   | [bluemonday](https://github.com/microcosm-cc/bluemonday) |
| Linting            | [golangci-lint](https://golangci-lint.run/)      |
| Containerization   | Docker + Docker Compose                          |

---

## Architecture

The project is organized using a clean / layered architecture. Each spot
feature is split into independent layers so that business logic, persistence and
transport concerns stay decoupled.

```
.
├── cmd/
│   └── main.go                 # Entry point: load config, init logger + DB, start server
├── internal/
│   ├── config/                 # Config structs + YAML files (local / docker)
│   ├── server/                 # Echo server bootstrap and route wiring
│   │   ├── server.go           # Server lifecycle (start, graceful shutdown)
│   │   └── handlers.go         # Maps groups: /api/v1/spots, /api/v1/health
│   └── spots/                  # The "spots" feature, layered:
│       ├── domain/             #   Spot / Geography / SpotList models
│       ├── repository/         #   SQL queries against PostGIS
│       ├── service/            #   Business logic / use-cases
│       └── delivery/http/      #   HTTP handlers + route definitions
├── pkg/                        # Reusable, app-agnostic helpers
│   ├── db/postgres/            # Postgres connection + auto-migrations
│   ├── httpErrors/             # Typed HTTP error responses
│   ├── logger/                 # zap logger wrapper
│   ├── sanitize/               # JSON sanitizing
│   └── utils/                  # Pagination, coordinate parsing, validation, http helpers
├── migrations/                 # SQL schema + seed data (up / down)
├── docker/Dockerfile           # Multi-stage build (golang → scratch)
├── docker-compose.yaml         # API + PostGIS services
└── Makefile                    # Common dev commands
```

**Request flow:**

```
HTTP request
   → delivery/http (handler: parse pagination + coordinate from query)
      → service (validate coordinate, apply use-case)
         → repository (run PostGIS spatial query)
            → PostgreSQL / PostGIS
```

---

## How the spatial query works

Coordinates are stored in the `spots.coordinates` column as a PostGIS `geometry`
in **SRID 4326 (WGS84 lat/lon)**.

To search "within range", the repository:

1. Projects the query point from **SRID 4326 → SRID 2163** (US National Atlas
   Equal Area, a metric projection) so distances can be reasoned about in meters.
2. Builds a bounding box (`ST_MakeEnvelope`) of `±range` around the point and
   projects it back to 4326.
3. Selects every spot whose geometry **`ST_Intersects`** that box.
4. Orders results by distance to the point using the KNN operator
   (`coordinates <-> ST_SetSRID(ST_MakePoint(...), 4326)`), nearest first.
5. Applies `OFFSET` / `LIMIT` for pagination.

A spatial index (`spots_coordinates_idx`) on `coordinates` keeps these lookups
fast. The full queries live in `internal/spots/repository/sql.go`.

---

## API reference

Base path: **`/api/v1`**

### `GET /spots`

Return spots near a coordinate, ordered by distance.

**Query parameters**

| Param   | Type    | Required | Default | Description                                          |
| ------- | ------- | -------- | ------- | ---------------------------------------------------- |
| `lat`   | float   | yes      | —       | Latitude of the search point                         |
| `lon`   | float   | yes      | —       | Longitude of the search point                        |
| `range` | int     | no       | `50`    | Search range (bounding-box half-width, in meters)    |
| `page`  | int     | no       | `0`     | Page number (1-based for offset calculation)         |
| `size`  | int     | no       | `10`    | Page size (number of results per page)               |

**Example**

```bash
curl "http://localhost:5001/api/v1/spots?lat=51.51&lon=-0.12&range=500&page=1&size=2"
```

**Response** `200 OK`

```json
{
  "total_count": 17,
  "total_pages": 9,
  "page": 1,
  "size": 2,
  "has_more": true,
  "spots": [
    {
      "spot_id": "00277153-a04b-437c-80b5-8d305947a256",
      "spot_name": "Balthazar",
      "website": "balthazarlondon.com",
      "coordinates": {
        "latitude": 51.512,
        "longitude": -0.122
      },
      "description": "In the heart of Covent Garden, Balthazar is open all day...",
      "rating": 3.68
    }
  ]
}
```

> Values above are illustrative; actual results depend on the seeded dataset.

### `GET /health`

Liveness probe.

```bash
curl "http://localhost:5001/api/v1/health"
# {"status":"OK"}
```

---

## Configuration

Configuration is loaded by **viper**. The active file is selected by the
`config` environment variable:

- `config` unset / anything → `internal/config/config-local.yml`
- `config=docker` → `config-docker.yml` (used inside the container)

Environment variables can override individual values (`AutomaticEnv` is enabled).

| Setting              | Local (`config-local.yml`) | Docker (`config-docker.yml`) |
| -------------------- | -------------------------- | ---------------------------- |
| Server port          | `:5001`                    | `:5000`                      |
| Postgres host        | `localhost`                | `postgesql` (compose service)|
| Postgres db          | `vessels_db`               | `vessels_db`                 |
| Postgres user / pass | `postgres` / `postgres`    | `postgres` / `postgres`      |
| Log encoding         | `json`                     | `console`                    |

Database migrations (schema + seed data) are **run automatically on startup** by
`pkg/db/postgres`. They are read from `file:///migrations`, which is why the
Docker image copies the `migrations/` directory to `/migrations`.

---

## Getting started

### Prerequisites

- [Docker](https://www.docker.com/) & Docker Compose, **or**
- Go 1.18+ and a local PostgreSQL with the PostGIS extension.

### Run with Docker Compose (recommended)

This brings up both the API and a PostGIS database, builds the image, and runs
migrations on first boot.

```bash
make docker_start
# equivalent to:
# docker-compose -f docker-compose.yaml up --build
```

- API container port `5000` is published on the host as **`4000`**:
  → `http://localhost:4000/api/v1/health`
- PostGIS is published on host `5432`.

Stop and clean up:

```bash
make down-local   # stop & remove all running containers
make clean        # docker system prune
```

### Run locally

1. Start a PostGIS-enabled PostgreSQL. The quickest way is to use the database
   service from compose:

   ```bash
   docker-compose up postgesql
   ```

2. Make sure the database referenced in `internal/config/config-local.yml`
   exists (see [caveats](#notes--known-caveats) below).

3. Apply migrations. They run automatically at startup, but you can also use the
   [`migrate`](https://github.com/golang-migrate/migrate) CLI:

   ```bash
   make migrate_up
   ```

4. Run the API:

   ```bash
   make run          # go run ./cmd/main.go
   # or build a binary:
   make build        # go build ./cmd/main.go
   ```

   The server listens on `:5001` locally → `http://localhost:5001/api/v1/spots`.

---

## Make targets

| Command              | What it does                                            |
| -------------------- | ------------------------------------------------------- |
| `make run`           | Run the API (`go run ./cmd/main.go`)                    |
| `make build`         | Build the server binary                                 |
| `make test`          | Run tests with coverage (`go test -cover ./...`)        |
| `make run-linter`    | Run `golangci-lint`                                     |
| `make docker_start`  | Build & start the full stack via Docker Compose         |
| `make migrate_up`    | Apply DB migrations via the `migrate` CLI               |
| `make migrate_down`  | Roll back the last migration                            |
| `make version`       | Show current migration version                          |
| `make down-local`    | Stop & remove all containers                            |
| `make clean`         | `docker system prune`                                   |
| `make tidy`          | `go mod tidy` + vendor                                  |

> The `migrate` targets expect a local database named `spots_db` at
> `postgres://postgres:postgres@localhost:5432/spots_db`.

---

## Notes & known caveats

A few configuration inconsistencies exist in the repo that are worth being aware
of before running:

- **Database name mismatch.** The config files (`config-local.yml` /
  `config-docker.yml`) point at `vessels_db`, while `docker-compose.yaml`
  provisions a database named `spots_db`, and the `Makefile` migration targets
  also use `spots_db`. To run successfully, align these — e.g. set
  `PostgresqlDbname: spots_db` in the config, or create a `vessels_db` database.
- **Migrations path.** Auto-migrations read from the absolute path
  `file:///migrations`. This resolves correctly inside the Docker image (where
  `migrations/` is copied to `/migrations`) but not for a plain `go run` on the
  host — use the `make migrate_up` CLI target locally instead.
- **Seed size.** The seed migration inserts ~12,000 rows and is several MB, so
  the very first startup / migration can take a little while.
</content>
</invoke>
