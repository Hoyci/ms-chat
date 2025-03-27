import { IMessage } from "@store/message";
import { CONTACTS } from "mocks/contacts";
import { create } from "zustand";

export type IRoom = {
  id: number;
  name: string;
  avatar?: string | null;
  messages: IMessage[];
};

type ContactStore = {
  contacts: IRoom[];
  selectedContact: IRoom | null;
  setSelectedContact: (contact: IRoom | null) => void;
  addContact: (contact: IRoom) => void;
  updateContact: (
    id: number,
    updates: Partial<IRoom> | ((contact: IRoom) => Partial<IRoom>)
  ) => void;
};

export const useContactStore = create<ContactStore>((set) => ({
  contacts: CONTACTS,
  selectedContact: null,
  setSelectedContact: (contact) => set({ selectedContact: contact }),
  addContact: (contact: IRoom) =>
    set((state) => ({ contacts: [...state.contacts, contact] })),
  updateContact: (id, updates) =>
    set((state) => {
      const updatedContacts = state.contacts.map((contact) =>
        contact.id === id
          ? {
              ...contact,
              ...(typeof updates === "function" ? updates(contact) : updates),
            }
          : contact
      );

      return {
        contacts: updatedContacts,
        selectedContact:
          updatedContacts.find((c) => c.id === state.selectedContact?.id) ||
          null,
      };
    }),
}));
