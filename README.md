# Budgeting app

Full-stack personal budgeting application: a **Go** REST API with **PostgreSQL**, and a **React** SPA for expenses, recurring expenses, monthly budgets, and a dashboard with charts.

| Area        | Location     | Stack |
|------------|--------------|--------|
| API        | `backend/`   | Go (chi), JWT auth, sqlx, migrations |
| Web client | `frontend/`  | React 19, Vite 8, TypeScript, TanStack Query, Tailwind v4 |

HTTP API reference: [`backend/API.md`](backend/API.md).

---

## Prerequisites

- **Go** — use the toolchain version declared in [`backend/go.mod`](backend/go.mod).
- **Node.js** — **20+** recommended (for the Vite frontend).
- **PostgreSQL** — **14+** recommended; empty database for local dev.
- **[golang-migrate](https://github.com/golang-migrate/migrate)** CLI — used by the backend `Makefile` for `migrate-up` / `migrate-down` / `migrate-create`.

  Install examples:

  ```bash
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

  Ensure the `migrate` binary is on your `PATH`.

---

## Repository layout

```
budgeting-app/
├── backend/           # Go API (cmd/api, handlers, services, repos, migrations)
│   ├── .env.example   # template for backend/.env (copy and edit)
│   ├── Makefile       # run, migrate-up, migrate-down, migrate-create
│   └── API.md         # endpoint documentation
├── frontend/          # React + Vite SPA
│   ├── .env.example   # VITE_API_URL template
│   └── package.json   # npm scripts
└── README.md          # this file
```

---

## Backend setup

### 1. Create a PostgreSQL database

Create a database and user as you prefer, for example:

```sql
CREATE DATABASE budgeting;
```

### 2. Environment variables

The API loads `.env` from the **current working directory** (via `godotenv`). When using `make run`, run commands from `backend/` so a `backend/.env` file is picked up.

From `backend/`, copy the template and edit secrets and URLs:

```bash
cp .env.example .env
```

Reference for each variable (comments and defaults are in `.env.example`):

| Variable | Required | Description |
|----------|----------|-------------|
| `PORT` | yes | HTTP listen port (e.g. `3000` → server on `:3000`). |
| `DATABASE_URL` | yes | PostgreSQL URL, e.g. `postgres://user:pass@localhost:5432/budgeting?sslmode=disable`. |
| `JWT_SECRET` | yes | Secret for signing access tokens (use a long random string in production). |
| `EMAIL_VERIFICATION_BASE_URL` | yes | Base URL for verification links (e.g. `http://localhost:5173` or your app origin). |
| `JWT_EXPIRY` | no | Token lifetime (default `24h`). |
| `EMAIL_VERIFICATION_TOKEN_TTL` | no | Verification token TTL (default `48h`). |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `SMTP_FROM` | yes for email | Used to send verification email on registration. |
| `CRON_SCHEDULE` | no | Cron expression for recurring-expense generator (default `5 0 * * *` — 00:05 UTC daily). |
| `CORS_ALLOWED_ORIGINS` | no | Comma-separated origins; default includes `http://localhost:5173` and `http://localhost:3001`. |

For local email testing, point SMTP at a catcher such as [MailHog](https://github.com/mailhog/MailHog) or [Mailpit](https://github.com/axllent/mailpit) and set `SMTP_FROM` accordingly.

### 3. Run migrations

From `backend/` (so `Makefile` can `-include .env` and export `DATABASE_URL`):

```bash
cd backend
make migrate-up
```

To roll back one step:

```bash
make migrate-down
```

To add a new migration (requires `migrate` CLI):

```bash
make migrate-create your_migration_name
```

### 4. Run the API

From `backend/`:

```bash
make run
```

Equivalent:

```bash
go run ./cmd/api .
```

The server listens on `:{PORT}` (e.g. `http://localhost:3000`). Health check: `GET http://localhost:{PORT}/api/health`.

### Other backend commands

| Command | Description |
|---------|-------------|
| `go build -o bin/api ./cmd/api` | Build a binary (output path is your choice). |
| `go test ./...` | Run all Go tests. |

---

## Frontend setup

### 1. Install dependencies

```bash
cd frontend
npm install
```

### 2. Environment

Copy the example env file and adjust the API base URL if your backend port differs:

```bash
cp .env.example .env
```

`VITE_API_URL` must be the **full base path** of the API, including `/api`, for example:

```env
VITE_API_URL=http://localhost:3000/api
```

Ensure that origin (`http://localhost:5173` by default) is allowed by backend `CORS_ALLOWED_ORIGINS` (or use the backend defaults).

### 3. Development server

```bash
npm run dev
```

Default dev URL: **http://localhost:5173** (see [`frontend/vite.config.ts`](frontend/vite.config.ts)).

### Frontend npm scripts

| Script | Description |
|--------|-------------|
| `npm run dev` | Vite dev server with HMR. |
| `npm run build` | Typecheck (`tsc -b`) and production build to `frontend/dist`. |
| `npm run preview` | Serve the production build locally. |
| `npm run lint` | ESLint. |

More frontend-specific notes: [`frontend/README.md`](frontend/README.md).

---

## Typical local workflow

1. Start PostgreSQL.
2. Configure `backend/.env` and `frontend/.env`.
3. `cd backend && make migrate-up && make run`
4. In another terminal: `cd frontend && npm run dev`
5. Open the app URL, register (email verification depends on SMTP), then sign in and use expenses, recurring items, and budgets.

---

## Docker (API, frontend, PostgreSQL, Mailpit)

Requires [Docker Compose](https://docs.docker.com/compose/) (Compose V2).

```bash
cp docker-compose.env.example .env
# Edit `.env` (set a strong `JWT_SECRET` outside local experiments).
docker compose up --build
```

| Service    | URL |
|-----------|-----|
| Frontend  | http://localhost:5173 |
| API       | http://localhost:3000 |
| Mailpit UI (dev SMTP sink) | http://localhost:8025 |

The API container runs `golang-migrate` against `DATABASE_URL` on startup, then starts the server. PostgreSQL files live in the `postgres_data` volume. For Google sign-in, set `GOOGLE_OAUTH_CLIENT_ID` in `.env` and run `docker compose build frontend` (or `docker compose up --build`) so the value is baked into the static bundle.

If Compose fails to bind **port 3000** (for example because the API is already running via `make run`), stop that process or set **`API_PORT`** to a free port (e.g. `3001`) in `.env`, set **`VITE_API_URL`** to the same host URL (e.g. `http://localhost:3001/api`), then run **`docker compose up --build`** so the frontend image picks up the new API URL.

---

## Features (high level)

- **Auth** — Register, login, JWT-protected routes, email verification (SMTP).
- **Expenses** — CRUD, filters (search, date range, categories, amounts), pagination, `sort_by`: `date_desc`, `date_asc`, `amount_desc`, `amount_asc`.
- **Expense categories** — Listed from `/api/expense-categories` (public).
- **Recurring expenses** — CRUD; server cron generates matching expenses.
- **Monthly budgets** — Budgets per month and per category; dashboard charts consume the same API.

---

## Documentation

- **REST API** — [`backend/API.md`](backend/API.md) (paths, bodies, query params, examples).
