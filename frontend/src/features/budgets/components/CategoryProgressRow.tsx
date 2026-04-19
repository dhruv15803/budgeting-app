import { motion } from "framer-motion"
import { Progress } from "@/components/ui/progress"
import { cn, formatCurrency } from "@/lib/utils"
import type { BudgetCategory } from "@/types/api"

interface Props {
  row: BudgetCategory
  color?: string
}

export function CategoryProgressRow({ row, color }: Props) {
  const pct =
    row.allocated_amount > 0
      ? Math.min(100, (row.spent_amount / row.allocated_amount) * 100)
      : 0
  const overspent = row.remaining < 0
  const tone = overspent
    ? "bg-destructive"
    : pct >= 90
      ? "bg-amber-500"
      : "bg-emerald-500"

  return (
    <motion.div
      initial={{ opacity: 0, y: 4 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-2"
    >
      <div className="flex items-center justify-between gap-3">
        <div className="flex min-w-0 items-center gap-2">
          <span
            className="size-2.5 shrink-0 rounded-full"
            style={{ backgroundColor: color ?? "var(--muted-foreground)" }}
          />
          <span className="truncate text-sm font-medium">{row.category_name}</span>
        </div>
        <div className="flex items-baseline gap-2 text-sm tabular-nums">
          <span className={cn("font-mono font-medium", overspent && "text-destructive")}>
            {formatCurrency(row.spent_amount)}
          </span>
          <span className="text-muted-foreground text-xs">
            / {formatCurrency(row.allocated_amount)}
          </span>
        </div>
      </div>
      <Progress value={pct} indicatorClassName={tone} />
    </motion.div>
  )
}
