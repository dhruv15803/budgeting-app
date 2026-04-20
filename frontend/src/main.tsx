import { StrictMode } from "react"
import { createRoot } from "react-dom/client"
import { BrowserRouter } from "react-router"
import { QueryClientProvider } from "@tanstack/react-query"
import { ReactQueryDevtools } from "@tanstack/react-query-devtools"
import { GoogleOAuthProvider } from "@react-oauth/google"
import "./index.css"
import App from "./App"
import { ThemeProvider } from "@/components/layout/ThemeProvider"
import { AuthProvider } from "@/features/auth/AuthProvider"
import { queryClient } from "@/lib/queryClient"
import { Toaster } from "@/components/ui/sonner"
import { ErrorBoundary } from "@/components/common/ErrorBoundary"

const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined

function AppShell() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <App />
        <Toaster />
      </AuthProvider>
    </BrowserRouter>
  )
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ErrorBoundary>
      <ThemeProvider>
        <QueryClientProvider client={queryClient}>
          {googleClientId ? (
            <GoogleOAuthProvider clientId={googleClientId}>
              <AppShell />
            </GoogleOAuthProvider>
          ) : (
            <AppShell />
          )}
          {import.meta.env.DEV && <ReactQueryDevtools initialIsOpen={false} buttonPosition="bottom-left" />}
        </QueryClientProvider>
      </ThemeProvider>
    </ErrorBoundary>
  </StrictMode>
)
