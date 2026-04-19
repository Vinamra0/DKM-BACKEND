# DKM-BACKEND

Backend service for DKM website.

## Prerequisites

- Go 1.22+
- MongoDB reachable from `MONGO_URI` (default: `mongodb://127.0.0.1:27017`)

## One-time Setup

```bash
# from repo root
go mod download

# create local env file
cp .env.example .env
```

Update [.env](.env) with your local values.

## Environment Variables

Main keys used by the app:

- `PORT`
- `BACKEND_URL`
- `MONGO_URI`
- `MONGO_DB`
- `JWT_SECRET`
- `ADMIN_EMAIL`
- `ADMIN_PASSWORD`
- `ADMIN_NAME`
- `ALLOWED_ORIGINS`
- `ADMIN_TOKEN` (optional dev fallback)

Template: [.env.example](.env.example)

## Run MongoDB

### Option A: System MongoDB (if installed)

```bash
sudo systemctl start mongod
sudo systemctl status mongod
```

### Option B: User-local MongoDB binary

```bash
mkdir -p ~/.local/mongodb/bin ~/.local/mongodb/data ~/.local/mongodb/logs

~/.local/mongodb/bin/mongod \
	--dbpath ~/.local/mongodb/data \
	--bind_ip 127.0.0.1 \
	--port 27017 \
	--logpath ~/.local/mongodb/logs/mongod.log \
	--fork
```

Stop local mongod:

```bash
pkill -f "mongod --dbpath ~/.local/mongodb/data"
```

## Run the Project

### Start API

```bash
go run cmd/api/main.go
```

### Seed Sample Data

```bash
go run cmd/seed/main.go
```

## Build Commands

```bash
# build API binary
go build -o bin/api ./cmd/api

# build seed binary
go build -o bin/seed ./cmd/seed
```

Run built binaries:

```bash
./bin/api
./bin/seed
```

## Verification Commands

```bash
# health
curl -i http://localhost:3001/api/health

# login
curl -s -X POST http://localhost:3001/api/auth/login \
	-H 'Content-Type: application/json' \
	-d '{"email":"admin@example.com","password":"change-this-local-password"}'

# public data
curl -s http://localhost:3001/api/blogs
curl -s http://localhost:3001/api/products
curl -s http://localhost:3001/api/careers/public
```

Check listening ports:

```bash
ss -ltn | rg ':3001|:27017'
```

## Development Commands

```bash
# run all tests
go test ./...

# format source files
gofmt -w cmd internal

# vet code
go vet ./...
```

## Common Management Commands

```bash
# restart API quickly
pkill -f "cmd/api/main.go" || true
go run cmd/api/main.go

# rerun seed
go run cmd/seed/main.go

# tail local MongoDB log (user-local mode)
tail -f ~/.local/mongodb/logs/mongod.log
```

## Troubleshooting

### Mongo connection refused

1. Ensure MongoDB is running on the URI in [.env](.env).
2. Prefer `127.0.0.1` over `localhost` to avoid IPv6 resolution issues on some systems.

### API starts but DB operations fail

`mongo.Connect` is lazy; startup can succeed before the first DB operation. Validate with seed or a DB-backed endpoint.
