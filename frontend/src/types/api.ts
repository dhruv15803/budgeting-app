export interface ApiSuccess<T> {
  success: true
  message: string
  data: T
}

export interface ApiFailure {
  success: false
  message: string
  status_code?: number
}

export type ApiResponse<T> = ApiSuccess<T> | ApiFailure

export interface PageMeta {
  total: number
  page: number
  page_size: number
  total_pages: number
}
export type PaginatedList<TKey extends string, TItem> = PageMeta & {
  [K in TKey]: TItem[]
}

export interface User {
  id: number
  email: string
  username: string | null
  image_url: string | null
  role: string
  created_at: string
  updated_at: string | null
}

export interface Category {
  id: number
  category_name: string
  created_at: string
  updated_at: string | null
}

export interface Expense {
  id: number
  title: string
  description: string | null
  amount: number
  user_id: number
  category_id: number
  recurring_expense_id: number | null
  expense_date: string
  created_at: string
  updated_at: string | null
}

export type Frequency = "daily" | "weekly" | "monthly" | "yearly"

export interface RecurringExpense {
  id: number
  title: string
  description: string | null
  amount: number
  user_id: number
  category_id: number
  start_date: string
  end_date: string | null
  frequency: Frequency
  next_occurrence: string
  is_active: boolean
  created_at: string
  updated_at: string | null
}

export interface MonthlyBudget {
  id: number
  budget_month: string
  total_amount: number
  created_at: string
  updated_at: string | null
}

export interface BudgetCategory {
  id: number
  category_id: number
  category_name: string
  allocated_amount: number
  spent_amount: number
  remaining: number
}

export interface BudgetOverview extends MonthlyBudget {
  total_spent: number
  remaining: number
  categories: BudgetCategory[]
}

/** Values accepted by `GET /api/expenses?sort_by=…` */
export type ExpenseSortBy = "date_desc" | "date_asc" | "amount_desc" | "amount_asc"

export interface ExpenseListParams {
  page?: number
  page_size?: number
  sort_by?: ExpenseSortBy
  search?: string
  date_from?: string
  date_to?: string
  month?: string
  year?: string
  category_id?: number[]
  amount_min?: number
  amount_max?: number
}

export interface RecurringListParams {
  page?: number
  page_size?: number
  sort_by?: string
  search?: string
  category_id?: number[]
  frequency?: Frequency
  is_active?: boolean
}

export interface ExpensePayload {
  title: string
  description?: string | null
  amount: number
  category_id: number
  expense_date: string
}

export interface RecurringPayload {
  title: string
  description?: string | null
  amount: number
  category_id: number
  start_date: string
  end_date?: string | null
  frequency: Frequency
  is_active?: boolean
}
