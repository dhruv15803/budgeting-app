import { motion } from "framer-motion"
import { PiggyBankIcon } from "lucide-react"
import { Link } from "react-router"
import { ThemeToggle } from "@/components/layout/ThemeToggle"

interface AuthLayoutProps {
  title: string
  description?: string
  children: React.ReactNode
  footer?: React.ReactNode
}

export function AuthLayout({ title, description, children, footer }: AuthLayoutProps) {
  return (
    <div className="bg-background flex min-h-screen">
      <div className="relative hidden w-1/2 flex-col justify-between overflow-hidden bg-gradient-to-br from-violet-600 via-fuchsia-500 to-indigo-600 p-10 text-white lg:flex">
        <div className="pointer-events-none absolute -top-24 -right-24 size-96 rounded-full bg-white/10 blur-3xl" />
        <div className="pointer-events-none absolute -bottom-32 -left-24 size-[28rem] rounded-full bg-black/10 blur-3xl" />
        <div className="relative flex items-center gap-2 text-lg font-semibold">
          <PiggyBankIcon className="size-6" />
          Budget
        </div>
        <div className="relative space-y-4">
          <motion.h1
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="text-4xl leading-tight font-semibold tracking-tight"
          >
            Take control of your spending.
          </motion.h1>
          <motion.p
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.1 }}
            className="max-w-md text-white/80"
          >
            Track expenses, schedule recurring bills, and set monthly budgets — all in one simple,
            private workspace.
          </motion.p>
        </div>
        <div className="relative text-xs text-white/60">&copy; {new Date().getFullYear()} Budget</div>
      </div>

      <div className="relative flex w-full flex-col items-center justify-center p-6 lg:w-1/2">
        <div className="absolute top-4 right-4">
          <ThemeToggle />
        </div>
        <motion.div
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.25 }}
          className="w-full max-w-sm space-y-6"
        >
          <div className="flex items-center gap-2 lg:hidden">
            <div className="bg-primary text-primary-foreground flex size-9 items-center justify-center rounded-lg">
              <PiggyBankIcon className="size-5" />
            </div>
            <Link to="/" className="text-lg font-semibold">
              Budget
            </Link>
          </div>
          <div className="space-y-2">
            <h2 className="text-2xl font-semibold tracking-tight">{title}</h2>
            {description && <p className="text-muted-foreground text-sm">{description}</p>}
          </div>
          {children}
          {footer && <div className="text-muted-foreground text-center text-sm">{footer}</div>}
        </motion.div>
      </div>
    </div>
  )
}
