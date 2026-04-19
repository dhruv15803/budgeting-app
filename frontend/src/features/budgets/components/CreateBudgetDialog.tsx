import { useState } from "react"
import { Loader2Icon } from "lucide-react"
import { useNavigate } from "react-router"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { CurrencyInput } from "@/components/common/CurrencyInput"
import { MonthPicker } from "@/features/budgets/components/MonthPicker"
import { useCreateBudget } from "@/features/budgets/hooks"
import { monthKey } from "@/lib/utils"

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  defaultMonth?: string
}

export function CreateBudgetDialog({ open, onOpenChange, defaultMonth }: Props) {
  const [month, setMonth] = useState(defaultMonth ?? monthKey())
  const [amount, setAmount] = useState<number | null>(null)
  const createMut = useCreateBudget()
  const navigate = useNavigate()

  const submit = async () => {
    if (!amount || amount <= 0) return
    await createMut.mutateAsync({ budget_month: month, total_amount: amount })
    onOpenChange(false)
    navigate(`/budgets/${month}`)
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(o) => {
        onOpenChange(o)
        if (o) {
          setMonth(defaultMonth ?? monthKey())
          setAmount(null)
        }
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New monthly budget</DialogTitle>
          <DialogDescription>Set a total spending target for a given month.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <MonthPicker label="Month" value={month} onChange={setMonth} />
          <div className="grid gap-1.5">
            <Label>Total amount</Label>
            <CurrencyInput value={amount} onValueChange={setAmount} placeholder="3000.00" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={submit} disabled={!amount || createMut.isPending}>
            {createMut.isPending && <Loader2Icon className="animate-spin" />}
            Create
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
