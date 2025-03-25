import { IContact } from "@store/contactStore";

export const CONTACTS: IContact[] = [
    {
        id: 1,
        name: "John Doe",
        avatar: null,
        messages: [
            {
                id: 1,
                text: "Opa, bom dia",
                sendId: 1,
                timestamp: new Date().toLocaleDateString(),
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
                text: "E ai, meu faixa!",
                sendId: 1,
                timestamp: new Date().toLocaleDateString(),
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