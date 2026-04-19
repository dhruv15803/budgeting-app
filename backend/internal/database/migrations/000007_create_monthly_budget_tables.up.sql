

CREATE TABLE IF NOT EXISTS monthly_budgets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    budget_month DATE NOT NULL CHECK (budget_month = date_trunc('month', budget_month)::date),
    total_amount DECIMAL NOT NULL CHECK (total_amount >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    UNIQUE (user_id, budget_month)
);

CREATE INDEX IF NOT EXISTS idx_monthly_budgets_user_id ON monthly_budgets (user_id);

CREATE TABLE IF NOT EXISTS monthly_category_budgets (
    id SERIAL PRIMARY KEY,
    monthly_budget_id INTEGER NOT NULL REFERENCES monthly_budgets(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES expense_categories(id) ON DELETE CASCADE,
    allocated_amount DECIMAL NOT NULL CHECK (allocated_amount >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    UNIQUE (monthly_budget_id, category_id)
);

CREATE INDEX IF NOT EXISTS idx_monthly_category_budgets_monthly_budget_id ON monthly_category_budgets (monthly_budget_id);
