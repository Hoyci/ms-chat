import apiClient from "../axios";
import { CreateRoomPayload, IRoom } from "./types";

export const roomsService = {
  async getRooms(): Promise<IRoom[]> {
    const response = await apiClient.get("/rooms");
    return response.data;
  },

  async createRoom(payload: CreateRoomPayload): Promise<IRoom> {
    const response = await apiClient.post("/rooms", payload);
    return response.data;
  },
};
