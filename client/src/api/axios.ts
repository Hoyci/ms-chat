import axios from "axios";
import snakecaseKeys from "snakecase-keys";
import camelcaseKeys from "camelcase-keys";
import { useAuthStore } from "@store/authStore";
import { config } from "config";

const apiClient = axios.create({
  baseURL: config.API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  const accessToken = useAuthStore.getState().accessToken;
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

apiClient.interceptors.request.use((config) => {
  if (config.data) {
    config.data = snakecaseKeys(config.data, { deep: true });
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => {
    if (response.data) {
      response.data = camelcaseKeys(response.data, { deep: true });
    }
    return response;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(undefined, (error) => {
  if (error.response?.status === 401) {
    const refreshToken = useAuthStore.getState().refreshToken;
    if (refreshToken) {
      useAuthStore.getState().validateToken({ refreshToken });
    }
  }
  return Promise.reject(error);
});

export default apiClient;
