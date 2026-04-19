import { ArrowRightIcon, ReceiptIcon } from "lucide-react"
import { Link } from "react-router"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { CategoryBadge } from "@/components/common/CategoryBadge"
import { EmptyState } from "@/components/common/EmptyState"
import { useCategoryMap } from "@/features/categories/hooks"
import { colorForIndex } from "@/lib/constants"
import { formatCurrency, formatDate } from "@/lib/utils"
import type { Expense } from "@/types/api"

interface Props {
  expenses?: Expense[]
  isLoading?: boolean
}

export function RecentExpensesList({ expenses, isLoading }: Props) {
  const catMap = useCategoryMap()

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Recent expenses</CardTitle>
          <Button asChild variant="ghost" size="sm" className="gap-1">
            <Link to="/expenses">
              View all <ArrowRightIcon className="size-3.5" />
            </Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="space-y-2">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        ) : !expenses || expenses.length === 0 ? (
          <EmptyState icon={<ReceiptIcon className="size-5" />} title="Nothing logged yet" />
        ) : (
          <ul className="divide-y">
            {expenses.map((e) => (
              <li key={e.id} className="flex items-center justify-between gap-3 py-3">
                <div className="flex min-w-0 items-center gap-3">
                  <div className="bg-muted text-muted-foreground flex size-9 shrink-0 items-center justify-center rounded-md">
                    <ReceiptIcon className="size-4" />
                  </div>
                  <div className="min-w-0">
                    <div className="truncate text-sm font-medium">{e.title}</div>
                    <div className="text-muted-foreground flex items-center gap-2 text-xs">
                      <span>{formatDate(e.expense_date)}</span>
                      <span>•</span>
                      <CategoryBadge
                        name={catMap.get(e.category_id) ?? `#${e.category_id}`}
                        color={colorForIndex(e.category_id)}
                        className="px-1.5 py-0 text-[10px]"
                      />
                    </div>
                  </div>
                </div>
                <div className="font-mono font-medium tabular-nums">{formatCurrency(e.amount)}</div>
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  )
}
