import { IRoom } from "@api/rooms/types";

export const ROOMS: IRoom[] = [
  {
    id: 1,
    participants: [
      {
        id: 1,
        name: "Ruan Ribeiro",
        email: "ruan@email.com",
        statusMessage: "Apenas ligação",
        avatar: null,
      },
      {
        id: 2,
        name: "Mac Miller",
        email: "mac@email.com",
        avatar: "https://i.pravatar.cc/100?img=1",
      },
    ],
    messages: [
      {
        id: 1,
        room_id: 1,
        text: "Opa, bom dia",
        sendId: 1,
        timestamp: new Date("2025-03-20T15:58:06.024Z"),
        status: "delivered",
      },
    ],
  },
  {
    id: 2,
    participants: [
      {
        id: 1,
        name: "Ruan Ribeiro",
        email: "ruan@email.com",
        statusMessage: "Apenas ligação",
        avatar: null,
      },
      {
        id: 3,
        name: "Kendrick Lamar",
        email: "kdot@email.com",
        statusMessage: "Remember the fist time you came to the house",
        avatar: "https://i.pravatar.cc/100?img=2",
      },
    ],
    messages: [
      {
        id: 1,
        room_id: 2,
        text: "E ai, meu faixa!",
        sendId: 1,
        timestamp: new Date("2025-03-26T15:58:06.024Z"),
        status: "sent",
      },
    ],
  },
  {
    id: 3,
    participants: [
      {
        id: 1,
        name: "Ruan Ribeiro",
        email: "ruan@email.com",
        statusMessage: "Apenas ligação",
        avatar: null,
      },
      {
        id: 4,
        name: "Tim Henson",
        email: "timtim@email.com",
        avatar: "https://i.pravatar.cc/100?img=5",
      },
    ],
    messages: [],
  },
];
