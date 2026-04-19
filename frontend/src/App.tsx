import { lazy, Suspense } from "react"
import { Route, Routes } from "react-router"
import { AppShell } from "@/components/layout/AppShell"
import { RedirectIfAuthed, RequireAuth } from "@/features/auth/RequireAuth"
import { LoginPage } from "@/features/auth/pages/LoginPage"
import { RegisterPage } from "@/features/auth/pages/RegisterPage"
import { VerifyEmailPage } from "@/features/auth/pages/VerifyEmailPage"
import { NotFoundPage } from "@/features/misc/NotFoundPage"
import { Skeleton } from "@/components/ui/skeleton"

const DashboardPage = lazy(() =>
  import("@/features/dashboard/pages/DashboardPage").then((m) => ({ default: m.DashboardPage }))
)
const ExpensesPage = lazy(() =>
  import("@/features/expenses/pages/ExpensesPage").then((m) => ({ default: m.ExpensesPage }))
)
const RecurringPage = lazy(() =>
  import("@/features/recurring/pages/RecurringPage").then((m) => ({ default: m.RecurringPage }))
)
const BudgetsPage = lazy(() =>
  import("@/features/budgets/pages/BudgetsPage").then((m) => ({ default: m.BudgetsPage }))
)
const BudgetMonthPage = lazy(() =>
  import("@/features/budgets/pages/BudgetMonthPage").then((m) => ({ default: m.BudgetMonthPage }))
)

function RouteFallback() {
  return (
    <div className="space-y-4">
      <Skeleton className="h-10 w-64" />
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Skeleton className="h-28" />
        <Skeleton className="h-28" />
        <Skeleton className="h-28" />
        <Skeleton className="h-28" />
      </div>
      <Skeleton className="h-72" />
    </div>
  )
}

function Lazy({ children }: { children: React.ReactNode }) {
  return <Suspense fallback={<RouteFallback />}>{children}</Suspense>
}

function App() {
  return (
    <Routes>
      <Route
        path="/login"
        element={
          <RedirectIfAuthed>
            <LoginPage />
          </RedirectIfAuthed>
        }
      />
      <Route
        path="/register"
        element={
          <RedirectIfAuthed>
            <RegisterPage />
          </RedirectIfAuthed>
        }
      />
      <Route path="/verify-email" element={<VerifyEmailPage />} />

      <Route element={<RequireAuth />}>
        <Route element={<AppShell />}>
          <Route
            index
            element={
              <Lazy>
                <DashboardPage />
              </Lazy>
            }
          />
          <Route
            path="/expenses"
            element={
              <Lazy>
                <ExpensesPage />
              </Lazy>
            }
          />
          <Route
            path="/recurring"
            element={
              <Lazy>
                <RecurringPage />
              </Lazy>
            }
          />
          <Route
            path="/budgets"
            element={
              <Lazy>
                <BudgetsPage />
              </Lazy>
            }
          />
          <Route
            path="/budgets/:month"
            element={
              <Lazy>
                <BudgetMonthPage />
              </Lazy>
            }
          />
        </Route>
      </Route>

      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  )
}

export default App
