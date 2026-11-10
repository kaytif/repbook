# Workout Tracker

A React frontend backed by the existing Go API and PostgreSQL. Supports accounts, workout history, exercise and program management, and individual sets with reps and weight.

## Run locally

From the repository root:

```powershell
docker compose up -d --build
cd frontend
npm.cmd install
npm.cmd run dev
```

Open the Vite URL (normally http://localhost:5173) and create an account. PostgreSQL initializes the checked-in schema automatically on a new Docker volume. The rebuilt app uses the `workout_tracker` database. A pre-existing volume needs that database created and the initial migration applied once; do not apply the initial migration to an already initialized database.

The frontend proxies `/api` to Go on port 8080, including session cookies. The same proxy works with `npm.cmd run preview` after a build. The Go container serves the API; run Vite separately for the frontend.

## Using the app

1. Create exercises in the Exercise library.
2. Optionally create a program and add exercises to it.
3. Create a workout with a date and optional program.
4. Add exercises individually or use **Add exercises from [program]**.
5. Log individual sets with reps and weight. Each set saves immediately.

Reopen a workout from session history to see its saved exercises and sets. Workouts, exercises, and program names can be edited. Sets support logging and deletion, matching the backend; there is no set-edit endpoint. Removing a workout exercise removes its sets. Workout deletion uses the backend's existing soft deletion. Exercises referenced by existing records cannot be deleted until those references are removed.

Use one consistent weight unit: the backend stores a number without a unit. Blank reps or weight are unrecorded (`null`); zero remains zero. No analytics, timers, goals, or other new backend features are included.

## Checks

```powershell
go test ./...
go vet ./...
$env:TEST_DATABASE_URL='postgres://workout:devpassword@127.0.0.1:5433/workout_tracker?sslmode=disable'
go test -run TestRouterIntegration -count=1 -v
cd frontend
npm.cmd run build
npm.cmd run lint
npm.cmd run test:e2e
```

The PostgreSQL integration test creates and removes an isolated schema; it does not use existing application tables. Browser tests use the running local backend and create separate test accounts. They cover desktop and mobile layouts using installed Microsoft Edge; adjust `channel` in `frontend/playwright.config.ts` for another installed Chromium browser. Playwright starts a temporary Vite server if port 5173 is free. Screenshots and failure traces go to the ignored `frontend/test-results` directory.

## External database

Set `DATABASE_URL` in the ignored root `.env` file to use an external PostgreSQL database such as Neon. Docker Compose reads it automatically. Run `docker compose up -d app` after changing it. Without that setting, the app uses local PostgreSQL. Initialize an empty external database with `migrations/001_initial_schema.up.sql` before first use. Existing local accounts and workouts remain in the local database and are not copied automatically.
