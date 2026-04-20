import { GoogleLogin, type CredentialResponse } from "@react-oauth/google"
import { useLocation, useNavigate } from "react-router"
import { toast } from "sonner"
import { googleAuth } from "@/features/auth/api"
import { useAuth } from "@/features/auth/AuthProvider"

const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined

type GoogleSignInButtonProps = {
  /** Matches Google button copy for sign-in vs registration screens */
  variant?: "signin" | "signup"
}

export function GoogleSignInButton({ variant = "signin" }: GoogleSignInButtonProps) {
  const auth = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const from = (location.state as { from?: { pathname: string } } | null)?.from?.pathname ?? "/"

  if (!googleClientId) {
    return null
  }

  const onSuccess = async (cred: CredentialResponse) => {
    if (!cred.credential) {
      toast.error("Google sign-in did not return a credential.")
      return
    }
    try {
      const { token } = await googleAuth(cred.credential)
      auth.login(token)
      toast.success(variant === "signup" ? "Welcome!" : "Signed in with Google")
      navigate(from, { replace: true })
    } catch (err) {
      toast.error((err as Error).message)
    }
  }

  return (
    <div className="flex w-full justify-center [&>div]:!w-full">
      <GoogleLogin
        onSuccess={onSuccess}
        onError={() => toast.error("Google sign-in was cancelled or failed.")}
        theme="outline"
        size="large"
        width="100%"
        text={variant === "signup" ? "signup_with" : "signin_with"}
        shape="rectangular"
      />
    </div>
  )
}

export function isGoogleSignInConfigured(): boolean {
  return Boolean(googleClientId && googleClientId.length > 0)
}
