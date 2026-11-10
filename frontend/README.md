# Workout Tracker frontend

React + TypeScript + Vite, connected to the existing Go API. See the [project README](../README.md) for setup, supported workflows, and database notes.

```powershell
npm.cmd install
npm.cmd run dev
```

Start the backend with `docker compose up -d --build` from the project root first. Create an account in the browser. Both the dev server and built preview proxy `/api` to port 8080.

```powershell
npm.cmd run build
npm.cmd run lint
npm.cmd run test:e2e
```

Browser tests use installed Microsoft Edge and check desktop and mobile flows. They create separate test accounts on the running local backend. Test screenshots and failure traces are saved in the ignored `test-results` directory.
