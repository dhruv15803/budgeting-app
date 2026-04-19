import { Cell, Label, Legend, Pie, PieChart } from "recharts"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart"
import { EmptyState } from "@/components/common/EmptyState"
import { colorForIndex } from "@/lib/constants"
import { formatCurrency } from "@/lib/utils"
import type { BudgetCategory } from "@/types/api"

interface Props {
  categories: BudgetCategory[]
  totalSpent: number
  title?: string
  description?: string
}

export function SpendingByCategoryChart({
  categories,
  totalSpent,
  title = "Spending by category",
  description,
}: Props) {
  const data = categories
    .filter((c) => c.spent_amount > 0)
    .map((c, i) => ({
      key: `cat_${c.category_id}`,
      name: c.category_name,
      value: Number(c.spent_amount),
      fill: colorForIndex(i),
    }))

  const config: ChartConfig = Object.fromEntries(
    data.map((d) => [d.key, { label: d.name, color: d.fill }])
  )

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        {description && <CardDescription>{description}</CardDescription>}
      </CardHeader>
      <CardContent>
        {data.length === 0 ? (
          <EmptyState
            title="No spending yet"
            description="Once you log some expenses this month, the breakdown will show here."
          />
        ) : (
          <ChartContainer config={config} className="mx-auto aspect-square max-h-72">
            <PieChart>
              <ChartTooltip
                cursor={false}
                content={
                  <ChartTooltipContent
                    formatter={(value, name) => (
                      <div className="flex w-full items-center justify-between gap-4">
                        <span className="text-muted-foreground">{name}</span>
                        <span className="text-foreground font-mono font-medium tabular-nums">
                          {formatCurrency(value as number)}
                        </span>
                      </div>
                    )}
                  />
                }
              />
              <Pie
                data={data}
                dataKey="value"
                nameKey="name"
                innerRadius={60}
                outerRadius={95}
                strokeWidth={2}
                paddingAngle={2}
              >
                {data.map((d) => (
                  <Cell key={d.key} fill={d.fill} />
                ))}
                <Label
                  content={({ viewBox }) => {
                    if (!viewBox || !("cx" in viewBox)) return null
                    return (
                      <text
                        x={viewBox.cx}
                        y={viewBox.cy}
                        textAnchor="middle"
                        dominantBaseline="middle"
                      >
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy ?? 0) - 6}
                          className="fill-muted-foreground text-[10px] uppercase"
                        >
                          Total spent
                        </tspan>
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy ?? 0) + 14}
                          className="fill-foreground text-lg font-semibold"
                        >
                          {formatCurrency(totalSpent)}
                        </tspan>
                      </text>
                    )
                  }}
                />
              </Pie>
              <Legend
                iconType="circle"
                iconSize={8}
                wrapperStyle={{ fontSize: 12 }}
                formatter={(value) => <span className="text-muted-foreground">{value}</span>}
              />
            </PieChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  )
}
