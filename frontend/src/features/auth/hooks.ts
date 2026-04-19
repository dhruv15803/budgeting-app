import { useMutation } from "@tanstack/react-query"
import { login as apiLogin, register as apiRegister, verifyEmail as apiVerifyEmail } from "@/features/auth/api"
import type { ApiError } from "@/lib/api"
import type { LoginInput, RegisterInput, TokenResponse } from "@/features/auth/api"

export function useLogin() {
  return useMutation<TokenResponse, ApiError, LoginInput>({
    mutationFn: apiLogin,
  })
}

export function useRegister() {
  return useMutation<null, ApiError, RegisterInput>({
    mutationFn: apiRegister,
  })
}

export function useVerifyEmail() {
  return useMutation<TokenResponse, ApiError, string>({
    mutationFn: apiVerifyEmail,
  })
}
