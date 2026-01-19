<h1 align="center"> Fandogh 🌰 </h1>

<p align="center">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/1995parham-teaching/fandogh/ci.yaml?logo=github&style=for-the-badge">
  <img alt="Codecov" src="https://img.shields.io/codecov/c/github/1995parham-teaching/fandogh?logo=codecov&style=for-the-badge">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/1995parham-teaching/fandogh?logo=go&style=for-the-badge">
</p>

Fandogh is a rental home listing platform that allows users to register, list properties for rent, and browse available homes. Property owners can create listings with photos, descriptions, and rental details, while buyers can search and view available properties.

## Features

- User registration and JWT-based authentication
- Create, update, and browse home listings
- Photo upload support with S3-compatible storage (MinIO/SeaweedFS)
- Role-based access control (owner/admin permissions)
- Pagination for listing queries
- Distributed tracing with OpenTelemetry and Jaeger
- Prometheus metrics for monitoring

## Tech Stack

- **Language:** Go
- **Web Framework:** Echo v4
- **Database:** MongoDB
- **File Storage:** MinIO (S3-compatible)
- **Authentication:** JWT
- **Dependency Injection:** Uber/fx
- **Logging:** Uber/zap
- **Tracing:** OpenTelemetry + Jaeger

## Getting Started

### Prerequisites

- Go 1.25+
- Docker and Docker Compose

### Running Locally

1. Start the infrastructure services:

```bash
docker compose -f deployments/docker-compose.yml up -d
```

2. Copy and configure the settings:

```bash
cp configs/config.example.yml configs/config.yml
```

3. Run database migrations:

```bash
go run ./cmd/fandogh migrate
```

4. Start the server:

```bash
go run ./cmd/fandogh server
```

The API will be available at `http://localhost:1378`.

### Using Docker

```bash
docker build -t fandogh -f build/package/Dockerfile .
docker run -p 1378:1378 fandogh
```

## Configuration

Configuration can be set via YAML file or environment variables (prefix: `FANDOGH_`).

| Setting             | Environment Variable        | Default                     | Description            |
| ------------------- | --------------------------- | --------------------------- | ---------------------- |
| `logger.level`      | `FANDOGH_LOGGER_LEVEL`      | `info`                      | Log level              |
| `database.url`      | `FANDOGH_DATABASE_URL`      | `mongodb://localhost:27017` | MongoDB connection URL |
| `database.name`     | `FANDOGH_DATABASE_NAME`     | `fandogh`                   | Database name          |
| `fs.endpoint`       | `FANDOGH_FS_ENDPOINT`       | `localhost:8333`            | MinIO endpoint         |
| `fs.access_key`     | `FANDOGH_FS_ACCESS_KEY`     | -                           | MinIO access key       |
| `fs.secret_key`     | `FANDOGH_FS_SECRET_KEY`     | -                           | MinIO secret key       |
| `jwt.access_secret` | `FANDOGH_JWT_ACCESS_SECRET` | -                           | JWT signing secret     |

## Project Structure

```
fandogh/
├── cmd/fandogh/          # Application entry point
├── internal/
│   ├── cmd/              # CLI commands (server, migrate)
│   ├── config/           # Configuration management
│   ├── db/               # MongoDB connection
│   ├── fs/               # File storage (MinIO)
│   ├── http/
│   │   ├── handler/      # HTTP request handlers
│   │   ├── request/      # Request DTOs with validation
│   │   ├── response/     # Response DTOs
│   │   ├── jwt/          # JWT generation/validation
│   │   └── server/       # Echo server setup
│   ├── model/            # Domain models (User, Home)
│   └── store/            # Data access layer
├── deployments/          # Docker Compose files
├── build/package/        # Dockerfile
└── configs/              # Configuration examples
```

## APIs

### Authentication

#### Register

```bash
curl 127.0.0.1:1378/register -X POST \
  -H 'Content-Type: application/json' \
  -d '{ "email": "user@example.com", "name": "John Doe", "password": "123456" }'
```

> Note: The first registered user automatically becomes an admin.

#### Login

```bash
curl 127.0.0.1:1378/login -X POST \
  -H 'Content-Type: application/json' \
  -d '{ "email": "user@example.com", "password": "123456" }'
```

Response:

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "Email": "user@example.com",
  "Name": "John Doe"
}
```

### Home Listings

All home endpoints require the `Authorization: Bearer <token>` header.

#### Create Home

```bash
curl 127.0.0.1:1378/api/homes -X POST \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Cozy Apartment",
    "location": "Rome, Italy",
    "description": "A beautiful place in the city center",
    "peoples": 3,
    "room": "living room",
    "bed": "single",
    "rooms": 2,
    "bathrooms": 1,
    "smoking": false,
    "guest": true,
    "pet": false,
    "bills_included": true,
    "contract": "1 year",
    "security_deposit": 1000,
    "price": 800
  }'
```

#### List Homes

```bash
curl '127.0.0.1:1378/api/homes?skip=0&limit=10' \
  -H 'Authorization: Bearer <token>'
```

Response:

```json
{
  "homes": [...],
  "total": 25,
  "skip": 0,
  "limit": 10
}
```

#### Get Home by ID

```bash
curl 127.0.0.1:1378/api/homes/<id> \
  -H 'Authorization: Bearer <token>'
```

#### Update Home

Only the owner or an admin can update a listing.

```bash
curl 127.0.0.1:1378/api/homes/<id> -X PUT \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{ "price": 900 }'
```

### Health Check

```bash
curl 127.0.0.1:1378/healthz
```

## Development

### Running Tests

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

## Observability

- **Metrics:** Prometheus metrics available at port 8080 (when enabled)
- **Tracing:** Jaeger UI available at `http://localhost:16686` (when running with Docker Compose)
