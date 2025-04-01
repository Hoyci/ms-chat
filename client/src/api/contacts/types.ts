import { z } from "zod";

export const ContactSchema = z.object({
  id: z.number(),
  name: z.string(),
  statusMessage: z.string().nullish(),
  email: z.string().email("Invalid email address"),
  avatar: z.string().url().nullish(),
});

export const GetContactsSchema = z.object({
  contacts: z.array(ContactSchema),
});

export const CreateContactSchema = ContactSchema.omit({ id: true });

export const DeleteContactSchema = ContactSchema.omit({
  name: true,
  email: true,
  avatar: true,
});

export type IContact = z.infer<typeof ContactSchema>;
export type GetContactsResponse = z.infer<typeof GetContactsSchema>;
export type CreateContactPayload = z.infer<typeof CreateContactSchema>;
export type DeleteContactResponse = z.infer<typeof DeleteContactSchema>;
