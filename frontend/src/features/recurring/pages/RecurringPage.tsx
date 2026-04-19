import { useState } from "react"
import { PlusIcon, SearchIcon } from "lucide-react"
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
import { Card } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { DataTablePager } from "@/components/common/DataTablePager"
import { ErrorState } from "@/components/common/ErrorState"
import { PageHeader } from "@/components/common/PageHeader"
import { RecurringFormDialog } from "@/features/recurring/components/RecurringFormDialog"
import { RecurringTable } from "@/features/recurring/components/RecurringTable"
import { useDeleteRecurring, useRecurring } from "@/features/recurring/hooks"
import type { Frequency, RecurringExpense, RecurringListParams } from "@/types/api"
import { Loader2Icon } from "lucide-react"

const DEFAULT_PARAMS: RecurringListParams = {
  page: 1,
  page_size: 20,
  sort_by: "next_occurrence ASC",
}

export function RecurringPage() {
  const [params, setParams] = useState<RecurringListParams>(DEFAULT_PARAMS)
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<RecurringExpense | null>(null)
  const [deleting, setDeleting] = useState<RecurringExpense | null>(null)
  const { data, isLoading, isFetching, isError, error, refetch } = useRecurring(params)
  const delMut = useDeleteRecurring()

  return (
    <div className="space-y-6">
      <PageHeader
        title="Recurring expenses"
        description="Subscriptions and bills that regenerate on a schedule."
        actions={
          <Button
            onClick={() => {
              setEditing(null)
              setFormOpen(true)
            }}
            className="gap-2"
          >
            <PlusIcon className="size-4" />
            New recurring
          </Button>
        }
      />

      <Card className="gap-0 py-4">
        <div className="flex flex-wrap items-center gap-2 px-4">
          <div className="relative min-w-48 flex-1">
            <SearchIcon className="text-muted-foreground absolute top-1/2 left-3 size-4 -translate-y-1/2" />
            <Input
              value={params.search ?? ""}
              onChange={(e) =>
                setParams({ ...params, search: e.target.value || undefined, page: 1 })
              }
              placeholder="Search title or description"
              className="pl-9"
            />
          </div>

          <Select
            value={params.frequency ?? "all"}
            onValueChange={(v) =>
              setParams({
                ...params,
                frequency: v === "all" ? undefined : (v as Frequency),
                page: 1,
              })
            }
          >
            <SelectTrigger className="w-40">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All frequencies</SelectItem>
              <SelectItem value="daily">Daily</SelectItem>
              <SelectItem value="weekly">Weekly</SelectItem>
              <SelectItem value="monthly">Monthly</SelectItem>
              <SelectItem value="yearly">Yearly</SelectItem>
            </SelectContent>
          </Select>

          <Select
            value={params.is_active == null ? "all" : params.is_active ? "active" : "paused"}
            onValueChange={(v) =>
              setParams({
                ...params,
                is_active: v === "all" ? undefined : v === "active",
                page: 1,
              })
            }
          >
            <SelectTrigger className="w-32">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All status</SelectItem>
              <SelectItem value="active">Active</SelectItem>
              <SelectItem value="paused">Paused</SelectItem>
            </SelectContent>
          </Select>

          <Select
            value={params.sort_by ?? "next_occurrence ASC"}
            onValueChange={(v) => setParams({ ...params, sort_by: v, page: 1 })}
          >
            <SelectTrigger className="min-w-44">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="next_occurrence ASC">Next soonest</SelectItem>
              <SelectItem value="next_occurrence DESC">Next latest</SelectItem>
              <SelectItem value="amount DESC">Amount: high to low</SelectItem>
              <SelectItem value="amount ASC">Amount: low to high</SelectItem>
              <SelectItem value="title ASC">Title A–Z</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="mt-4 border-t">
          {isError ? (
            <div className="p-4">
              <ErrorState message={error?.message} onRetry={() => refetch()} />
            </div>
          ) : (
            <RecurringTable
              items={data?.recurring_expenses}
              isLoading={isLoading || isFetching}
              onEdit={(r) => {
                setEditing(r)
                setFormOpen(true)
              }}
              onDelete={(r) => setDeleting(r)}
            />
          )}
        </div>
        {data && data.total > 0 && (
          <div className="border-t px-2">
            <DataTablePager
              page={data.page}
              totalPages={data.total_pages}
              total={data.total}
              pageSize={data.page_size}
              onPageChange={(p) => setParams({ ...params, page: p })}
            />
          </div>
        )}
      </Card>

      <RecurringFormDialog
        open={formOpen}
        onOpenChange={(o) => {
          setFormOpen(o)
          if (!o) setEditing(null)
        }}
        recurring={editing}
      />

      <AlertDialog open={!!deleting} onOpenChange={(o) => !o && setDeleting(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete recurring expense?</AlertDialogTitle>
            <AlertDialogDescription>
              This will stop future generation for{" "}
              <span className="text-foreground font-medium">{deleting?.title}</span>. Existing expenses
              already generated will remain.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={delMut.isPending}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-white hover:bg-destructive/90"
              disabled={delMut.isPending}
              onClick={async (e) => {
                e.preventDefault()
                if (!deleting) return
                await delMut.mutateAsync(deleting.id)
                setDeleting(null)
              }}
            >
              {delMut.isPending && <Loader2Icon className="animate-spin" />}
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
