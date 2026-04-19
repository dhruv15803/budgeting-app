import { AlertTriangleIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

interface ErrorStateProps {
  title?: string
  message?: string
  onRetry?: () => void
  className?: string
}

export function ErrorState({
  title = "Something went wrong",
  message,
  onRetry,
  className,
}: ErrorStateProps) {
  return (
    <div
      className={cn(
        "border-destructive/20 bg-destructive/5 flex flex-col items-center gap-3 rounded-xl border px-6 py-10 text-center",
        className
      )}
    >
      <div className="bg-destructive/10 text-destructive flex size-10 items-center justify-center rounded-full">
        <AlertTriangleIcon className="size-5" />
      </div>
      <div className="space-y-1">
        <h3 className="font-semibold">{title}</h3>
        {message && <p className="text-muted-foreground text-sm">{message}</p>}
      </div>
      {onRetry && (
        <Button variant="outline" size="sm" onClick={onRetry}>
          Try again
        </Button>
      )}
    </div>
  )
}
