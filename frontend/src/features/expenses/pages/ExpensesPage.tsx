import { useState } from "react"
import { PlusIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { DataTablePager } from "@/components/common/DataTablePager"
import { ErrorState } from "@/components/common/ErrorState"
import { PageHeader } from "@/components/common/PageHeader"
import { DeleteExpenseDialog } from "@/features/expenses/components/DeleteExpenseDialog"
import { ExpenseFilters } from "@/features/expenses/components/ExpenseFilters"
import { ExpenseFormDialog } from "@/features/expenses/components/ExpenseFormDialog"
import { ExpenseTable } from "@/features/expenses/components/ExpenseTable"
import { useExpenses } from "@/features/expenses/hooks"
import type { Expense, ExpenseListParams } from "@/types/api"

const DEFAULT_PARAMS: ExpenseListParams = {
  page: 1,
  page_size: 20,
  sort_by: "date_desc",
}

export function ExpensesPage() {
  const [params, setParams] = useState<ExpenseListParams>(DEFAULT_PARAMS)
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<Expense | null>(null)
  const [deleting, setDeleting] = useState<Expense | null>(null)

  const { data, isLoading, isError, error, refetch, isFetching } = useExpenses(params)

  return (
    <div className="space-y-6">
      <PageHeader
        title="Expenses"
        description="Search, filter, and manage every expense you've logged."
        actions={
          <Button
            onClick={() => {
              setEditing(null)
              setFormOpen(true)
            }}
            className="gap-2"
          >
            <PlusIcon className="size-4" />
            Add expense
          </Button>
        }
      />

      <Card className="gap-0 py-4">
        <div className="px-4">
          <ExpenseFilters value={params} onChange={setParams} />
        </div>
        <div className="mt-4 border-t">
          {isError ? (
            <div className="p-4">
              <ErrorState message={error?.message} onRetry={() => refetch()} />
            </div>
          ) : (
            <ExpenseTable
              expenses={data?.expenses}
              isLoading={isLoading || isFetching}
              onEdit={(e) => {
                setEditing(e)
                setFormOpen(true)
              }}
              onDelete={(e) => setDeleting(e)}
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
              onPageChange={(p) => setParams((v) => ({ ...v, page: p }))}
            />
          </div>
        )}
      </Card>

      <ExpenseFormDialog
        open={formOpen}
        onOpenChange={(o) => {
          setFormOpen(o)
          if (!o) setEditing(null)
        }}
        expense={editing}
      />
      <DeleteExpenseDialog
        open={!!deleting}
        onOpenChange={(o) => {
          if (!o) setDeleting(null)
        }}
        expense={deleting}
      />
    </div>
  )
}
