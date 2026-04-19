CREATE UNIQUE INDEX IF NOT EXISTS uniq_expenses_recurring_date
    ON expenses (recurring_expense_id, expense_date)
    WHERE recurring_expense_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_recurring_expenses_due
    ON recurring_expenses (next_occurrence)
    WHERE is_active = TRUE;
