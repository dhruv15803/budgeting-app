import { useEffect } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { format } from "date-fns"
import { CalendarIcon, Loader2Icon } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
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
import { Textarea } from "@/components/ui/textarea"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Calendar } from "@/components/ui/calendar"
import { Switch } from "@/components/ui/switch"
import { CurrencyInput } from "@/components/common/CurrencyInput"
import { useCategories } from "@/features/categories/hooks"
import { useCreateRecurring, useUpdateRecurring } from "@/features/recurring/hooks"
import { recurringSchema, type RecurringFormValues } from "@/features/recurring/schemas"
import { FREQUENCY_LABELS } from "@/lib/constants"
import { cn, toDateInput } from "@/lib/utils"
import type { Frequency, RecurringExpense } from "@/types/api"

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  recurring?: RecurringExpense | null
}

export function RecurringFormDialog({ open, onOpenChange, recurring }: Props) {
  const isEdit = !!recurring
  const { data: categories = [] } = useCategories()
  const createMut = useCreateRecurring()
  const updateMut = useUpdateRecurring()

  const form = useForm<RecurringFormValues>({
    resolver: zodResolver(recurringSchema),
    defaultValues: {
      title: "",
      description: "",
      amount: undefined as unknown as number,
      category_id: undefined as unknown as number,
      frequency: "monthly",
      start_date: toDateInput(new Date()),
      end_date: "",
      is_active: true,
    },
  })

  useEffect(() => {
    if (!open) return
    if (recurring) {
      form.reset({
        title: recurring.title,
        description: recurring.description ?? "",
        amount: Number(recurring.amount),
        category_id: recurring.category_id,
        frequency: recurring.frequency,
        start_date: recurring.start_date,
        end_date: recurring.end_date ?? "",
        is_active: recurring.is_active,
      })
    } else {
      form.reset({
        title: "",
        description: "",
        amount: undefined as unknown as number,
        category_id: undefined as unknown as number,
        frequency: "monthly",
        start_date: toDateInput(new Date()),
        end_date: "",
        is_active: true,
      })
    }
  }, [recurring, open, form])

  const isPending = createMut.isPending || updateMut.isPending

  const onSubmit = async (values: RecurringFormValues) => {
    const payload = {
      title: values.title,
      description: values.description || null,
      amount: values.amount,
      category_id: values.category_id,
      frequency: values.frequency as Frequency,
      start_date: values.start_date,
      end_date: values.end_date ? values.end_date : null,
      is_active: values.is_active,
    }
    if (isEdit && recurring) {
      await updateMut.mutateAsync({ id: recurring.id, payload })
    } else {
      await createMut.mutateAsync(payload)
    }
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit recurring expense" : "New recurring expense"}</DialogTitle>
          <DialogDescription>
            Recurring expenses automatically generate a new expense on each scheduled occurrence.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="title"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Title</FormLabel>
                  <FormControl>
                    <Input placeholder="Netflix" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid gap-4 sm:grid-cols-2">
              <FormField
                control={form.control}
                name="amount"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Amount</FormLabel>
                    <FormControl>
                      <CurrencyInput
                        value={field.value ?? null}
                        onValueChange={(v) => field.onChange(v ?? undefined)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="frequency"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Frequency</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger className="w-full">
                          <SelectValue />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {(["daily", "weekly", "monthly", "yearly"] as const).map((f) => (
                          <SelectItem key={f} value={f}>
                            {FREQUENCY_LABELS[f]}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="category_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Category</FormLabel>
                  <Select
                    value={field.value ? String(field.value) : undefined}
                    onValueChange={(v) => field.onChange(Number(v))}
                  >
                    <FormControl>
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder="Select a category" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {categories.map((c) => (
                        <SelectItem key={c.id} value={String(c.id)}>
                          {c.category_name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid gap-4 sm:grid-cols-2">
              <FormField
                control={form.control}
                name="start_date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Start date</FormLabel>
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button
                          type="button"
                          variant="outline"
                          className={cn("justify-start font-normal", !field.value && "text-muted-foreground")}
                        >
                          <CalendarIcon />
                          {field.value ? format(new Date(field.value), "PPP") : "Pick a date"}
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0" align="start">
                        <Calendar
                          mode="single"
                          selected={field.value ? new Date(field.value) : undefined}
                          onSelect={(d) => d && field.onChange(toDateInput(d))}
                        />
                      </PopoverContent>
                    </Popover>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="end_date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>End date</FormLabel>
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button
                          type="button"
                          variant="outline"
                          className={cn("justify-start font-normal", !field.value && "text-muted-foreground")}
                        >
                          <CalendarIcon />
                          {field.value ? format(new Date(field.value), "PPP") : "Indefinite"}
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0" align="start">
                        <Calendar
                          mode="single"
                          selected={field.value ? new Date(field.value) : undefined}
                          onSelect={(d) => field.onChange(d ? toDateInput(d) : "")}
                        />
                        {field.value && (
                          <div className="border-t p-2">
                            <Button
                              variant="ghost"
                              size="sm"
                              className="w-full"
                              onClick={() => field.onChange("")}
                            >
                              Clear
                            </Button>
                          </div>
                        )}
                      </PopoverContent>
                    </Popover>
                    <FormDescription>Leave empty for indefinite.</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea placeholder="Optional" rows={2} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="is_active"
              render={({ field }) => (
                <FormItem className="flex-row items-center justify-between rounded-lg border p-3">
                  <div className="space-y-0.5">
                    <FormLabel>Active</FormLabel>
                    <FormDescription>Pause the schedule without deleting it.</FormDescription>
                  </div>
                  <FormControl>
                    <Switch checked={field.value} onCheckedChange={field.onChange} />
                  </FormControl>
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending && <Loader2Icon className="animate-spin" />}
                {isEdit ? "Save changes" : "Create"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
