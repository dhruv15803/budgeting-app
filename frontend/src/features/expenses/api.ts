import { api, stringifyParams } from "@/lib/api"
import type { Expense, ExpenseListParams, ExpensePayload, PaginatedList } from "@/types/api"

export async function listExpenses(
  params: ExpenseListParams
): Promise<PaginatedList<"expenses", Expense>> {
  const { data } = await api.get<PaginatedList<"expenses", Expense>>("/expenses", {
    params: stringifyParams(params as Record<string, unknown>),
    paramsSerializer: {
      indexes: null,
    },
  })
  return data
}

export async function createExpense(payload: ExpensePayload): Promise<Expense> {
  const { data } = await api.post<Expense>("/expenses", payload)
  return data
}

export async function updateExpense(id: number, payload: ExpensePayload): Promise<Expense> {
  const { data } = await api.put<Expense>(`/expenses/${id}`, payload)
  return data
}

export async function deleteExpense(id: number): Promise<void> {
  await api.delete(`/expenses/${id}`)
}
