import type { ExpenseListParams, RecurringListParams } from "@/types/api"

export const TOKEN_STORAGE_KEY = "budgeting.token"

export const LOGOUT_EVENT = "budgeting:logout"

export const qk = {
  me: ["auth", "me"] as const,
  categories: ["categories"] as const,
  expenses: (params: ExpenseListParams) => ["expenses", params] as const,
  expensesAll: ["expenses"] as const,
  recurring: (params: RecurringListParams) => ["recurring", params] as const,
  recurringAll: ["recurring"] as const,
  budgets: (params?: { page?: number; page_size?: number; year?: number }) =>
    ["budgets", params ?? {}] as const,
  budgetsAll: ["budgets"] as const,
  budgetMonth: (month: string) => ["budget", month] as const,
}

export const FREQUENCY_LABELS: Record<string, string> = {
  daily: "Daily",
  weekly: "Weekly",
  monthly: "Monthly",
  yearly: "Yearly",
}

export const CATEGORY_COLORS = [
  "var(--chart-1)",
  "var(--chart-2)",
  "var(--chart-3)",
  "var(--chart-4)",
  "var(--chart-5)",
  "#8b5cf6",
  "#06b6d4",
  "#f43f5e",
  "#22c55e",
  "#f97316",
  "#6366f1",
  "#ec4899",
]

export function colorForIndex(i: number) {
  return CATEGORY_COLORS[i % CATEGORY_COLORS.length]
}
