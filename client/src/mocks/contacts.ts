import { IRoom } from "@store/roomStore";

export const CONTACTS: IRoom[] = [
    {
        id: 1,
        name: "John Doe",
        avatar: null,
        messages: [
            {
                id: 1,
                room_id: 1,
                text: "Opa, bom dia",
                sendId: 1,
                timestamp: new Date("2025-03-20T15:58:06.024Z").toISOString(),
                status: "delivered",
            }
        ]
      },
      {
        id: 2,
        name: "Mac Miller",
        avatar: "https://i.pravatar.cc/100?img=1",
        messages: [
            {
                id: 1,
                room_id: 2,
                text: "E ai, meu faixa!",
                sendId: 1,
                timestamp: new Date("2025-03-26T15:58:06.024Z").toISOString(),
                status: "sent",
            }
        ]
      },
      {
        id: 3,
        name: "Kendrick Lamar",
        avatar: "https://i.pravatar.cc/100?img=2",
        messages: []
      },
]