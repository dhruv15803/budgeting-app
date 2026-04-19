import { NavLink } from "react-router"
import {
  LayoutDashboardIcon,
  ReceiptIcon,
  RepeatIcon,
  WalletIcon,
  PiggyBankIcon,
} from "lucide-react"
import { cn } from "@/lib/utils"

interface NavItem {
  to: string
  label: string
  icon: React.ComponentType<{ className?: string }>
  end?: boolean
}

export const NAV_ITEMS: NavItem[] = [
  { to: "/", label: "Dashboard", icon: LayoutDashboardIcon, end: true },
  { to: "/expenses", label: "Expenses", icon: ReceiptIcon },
  { to: "/recurring", label: "Recurring", icon: RepeatIcon },
  { to: "/budgets", label: "Budgets", icon: WalletIcon },
]

export function SidebarBrand({ compact = false }: { compact?: boolean }) {
  return (
    <div className="flex items-center gap-2 px-4 py-5">
      <div className="bg-primary text-primary-foreground flex size-9 items-center justify-center rounded-lg shadow-sm">
        <PiggyBankIcon className="size-5" />
      </div>
      {!compact && (
        <div className="leading-tight">
          <div className="text-foreground font-semibold">Budget</div>
          <div className="text-muted-foreground text-xs">Personal finance</div>
        </div>
      )}
    </div>
  )
}

export function SidebarNav({ onNavigate }: { onNavigate?: () => void }) {
  return (
    <nav className="flex flex-col gap-1 px-3 py-2">
      {NAV_ITEMS.map((item) => (
        <NavLink
          key={item.to}
          to={item.to}
          end={item.end}
          onClick={onNavigate}
          className={({ isActive }) =>
            cn(
              "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
              "hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
              isActive
                ? "bg-sidebar-accent text-sidebar-accent-foreground shadow-sm"
                : "text-muted-foreground"
            )
          }
        >
          <item.icon className="size-4" />
          {item.label}
        </NavLink>
      ))}
    </nav>
  )
}

export function Sidebar() {
  return (
    <aside className="bg-sidebar text-sidebar-foreground hidden h-full w-60 shrink-0 flex-col border-r md:flex">
      <SidebarBrand />
      <SidebarNav />
      <div className="text-muted-foreground mt-auto px-4 py-4 text-xs">
        v1.0 &middot; Budgeting App
      </div>
    </aside>
  )
}
