import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import {
  bulkSetCategoryAllocations,
  createBudget,
  deleteBudget,
  deleteCategoryAllocation,
  getBudgetOverview,
  listBudgets,
  updateBudgetTotal,
  upsertCategoryAllocation,
  type BudgetListParams,
} from "@/features/budgets/api"
import type { ApiError } from "@/lib/api"
import { qk } from "@/lib/constants"

export function useBudgets(params: BudgetListParams = {}) {
  return useQuery({
    queryKey: qk.budgets(params),
    queryFn: () => listBudgets(params),
    placeholderData: keepPreviousData,
  })
}

export function useBudgetOverview(month: string | undefined) {
  return useQuery({
    queryKey: qk.budgetMonth(month ?? "none"),
    queryFn: () => getBudgetOverview(month as string),
    enabled: !!month,
    retry: (failureCount, error) => {
      const err = error as ApiError
      if (err.status === 404) return false
      return failureCount < 1
    },
  })
}

function invalidateBudgetFor(qc: ReturnType<typeof useQueryClient>, month?: string) {
  qc.invalidateQueries({ queryKey: qk.budgetsAll })
  if (month) qc.invalidateQueries({ queryKey: qk.budgetMonth(month) })
}

export function useCreateBudget() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: createBudget,
    onSuccess: (_d, vars) => {
      invalidateBudgetFor(qc, vars.budget_month)
      toast.success("Budget created")
    },
    onError: (err: ApiError) => toast.error(err.message),
  })
}

export function useUpdateBudgetTotal() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ month, total }: { month: string; total: number }) =>
      updateBudgetTotal(month, total),
    onSuccess: (_d, vars) => {
      invalidateBudgetFor(qc, vars.month)
      toast.success("Budget updated")
    },
    onError: (err: ApiError) => toast.error(err.message),
  })
}

export function useDeleteBudget() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (month: string) => deleteBudget(month),
    onSuccess: (_d, month) => {
      invalidateBudgetFor(qc, month)
      toast.success("Budget deleted")
    },
    onError: (err: ApiError) => toast.error(err.message),
  })
}

export function useUpsertAllocation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      month,
      ...input
    }: {
      month: string
      category_id: number
      allocated_amount: number
    }) => upsertCategoryAllocation(month, input),
    onSuccess: (_d, vars) => invalidateBudgetFor(qc, vars.month),
    onError: (err: ApiError) => toast.error(err.message),
  })
}

export function useBulkSetAllocations() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      month,
      categories,
    }: {
      month: string
      categories: Array<{ category_id: number; allocated_amount: number }>
    }) => bulkSetCategoryAllocations(month, categories),
    onSuccess: (_d, vars) => {
      invalidateBudgetFor(qc, vars.month)
      toast.success("Allocations saved")
    },
    onError: (err: ApiError) => toast.error(err.message),
  })
}

export function useDeleteAllocation() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ month, categoryId }: { month: string; categoryId: number }) =>
      deleteCategoryAllocation(month, categoryId),
    onSuccess: (_d, vars) => invalidateBudgetFor(qc, vars.month),
    onError: (err: ApiError) => toast.error(err.message),
  })
}
