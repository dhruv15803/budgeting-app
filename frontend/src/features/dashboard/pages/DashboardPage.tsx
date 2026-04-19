import { useMemo } from "react"
import { useQueries, useQuery } from "@tanstack/react-query"
import { PlusIcon, SparklesIcon, WalletIcon } from "lucide-react"
import { Link } from "react-router"
import { Button } from "@/components/ui/button"
import { PageHeader } from "@/components/common/PageHeader"
import { EmptyState } from "@/components/common/EmptyState"
import { KpiCards } from "@/features/dashboard/components/KpiCards"
import { MonthlyTrendChart } from "@/features/dashboard/components/MonthlyTrendChart"
import { BudgetProgressStrip } from "@/features/dashboard/components/BudgetProgressStrip"
import { RecentExpensesList } from "@/features/dashboard/components/RecentExpensesList"
import { SpendingByCategoryChart } from "@/features/budgets/components/SpendingByCategoryChart"
import { useBudgetOverview } from "@/features/budgets/hooks"
import { listExpenses } from "@/features/expenses/api"
import { listRecurring } from "@/features/recurring/api"
import { useCategoryMap } from "@/features/categories/hooks"
import { qk } from "@/lib/constants"
import { formatMonthLabel, monthKey, parseMonthKey } from "@/lib/utils"
import { useAuth } from "@/features/auth/AuthProvider"
import type { BudgetCategory } from "@/types/api"

function lastNMonths(n: number): string[] {
  const months: string[] = []
  const d = new Date()
  d.setDate(1)
  for (let i = n - 1; i >= 0; i--) {
    const m = new Date(d.getFullYear(), d.getMonth() - i, 1)
    months.push(monthKey(m))
  }
  return months
}

export function DashboardPage() {
  const { user } = useAuth()
  const currentMonth = monthKey()
  const prevMonth = monthKey(new Date(new Date().setMonth(new Date().getMonth() - 1)))

  const categoryMap = useCategoryMap()

  const budget = useBudgetOverview(currentMonth)

  const currentMonthExpenses = useQuery({
    queryKey: qk.expenses({ page: 1, page_size: 200, month: currentMonth, sort_by: "date_desc" }),
    queryFn: () =>
      listExpenses({
        page: 1,
        page_size: 200,
        month: currentMonth,
        sort_by: "date_desc",
      }),
  })

  const prevMonthExpenses = useQuery({
    queryKey: qk.expenses({ page: 1, page_size: 200, month: prevMonth, sort_by: "date_desc" }),
    queryFn: () =>
      listExpenses({ page: 1, page_size: 200, month: prevMonth, sort_by: "date_desc" }),
  })

  const recent = useQuery({
    queryKey: qk.expenses({ page: 1, page_size: 5, sort_by: "date_desc" }),
    queryFn: () => listExpenses({ page: 1, page_size: 5, sort_by: "date_desc" }),
  })

  const activeRecurring = useQuery({
    queryKey: qk.recurring({ page: 1, page_size: 100, is_active: true }),
    queryFn: () => listRecurring({ page: 1, page_size: 100, is_active: true }),
  })

  const trendMonths = useMemo(() => lastNMonths(6), [])
  const trendQueries = useQueries({
    queries: trendMonths.map((m) => ({
      queryKey: qk.expenses({ page: 1, page_size: 1000, month: m, sort_by: "date_desc" }),
      queryFn: () =>
        listExpenses({ page: 1, page_size: 1000, month: m, sort_by: "date_desc" }),
      staleTime: 60_000,
    })),
  })

  const trendData = useMemo(
    () =>
      trendMonths.map((m, i) => {
        const q = trendQueries[i]
        const total = q.data?.expenses.reduce((s, e) => s + Number(e.amount), 0) ?? 0
        return {
          month: m,
          label: new Intl.DateTimeFormat("en-US", { month: "short" }).format(parseMonthKey(m)),
          total,
          isCurrent: m === currentMonth,
        }
      }),
    [trendMonths, trendQueries, currentMonth]
  )

  const monthSpent = currentMonthExpenses.data?.expenses.reduce(
    (s, e) => s + Number(e.amount),
    0
  )
  const prevMonthSpent = prevMonthExpenses.data?.expenses.reduce(
    (s, e) => s + Number(e.amount),
    0
  )

  // Sum of active recurring due in next 30 days
  const now = new Date()
  const thirtyOut = new Date(now)
  thirtyOut.setDate(thirtyOut.getDate() + 30)
  const upcomingTotal =
    activeRecurring.data?.recurring_expenses
      .filter((r) => {
        const next = new Date(r.next_occurrence)
        return next >= now && next <= thirtyOut
      })
      .reduce((s, r) => s + Number(r.amount), 0) ?? 0

  // Build a spending-by-category list for the donut (fallback when no budget exists)
  const fallbackCategories: BudgetCategory[] = useMemo(() => {
    const items = currentMonthExpenses.data?.expenses ?? []
    const totals = new Map<number, number>()
    for (const e of items) {
      totals.set(e.category_id, (totals.get(e.category_id) ?? 0) + Number(e.amount))
    }
    return Array.from(totals.entries()).map(([catId, spent]) => ({
      id: catId,
      category_id: catId,
      category_name: categoryMap.get(catId) ?? `#${catId}`,
      allocated_amount: 0,
      spent_amount: spent,
      remaining: -spent,
    }))
  }, [currentMonthExpenses.data, categoryMap])

  const greeting = (() => {
    const h = new Date().getHours()
    if (h < 12) return "Good morning"
    if (h < 18) return "Good afternoon"
    return "Good evening"
  })()
  const name = user?.username || user?.email?.split("@")[0] || "there"

  const chartCategories = budget.data?.categories.length
    ? budget.data.categories
    : fallbackCategories
  const chartTotalSpent = budget.data?.total_spent ?? monthSpent ?? 0

  return (
    <div className="space-y-6">
      <PageHeader
        title={`${greeting}, ${name}`}
        description={`Here's a snapshot of your spending for ${formatMonthLabel(currentMonth)}.`}
        actions={
          <Button asChild className="gap-2">
            <Link to="/expenses">
              <PlusIcon className="size-4" /> Log expense
            </Link>
          </Button>
        }
      />

      <KpiCards
        monthLabel={new Intl.DateTimeFormat("en-US", { month: "short" }).format(new Date())}
        monthSpent={monthSpent}
        monthSpentLoading={currentMonthExpenses.isLoading}
        budgetRemaining={budget.data?.remaining ?? null}
        budgetLoading={budget.isLoading}
        upcomingRecurring={upcomingTotal}
        expenseCount={currentMonthExpenses.data?.expenses.length ?? 0}
        prevMonthSpent={prevMonthSpent}
      />

      {!budget.isLoading && !budget.data && (
        <EmptyState
          icon={<SparklesIcon className="size-5" />}
          title={`No budget for ${formatMonthLabel(currentMonth)}`}
          description="Set a monthly target to unlock allocations, progress, and category breakdowns."
          action={
            <Button asChild className="gap-2">
              <Link to="/budgets">
                <WalletIcon className="size-4" /> Set up this month's budget
              </Link>
            </Button>
          }
        />
      )}

      <div className="grid gap-6 lg:grid-cols-5">
        <div className="space-y-6 lg:col-span-3">
          <MonthlyTrendChart
            data={trendData}
            isLoading={trendQueries.some((q) => q.isLoading)}
          />
          {budget.data ? (
            <BudgetProgressStrip overview={budget.data} />
          ) : null}
        </div>
        <div className="space-y-6 lg:col-span-2">
          <SpendingByCategoryChart
            categories={chartCategories}
            totalSpent={chartTotalSpent}
            description={formatMonthLabel(currentMonth)}
          />
          <RecentExpensesList
            expenses={recent.data?.expenses}
            isLoading={recent.isLoading}
          />
        </div>
      </div>
    </div>
  )
}
