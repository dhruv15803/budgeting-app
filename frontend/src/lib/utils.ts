import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const currencyFormatter = new Intl.NumberFormat("en-US", {
  style: "currency",
  currency: "USD",
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
})

export function formatCurrency(value: number | string | null | undefined) {
  const num = typeof value === "string" ? Number(value) : value
  if (num == null || Number.isNaN(num)) return "$0.00"
  return currencyFormatter.format(num)
}

export function formatCurrencyCompact(value: number | string | null | undefined) {
  const num = typeof value === "string" ? Number(value) : value
  if (num == null || Number.isNaN(num)) return "$0"
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    notation: "compact",
    maximumFractionDigits: 1,
  }).format(num)
}

export function formatDate(date: string | Date, opts?: Intl.DateTimeFormatOptions) {
  const d = typeof date === "string" ? new Date(date) : date
  if (Number.isNaN(d.getTime())) return ""
  return new Intl.DateTimeFormat("en-US", opts ?? { year: "numeric", month: "short", day: "numeric" }).format(d)
}

export function monthKey(date: Date = new Date()): string {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, "0")
  return `${y}-${m}`
}

export function parseMonthKey(key: string): Date {
  const [y, m] = key.split("-").map(Number)
  return new Date(y, (m ?? 1) - 1, 1)
}

export function formatMonthLabel(key: string): string {
  const d = parseMonthKey(key)
  return new Intl.DateTimeFormat("en-US", { month: "long", year: "numeric" }).format(d)
}

export function toDateInput(date: Date): string {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, "0")
  const d = String(date.getDate()).padStart(2, "0")
  return `${y}-${m}-${d}`
}

export function debounce<T extends (...args: never[]) => void>(fn: T, wait = 300) {
  let timer: ReturnType<typeof setTimeout> | undefined
  return (...args: Parameters<T>) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => fn(...args), wait)
  }
}
