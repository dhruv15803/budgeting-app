import { useState } from "react"
import { motion } from "framer-motion"
import { ArrowRightIcon, CalendarIcon, PlusIcon, WalletIcon } from "lucide-react"
import { Link } from "react-router"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { EmptyState } from "@/components/common/EmptyState"
import { ErrorState } from "@/components/common/ErrorState"
import { PageHeader } from "@/components/common/PageHeader"
import { CreateBudgetDialog } from "@/features/budgets/components/CreateBudgetDialog"
import { useBudgets } from "@/features/budgets/hooks"
import { formatCurrency, formatMonthLabel } from "@/lib/utils"

export function BudgetsPage() {
  const [createOpen, setCreateOpen] = useState(false)
  const { data, isLoading, isError, error, refetch } = useBudgets({ page: 1, page_size: 60 })

  return (
    <div className="space-y-6">
      <PageHeader
        title="Budgets"
        description="A budget for every month. Open one to see breakdowns and allocations."
        actions={
          <Button onClick={() => setCreateOpen(true)} className="gap-2">
            <PlusIcon className="size-4" />
            New budget
          </Button>
        }
      />

      {isError ? (
        <ErrorState message={error?.message} onRetry={() => refetch()} />
      ) : isLoading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-36 w-full" />
          ))}
        </div>
      ) : !data || data.budgets.length === 0 ? (
        <EmptyState
          icon={<WalletIcon className="size-5" />}
          title="No budgets yet"
          description="Create your first monthly budget to start allocating and tracking spending."
          action={
            <Button onClick={() => setCreateOpen(true)} className="gap-2">
              <PlusIcon className="size-4" />
              Create budget
            </Button>
          }
        />
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {data.budgets.map((b, idx) => (
            <motion.div
              key={b.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.2, delay: Math.min(idx * 0.04, 0.3) }}
            >
              <Link
                to={`/budgets/${b.budget_month}`}
                className="focus-visible:ring-ring block rounded-xl transition focus-visible:ring-2 focus-visible:ring-offset-2"
              >
                <Card className="hover:border-primary/40 hover:shadow-primary/5 group h-full transition-all hover:shadow-md">
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <div className="text-muted-foreground flex items-center gap-2 text-xs uppercase">
                        <CalendarIcon className="size-3.5" />
                        {formatMonthLabel(b.budget_month)}
                      </div>
                      <ArrowRightIcon className="text-muted-foreground group-hover:text-foreground size-4 transition-colors" />
                    </div>
                    <CardTitle className="text-2xl tabular-nums">
                      {formatCurrency(b.total_amount)}
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-muted-foreground text-xs">
                      Click to view allocations and spending
                    </div>
                  </CardContent>
                </Card>
              </Link>
            </motion.div>
          ))}
        </div>
      )}

      <CreateBudgetDialog open={createOpen} onOpenChange={setCreateOpen} />
    </div>
  )
}
