import { roomsService } from "@api/rooms/service";
import { CreateRoomPayload, IRoom } from "@api/rooms/types";
import { AxiosError } from "axios";
import { ROOMS } from "mocks/rooms";
import { create } from "zustand";
import { persist } from "zustand/middleware";

type RoomStore = {
  rooms: IRoom[];
  loading: boolean;
  error: string | null;
  selectedRoom: IRoom | null;
  setSelectedRoom: (room: IRoom | null) => void;
  addRoom: (room: CreateRoomPayload) => void;
  updateRoom: (
    id: number,
    updates: Partial<IRoom> | ((room: IRoom) => Partial<IRoom>)
  ) => void;
};

export const useRoomStore = create<RoomStore>()(
  persist(
    (set, get) => ({
      rooms: ROOMS,
      loading: false,
      error: null,
      selectedRoom: null,
      setSelectedRoom: (room) => set({ selectedRoom: room }),
      addRoom: async (payload: CreateRoomPayload) => {
        try {
          const room = await roomsService.createRoom(payload);
          set({ rooms: [...get().rooms, room], loading: false });
        } catch (error) {
          if (error instanceof AxiosError) {
            set({ error: error.response?.data.error, loading: false });
          } else {
            set({ error: "Failed to fetch contacts", loading: false });
          }
        }
      },
      updateRoom: (id, updates) =>
        set((state) => {
          const updatedRooms = state.rooms.map((room) =>
            room.id === id
              ? {
                  ...room,
                  ...(typeof updates === "function" ? updates(room) : updates),
                }
              : room
          );

          return {
            rooms: updatedRooms,
            selectedRoom:
              updatedRooms.find((c) => c.id === state.selectedRoom?.id) || null,
          };
        }),
    }),
    {
      name: "room-store",
      partialize: (state) => ({ rooms: state.rooms }),
    }
  )
);
