import { useEffect, useMemo, useState } from "react"
import { Loader2Icon, PlusIcon, Trash2Icon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { CurrencyInput } from "@/components/common/CurrencyInput"
import { EmptyState } from "@/components/common/EmptyState"
import { useCategories } from "@/features/categories/hooks"
import { useBulkSetAllocations } from "@/features/budgets/hooks"
import type { BudgetOverview, Category } from "@/types/api"

interface Props {
  overview: BudgetOverview
}

interface Row {
  id: string
  category_id: number | null
  allocated_amount: number | null
}

export function CategoryAllocationsEditor({ overview }: Props) {
  const { data: categories = [] } = useCategories()
  const bulkMut = useBulkSetAllocations()

  const initialRows = useMemo<Row[]>(
    () =>
      overview.categories.map((c) => ({
        id: `a_${c.id}`,
        category_id: c.category_id,
        allocated_amount: Number(c.allocated_amount),
      })),
    [overview.categories]
  )

  const [rows, setRows] = useState<Row[]>(initialRows)

  useEffect(() => {
    setRows(initialRows)
  }, [initialRows])

  const usedIds = new Set(rows.map((r) => r.category_id).filter((v): v is number => v != null))

  const availableForRow = (current: number | null): Category[] =>
    categories.filter((c) => !usedIds.has(c.id) || c.id === current)

  const addRow = () =>
    setRows((rs) => [...rs, { id: `row_${Date.now()}_${rs.length}`, category_id: null, allocated_amount: null }])

  const removeRow = (id: string) => setRows((rs) => rs.filter((r) => r.id !== id))

  const totalAllocated = rows.reduce((sum, r) => sum + (r.allocated_amount ?? 0), 0)

  const save = async () => {
    const payload = rows
      .filter((r) => r.category_id != null && r.allocated_amount != null && r.allocated_amount >= 0)
      .map((r) => ({
        category_id: r.category_id as number,
        allocated_amount: r.allocated_amount as number,
      }))
    await bulkMut.mutateAsync({ month: overview.budget_month, categories: payload })
  }

  const hasChanges = JSON.stringify(rows) !== JSON.stringify(initialRows)

  return (
    <Card>
      <CardHeader>
        <CardTitle>Category allocations</CardTitle>
        <CardDescription>
          Split your monthly budget across categories. Unallocated expenses still count toward the
          total.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {rows.length === 0 ? (
          <EmptyState
            title="No allocations yet"
            description="Add your first category allocation to track spending per bucket."
            action={
              <Button onClick={addRow} variant="outline" className="gap-2">
                <PlusIcon className="size-4" />
                Add allocation
              </Button>
            }
          />
        ) : (
          <div className="space-y-2">
            {rows.map((row) => {
              const options = availableForRow(row.category_id)
              return (
                <div key={row.id} className="flex flex-wrap items-center gap-2">
                  <Select
                    value={row.category_id ? String(row.category_id) : undefined}
                    onValueChange={(v) =>
                      setRows((rs) =>
                        rs.map((r) => (r.id === row.id ? { ...r, category_id: Number(v) } : r))
                      )
                    }
                  >
                    <SelectTrigger className="min-w-44 flex-1">
                      <SelectValue placeholder="Pick a category" />
                    </SelectTrigger>
                    <SelectContent>
                      {options.map((c) => (
                        <SelectItem key={c.id} value={String(c.id)}>
                          {c.category_name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <div className="w-36">
                    <CurrencyInput
                      value={row.allocated_amount}
                      onValueChange={(v) =>
                        setRows((rs) =>
                          rs.map((r) => (r.id === row.id ? { ...r, allocated_amount: v } : r))
                        )
                      }
                    />
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="size-9"
                    onClick={() => removeRow(row.id)}
                    aria-label="Remove"
                  >
                    <Trash2Icon />
                  </Button>
                </div>
              )
            })}
          </div>
        )}

        <div className="flex flex-wrap items-center justify-between gap-2 border-t pt-4">
          <div className="text-muted-foreground text-sm">
            Allocated total:{" "}
            <span className="text-foreground font-mono font-medium tabular-nums">
              {new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(
                totalAllocated
              )}
            </span>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={addRow}
              disabled={usedIds.size >= categories.length}
              className="gap-2"
            >
              <PlusIcon className="size-4" />
              Add
            </Button>
            <Button
              onClick={save}
              disabled={!hasChanges || bulkMut.isPending || rows.some((r) => r.category_id == null)}
            >
              {bulkMut.isPending && <Loader2Icon className="animate-spin" />}
              Save allocations
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
