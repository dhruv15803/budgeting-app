import { useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { Link } from "react-router"
import { CheckCircle2Icon, EyeIcon, EyeOffIcon, Loader2Icon, MailIcon } from "lucide-react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { AuthLayout } from "@/features/auth/components/AuthLayout"
import { GoogleSignInButton, isGoogleSignInConfigured } from "@/features/auth/components/GoogleSignInButton"
import { useRegister } from "@/features/auth/hooks"
import { registerSchema, type RegisterValues } from "@/features/auth/schemas"

export function RegisterPage() {
  const form = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { email: "", password: "", username: "" },
  })
  const { mutateAsync, isPending } = useRegister()
  const [done, setDone] = useState<string | null>(null)
  const [showPassword, setShowPassword] = useState(false)

  const onSubmit = async (values: RegisterValues) => {
    try {
      await mutateAsync({
        email: values.email,
        password: values.password,
        username: values.username?.trim() || undefined,
      })
      setDone(values.email)
    } catch (err) {
      toast.error((err as Error).message)
    }
  }

  if (done) {
    return (
      <AuthLayout
        title="Check your inbox"
        footer={
          <>
            Wrong email?{" "}
            <button
              type="button"
              className="text-primary font-medium hover:underline"
              onClick={() => setDone(null)}
            >
              Go back
            </button>
          </>
        }
      >
        <div className="border-primary/20 bg-primary/5 flex flex-col items-center gap-3 rounded-xl border p-6 text-center">
          <div className="bg-primary/15 text-primary flex size-12 items-center justify-center rounded-full">
            <MailIcon className="size-6" />
          </div>
          <div className="space-y-1">
            <div className="flex items-center justify-center gap-2 font-semibold">
              <CheckCircle2Icon className="size-4 text-emerald-500" /> Account created
            </div>
            <p className="text-muted-foreground text-sm">
              We sent a verification link to <span className="text-foreground font-medium">{done}</span>. Open
              it to activate your account.
            </p>
          </div>
          <Button asChild variant="outline" className="mt-2 w-full">
            <Link to="/login">Back to sign in</Link>
          </Button>
        </div>
      </AuthLayout>
    )
  }

  return (
    <AuthLayout
      title="Create account"
      description="Get started tracking your spending in under a minute."
      footer={
        <>
          Already have an account?{" "}
          <Link to="/login" className="text-primary font-medium hover:underline">
            Sign in
          </Link>
        </>
      }
    >
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input type="email" placeholder="you@example.com" autoComplete="email" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Username</FormLabel>
                <FormControl>
                  <Input placeholder="alice" autoComplete="username" {...field} />
                </FormControl>
                <FormDescription>Optional — shown in your profile.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password</FormLabel>
                <FormControl>
                  <div className="relative">
                    <Input
                      type={showPassword ? "text" : "password"}
                      placeholder="At least 6 characters"
                      autoComplete="new-password"
                      className="pr-10"
                      {...field}
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword((s) => !s)}
                      className="text-muted-foreground hover:text-foreground absolute top-1/2 right-2 -translate-y-1/2 rounded p-1"
                      aria-label={showPassword ? "Hide password" : "Show password"}
                    >
                      {showPassword ? <EyeOffIcon className="size-4" /> : <EyeIcon className="size-4" />}
                    </button>
                  </div>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit" className="w-full" disabled={isPending}>
            {isPending && <Loader2Icon className="animate-spin" />}
            Create account
          </Button>
          {isGoogleSignInConfigured() ? (
            <>
              <div className="flex items-center gap-3 py-1">
                <div className="bg-border h-px flex-1" />
                <span className="text-muted-foreground text-xs">or</span>
                <div className="bg-border h-px flex-1" />
              </div>
              <GoogleSignInButton variant="signup" />
            </>
          ) : null}
        </form>
      </Form>
    </AuthLayout>
  )
}
