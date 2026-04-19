import { useLocation } from "react-router"
import { MobileNav } from "@/components/layout/MobileNav"
import { ThemeToggle } from "@/components/layout/ThemeToggle"
import { UserMenu } from "@/components/layout/UserMenu"
import { NAV_ITEMS } from "@/components/layout/Sidebar"

function usePageTitle() {
  const { pathname } = useLocation()
  if (pathname.startsWith("/budgets/")) return "Budget detail"
  const match = NAV_ITEMS.find((n) => (n.end ? pathname === n.to : pathname.startsWith(n.to)))
  return match?.label ?? "Budget"
}

export function Topbar() {
  const title = usePageTitle()
  return (
    <header className="bg-background/80 sticky top-0 z-30 flex h-14 items-center gap-2 border-b px-3 backdrop-blur md:px-6">
      <MobileNav />
      <h1 className="text-base font-semibold md:text-lg">{title}</h1>
      <div className="ml-auto flex items-center gap-1">
        <ThemeToggle />
        <UserMenu />
      </div>
    </header>
  )
}
