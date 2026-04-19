import { api } from "@/lib/api"
import type { Category } from "@/types/api"

export async function fetchCategories(): Promise<Category[]> {
  const { data } = await api.get<{ categories: Category[] } | Category[]>("/expense-categories")
  if (Array.isArray(data)) return data
  return data.categories ?? []
}
