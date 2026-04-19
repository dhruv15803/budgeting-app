import { Link } from "react-router"
import { Button } from "@/components/ui/button"

export function NotFoundPage() {
  return (
    <div className="bg-background flex h-full min-h-screen flex-col items-center justify-center gap-4 p-8 text-center">
      <div className="text-primary text-6xl font-bold">404</div>
      <h1 className="text-xl font-semibold">Page not found</h1>
      <p className="text-muted-foreground max-w-sm text-sm">
        The page you're looking for doesn't exist or has been moved.
      </p>
      <Button asChild>
        <Link to="/">Back to dashboard</Link>
      </Button>
    </div>
  )
}
