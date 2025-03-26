import PendingIcon from "@assets/pending.svg?react";
import SentIcon from "@assets/sent.svg?react";
import DeliveredIcon from "@assets/delivered.svg?react";

export const messageStatus = {
  pending: PendingIcon,
  sent: SentIcon,
  delivered: DeliveredIcon,
};

type StatusType = keyof typeof messageStatus;

export type IMessage = {
  id: number;
  text: string;
  sendId: number;
  timestamp: string;
  status: StatusType;
};
