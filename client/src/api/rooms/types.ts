import { z } from "zod";
import { ContactSchema } from "@api/contacts/types";
import PendingIcon from "@assets/pending.svg?react";
import SentIcon from "@assets/sent.svg?react";
import DeliveredIcon from "@assets/delivered.svg?react";

export const messageStatus = {
  pending: PendingIcon,
  sent: SentIcon,
  delivered: DeliveredIcon,
} as const;

const statusSchema = z.enum(["pending", "sent", "delivered"]);

export const MessageSchema = z.object({
  id: z.number(),
  text: z.string(),
  room_id: z.number(),
  sendId: z.number(),
  timestamp: z.date(),
  status: statusSchema,
});

export const RoomSchema = z.object({
  id: z.number(),
  participants: z.array(ContactSchema),
  messages: z.array(MessageSchema),
});

export const CreateRoomSchema = z.object({
  participants: z.array(z.number()),
});

export type IMessage = z.infer<typeof MessageSchema>;
export type IRoom = z.infer<typeof RoomSchema>;
export type CreateRoomPayload = z.infer<typeof CreateRoomSchema>;
