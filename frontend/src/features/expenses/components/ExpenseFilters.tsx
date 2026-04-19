import { useEffect, useRef, useState } from "react"
import { format } from "date-fns"
import { CalendarIcon, FilterIcon, SearchIcon, XIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Calendar } from "@/components/ui/calendar"
import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { useCategories } from "@/features/categories/hooks"
import { cn, toDateInput } from "@/lib/utils"
import type { ExpenseListParams, ExpenseSortBy } from "@/types/api"

interface ExpenseFiltersProps {
  value: ExpenseListParams
  onChange: (next: ExpenseListParams) => void
}

const SORT_OPTIONS: { value: ExpenseSortBy; label: string }[] = [
  { value: "date_desc", label: "Newest first" },
  { value: "date_asc", label: "Oldest first" },
  { value: "amount_desc", label: "Amount: high to low" },
  { value: "amount_asc", label: "Amount: low to high" },
]

export function ExpenseFilters({ value, onChange }: ExpenseFiltersProps) {
  const [search, setSearch] = useState(value.search ?? "")
  const searchTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)
  const { data: categories = [] } = useCategories()

  useEffect(() => {
    setSearch(value.search ?? "")
  }, [value.search])

  const handleSearch = (v: string) => {
    setSearch(v)
    if (searchTimer.current) clearTimeout(searchTimer.current)
    searchTimer.current = setTimeout(() => {
      onChange({ ...value, search: v || undefined, page: 1 })
    }, 300)
  }

  const toggleCategory = (id: number, checked: boolean) => {
    const cur = new Set(value.category_id ?? [])
    if (checked) cur.add(id)
    else cur.delete(id)
    onChange({
      ...value,
      category_id: cur.size ? Array.from(cur) : undefined,
      page: 1,
    })
  }

  const selectedCount = value.category_id?.length ?? 0
  const dateFrom = value.date_from ? new Date(value.date_from) : undefined
  const dateTo = value.date_to ? new Date(value.date_to) : undefined
  const activeCount =
    (value.search ? 1 : 0) +
    (value.category_id?.length ? 1 : 0) +
    (value.date_from || value.date_to ? 1 : 0) +
    (value.amount_min != null ? 1 : 0) +
    (value.amount_max != null ? 1 : 0)

  return (
    <div className="flex flex-wrap items-center gap-2">
      <div className="relative min-w-48 flex-1">
        <SearchIcon className="text-muted-foreground absolute top-1/2 left-3 size-4 -translate-y-1/2" />
        <Input
          value={search}
          onChange={(e) => handleSearch(e.target.value)}
          placeholder="Search title or description"
          className="pl-9"
        />
      </div>

      <Popover>
        <PopoverTrigger asChild>
          <Button variant="outline" className="gap-2">
            <FilterIcon className="size-4" />
            Categories
            {selectedCount > 0 && (
              <Badge variant="secondary" className="ml-1 px-1.5 tabular-nums">
                {selectedCount}
              </Badge>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-64 p-0" align="start">
          <div className="max-h-72 overflow-y-auto p-2">
            {categories.length === 0 && (
              <div className="text-muted-foreground p-3 text-center text-sm">No categories</div>
            )}
            {categories.map((c) => {
              const checked = value.category_id?.includes(c.id) ?? false
              return (
                <Label
                  key={c.id}
                  className="hover:bg-accent flex cursor-pointer items-center gap-3 rounded-md px-2 py-1.5 text-sm font-normal"
                >
                  <Checkbox checked={checked} onCheckedChange={(v) => toggleCategory(c.id, !!v)} />
                  <span className="truncate">{c.category_name}</span>
                </Label>
              )
            })}
          </div>
          {selectedCount > 0 && (
            <div className="border-t p-2">
              <Button
                variant="ghost"
                size="sm"
                className="w-full justify-center"
                onClick={() => onChange({ ...value, category_id: undefined, page: 1 })}
              >
                Clear
              </Button>
            </div>
          )}
        </PopoverContent>
      </Popover>

      <Popover>
        <PopoverTrigger asChild>
          <Button variant="outline" className="gap-2">
            <CalendarIcon className="size-4" />
            {dateFrom || dateTo ? (
              <span className="text-xs">
                {dateFrom ? format(dateFrom, "MMM d") : "…"} – {dateTo ? format(dateTo, "MMM d") : "…"}
              </span>
            ) : (
              "Date range"
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="start">
          <Calendar
            mode="range"
            numberOfMonths={2}
            selected={{ from: dateFrom, to: dateTo }}
            onSelect={(range) => {
              onChange({
                ...value,
                date_from: range?.from ? toDateInput(range.from) : undefined,
                date_to: range?.to ? toDateInput(range.to) : undefined,
                month: undefined,
                year: undefined,
                page: 1,
              })
            }}
          />
          {(dateFrom || dateTo) && (
            <div className="border-t p-2">
              <Button
                variant="ghost"
                size="sm"
                className="w-full"
                onClick={() =>
                  onChange({ ...value, date_from: undefined, date_to: undefined, page: 1 })
                }
              >
                Clear range
              </Button>
            </div>
          )}
        </PopoverContent>
      </Popover>

      <Popover>
        <PopoverTrigger asChild>
          <Button variant="outline" className="gap-2">
            Amount
            {(value.amount_min != null || value.amount_max != null) && (
              <Badge variant="secondary" className="px-1.5">•</Badge>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-64" align="start">
          <div className="space-y-3">
            <div className="grid gap-1.5">
              <Label className="text-xs">Min</Label>
              <Input
                type="number"
                inputMode="decimal"
                placeholder="0"
                value={value.amount_min ?? ""}
                onChange={(e) =>
                  onChange({
                    ...value,
                    amount_min: e.target.value === "" ? undefined : Number(e.target.value),
                    page: 1,
                  })
                }
              />
            </div>
            <div className="grid gap-1.5">
              <Label className="text-xs">Max</Label>
              <Input
                type="number"
                inputMode="decimal"
                placeholder="No limit"
                value={value.amount_max ?? ""}
                onChange={(e) =>
                  onChange({
                    ...value,
                    amount_max: e.target.value === "" ? undefined : Number(e.target.value),
                    page: 1,
                  })
                }
              />
            </div>
          </div>
        </PopoverContent>
      </Popover>

      <Select
        value={value.sort_by ?? "date_desc"}
        onValueChange={(v) =>
          onChange({ ...value, sort_by: v as ExpenseSortBy, page: 1 })
        }
      >
        <SelectTrigger className={cn("min-w-40")}>
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {SORT_OPTIONS.map((o) => (
            <SelectItem key={o.value} value={o.value}>
              {o.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      {activeCount > 0 && (
        <Button
          variant="ghost"
          size="sm"
          onClick={() =>
            onChange({
              page: 1,
              page_size: value.page_size,
              sort_by: value.sort_by,
            })
          }
          className="gap-1"
        >
          <XIcon className="size-4" /> Reset
        </Button>
      )}
    </div>
  )
}
