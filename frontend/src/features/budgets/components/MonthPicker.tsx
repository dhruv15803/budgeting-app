import { ChevronLeftIcon, ChevronRightIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { formatMonthLabel, monthKey, parseMonthKey } from "@/lib/utils"

interface MonthPickerProps {
  value: string
  onChange: (v: string) => void
  label?: string
  min?: string
  max?: string
}

export function MonthPicker({ value, onChange, label, min, max }: MonthPickerProps) {
  const shift = (delta: number) => {
    const d = parseMonthKey(value)
    d.setMonth(d.getMonth() + delta)
    onChange(monthKey(d))
  }

  return (
    <div className="grid gap-1.5">
      {label && <Label>{label}</Label>}
      <div className="flex items-center gap-2">
        <Button
          type="button"
          variant="outline"
          size="icon"
          className="size-9"
          onClick={() => shift(-1)}
          aria-label="Previous month"
        >
          <ChevronLeftIcon />
        </Button>
        <Input
          type="month"
          value={value}
          min={min}
          max={max}
          onChange={(e) => onChange(e.target.value || monthKey())}
          className="flex-1 text-center tabular-nums"
        />
        <Button
          type="button"
          variant="outline"
          size="icon"
          className="size-9"
          onClick={() => shift(1)}
          aria-label="Next month"
        >
          <ChevronRightIcon />
        </Button>
      </div>
      <div className="text-muted-foreground text-xs">{formatMonthLabel(value)}</div>
    </div>
  )
}
