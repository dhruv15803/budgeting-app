import { Navigate, Outlet, useLocation } from "react-router"
import { useAuth } from "@/features/auth/AuthProvider"
import { Skeleton } from "@/components/ui/skeleton"

export function RequireAuth() {
  const { token, user, isLoadingUser } = useAuth()
  const location = useLocation()

  if (!token) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  if (isLoadingUser && !user) {
    return (
      <div className="flex h-full w-full items-center justify-center p-8">
        <div className="w-full max-w-md space-y-4">
          <Skeleton className="h-8 w-40" />
          <Skeleton className="h-24 w-full" />
          <Skeleton className="h-24 w-full" />
        </div>
      </div>
    )
  }

  return <Outlet />
}

export function RedirectIfAuthed({ children }: { children: React.ReactNode }) {
  const { token } = useAuth()
  if (token) return <Navigate to="/" replace />
  return <>{children}</>
}
