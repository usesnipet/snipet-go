# Snipet

Go backend for orchestrating AI agents with multi-client support, conversation sessions, knowledge bases, and a configurable runtime.

Snipet exposes a REST API that lets external applications (websites, mobile apps, Discord, WhatsApp, etc.) integrate with agents featuring personas, language models, tools, and knowledge sources.

## Features

- **Multi-tenant by client** — each client has its own users, sessions, and configuration
- **Flexible authentication** — JWT, webhook, and anonymous users via configurable providers
- **Configurable agents** — persona, LLM models, tools, and knowledge base bindings
- **Knowledge** — management of static or semi-static data sources (documents, RAG, vector stores, etc.)
- **Sessions** — conversation history and dynamic state
- **API Keys** — authentication for internal services and integrations
- **Runtime pipeline** — architecture ready for planner, context builder, executor, and state manager (see [docs/ai.md](docs/ai.md))

## Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.26+ |
| HTTP | [chi](https://github.com/go-chi/chi) |
| ORM | GORM |
| Database | PostgreSQL 17 |
| Migrations | golang-migrate + Atlas |
| Auth | JWT, OIDC, webhooks |
| LLM | Google GenAI (provider) |

## Project structure

```
snipet/
├── cmd/api/              # Application entrypoint
├── config/               # Environment-based configuration
├── docs/                 # Architecture documentation
├── internal/
│   ├── api/              # HTTP helpers (parser, response, serve)
│   ├── auth/             # JWT, API keys, hashing
│   ├── bootstrap/        # Dependency wiring and routes
│   ├── infra/            # Database, cache
│   ├── model/            # GORM entities
│   ├── module/           # HTTP modules (handler, service, dto)
│   ├── provider/         # Provider abstractions (LLM)
│   ├── repository/       # Persistence layer
│   └── runtime/          # Source and index drivers
├── migrations/           # SQL migrations
├── providers/            # External provider implementations
└── docker-compose.yml
```

## Prerequisites

- Go 1.26+
- Docker and Docker Compose (for local PostgreSQL)
- [Atlas CLI](https://atlasgo.io/) (optional, to generate migrations from models)
- [Air](https://github.com/air-verse/air) (optional, hot reload in development)

## Quick start

### 1. Start the database

```bash
docker compose up -d postgres
```

PostgreSQL will be available at `localhost:8853` with default credentials (`postgres` / `postgres`, database `template`).

### 2. Configure environment variables

Create a `.env` file in the project root:

```env
ENV=development
LOG_LEVEL=debug

DB_URL=postgres://postgres:postgres@localhost:8853/template?sslmode=disable
DB_AUTO_MIGRATE=true

SERVER_PORT=8852

AUTH_JWT_SECRET=change-me-in-development
```

Variables are loaded automatically from a `.env` file found in the current directory or any parent directory.

### 3. Run the API

```bash
# Development with hot reload
make dev

# Or build + run manually
make build
./tmp/api
```

The API will be available at `http://localhost:8852/api`.

### 4. Full Docker Compose

To start both the API and PostgreSQL:

```bash
docker compose up -d
```

The API will be available at `http://localhost:8852`.

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ENV` | `development` | Runtime environment |
| `LOG_LEVEL` | `info` | Log level (`debug`, `info`, `warn`, `error`) |
| `SERVER_PORT` | `8852` | HTTP port |
| `SERVER_SHUTDOWN_TIMEOUT` | `15s` | Graceful shutdown timeout |
| `DB_URL` | — | PostgreSQL connection string |
| `DB_AUTO_MIGRATE` | `true` | Run migrations on startup |
| `DB_AUTO_CREATE` | `true` | Create the database if it does not exist |
| `DB_MIGRATION_DIR` | `migrations` | Migrations directory |
| `AUTH_JWT_SECRET` | — | Secret for JWT signing |
| `AUTH_JWT_EXPIRATION` | `15m` | JWT expiration time |
| `AUTH_JWT_ISSUER` | `https://snipet.cloud` | JWT issuer |
| `AUTH_JWT_AUDIENCE` | `https://snipet.cloud` | JWT audience |
| `DEV_PROXY` | `http://localhost:5173` | Development proxy for frontend |

## API

All routes are prefixed with `/api`.

### Authentication

| Method | Route | Description |
|--------|-------|-------------|
| `POST` | `/client/{client_code}/authenticate/{provider_name}` | Authenticate via provider (webhook, OIDC) |
| `POST` | `/client/{client_code}/authenticate/anonymous` | Create anonymous user session |

### Main resources

| Resource | Base route | Auth |
|----------|------------|------|
| Clients | `/client` | API Key |
| Users | `/client/{client_code}/user` | API Key / JWT |
| Agents | `/agent` | API Key |
| Sessions | `/client/{client_code}/session` | JWT / API Key |
| API Keys | `/api-key` | API Key |
| Knowledge | `/knowledge` | API Key |
| Knowledge Items | `/knowledge/{knowledge_id}/items` | API Key |
| Knowledge Index | `/knowledge/{knowledge_id}/index` | API Key |
| Indexed Items | `/knowledge/{knowledge_id}/index/{index_id}/items` | API Key |

For details on client integration and authentication flows, see [docs/client.md](docs/client.md).

## Development

### Make commands

```bash
make install      # Install dependencies
make test         # Run tests
make dev          # Hot reload with Air
make build        # Development build
make build-prod   # Optimized production build
make mocks        # Generate mocks with mockery
make db-generate <name>  # Generate migration via Atlas from GORM models
```

### Migrations

SQL migrations live in `migrations/` and are applied automatically on startup when `DB_AUTO_MIGRATE=true`.

To generate a new migration from models:

```bash
make db-generate add_new_feature
```

Requires Atlas CLI installed and Docker available (Atlas uses a temporary PostgreSQL container).

### Tests

```bash
go test ./...
```

### Mocks

Repository mocks are generated with [mockery](https://github.com/vektra/mockery) and live in `internal/repository/mocks/`.

```bash
make mocks
```

## Architecture

Snipet separates two core concepts:

- **Knowledge** — static or semi-static information queried by the AI (documents, RAG, vector stores, etc.)
- **State** — dynamic conversation data (history, summaries, variables, preferences)

The runtime executes each request through a pipeline: Planner → Context Builder → Executor → State Manager.

For the full AI architecture overview, see [docs/ai.md](docs/ai.md).

## License

This project is licensed under the [Apache License 2.0](LICENSE).
