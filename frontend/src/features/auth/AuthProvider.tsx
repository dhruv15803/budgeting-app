import * as React from "react"
import { useQuery, useQueryClient } from "@tanstack/react-query"
import { fetchMe } from "@/features/auth/api"
import { LOGOUT_EVENT, TOKEN_STORAGE_KEY, qk } from "@/lib/constants"
import type { User } from "@/types/api"

interface AuthContextValue {
  token: string | null
  user: User | null
  isLoadingUser: boolean
  login: (token: string) => void
  logout: () => void
}

const AuthContext = React.createContext<AuthContextValue | null>(null)

export function useAuth() {
  const ctx = React.useContext(AuthContext)
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>")
  return ctx
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = React.useState<string | null>(() =>
    typeof window !== "undefined" ? localStorage.getItem(TOKEN_STORAGE_KEY) : null
  )
  const qc = useQueryClient()

  const meQuery = useQuery({
    queryKey: qk.me,
    queryFn: fetchMe,
    enabled: !!token,
    staleTime: 5 * 60_000,
    retry: false,
  })

  const login = React.useCallback(
    (nextToken: string) => {
      localStorage.setItem(TOKEN_STORAGE_KEY, nextToken)
      setToken(nextToken)
      qc.invalidateQueries({ queryKey: qk.me })
    },
    [qc]
  )

  const logout = React.useCallback(() => {
    localStorage.removeItem(TOKEN_STORAGE_KEY)
    setToken(null)
    qc.clear()
  }, [qc])

  React.useEffect(() => {
    const handler = () => {
      setToken(null)
      qc.clear()
    }
    window.addEventListener(LOGOUT_EVENT, handler)
    return () => window.removeEventListener(LOGOUT_EVENT, handler)
  }, [qc])

  const value = React.useMemo<AuthContextValue>(
    () => ({
      token,
      user: meQuery.data ?? null,
      isLoadingUser: meQuery.isLoading,
      login,
      logout,
    }),
    [token, meQuery.data, meQuery.isLoading, login, logout]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
