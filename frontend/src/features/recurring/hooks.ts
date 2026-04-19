import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import {
  createRecurring,
  deleteRecurring,
  listRecurring,
  updateRecurring,
} from "@/features/recurring/api"
import type { ApiError } from "@/lib/api"
import { qk } from "@/lib/constants"
import type { RecurringExpense, RecurringListParams, RecurringPayload } from "@/types/api"

export function useRecurring(params: RecurringListParams) {
  return useQuery({
    queryKey: qk.recurring(params),
    queryFn: () => listRecurring(params),
    placeholderData: keepPreviousData,
  })
}

export function useCreateRecurring() {
  const qc = useQueryClient()
  return useMutation<RecurringExpense, ApiError, RecurringPayload>({
    mutationFn: createRecurring,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: qk.recurringAll })
      toast.success("Recurring expense created")
    },
    onError: (err) => toast.error(err.message),
  })
}

export function useUpdateRecurring() {
  const qc = useQueryClient()
  return useMutation<RecurringExpense, ApiError, { id: number; payload: RecurringPayload }>({
    mutationFn: ({ id, payload }) => updateRecurring(id, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: qk.recurringAll })
      toast.success("Recurring expense updated")
    },
    onError: (err) => toast.error(err.message),
  })
}

export function useDeleteRecurring() {
  const qc = useQueryClient()
  return useMutation<void, ApiError, number>({
    mutationFn: deleteRecurring,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: qk.recurringAll })
      toast.success("Recurring expense deleted")
    },
    onError: (err) => toast.error(err.message),
  })
}
