# iRegistro

Electronic School Register - Go Backend with Clean Architecture.

## Requirements

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and adjust values:
    ```sh
    cp .env.example .env
    ```
3. Initialize the database (ensure Postgres is running):
    ```sh
    go run cmd/migrate/main.go
    ```
4. Run the application:
    ```sh
    make run
    ```

## Development

- Run tests: `make test`
- Build binary: `make build`
- Run lint: `make lint`

## Architecture

Project follows Domain-Driven Design (DDD) and Clean Architecture.

- `cmd/`: Entry points
- `internal/domain/`: Enterprise business rules (Entities)
- `internal/application/`: Application business rules (Use Cases)
- `internal/presentation/`: Interface adapters (HTTP Handlers)
- `internal/infrastructure/`: Frameworks and drivers (DB, External Services)

## API endpoints

- `GET /health` - Health check
