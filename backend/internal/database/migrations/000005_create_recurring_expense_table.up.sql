

CREATE TYPE recurring_frequency AS ENUM ('daily','weekly','monthly','yearly');

CREATE TABLE IF NOT EXISTS recurring_expenses (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    amount DECIMAL NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES expense_categories(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE,
    frequency recurring_frequency NOT NULL DEFAULT 'monthly',
    next_occurrence DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);
