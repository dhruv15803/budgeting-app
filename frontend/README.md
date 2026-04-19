# Budget — Frontend

A modern React 19 + TypeScript + Tailwind v4 + shadcn/ui frontend for the Go budgeting API.
It handles authentication, expenses, recurring expenses, monthly budgets, and a dashboard
with charts.

## Stack

- **Vite 8** + **React 19** + **TypeScript**
- **Tailwind CSS v4** (`@tailwindcss/vite`) + **shadcn/ui** (Radix primitives, New York style)
- **TanStack Query v5** + **Axios** for data fetching & caching
- **React Router v7** for routing
- **React Hook Form** + **Zod** for forms & validation
- **Recharts** (via shadcn Charts) for visualizations
- **framer-motion** for page transitions and list entrance animations
- **next-themes** + `sonner` for theming and toasts

## Getting started

### 1. Install dependencies

```bash
npm install
```

### 2. Configure the API URL

Copy `.env.example` → `.env` and point it at your Go backend:

```
VITE_API_URL=http://localhost:3000/api
```

The backend CORS allow-list already includes `http://localhost:5173`, which is this app's
default dev port.

### 3. Run the dev server

```bash
npm run dev
```

Open http://localhost:5173.

### 4. Build for production

```bash
npm run build
npm run preview
```

## Project layout

```
src/
  main.tsx                     # Providers: ErrorBoundary, Theme, QueryClient, Router, Auth, Toaster
  App.tsx                      # Route tree (public + protected + AppShell)
  index.css                    # Tailwind v4 import + shadcn CSS vars (light + dark)
  lib/
    api.ts                     # axios instance, JWT + envelope + 401 interceptors
    queryClient.ts             # TanStack Query defaults
    constants.ts               # query keys, chart palette
    utils.ts                   # cn(), formatCurrency/Date, monthKey helpers
  types/api.ts                 # Shared API types
  components/
    ui/                        # shadcn primitives (Button, Dialog, Table, Select, Chart, …)
    layout/                    # AppShell, Sidebar, Topbar, MobileNav, ThemeToggle, UserMenu
    common/                    # PageHeader, EmptyState, ErrorState, DataTablePager, CurrencyInput, CategoryBadge, ErrorBoundary
  features/
    auth/                      # login, register, verify-email, AuthProvider, RequireAuth
    categories/                # cached /expense-categories query
    expenses/                  # list, filters, create/edit/delete, pagination
    recurring/                 # recurring expenses + next-occurrence hints + active toggle
    budgets/                   # month list + month detail (overview, donut, allocations editor)
    dashboard/                 # KPI cards, 6-month trend, budget strip, recent expenses
    misc/                      # 404
```

## Feature tour

- **Auth** — JWT stored in `localStorage` (`budgeting.token`). Request interceptor attaches it;
  401 responses clear the token and broadcast a logout event that the `AuthProvider` listens for.
- **Envelope unwrapping** — the axios response interceptor reads `{ success, message, data }`,
  rejects failures with a typed `ApiError`, and unwraps `data` so consumers get clean payloads.
- **Expenses** — sticky filter bar (debounced search, multi-select categories, calendar range,
  amount min/max, sort). Paginated table with row actions, shared dialog for create/edit.
  Mutations invalidate expenses **and** budgets so totals refresh automatically.
- **Recurring** — inline active/paused switch, frequency badge, next-occurrence hint
  (`formatDistanceToNowStrict`). Same dialog pattern as expenses with start/end dates.
- **Budgets** — grid of month cards; month detail has editable total, donut chart, a bulk
  category allocations editor (PUT `/budgets/:month/categories`), and per-category progress
  rows with green → amber → red tones.
- **Dashboard** — greeting, KPI cards (spent, remaining, upcoming recurring, count),
  6-month trend chart (client-side via parallel `useQueries`), spending-by-category donut,
  budget progress strip, and a recent-expenses list. Empty-first-run CTA if no budget exists.

## Theming

- `next-themes` manages `class="dark"` on `<html>`. The toggle is in the top bar.
- Colors are driven entirely by CSS custom properties in `src/index.css`, so Recharts and
  shadcn components switch themes in lockstep.
- Respects `prefers-reduced-motion` (see `src/index.css`).

## Scripts

| Command            | What it does                         |
| ------------------ | ------------------------------------ |
| `npm run dev`      | Start the Vite dev server on `:5173` |
| `npm run build`    | Type-check + production build        |
| `npm run preview`  | Preview the production build         |
| `npm run lint`     | Run ESLint                           |
