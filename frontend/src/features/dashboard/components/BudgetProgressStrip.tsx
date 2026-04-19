import { ArrowRightIcon } from "lucide-react"
import { Link } from "react-router"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { CategoryProgressRow } from "@/features/budgets/components/CategoryProgressRow"
import { colorForIndex } from "@/lib/constants"
import { cn, formatCurrency, formatMonthLabel } from "@/lib/utils"
import type { BudgetOverview } from "@/types/api"

interface Props {
  overview: BudgetOverview
}

export function BudgetProgressStrip({ overview }: Props) {
  const pct =
    overview.total_amount > 0
      ? Math.min(100, (overview.total_spent / overview.total_amount) * 100)
      : 0

  const topCategories = [...overview.categories]
    .filter((c) => c.allocated_amount > 0)
    .sort(
      (a, b) =>
        (b.spent_amount / b.allocated_amount) - (a.spent_amount / a.allocated_amount)
    )
    .slice(0, 3)

  const tone =
    pct >= 100 ? "bg-destructive" : pct >= 90 ? "bg-amber-500" : pct >= 70 ? "bg-amber-400" : "bg-emerald-500"

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Budget progress</CardTitle>
            <CardDescription>{formatMonthLabel(overview.budget_month)}</CardDescription>
          </div>
          <Button asChild variant="ghost" size="sm" className="gap-1">
            <Link to={`/budgets/${overview.budget_month}`}>
              Open <ArrowRightIcon className="size-3.5" />
            </Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-5">
        <div>
          <div className="mb-2 flex items-baseline justify-between">
            <div className="text-muted-foreground text-xs uppercase">Overall</div>
            <div className="text-sm tabular-nums">
              <span className={cn("font-mono font-medium", overview.remaining < 0 && "text-destructive")}>
                {formatCurrency(overview.total_spent)}
              </span>
              <span className="text-muted-foreground">
                {" "}
                / {formatCurrency(overview.total_amount)}
              </span>
            </div>
          </div>
          <Progress value={pct} indicatorClassName={tone} />
        </div>

        {topCategories.length > 0 && (
          <div className="space-y-4">
            <div className="text-muted-foreground text-xs tracking-wide uppercase">
              Top categories
            </div>
            {topCategories.map((c, i) => (
              <CategoryProgressRow key={c.category_id} row={c} color={colorForIndex(i)} />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
