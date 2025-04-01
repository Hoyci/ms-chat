import PendingIcon from "@assets/pending.svg?react";
import SentIcon from "@assets/sent.svg?react";
import DeliveredIcon from "@assets/delivered.svg?react";
import { format } from "date-fns";
import { IMessage } from "@api/rooms/types";

function Message({ message }: { message: IMessage }) {
  return (
    <div
      className={`max-w-xs p-2 rounded-lg flex flex-col ${
        message.sendId === 1
          ? "bg-secondary-100 text-white self-end"
          : "bg-primary-100 text-white self-start"
      }`}
    >
      <p className="break-words">{message.text}</p>

      <div className="flex justify-end items-center gap-1 mt-1">
        <span className="text-xs opacity-80">
          {format(message.timestamp, "HH:mm")}
        </span>

        {message.status === "pending" ? (
          <PendingIcon className="w-3 h-3 fill-current" />
        ) : message.status === "sent" ? (
          <SentIcon className="w-3 h-3 fill-current" />
        ) : message.status === "delivered" ? (
          <DeliveredIcon className="w-3 h-3 fill-current" />
        ) : (
          <DeliveredIcon className="w-3 h-3 fill-current text-blue-400" />
        )}
      </div>
    </div>
  );
}

export default Message;
