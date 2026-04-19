import { api, stringifyParams } from "@/lib/api"
import type {
  PaginatedList,
  RecurringExpense,
  RecurringListParams,
  RecurringPayload,
} from "@/types/api"

export async function listRecurring(
  params: RecurringListParams
): Promise<PaginatedList<"recurring_expenses", RecurringExpense>> {
  const { data } = await api.get<PaginatedList<"recurring_expenses", RecurringExpense>>(
    "/recurring-expenses",
    {
      params: stringifyParams(params as Record<string, unknown>),
      paramsSerializer: { indexes: null },
    }
  )
  return data
}

export async function createRecurring(payload: RecurringPayload): Promise<RecurringExpense> {
  const { data } = await api.post<RecurringExpense>("/recurring-expenses", payload)
  return data
}

export async function updateRecurring(
  id: number,
  payload: RecurringPayload
): Promise<RecurringExpense> {
  const { data } = await api.put<RecurringExpense>(`/recurring-expenses/${id}`, payload)
  return data
}

export async function deleteRecurring(id: number): Promise<void> {
  await api.delete(`/recurring-expenses/${id}`)
}
