import { useState } from "react"
import { useNavigate, useParams } from "react-router"
import { ArrowLeftIcon, Loader2Icon, Trash2Icon, WalletIcon } from "lucide-react"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { EmptyState } from "@/components/common/EmptyState"
import { ErrorState } from "@/components/common/ErrorState"
import { PageHeader } from "@/components/common/PageHeader"
import { BudgetOverviewCard } from "@/features/budgets/components/BudgetOverviewCard"
import { CategoryAllocationsEditor } from "@/features/budgets/components/CategoryAllocationsEditor"
import { CategoryProgressRow } from "@/features/budgets/components/CategoryProgressRow"
import { CreateBudgetDialog } from "@/features/budgets/components/CreateBudgetDialog"
import { SpendingByCategoryChart } from "@/features/budgets/components/SpendingByCategoryChart"
import { useBudgetOverview, useDeleteBudget } from "@/features/budgets/hooks"
import type { ApiError } from "@/lib/api"
import { colorForIndex } from "@/lib/constants"
import { formatMonthLabel } from "@/lib/utils"

export function BudgetMonthPage() {
  const { month } = useParams<{ month: string }>()
  const navigate = useNavigate()
  const { data, isLoading, isError, error, refetch } = useBudgetOverview(month)
  const [delOpen, setDelOpen] = useState(false)
  const [createOpen, setCreateOpen] = useState(false)
  const deleteMut = useDeleteBudget()

  const apiErr = error as ApiError | null
  const isNotFound = apiErr?.status === 404

  if (!month) return null

  return (
    <div className="space-y-6">
      <div>
        <Button variant="ghost" size="sm" asChild className="-ml-2 mb-2 gap-1">
          <button type="button" onClick={() => navigate("/budgets")}>
            <ArrowLeftIcon className="size-4" /> Back to budgets
          </button>
        </Button>
        <PageHeader
          title={formatMonthLabel(month)}
          description="Monthly budget overview, spending by category, and per-category allocations."
          actions={
            data && (
              <Button variant="outline" onClick={() => setDelOpen(true)} className="gap-2">
                <Trash2Icon className="size-4" /> Delete
              </Button>
            )
          }
        />
      </div>

      {isLoading ? (
        <div className="space-y-4">
          <Skeleton className="h-40 w-full" />
          <div className="grid gap-4 lg:grid-cols-2">
            <Skeleton className="h-80 w-full" />
            <Skeleton className="h-80 w-full" />
          </div>
        </div>
      ) : isError && isNotFound ? (
        <EmptyState
          icon={<WalletIcon className="size-5" />}
          title={`No budget for ${formatMonthLabel(month)}`}
          description="Create a monthly budget to start tracking and allocating."
          action={<Button onClick={() => setCreateOpen(true)}>Create budget for this month</Button>}
        />
      ) : isError ? (
        <ErrorState message={error?.message} onRetry={() => refetch()} />
      ) : data ? (
        <div className="space-y-6">
          <BudgetOverviewCard overview={data} />

          <div className="grid gap-6 lg:grid-cols-5">
            <div className="lg:col-span-3">
              <CategoryAllocationsEditor overview={data} />
            </div>
            <div className="lg:col-span-2">
              <SpendingByCategoryChart categories={data.categories} totalSpent={data.total_spent} />
            </div>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Category progress</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {data.categories.length === 0 ? (
                <EmptyState
                  title="No category allocations"
                  description="Add allocations above to see per-category progress."
                />
              ) : (
                data.categories.map((c, i) => (
                  <CategoryProgressRow key={c.category_id} row={c} color={colorForIndex(i)} />
                ))
              )}
            </CardContent>
          </Card>
        </div>
      ) : null}

      <AlertDialog open={delOpen} onOpenChange={setDelOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete this budget?</AlertDialogTitle>
            <AlertDialogDescription>
              Removes the monthly budget for {formatMonthLabel(month)} along with all category
              allocations. Existing expenses are not affected.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteMut.isPending}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-white hover:bg-destructive/90"
              disabled={deleteMut.isPending}
              onClick={async (e) => {
                e.preventDefault()
                await deleteMut.mutateAsync(month)
                setDelOpen(false)
                navigate("/budgets")
              }}
            >
              {deleteMut.isPending && <Loader2Icon className="animate-spin" />}
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <CreateBudgetDialog open={createOpen} onOpenChange={setCreateOpen} defaultMonth={month} />
    </div>
  )
}
