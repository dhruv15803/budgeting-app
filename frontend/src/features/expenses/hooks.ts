import { useMutation, useQuery, useQueryClient, keepPreviousData } from "@tanstack/react-query"
import { toast } from "sonner"
import { createExpense, deleteExpense, listExpenses, updateExpense } from "@/features/expenses/api"
import type { ApiError } from "@/lib/api"
import { qk } from "@/lib/constants"
import type { Expense, ExpenseListParams, ExpensePayload } from "@/types/api"

export function useExpenses(params: ExpenseListParams) {
  return useQuery({
    queryKey: qk.expenses(params),
    queryFn: () => listExpenses(params),
    placeholderData: keepPreviousData,
  })
}

function invalidateExpensesAndBudgets(qc: ReturnType<typeof useQueryClient>) {
  qc.invalidateQueries({ queryKey: qk.expensesAll })
  qc.invalidateQueries({ queryKey: qk.budgetsAll })
}

export function useCreateExpense() {
  const qc = useQueryClient()
  return useMutation<Expense, ApiError, ExpensePayload>({
    mutationFn: createExpense,
    onSuccess: () => {
      invalidateExpensesAndBudgets(qc)
      toast.success("Expense created")
    },
    onError: (err) => toast.error(err.message),
  })
}

export function useUpdateExpense() {
  const qc = useQueryClient()
  return useMutation<Expense, ApiError, { id: number; payload: ExpensePayload }>({
    mutationFn: ({ id, payload }) => updateExpense(id, payload),
    onSuccess: () => {
      invalidateExpensesAndBudgets(qc)
      toast.success("Expense updated")
    },
    onError: (err) => toast.error(err.message),
  })
}

export function useDeleteExpense() {
  const qc = useQueryClient()
  return useMutation<void, ApiError, number>({
    mutationFn: deleteExpense,
    onSuccess: () => {
      invalidateExpensesAndBudgets(qc)
      toast.success("Expense deleted")
    },
    onError: (err) => toast.error(err.message),
  })
}
