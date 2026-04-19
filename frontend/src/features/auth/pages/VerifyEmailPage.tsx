import { useEffect, useRef } from "react"
import { Link, useNavigate, useSearchParams } from "react-router"
import { AlertTriangleIcon, CheckCircle2Icon, Loader2Icon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { AuthLayout } from "@/features/auth/components/AuthLayout"
import { useVerifyEmail } from "@/features/auth/hooks"
import { useAuth } from "@/features/auth/AuthProvider"

export function VerifyEmailPage() {
  const [params] = useSearchParams()
  const token = params.get("token")
  const verify = useVerifyEmail()
  const auth = useAuth()
  const navigate = useNavigate()
  const fired = useRef(false)

  useEffect(() => {
    if (!token || fired.current) return
    fired.current = true
    verify.mutate(token, {
      onSuccess: ({ token: jwt }) => {
        auth.login(jwt)
        setTimeout(() => navigate("/", { replace: true }), 900)
      },
    })
  }, [token, verify, auth, navigate])

  if (!token) {
    return (
      <AuthLayout title="Invalid verification link">
        <div className="border-destructive/20 bg-destructive/5 flex flex-col items-center gap-3 rounded-xl border p-6 text-center">
          <div className="bg-destructive/10 text-destructive flex size-10 items-center justify-center rounded-full">
            <AlertTriangleIcon className="size-5" />
          </div>
          <p className="text-muted-foreground text-sm">
            This link is missing the verification token. Please use the link from the email we sent you.
          </p>
          <Button asChild variant="outline" className="w-full">
            <Link to="/login">Back to sign in</Link>
          </Button>
        </div>
      </AuthLayout>
    )
  }

  return (
    <AuthLayout title="Verifying your email">
      {verify.isPending && (
        <div className="flex flex-col items-center gap-3 py-6">
          <Loader2Icon className="text-primary size-8 animate-spin" />
          <p className="text-muted-foreground text-sm">Hang tight while we confirm your email…</p>
        </div>
      )}
      {verify.isSuccess && (
        <div className="border-primary/20 bg-primary/5 flex flex-col items-center gap-3 rounded-xl border p-6 text-center">
          <div className="flex size-12 items-center justify-center rounded-full bg-emerald-500/15 text-emerald-500">
            <CheckCircle2Icon className="size-6" />
          </div>
          <div className="space-y-1">
            <div className="font-semibold">Email verified</div>
            <p className="text-muted-foreground text-sm">Signing you in…</p>
          </div>
        </div>
      )}
      {verify.isError && (
        <div className="border-destructive/20 bg-destructive/5 flex flex-col items-center gap-3 rounded-xl border p-6 text-center">
          <div className="bg-destructive/10 text-destructive flex size-10 items-center justify-center rounded-full">
            <AlertTriangleIcon className="size-5" />
          </div>
          <div className="space-y-1">
            <div className="font-semibold">Verification failed</div>
            <p className="text-muted-foreground text-sm">{verify.error.message}</p>
          </div>
          <Button asChild variant="outline" className="w-full">
            <Link to="/login">Back to sign in</Link>
          </Button>
        </div>
      )}
    </AuthLayout>
  )
}
