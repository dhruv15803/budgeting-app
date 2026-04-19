import { motion } from "framer-motion"
import { formatDistanceToNowStrict } from "date-fns"
import { MoreHorizontalIcon, PencilIcon, Trash2Icon } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { CategoryBadge } from "@/components/common/CategoryBadge"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { EmptyState } from "@/components/common/EmptyState"
import { Skeleton } from "@/components/ui/skeleton"
import { Switch } from "@/components/ui/switch"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useCategoryMap } from "@/features/categories/hooks"
import { useUpdateRecurring } from "@/features/recurring/hooks"
import { FREQUENCY_LABELS, colorForIndex } from "@/lib/constants"
import { formatCurrency, formatDate } from "@/lib/utils"
import type { RecurringExpense } from "@/types/api"

interface Props {
  items: RecurringExpense[] | undefined
  isLoading: boolean
  onEdit: (item: RecurringExpense) => void
  onDelete: (item: RecurringExpense) => void
}

export function RecurringTable({ items, isLoading, onEdit, onDelete }: Props) {
  const catMap = useCategoryMap()
  const updateMut = useUpdateRecurring()

  if (isLoading && !items) {
    return (
      <div className="space-y-2 p-4">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-12 w-full" />
        ))}
      </div>
    )
  }

  if (!items || items.length === 0) {
    return (
      <EmptyState
        title="No recurring expenses yet"
        description="Add subscriptions or regular bills so they're logged automatically."
      />
    )
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Title</TableHead>
          <TableHead className="hidden md:table-cell">Category</TableHead>
          <TableHead className="w-24">Frequency</TableHead>
          <TableHead className="w-32 hidden sm:table-cell">Next</TableHead>
          <TableHead className="w-28 text-right">Amount</TableHead>
          <TableHead className="w-20 text-center">Active</TableHead>
          <TableHead className="w-12" />
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((r, idx) => {
          const catName = catMap.get(r.category_id) ?? `#${r.category_id}`
          const next = new Date(r.next_occurrence)
          const relative = formatDistanceToNowStrict(next, { addSuffix: true })
          return (
            <motion.tr
              key={r.id}
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.15, delay: Math.min(idx * 0.02, 0.2) }}
              className="hover:bg-muted/50 border-b transition-colors"
            >
              <TableCell className="font-medium">
                <div className="flex flex-col">
                  <span>{r.title}</span>
                  {r.description && (
                    <span className="text-muted-foreground truncate text-xs">{r.description}</span>
                  )}
                </div>
              </TableCell>
              <TableCell className="hidden md:table-cell">
                <CategoryBadge name={catName} color={colorForIndex(r.category_id)} />
              </TableCell>
              <TableCell>
                <Badge variant="outline">{FREQUENCY_LABELS[r.frequency]}</Badge>
              </TableCell>
              <TableCell className="hidden sm:table-cell">
                <div className="flex flex-col">
                  <span className="text-xs tabular-nums">{formatDate(r.next_occurrence)}</span>
                  <span className="text-muted-foreground text-xs">{relative}</span>
                </div>
              </TableCell>
              <TableCell className="text-right font-mono font-medium tabular-nums">
                {formatCurrency(r.amount)}
              </TableCell>
              <TableCell className="text-center">
                <Switch
                  checked={r.is_active}
                  aria-label={r.is_active ? "Pause" : "Activate"}
                  onCheckedChange={(checked) =>
                    updateMut.mutate({
                      id: r.id,
                      payload: {
                        title: r.title,
                        description: r.description,
                        amount: r.amount,
                        category_id: r.category_id,
                        frequency: r.frequency,
                        start_date: r.start_date,
                        end_date: r.end_date,
                        is_active: checked,
                      },
                    })
                  }
                />
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon" className="size-8" aria-label="Actions">
                      <MoreHorizontalIcon />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => onEdit(r)}>
                      <PencilIcon /> Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem variant="destructive" onClick={() => onDelete(r)}>
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
