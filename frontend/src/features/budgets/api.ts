import { api, stringifyParams } from "@/lib/api"
import type {
  BudgetOverview,
  MonthlyBudget,
  PaginatedList,
} from "@/types/api"

export interface BudgetListParams {
  page?: number
  page_size?: number
  sort_by?: string
  year?: number
}

export async function listBudgets(
  params: BudgetListParams = {}
): Promise<PaginatedList<"budgets", MonthlyBudget>> {
  const { data } = await api.get<PaginatedList<"budgets", MonthlyBudget>>("/budgets", {
    params: stringifyParams(params as Record<string, unknown>),
  })
  return data
}

export async function createBudget(input: {
  budget_month: string
  total_amount: number
}): Promise<MonthlyBudget> {
  const { data } = await api.post<MonthlyBudget>("/budgets", input)
  return data
}

export async function getBudgetOverview(month: string): Promise<BudgetOverview> {
  const { data } = await api.get<BudgetOverview>(`/budgets/${month}`)
  return data
}

export async function updateBudgetTotal(
  month: string,
  total_amount: number
): Promise<MonthlyBudget> {
  const { data } = await api.put<MonthlyBudget>(`/budgets/${month}`, { total_amount })
  return data
}

export async function deleteBudget(month: string): Promise<void> {
  await api.delete(`/budgets/${month}`)
}

export async function upsertCategoryAllocation(
  month: string,
  input: { category_id: number; allocated_amount: number }
): Promise<{ id: number; category_id: number; allocated_amount: number }> {
  const { data } = await api.post<{ id: number; category_id: number; allocated_amount: number }>(
    `/budgets/${month}/categories`,
    input
  )
  return data
}

export async function bulkSetCategoryAllocations(
  month: string,
  categories: Array<{ category_id: number; allocated_amount: number }>
): Promise<void> {
  await api.put(`/budgets/${month}/categories`, { categories })
}

export async function deleteCategoryAllocation(
  month: string,
  categoryId: number
): Promise<void> {
  await api.delete(`/budgets/${month}/categories/${categoryId}`)
}
