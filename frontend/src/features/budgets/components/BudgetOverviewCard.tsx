import { useState } from "react"
import { motion } from "framer-motion"
import { CheckIcon, Loader2Icon, PencilIcon, XIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { CurrencyInput } from "@/components/common/CurrencyInput"
import { useUpdateBudgetTotal } from "@/features/budgets/hooks"
import { cn, formatCurrency, formatMonthLabel } from "@/lib/utils"
import type { BudgetOverview } from "@/types/api"

interface Props {
  overview: BudgetOverview
}

export function BudgetOverviewCard({ overview }: Props) {
  const [editing, setEditing] = useState(false)
  const [amount, setAmount] = useState<number | null>(overview.total_amount)
  const updateMut = useUpdateBudgetTotal()

  const pct = overview.total_amount > 0
    ? Math.min(100, (overview.total_spent / overview.total_amount) * 100)
    : 0

  const tone =
    pct >= 100 ? "destructive" : pct >= 90 ? "warning" : pct >= 70 ? "caution" : "ok"

  const indicatorClass =
    tone === "destructive"
      ? "bg-destructive"
      : tone === "warning"
        ? "bg-amber-500"
        : tone === "caution"
          ? "bg-amber-400"
          : "bg-emerald-500"

  return (
    <Card className="overflow-hidden">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
            {formatMonthLabel(overview.budget_month)}
          </CardTitle>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid gap-4 md:grid-cols-3">
          <div className="space-y-1">
            <div className="text-muted-foreground text-xs uppercase">Total budget</div>
            {editing ? (
              <div className="flex items-center gap-2">
                <CurrencyInput value={amount} onValueChange={setAmount} className="h-9" />
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-9"
                  disabled={!amount || amount <= 0 || updateMut.isPending}
                  onClick={async () => {
                    if (!amount) return
                    await updateMut.mutateAsync({ month: overview.budget_month, total: amount })
                    setEditing(false)
                  }}
                  aria-label="Save"
                >
                  {updateMut.isPending ? <Loader2Icon className="animate-spin" /> : <CheckIcon />}
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-9"
                  onClick={() => {
                    setEditing(false)
                    setAmount(overview.total_amount)
                  }}
                  aria-label="Cancel"
                >
                  <XIcon />
                </Button>
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <motion.div
                  key={overview.total_amount}
                  initial={{ opacity: 0, y: 4 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="text-2xl font-semibold tabular-nums"
                >
                  {formatCurrency(overview.total_amount)}
                </motion.div>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-8"
                  onClick={() => setEditing(true)}
                  aria-label="Edit total"
                >
                  <PencilIcon />
                </Button>
              </div>
            )}
          </div>
          <div className="space-y-1">
            <div className="text-muted-foreground text-xs uppercase">Spent</div>
            <div className="text-2xl font-semibold tabular-nums">
              {formatCurrency(overview.total_spent)}
            </div>
            <div className="text-muted-foreground text-xs">{pct.toFixed(1)}% of budget</div>
          </div>
          <div className="space-y-1">
            <div className="text-muted-foreground text-xs uppercase">Remaining</div>
            <div
              className={cn(
                "text-2xl font-semibold tabular-nums",
                overview.remaining < 0 && "text-destructive"
              )}
            >
              {formatCurrency(overview.remaining)}
            </div>
          </div>
        </div>

        <div>
          <Progress value={pct} indicatorClassName={indicatorClass} />
        </div>
      </CardContent>
    </Card>
  )
}
