import { z } from "zod"

export const recurringSchema = z
  .object({
    title: z.string().trim().min(1, "Required").max(120, "Too long"),
    description: z.string().trim().max(500, "Too long").optional(),
    amount: z.number().positive("Must be greater than 0"),
    category_id: z.number().int().positive("Pick a category"),
    frequency: z.enum(["daily", "weekly", "monthly", "yearly"]),
    start_date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "Invalid date"),
    end_date: z
      .string()
      .regex(/^\d{4}-\d{2}-\d{2}$/, "Invalid date")
      .optional()
      .or(z.literal(""))
      .nullable(),
    is_active: z.boolean(),
  })
  .refine(
    (v) => !v.end_date || v.end_date === "" || v.end_date >= v.start_date,
    { message: "End date must be after start", path: ["end_date"] }
  )

export type RecurringFormValues = z.infer<typeof recurringSchema>
