# HSE Backend (Go + Gin + database/sql + MySQL)

## 1. Prerequisites
- Go 1.22+ (`go version` to check)
- MySQL server running
- Docker (recommended for running migrations - see below)

## 2. Install Go dependencies

```bash
cd hse-backend-go
go mod tidy
```

## 3. Database migrations (golang-migrate)

Migrations live in `migrations/`, named like `000001_description.up.sql` /
`000001_description.down.sql`. This project uses the
[golang-migrate](https://github.com/golang-migrate/migrate) tool - it keeps
track of which migrations have already run (in a `schema_migrations` table),
so you never have to remember that manually again.

### Option A - Docker (recommended, no local install needed)

Migrations run **automatically** every time you do:
```bash
docker compose up -d --build
```
A one-shot `migrate` service applies any pending migrations before `backend`
starts. You don't need to do anything else - this works the same whether
it's the first time (fresh database) or the 20th time (a few new
migrations since last time).

To run migrations manually via Docker without starting the whole stack:
```bash
docker compose run --rm migrate -path=/migrations -database "mysql://root:${MYSQL_ROOT_PASSWORD}@tcp(mysql:3306)/hse_db" up
```

### Option B - migrate CLI installed locally

Install the CLI once:
```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migrations against your local MySQL:
```bash
migrate -path hse-backend-go/migrations -database "mysql://root:YOUR_PASSWORD@tcp(127.0.0.1:3306)/hse_db" up
```

Other useful commands:
```bash
# Roll back the last migration
migrate -path hse-backend-go/migrations -database "mysql://..." down 1

# Check current version
migrate -path hse-backend-go/migrations -database "mysql://..." version

# Create a new empty migration pair (up + down)
migrate create -ext sql -dir hse-backend-go/migrations -seq add_something_new
```

## 4. Configure connection

```bash
cp .env.example .env
```
Edit `DB_DSN` and the other values with your real credentials.

## 5. Run the server

```bash
go run main.go
```
Runs at `http://localhost:3000` - matches the `/api` proxy in the
frontend's `vite.config.js`.

## Structure

```
hse-backend-go/
  main.go                     -> entry point
  db/db.go                    -> MySQL connection
  models/                     -> data structs
  handlers/                   -> endpoint logic (plain SQL, no ORM)
  middleware/auth.go          -> JWT auth middleware
  utils/jwt.go                -> JWT generate/parse
  migrations/                 -> golang-migrate up/down pairs
  .env.example
  go.mod
```

## Main endpoints

| Method | Path                               | Purpose                        |
|--------|--------------------------------------|---------------------------------|
| POST   | /api/sign-in                        | Login                          |
| POST   | /api/sign-up                        | Public registration            |
| GET    | /api/verify-email                   | Activate a Pending account     |
| POST   | /api/oauth/microsoft                | Microsoft login callback       |
| GET    | /api/menu                           | Menu tree for the logged-in user |
| POST   | /api/organizations                  | Create organization profile    |
| GET    | /api/organizations                  | List (for ANP HSE review)      |
| PATCH  | /api/organizations/:id/decision     | Approve/reject organization    |
| GET/POST/PATCH | /api/hse/submissions/...    | HSE certification submissions  |
