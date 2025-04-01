import { useEffect } from "react";
import SocketService from "@services/socket";
import { config } from "config";
import { useAuthStore } from "@store/authStore";

export function useSocket() {
  const { accessToken } = useAuthStore();
  useEffect(() => {
    SocketService.connect(`${config.API_URL}/ws?token=${accessToken}`);

    return () => {
      SocketService.disconnect();
    };
  }, []);
}
