import { motion } from "framer-motion"
import {
  TrendingDownIcon,
  TrendingUpIcon,
  WalletIcon,
  RepeatIcon,
  ReceiptIcon,
} from "lucide-react"
import { Card, CardContent } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { cn, formatCurrency } from "@/lib/utils"

interface Kpi {
  label: string
  value: number
  icon: React.ComponentType<{ className?: string }>
  accent: string
  subtext?: React.ReactNode
  format?: "currency" | "number"
}

interface Props {
  monthLabel: string
  monthSpent?: number
  monthSpentLoading?: boolean
  budgetRemaining?: number | null
  budgetLoading?: boolean
  upcomingRecurring?: number
  expenseCount?: number
  prevMonthSpent?: number
}

export function KpiCards({
  monthLabel,
  monthSpent,
  monthSpentLoading,
  budgetRemaining,
  budgetLoading,
  upcomingRecurring = 0,
  expenseCount = 0,
  prevMonthSpent,
}: Props) {
  const diff =
    monthSpent != null && prevMonthSpent != null && prevMonthSpent > 0
      ? ((monthSpent - prevMonthSpent) / prevMonthSpent) * 100
      : null
  const up = diff != null && diff > 0

  const items: Kpi[] = [
    {
      label: `Spent in ${monthLabel}`,
      value: monthSpent ?? 0,
      icon: ReceiptIcon,
      accent: "from-violet-500/20 to-violet-500/0",
      subtext:
        diff != null ? (
          <div className={cn("flex items-center gap-1 text-xs", up ? "text-destructive" : "text-emerald-500")}>
            {up ? <TrendingUpIcon className="size-3" /> : <TrendingDownIcon className="size-3" />}
            {Math.abs(diff).toFixed(1)}% vs last month
          </div>
        ) : (
          <div className="text-muted-foreground text-xs">No prior month to compare</div>
        ),
      format: "currency",
    },
    {
      label: "Budget remaining",
      value: budgetRemaining ?? 0,
      icon: WalletIcon,
      accent: "from-emerald-500/20 to-emerald-500/0",
      subtext: (
        <div className="text-muted-foreground text-xs">
          {budgetRemaining == null ? "No budget set" : "for this month"}
        </div>
      ),
      format: "currency",
    },
    {
      label: "Upcoming recurring",
      value: upcomingRecurring,
      icon: RepeatIcon,
      accent: "from-fuchsia-500/20 to-fuchsia-500/0",
      subtext: <div className="text-muted-foreground text-xs">Next 30 days</div>,
      format: "currency",
    },
    {
      label: "Expenses logged",
      value: expenseCount,
      icon: ReceiptIcon,
      accent: "from-indigo-500/20 to-indigo-500/0",
      subtext: <div className="text-muted-foreground text-xs">This month</div>,
      format: "number",
    },
  ]

  return (
    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      {items.map((k, i) => (
        <motion.div
          key={k.label}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.25, delay: i * 0.05 }}
        >
          <Card className="relative overflow-hidden">
            <div
              className={cn(
                "pointer-events-none absolute inset-0 bg-gradient-to-br opacity-60",
                k.accent
              )}
            />
            <CardContent className="relative">
              <div className="flex items-center justify-between">
                <div className="text-muted-foreground text-xs font-medium tracking-wide uppercase">
                  {k.label}
                </div>
                <div className="bg-background/80 flex size-8 items-center justify-center rounded-md border">
                  <k.icon className="size-4" />
                </div>
              </div>
              <div className="mt-4 space-y-1">
                {(k.label.startsWith("Spent") && monthSpentLoading) ||
                (k.label === "Budget remaining" && budgetLoading) ? (
                  <Skeleton className="h-8 w-32" />
                ) : (
                  <div className="text-2xl font-semibold tabular-nums">
                    {k.format === "number"
                      ? k.value.toLocaleString()
                      : formatCurrency(k.value)}
                  </div>
                )}
                {k.subtext}
              </div>
            </CardContent>
          </Card>
        </motion.div>
      ))}
    </div>
  )
}
