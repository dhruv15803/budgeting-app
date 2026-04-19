import { useQuery } from "@tanstack/react-query"
import { fetchCategories } from "@/features/categories/api"
import { qk } from "@/lib/constants"

export function useCategories() {
  return useQuery({
    queryKey: qk.categories,
    queryFn: fetchCategories,
    staleTime: Infinity,
    gcTime: Infinity,
  })
}

export function useCategoryMap() {
  const { data } = useCategories()
  const map = new Map<number, string>()
  for (const c of data ?? []) map.set(c.id, c.category_name)
  return map
}
