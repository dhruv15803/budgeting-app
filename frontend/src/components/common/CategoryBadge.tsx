import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

interface CategoryBadgeProps {
  name: string
  color?: string
  className?: string
}

export function CategoryBadge({ name, color, className }: CategoryBadgeProps) {
  return (
    <Badge variant="outline" className={cn("gap-1.5", className)}>
      <span
        className="size-2 rounded-full"
        style={{ backgroundColor: color ?? "var(--muted-foreground)" }}
        aria-hidden="true"
      />
      <span className="truncate">{name}</span>
    </Badge>
  )
}
