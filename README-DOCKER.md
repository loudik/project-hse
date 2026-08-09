# Running the HSE App with Docker Compose

## 1. Folder layout

`docker-compose.yml` sits at the project root, alongside the two project
folders:

```
PROJECT_HSE/
  docker-compose.yml
  .env
  hse-backend-go/
  hse-frontend-react/
```

## 2. Prepare .env

```bash
cp .env.example .env
```
Set `MYSQL_ROOT_PASSWORD` to a strong password.

Also make sure `hse-backend-go/.env` exists (`cp hse-backend-go/.env.example hse-backend-go/.env`
then fill in `JWT_SECRET` and the `MS_*` Microsoft OAuth values).

## 3. Run everything

```bash
docker compose up -d --build
```

This will:
1. Start MySQL (creates the `hse_db` database via `MYSQL_DATABASE`)
2. Run the **`migrate`** service - applies every migration in
   `hse-backend-go/migrations/` that hasn't run yet (tracked automatically,
   safe to run repeatedly - see `hse-backend-go/README.md` for details)
3. Start the Go backend, after `migrate` finishes successfully
4. Start the frontend (Nginx serving the React build + reverse-proxying
   `/api` to the backend)

You do **not** need to run migrations manually when using Docker - this
happens every time you `docker compose up`, whether it's the first time or
the fiftieth.

## 4. Access the app

Open `http://localhost` - the frontend calls `/api/...` which Nginx
proxies to the backend, so both are served from the same origin.

## Useful commands

```bash
docker compose logs -f              # all services
docker compose logs -f backend      # backend only
docker compose logs -f migrate      # check migration output/errors
docker compose down                 # stop everything
docker compose down -v              # stop + wipe MySQL data (full reset)
docker compose up -d --build backend   # rebuild & restart backend only
```

## Running a migration manually (without starting the whole stack)

```bash
docker compose run --rm migrate -path=/migrations -database "mysql://root:${MYSQL_ROOT_PASSWORD}@tcp(mysql:3306)/hse_db" up
```

## HTTPS / a real domain

This setup is HTTP-only on port 80. For production with a domain + SSL,
add a reverse proxy in front (e.g. Caddy or Nginx with Let's Encrypt) -
ask when you're ready for that step.
