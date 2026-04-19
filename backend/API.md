# Budgeting App — Backend API Reference

> **Base URL:** `http://localhost:3000/api`  
> All dates are `YYYY-MM-DD` strings. All month parameters are `YYYY-MM` strings.  
> All responses are JSON with `Content-Type: application/json`.

---

## Table of Contents

1. [Response Envelope](#response-envelope)
2. [Authentication](#authentication)
3. [Expense Categories](#expense-categories)
4. [Expenses](#expenses)
5. [Recurring Expenses](#recurring-expenses)
6. [Budgets](#budgets)
7. [Error Reference](#error-reference)

---

## Response Envelope

Every response follows one of two shapes:

**Success**
```json
{
  "success": true,
  "message": "human-readable message",
  "data": { ... }
}
```

**Error**
```json
{
  "success": false,
  "message": "what went wrong",
  "status_code": 400
}
```

---

## Authentication

All protected endpoints require a Bearer token in the `Authorization` header:

```
Authorization: Bearer <token>
```

Tokens are issued by `/auth/login` and `/auth/verify-email`. There is no refresh-token mechanism; re-login when the token expires.

---

### `POST /auth/register`

Create a new user account. A verification email is sent; the account cannot log in until verified.

**Request body**
```json
{
  "email": "user@example.com",
  "password": "secret123",
  "username": "alice"
}
```

| Field | Type | Required | Notes |
|---|---|---|---|
| `email` | string | yes | Must be unique |
| `password` | string | yes | |
| `username` | string | no | |

**Response `201`**
```json
{
  "success": true,
  "message": "Registration successful. Check your email to verify your account.",
  "data": null
}
```

**Errors**

| Status | Condition |
|---|---|
| `409 Conflict` | Email already registered |
| `400 Bad Request` | Validation failure |

---

### `POST /auth/login`

Authenticate and receive a JWT.

**Request body**
```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

**Response `200`**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "<jwt>"
  }
}
```

**Errors**

| Status | Condition |
|---|---|
| `401 Unauthorized` | Wrong email or password |
| `403 Forbidden` | Email not verified |

---

### `GET /auth/verify-email?token=<token>`

Verify a user's email address via the link emailed on registration. Returns a JWT so the user is immediately logged in after verification.

**Query parameter:** `token` (string, required)

**Response `200`**
```json
{
  "success": true,
  "message": "Email verified",
  "data": {
    "token": "<jwt>"
  }
}
```

**Errors**

| Status | Condition |
|---|---|
| `400 Bad Request` | Token invalid or expired |

---

### `GET /auth/me` 🔒

Return the authenticated user's profile.

**Response `200`**
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "username": "alice",
    "image_url": null,
    "role": "user",
    "created_at": "2026-01-15T10:00:00Z",
    "updated_at": null
  }
}
```

---

## Expense Categories

### `GET /api/expense-categories`

Public endpoint — no authentication required. Returns all categories sorted alphabetically.

**Response `200`**
```json
{
  "success": true,
  "message": "categories fetched successfully",
  "categories": [
    { "id": 1, "category_name": "Food & Dining", "created_at": "...", "updated_at": null },
    { "id": 2, "category_name": "Housing",        "created_at": "...", "updated_at": null }
  ]
}
```

> **Tip:** Fetch this once on app load and cache it. Use the `id` values when creating or filtering expenses.

---

## Expenses

All expense endpoints require authentication (`Authorization: Bearer <token>`).

### Expense object

```json
{
  "id": 42,
  "title": "Grocery run",
  "description": "Weekly shop",
  "amount": 85.50,
  "user_id": 1,
  "category_id": 3,
  "recurring_expense_id": null,
  "expense_date": "2026-04-15",
  "created_at": "2026-04-15T12:00:00Z",
  "updated_at": null
}
```

`recurring_expense_id` is non-null only for expenses generated automatically by the recurring expense scheduler.

---

### `POST /api/expenses` 🔒

Create a new expense.

**Request body**
```json
{
  "title": "Grocery run",
  "description": "Weekly shop",
  "amount": 85.50,
  "category_id": 3,
  "expense_date": "2026-04-15"
}
```

| Field | Type | Required | Notes |
|---|---|---|---|
| `title` | string | yes | |
| `description` | string | no | |
| `amount` | number | yes | Must be > 0 |
| `category_id` | integer | yes | Must exist in `expense_categories` |
| `expense_date` | string | yes | `YYYY-MM-DD` |

**Response `201`**
```json
{
  "success": true,
  "message": "Expense created successfully",
  "data": { ...expense object... }
}
```

---

### `GET /api/expenses` 🔒

List the authenticated user's expenses with optional filtering and pagination.

**Query parameters**

| Parameter | Type | Default | Description |
|---|---|---|---|
| `page` | integer | `1` | Page number |
| `page_size` | integer | `20` | Items per page |
| `sort_by` | string | `expense_date DESC` | Column + direction, e.g. `amount ASC`, `title ASC`, `expense_date ASC` |
| `search` | string | — | Full-text search on `title` and `description` (case-insensitive) |
| `date_from` | string | — | `YYYY-MM-DD` — start of range (inclusive) |
| `date_to` | string | — | `YYYY-MM-DD` — end of range (inclusive) |
| `month` | string | — | `YYYY-MM` shorthand — expands to first/last day of the month. Ignored if `date_from`/`date_to` are set |
| `year` | string | — | `YYYY` shorthand — expands to Jan 1 – Dec 31. Ignored if `date_from`/`date_to` or `month` are set |
| `category_id` | integer | — | Repeatable: `?category_id=1&category_id=3` |
| `amount_min` | number | — | Minimum amount (inclusive) |
| `amount_max` | number | — | Maximum amount (inclusive) |

**Response `200`**
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "expenses": [ ...expense objects... ],
    "total": 87,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

**Example requests**
```
GET /api/expenses?month=2026-04
GET /api/expenses?year=2026&category_id=1&category_id=5
GET /api/expenses?search=grocery&amount_min=10&amount_max=200
GET /api/expenses?date_from=2026-01-01&date_to=2026-03-31&sort_by=amount+DESC
```

---

### `PUT /api/expenses/{id}` 🔒

Replace all fields on an existing expense. Only the owner can update.

**URL parameter:** `id` (integer)

**Request body** — same shape as create:
```json
{
  "title": "Grocery run",
  "description": "Updated note",
  "amount": 90.00,
  "category_id": 3,
  "expense_date": "2026-04-15"
}
```

**Response `200`**
```json
{
  "success": true,
  "message": "Expense updated successfully",
  "data": { ...expense object... }
}
```

**Errors**

| Status | Condition |
|---|---|
| `404 Not Found` | Expense does not exist |
| `403 Forbidden` | Expense belongs to another user |

---

### `DELETE /api/expenses/{id}` 🔒

Delete an expense. Only the owner can delete.

**URL parameter:** `id` (integer)

**Response `200`**
```json
{
  "success": true,
  "message": "Expense deleted successfully"
}
```

**Errors**

| Status | Condition |
|---|---|
| `404 Not Found` | Expense does not exist |
| `403 Forbidden` | Expense belongs to another user |

---

## Recurring Expenses

All recurring expense endpoints require authentication. The scheduler automatically creates a matching entry in `expenses` each time `next_occurrence` is reached (daily cron, with catch-up on startup).

### Recurring expense object

```json
{
  "id": 7,
  "title": "Netflix",
  "description": null,
  "amount": 15.99,
  "user_id": 1,
  "category_id": 5,
  "start_date": "2026-01-01",
  "end_date": null,
  "frequency": "monthly",
  "next_occurrence": "2026-05-01",
  "is_active": true,
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": null
}
```

**Frequency values:** `daily` · `weekly` · `monthly` · `yearly`

---

### `POST /api/recurring-expenses` 🔒

**Request body**
```json
{
  "title": "Netflix",
  "description": null,
  "amount": 15.99,
  "category_id": 5,
  "start_date": "2026-01-01",
  "end_date": null,
  "frequency": "monthly"
}
```

| Field | Type | Required | Notes |
|---|---|---|---|
| `title` | string | yes | |
| `description` | string | no | |
| `amount` | number | yes | Must be > 0 |
| `category_id` | integer | yes | |
| `start_date` | string | yes | `YYYY-MM-DD` |
| `end_date` | string | no | `YYYY-MM-DD` or `null` — leave null for indefinite |
| `frequency` | string | yes | `daily` \| `weekly` \| `monthly` \| `yearly` |

**Response `201`**
```json
{
  "success": true,
  "message": "Recurring expense created successfully",
  "data": { ...recurring expense object... }
}
```

---

### `GET /api/recurring-expenses` 🔒

**Query parameters**

| Parameter | Type | Default | Description |
|---|---|---|---|
| `page` | integer | `1` | |
| `page_size` | integer | `20` | |
| `sort_by` | string | `next_occurrence ASC` | |
| `search` | string | — | Search title/description |
| `category_id` | integer | — | Repeatable |
| `frequency` | string | — | Filter by frequency value |
| `is_active` | boolean | — | `true` or `false` |

**Response `200`**
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "recurring_expenses": [ ...recurring expense objects... ],
    "total": 12,
    "page": 1,
    "page_size": 20,
    "total_pages": 1
  }
}
```

---

### `GET /api/recurring-expenses/{id}` 🔒

Get a single recurring expense by ID.

**Response `200`**
```json
{
  "success": true,
  "message": "ok",
  "data": { ...recurring expense object... }
}
```

**Errors:** `404` not found, `403` forbidden.

---

### `PUT /api/recurring-expenses/{id}` 🔒

Update a recurring expense. Changing `start_date` or `frequency` recalculates `next_occurrence`. Use `is_active: false` to pause the schedule without deleting.

**Request body**
```json
{
  "title": "Netflix Premium",
  "description": null,
  "amount": 22.99,
  "category_id": 5,
  "start_date": "2026-01-01",
  "end_date": null,
  "frequency": "monthly",
  "is_active": true
}
```

`is_active` defaults to `true` if omitted.

**Response `200`**
```json
{
  "success": true,
  "message": "Recurring expense updated successfully",
  "data": { ...recurring expense object... }
}
```

---

### `DELETE /api/recurring-expenses/{id}` 🔒

Delete a recurring expense (does **not** delete already-generated expense records).

**Response `200`**
```json
{
  "success": true,
  "message": "Recurring expense deleted successfully"
}
```

---

## Budgets

All budget endpoints require authentication. Budgets are per-user and per-month. A budget can optionally have per-category allocations.

### Monthly budget object

```json
{
  "id": 3,
  "budget_month": "2026-04",
  "total_amount": 3000.00,
  "created_at": "2026-04-01T00:00:00Z",
  "updated_at": null
}
```

### Budget overview object (returned by `GET /api/budgets/{month}`)

```json
{
  "id": 3,
  "budget_month": "2026-04",
  "total_amount": 3000.00,
  "total_spent": 1240.50,
  "remaining": 1759.50,
  "categories": [
    {
      "id": 11,
      "category_id": 1,
      "category_name": "Food & Dining",
      "allocated_amount": 500.00,
      "spent_amount": 312.75,
      "remaining": 187.25
    }
  ],
  "created_at": "2026-04-01T00:00:00Z",
  "updated_at": null
}
```

---

### `POST /api/budgets` 🔒

Create a monthly budget. Only one budget per month per user is allowed.

**Request body**
```json
{
  "budget_month": "2026-04",
  "total_amount": 3000.00
}
```

**Response `201`**
```json
{
  "success": true,
  "message": "Budget created successfully",
  "data": { ...monthly budget object... }
}
```

**Errors:** `400` if the month already has a budget.

---

### `GET /api/budgets` 🔒

List all of the user's budgets.

**Query parameters**

| Parameter | Type | Default | Description |
|---|---|---|---|
| `page` | integer | `1` | |
| `page_size` | integer | `20` | |
| `sort_by` | string | `budget_month DESC` | |
| `year` | integer | — | Filter to a specific calendar year |

**Response `200`**
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "budgets": [ ...monthly budget objects... ],
    "total": 6,
    "page": 1,
    "page_size": 20,
    "total_pages": 1
  }
}
```

---

### `GET /api/budgets/{month}` 🔒

Get the full budget overview for a month — total allocated vs. spent, plus per-category breakdown.

**URL parameter:** `month` — `YYYY-MM`

**Response `200`**
```json
{
  "success": true,
  "message": "ok",
  "data": { ...budget overview object... }
}
```

**Errors:** `404` if no budget exists for that month.

---

### `PUT /api/budgets/{month}` 🔒

Update the total budget amount for a month.

**URL parameter:** `month` — `YYYY-MM`

**Request body**
```json
{
  "total_amount": 3500.00
}
```

**Response `200`**
```json
{
  "success": true,
  "message": "Budget updated successfully",
  "data": { ...monthly budget object... }
}
```

---

### `DELETE /api/budgets/{month}` 🔒

Delete a monthly budget and all its category allocations.

**URL parameter:** `month` — `YYYY-MM`

**Response `200`**
```json
{
  "success": true,
  "message": "Budget deleted successfully"
}
```

---

### `POST /api/budgets/{month}/categories` 🔒

Set (or update) the allocation for a single category within a monthly budget. This is an **upsert** — safe to call multiple times.

**URL parameter:** `month` — `YYYY-MM`

**Request body**
```json
{
  "category_id": 1,
  "allocated_amount": 500.00
}
```

**Response `200`**
```json
{
  "success": true,
  "message": "Category budget set successfully",
  "data": {
    "id": 11,
    "category_id": 1,
    "allocated_amount": 500.00
  }
}
```

**Errors:** `404` if no budget exists for that month.

---

### `PUT /api/budgets/{month}/categories` 🔒

**Bulk replace** all category allocations for a month in a single atomic transaction. Existing allocations not in the payload are deleted; entries in the payload are upserted.

**URL parameter:** `month` — `YYYY-MM`

**Request body**
```json
{
  "categories": [
    { "category_id": 1, "allocated_amount": 500.00 },
    { "category_id": 2, "allocated_amount": 800.00 },
    { "category_id": 5, "allocated_amount": 200.00 }
  ]
}
```

**Response `200`**
```json
{
  "success": true,
  "message": "Category budgets updated successfully"
}
```

---

### `DELETE /api/budgets/{month}/categories/{category_id}` 🔒

Remove the allocation for a single category from a monthly budget.

**URL parameters:** `month` — `YYYY-MM`, `category_id` — integer

**Response `200`**
```json
{
  "success": true,
  "message": "Category budget removed successfully"
}
```

---

## Error Reference

| HTTP Status | Meaning |
|---|---|
| `400 Bad Request` | Malformed JSON, missing required field, invalid date format, business rule violation |
| `401 Unauthorized` | Missing or invalid JWT |
| `403 Forbidden` | Valid JWT but resource belongs to another user, or email not verified |
| `404 Not Found` | Resource does not exist |
| `409 Conflict` | Duplicate resource (e.g. email already registered, budget already exists for that month) |
| `500 Internal Server Error` | Unexpected server-side failure |

---

## Quick-Start Checklist for the Frontend

1. **On app load** — call `GET /api/expense-categories` and store the list locally (used as dropdown data for create/filter forms).
2. **Auth flow** — `POST /auth/register` → user verifies email → `POST /auth/login` → store JWT in memory or `localStorage`.
3. **Attach JWT** — every request to a 🔒 endpoint needs `Authorization: Bearer <token>`.
4. **Pagination** — use `page` / `page_size` query params; read `total_pages` from the response to render a pager.
5. **Month dashboard** — call `GET /api/budgets/YYYY-MM` for the full spending overview; it already includes per-category breakdowns.
6. **CORS** — the server allows `http://localhost:5173` and `http://localhost:3001` by default. Override with the `CORS_ALLOWED_ORIGINS` environment variable in production (comma-separated).
