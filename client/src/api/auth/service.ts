import apiClient from "../axios";
import {
  LoginPayload,
  SignupPayload,
  AuthResponse,
  RefreshTokenPayload,
} from "./types";

export const authService = {
  async login(payload: LoginPayload): Promise<AuthResponse> {
    const response = await apiClient.post("/auth", payload);
    return response.data;
  },

  async signup(payload: SignupPayload): Promise<AuthResponse> {
    const response = await apiClient.post("/users", payload);
    return response.data;
  },

  async validateToken(payload: RefreshTokenPayload): Promise<AuthResponse> {
    const response = await apiClient.post("/auth/refresh", payload);
    return response.data;
  },
};
