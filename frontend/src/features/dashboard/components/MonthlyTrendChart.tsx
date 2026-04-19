import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart"
import { Skeleton } from "@/components/ui/skeleton"
import { formatCurrency, formatCurrencyCompact } from "@/lib/utils"

interface Props {
  data: { month: string; label: string; total: number; isCurrent?: boolean }[]
  isLoading?: boolean
}

export function MonthlyTrendChart({ data, isLoading }: Props) {
  const config: ChartConfig = {
    total: { label: "Spent", color: "var(--chart-1)" },
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>6-month trend</CardTitle>
        <CardDescription>Total monthly expenses</CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <Skeleton className="h-48 w-full" />
        ) : (
          <ChartContainer config={config} className="h-52 w-full">
            <BarChart data={data} margin={{ top: 8, right: 8, bottom: 4, left: -12 }}>
              <CartesianGrid vertical={false} strokeDasharray="3 3" />
              <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} />
              <YAxis
                tickLine={false}
                axisLine={false}
                tickFormatter={(v) => formatCurrencyCompact(v)}
                width={60}
              />
              <ChartTooltip
                cursor={{ fill: "var(--muted)", opacity: 0.3 }}
                content={
                  <ChartTooltipContent
                    formatter={(value) => (
                      <div className="flex w-full items-center justify-between gap-6">
                        <span className="text-muted-foreground">Spent</span>
                        <span className="font-mono font-medium tabular-nums">
                          {formatCurrency(value as number)}
                        </span>
                      </div>
                    )}
                  />
                }
              />
              <Bar dataKey="total" fill="var(--color-total)" radius={[6, 6, 0, 0]} />
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  )
}
