import { create } from "zustand";
import { persist } from "zustand/middleware";
import { authService } from "@api/auth/service";
import type {
  AuthResponse,
  LoginPayload,
  RefreshTokenPayload,
  SignupPayload,
} from "@api/auth/types";
import { AxiosError } from "axios";
import socket from "@services/socket";

type AuthState = {
  user: AuthResponse["user"] | null;
  accessToken: string | null;
  refreshToken: string | null;
  loading: boolean;
  error: string | null;
  login: (payload: LoginPayload) => Promise<void>;
  signup: (payload: SignupPayload) => Promise<void>;
  logout: () => void;
  validateToken: (payload: RefreshTokenPayload) => Promise<boolean>;
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      loading: false,
      error: null,

      login: async (payload) => {
        set({ loading: true, error: null });
        try {
          const data = await authService.login(payload);
          set({
            user: data.user,
            accessToken: data.accessToken,
            refreshToken: data.refreshToken,
            loading: false,
          });
        } catch (error) {
          console.log(error);
          if (error instanceof AxiosError) {
            set({
              error: error.response?.data.error,
              loading: false,
            });
          } else {
            set({
              error: error instanceof Error ? error.message : "Login failed",
              loading: false,
            });
          }
        }
      },

      signup: async (payload) => {
        set({ loading: true, error: null });
        try {
          const data = await authService.signup(payload);
          set({
            user: data.user,
            accessToken: data.accessToken,
            refreshToken: data.refreshToken,
            loading: false,
          });
        } catch (error) {
          if (error instanceof AxiosError) {
            set({
              error: error.response?.data.error,
              loading: false,
            });
          } else {
            set({
              error: error instanceof Error ? error.message : "Login failed",
              loading: false,
            });
          }
        }
      },

      logout: () => {
        try {
          socket.disconnect();
          set({ user: null, accessToken: null, refreshToken: null });
        } catch (error) {
          console.log(error);
        }
      },

      validateToken: async (payload) => {
        set({ loading: true });
        try {
          const data = await authService.validateToken(payload);
          set({ user: data.user, accessToken: data.refreshToken });
          return true;
        } catch (error) {
          console.log(error);
          set({ user: null, accessToken: null, refreshToken: null });
          return false;
        } finally {
          set({ loading: false });
        }
      },
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        accessToken: state.accessToken,
        user: state.user,
      }),
    }
  )
);
