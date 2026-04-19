import * as React from "react"
import { Input } from "@/components/ui/input"
import { cn } from "@/lib/utils"

interface CurrencyInputProps
  extends Omit<React.ComponentProps<"input">, "value" | "onChange" | "type" | "prefix"> {
  value: number | null | undefined
  onValueChange: (value: number | null) => void
  currencySymbol?: string
}

export function CurrencyInput({
  value,
  onValueChange,
  currencySymbol = "$",
  className,
  ...props
}: CurrencyInputProps) {
  const [text, setText] = React.useState(value == null ? "" : String(value))

  React.useEffect(() => {
    if (value == null) {
      setText("")
      return
    }
    const asNum = Number(text)
    if (Number.isFinite(asNum) && asNum === value) return
    setText(String(value))
  }, [value, text])

  return (
    <div className="relative">
      <span className="text-muted-foreground pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-sm">
        {currencySymbol}
      </span>
      <Input
        inputMode="decimal"
        className={cn("pl-7", className)}
        value={text}
        onChange={(e) => {
          const raw = e.target.value.replace(/[^0-9.]/g, "")
          const parts = raw.split(".")
          const normalized =
            parts.length > 2 ? `${parts[0]}.${parts.slice(1).join("")}` : raw
          setText(normalized)
          if (normalized === "" || normalized === ".") {
            onValueChange(null)
          } else {
            const n = Number(normalized)
            onValueChange(Number.isFinite(n) ? n : null)
          }
        }}
        onBlur={() => {
          if (text === "") return
          const n = Number(text)
          if (Number.isFinite(n)) setText(n.toFixed(2))
        }}
        {...props}
      />
    </div>
  )
}
