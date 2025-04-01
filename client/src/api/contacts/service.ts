import apiClient from "../axios";
import {
  IContact,
  CreateContactPayload,
  DeleteContactResponse,
  GetContactsResponse,
} from "./types";

export const contactsService = {
  async getContacts(): Promise<GetContactsResponse> {
    const response = await apiClient.get("/contacts");
    return response.data;
  },

  async createContact(payload: CreateContactPayload): Promise<IContact> {
    const response = await apiClient.post("/contacts", payload);
    return response.data;
  },

  async deleteContact(id: IContact["id"]): Promise<DeleteContactResponse> {
    const response = await apiClient.delete(`/contacts/${id}`);
    return response.data;
  },
};
