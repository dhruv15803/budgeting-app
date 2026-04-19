import { api } from "@/lib/api"
import type { User } from "@/types/api"

export interface LoginInput {
  email: string
  password: string
}

export interface RegisterInput {
  email: string
  password: string
  username?: string
}

export interface TokenResponse {
  token: string
}

export async function login(input: LoginInput): Promise<TokenResponse> {
  const { data } = await api.post<TokenResponse>("/auth/login", input)
  return data
}

export async function register(input: RegisterInput): Promise<null> {
  const { data } = await api.post<null>("/auth/register", input)
  return data
}

export async function verifyEmail(token: string): Promise<TokenResponse> {
  const { data } = await api.get<TokenResponse>("/auth/verify-email", {
    params: { token },
  })
  return data
}

export async function fetchMe(): Promise<User> {
  const { data } = await api.get<User>("/auth/me")
  return data
}
