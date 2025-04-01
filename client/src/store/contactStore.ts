import { create } from "zustand";
import { persist } from "zustand/middleware";
// import { contactService } from "@api/contacts/service";
// import type { Contact } from "@api/contacts/types";
import { AxiosError } from "axios";
import { contactsService } from "@api/contacts/service";
import { CreateContactPayload, IContact } from "@api/contacts/types";
import { CONTACTS } from "mocks/contacts";

type ContactsState = {
  contacts: IContact[];
  loading: boolean;
  error: string | null;
  fetchContacts: () => Promise<void>;
  addContact: (contact: CreateContactPayload) => Promise<void>;
  deleteContact: (id: number) => Promise<void>;
};

export const useContactsStore = create<ContactsState>()(
  persist(
    (set, get) => ({
      contacts: CONTACTS,
      loading: false,
      error: null,

      fetchContacts: async () => {
        set({ loading: true, error: null });
        try {
          const { contacts } = await contactsService.getContacts();
          set({ contacts, loading: false });
        } catch (error) {
          if (error instanceof AxiosError) {
            set({ error: error.response?.data.error, loading: false });
          } else {
            set({ error: "Failed to fetch contacts", loading: false });
          }
        }
      },

      addContact: async (contact) => {
        set({ loading: true, error: null });
        try {
          const newContact = await contactsService.createContact(contact);
          set({ contacts: [...get().contacts, newContact], loading: false });
        } catch (error) {
          if (error instanceof AxiosError) {
            set({ error: error.response?.data.error, loading: false });
          } else {
            set({ error: "Failed to add contact", loading: false });
          }
        }
      },

      deleteContact: async (id) => {
        set({ loading: true, error: null });
        try {
          await contactsService.deleteContact(id);
          set({
            contacts: get().contacts.filter((c) => c.id !== id),
            loading: false,
          });
        } catch (error) {
          if (error instanceof AxiosError) {
            set({ error: error.response?.data.error, loading: false });
          } else {
            set({ error: "Failed to delete contact", loading: false });
          }
        }
      },
    }),
    {
      name: "contacts-storage",
      partialize: (state) => ({ contacts: state.contacts }),
    }
  )
);
