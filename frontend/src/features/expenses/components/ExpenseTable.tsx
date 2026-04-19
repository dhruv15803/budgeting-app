import { motion } from "framer-motion"
import { MoreHorizontalIcon, PencilIcon, RepeatIcon, Trash2Icon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { CategoryBadge } from "@/components/common/CategoryBadge"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"
import { EmptyState } from "@/components/common/EmptyState"
import { useCategoryMap } from "@/features/categories/hooks"
import { colorForIndex } from "@/lib/constants"
import { formatCurrency, formatDate } from "@/lib/utils"
import type { Expense } from "@/types/api"

interface ExpenseTableProps {
  expenses: Expense[] | undefined
  isLoading: boolean
  onEdit: (expense: Expense) => void
  onDelete: (expense: Expense) => void
}

export function ExpenseTable({ expenses, isLoading, onEdit, onDelete }: ExpenseTableProps) {
  const catMap = useCategoryMap()

  if (isLoading && !expenses) {
    return (
      <div className="space-y-2 p-4">
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-12 w-full" />
        ))}
      </div>
    )
  }

  if (!expenses || expenses.length === 0) {
    return (
      <EmptyState
        title="No expenses yet"
        description="Add your first expense to start tracking your spending."
      />
    )
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-28">Date</TableHead>
          <TableHead>Title</TableHead>
          <TableHead className="hidden md:table-cell">Category</TableHead>
          <TableHead className="hidden lg:table-cell">Description</TableHead>
          <TableHead className="w-28 text-right">Amount</TableHead>
          <TableHead className="w-12" />
        </TableRow>
      </TableHeader>
      <TableBody>
        {expenses.map((e, idx) => {
          const catName = catMap.get(e.category_id) ?? `#${e.category_id}`
          return (
            <motion.tr
              key={e.id}
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.15, delay: Math.min(idx * 0.02, 0.2) }}
              className="hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors"
            >
              <TableCell className="text-muted-foreground text-xs tabular-nums">
                {formatDate(e.expense_date)}
              </TableCell>
              <TableCell className="font-medium">
                <div className="flex items-center gap-2">
                  <span className="truncate">{e.title}</span>
                  {e.recurring_expense_id && (
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <span className="text-muted-foreground">
                          <RepeatIcon className="size-3.5" />
                        </span>
                      </TooltipTrigger>
                      <TooltipContent>Generated from recurring</TooltipContent>
                    </Tooltip>
                  )}
                </div>
              </TableCell>
              <TableCell className="hidden md:table-cell">
                <CategoryBadge name={catName} color={colorForIndex(e.category_id)} />
              </TableCell>
              <TableCell className="text-muted-foreground hidden max-w-xs truncate lg:table-cell">
                {e.description ?? "—"}
              </TableCell>
              <TableCell className="text-right font-mono font-medium tabular-nums">
                {formatCurrency(e.amount)}
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon" className="size-8" aria-label="Actions">
                      <MoreHorizontalIcon />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => onEdit(e)}>
                      <PencilIcon /> Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem variant="destructive" onClick={() => onDelete(e)}>
                      <Trash2Icon /> Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </motion.tr>
          )
        })}
      </TableBody>
    </Table>
  )
}
