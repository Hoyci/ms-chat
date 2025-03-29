import { z } from "zod";

export const LoginSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(6, "Password must be at least 6 characters long"),
});

export const SignupSchema = z
  .object({
    username: z.string(),
    email: z.string().email("Invalid email address"),
    password: z.string().min(6, "Password must be at least 6 characters long"),
    confirmPassword: z
      .string()
      .min(6, "Password must be at least 6 characters long"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

export const RefreshTokenSchema = z.object({
  refreshToken: z.string().jwt({
    alg: "RS256",
    message: "Invalid refresh JWT token",
  }),
});

export const UserSchema = z.object({
  id: z.number(),
  username: z.string(),
  email: z.string(),
  createdAt: z.string(),
  deletedAt: z.string().optional(),
  updatedAt: z.string().optional(),
});

export const AuthResponseSchema = z.object({
  user: UserSchema,
  accessToken: z.string().jwt({
    alg: "RS256",
    message: "Invalid refresh JWT token",
  }),
  refreshToken: z.string().jwt({
    alg: "RS256",
    message: "Invalid refresh JWT token",
  }),
});

export type LoginPayload = z.infer<typeof LoginSchema>;
export type SignupPayload = z.infer<typeof SignupSchema>;
export type RefreshTokenPayload = z.infer<typeof RefreshTokenSchema>;
export type UserPayload = z.infer<typeof UserSchema>;
export type AuthResponse = z.infer<typeof AuthResponseSchema>;
