# TK-BE

A backend service for a Project Management platform.

This project provides APIs for authentication, user profiles, organization management, and organization membership workflows. It is designed with a layered/hexagonal architecture so business logic stays independent from frameworks and infrastructure.

## Project Overview

TK-BE is focused on multi-organization project management scenarios, where users can:

- Sign in using basic login or Google social login.
- Manage their own profile and primary organization.
- Create and manage organizations.
- Invite and manage organization members.

The codebase is organized to keep domain logic clean, testable, and portable.

## Tech Stack

- Language: Go (Go 1.25+)
- HTTP Framework: Fiber v2
- Database: PostgreSQL (via GORM)
- Cache: Redis
- Authentication: JWT + Google token verification
- Configuration: Environment variables (`.env`) via `godotenv`
- Containerization: Docker + Docker Compose

## Architecture

The project follows a clean/hexagonal style split into:

- `core/domain`: domain entities and request models
- `core/ports`: interfaces for inbound and outbound dependencies
- `core/services`: use cases / business logic
- `internal/adapters`: framework-specific and infrastructure adapters
	- HTTP handlers and routes
	- middleware
	- external services (Google verifier)
	- storage adapters (repositories, ORM, cache)

This separation makes it easier to:

- Change frameworks or storage implementations with minimal impact on business rules.
- Keep use cases independent from transport (HTTP) and persistence details.

## Project Structure

```text
cmd/
	main.go                       # application bootstrap and dependency wiring

internal/
	adapters/
		config/                     # environment config loader
		google/                     # Google token verifier adapter
		middleware/                 # JWT middleware
		routes/                     # Fiber route registration + HTTP handlers
		storage/
			cache/                    # Redis adapter
			orm/                      # GORM/PostgreSQL setup + models
			repository/               # repository implementations (output adapters)
	core/
		domain/                     # domain models
		ports/                      # input/output interfaces
		services/                   # business services (use cases)
	utils/                        # shared utilities (jwt, bcrypt, struct helpers)

docker/
	Dockerfile
	docker-compose.yml
```

## Running the Project

### 1. Prerequisites

- Go 1.25+
- PostgreSQL
- Redis
- Docker (optional, recommended for local setup)

### 2. Environment Variables

Create a `.env` file at the project root.

Minimum commonly-used variables:

```env
APP_NAME=TK-BE
APP_VERSION=1.0.0
API_PORT=3000
ALLOWED_CREDENTIAL_ORIGINS=*

JWT_SECRET_KEY=change_me
JWT_EXPIRE_DAYS_COUNT=7
JWT_ISSUER=tk-be
JWT_SUBJECT=auth
JWT_SIGNING_METHOD=HS256

POSTGRE_URI=postgres://aprilpollo:aprilpollo@localhost:5432/aprilpollo?sslmode=disable
POSTGRE_URI_MIGRATION=postgres://aprilpollo:aprilpollo@localhost:5432/aprilpollo?sslmode=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

OAUTH_GOOGLE_CLIENT_ID=your_google_client_id
```

### 3. Run with Docker Compose

From project root:

```bash
docker compose -f docker/docker-compose.yml up --build
```

The API will be available at `http://localhost:3000`.

### 4. Run Locally (without Docker)

```bash
go mod download
go run ./cmd/main.go
```

## Development Notes

- Keep business rules inside `core/services` and `core/domain`.
- Add new external integrations under `internal/adapters`.
- Prefer depending on interfaces in `core/ports` for better testability.

## License

[MIT](https://choosealicense.com/licenses/mit/) © 2026 p.phonsing_