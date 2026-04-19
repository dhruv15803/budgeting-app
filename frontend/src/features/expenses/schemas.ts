import { z } from "zod"

export const expenseSchema = z.object({
  title: z.string().trim().min(1, "Required").max(120, "Too long"),
  description: z.string().trim().max(500, "Too long").optional(),
  amount: z.number().positive("Must be greater than 0"),
  category_id: z.number().int().positive("Pick a category"),
  expense_date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "Invalid date"),
})

export type ExpenseFormValues = z.infer<typeof expenseSchema>
