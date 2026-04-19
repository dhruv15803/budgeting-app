


CREATE TABLE IF NOT EXISTS expenses ( 
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    amount DECIMAL NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES expense_categories(id) ON DELETE CASCADE,
    recurring_expense_id INTEGER REFERENCES recurring_expenses(id) ON DELETE CASCADE,
    expense_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);
