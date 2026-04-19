import axios, { AxiosError, type AxiosInstance } from "axios"
import { LOGOUT_EVENT, TOKEN_STORAGE_KEY } from "@/lib/constants"

export class ApiError extends Error {
  status?: number
  constructor(message: string, status?: number) {
    super(message)
    this.name = "ApiError"
    this.status = status
  }

  static fromAxios(err: unknown): ApiError {
    if (axios.isAxiosError(err)) {
      const ax = err as AxiosError<{ success?: boolean; message?: string; status_code?: number }>
      const status = ax.response?.status
      const body = ax.response?.data
      if (body && typeof body === "object" && "message" in body && typeof body.message === "string") {
        return new ApiError(body.message, status)
      }
      if (ax.code === "ERR_NETWORK") {
        return new ApiError("Can't reach the server. Is the backend running?", status)
      }
      return new ApiError(ax.message || "Request failed", status)
    }
    if (err instanceof Error) return new ApiError(err.message)
    return new ApiError("Unexpected error")
  }
}

const baseURL = (import.meta.env.VITE_API_URL as string | undefined) ?? "http://localhost:3000/api"

export const api: AxiosInstance = axios.create({
  baseURL,
  headers: { "Content-Type": "application/json" },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_STORAGE_KEY)
  if (token) {
    config.headers = config.headers ?? {}
    ;(config.headers as Record<string, string>).Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (res) => {
    const body = res.data
    if (body && typeof body === "object") {
      if ("success" in body) {
        if (body.success === false) {
          return Promise.reject(
            new ApiError(body.message ?? "Request failed", body.status_code ?? res.status)
          )
        }
        if ("data" in body) {
          return { ...res, data: body.data }
        }
        const { success: _s, message: _m, ...rest } = body as Record<string, unknown>
        return { ...res, data: rest }
      }
    }
    return res
  },
  (err) => {
    const apiErr = ApiError.fromAxios(err)
    if (apiErr.status === 401) {
      localStorage.removeItem(TOKEN_STORAGE_KEY)
      window.dispatchEvent(new Event(LOGOUT_EVENT))
    }
    return Promise.reject(apiErr)
  }
)

export function stringifyParams(
  params: Record<string, unknown>
): Record<string, string | number | boolean | Array<string | number>> {
  const out: Record<string, string | number | boolean | Array<string | number>> = {}
  for (const [k, v] of Object.entries(params)) {
    if (v == null || v === "" || (Array.isArray(v) && v.length === 0)) continue
    out[k] = v as string | number | boolean | Array<string | number>
  }
  return out
}
