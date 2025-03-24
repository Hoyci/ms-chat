import { IMessage } from '@types/message';
import { create } from 'zustand';

export type IContact = {
    id: number;
    name: string;
    avatar?: string | null;
    lastMessage: IMessage
};

type ContactStore = {
    contacts: IContact[];
    selectedContact: IContact | null;
    setSelectedContact: (contact: IContact | null) => void;
    addContact: (contact: IContact) => void;
    updateContact: (id: number, updates: Partial<IContact>) => void;
};
  

export const useContactStore = create<ContactStore>((set) => ({
    contacts: [
      {
        id: 1,
        name: "John Doe",
        avatar: null,
        lastMessage: {
          id: 1,
          text: "Opa, bom dia",
          sendId: 1,
          timestamp: new Date().toLocaleDateString(),
          status: "sent",
        }
      },
      {
        id: 2,
        name: "Mac Miller",
        avatar: "https://i.pravatar.cc/100?img=1",
        lastMessage: {
          id: 1,
          text: "Opa, bom dia",
          sendId: 1,
          timestamp: new Date().toLocaleDateString(),
          status: "delivered",
        }
      },
      {
        id: 3,
        name: "Kendrick Lamar",
        avatar: "https://i.pravatar.cc/100?img=2",
        lastMessage: {
          id: 1,
          text: "Opa, bom dia",
          sendId: 1,
          timestamp: new Date().toLocaleDateString(),
          status: "pending",
        }
      },
    ],
    selectedContact: null,
    setSelectedContact: (contact) => set({ selectedContact: contact }),
    addContact: (contact) => 
      set((state) => ({ contacts: [...state.contacts, contact] })),
    updateContact: (id, updates) =>
      set((state) => ({
        contacts: state.contacts.map(contact =>
          contact.id === id ? { ...contact, ...updates } : contact
        )
      })),
  }));