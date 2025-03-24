import React, { useState } from "react";
import UserIcon from "@assets/user.svg?react";
import PendingIcon from "@assets/pending.svg?react";
import SentIcon from "@assets/sent.svg?react";
import DeliveredIcon from "@assets/delivered.svg?react";

const messageStatus = {
  pending: PendingIcon,
  sent: SentIcon,
  delivered: DeliveredIcon,
};

type StatusType = keyof typeof messageStatus;

interface Contact {
  id: number;
  name: string;
  message: string;
  date: string;
  avatar?: string;
  status?: StatusType;
}

const contacts: Contact[] = [
  {
    id: 1,
    name: "Alice Santos",
    message: "Oi, tudo bem?",
    date: "14:30",
    avatar: "https://i.pravatar.cc/1000?img=1",
    status: "pending",
  },
  {
    id: 2,
    name: "Bruno Lima",
    message: "Enviado o relatório!",
    date: "Ontem",
    status: "sent",
  },
  {
    id: 3,
    name: "Carlos Oliveira",
    message: "Vamos marcar a reunião?",
    date: "Segunda-feira",
    avatar: "https://i.pravatar.cc/1000?img=3",
    status: "delivered",
  },
];

function ScrollArea() {
  const [selectedContact, setSelectedContact] = useState<number | null>(null);

  return (
    <div className="scroll-area flex flex-col gap-[0.75px] w-full h-[calc(100vh-160px)] overflow-y-auto mt-2">
      {contacts.length ? (
        contacts.map((contact: Contact) => (
          <div
            key={contact.id}
            className={`flex items-center gap-3 pl-3 cursor-pointer 
            transition-all duration-100 ease-in
            ${
              selectedContact === contact.id
                ? "bg-primary-50 opacity-100"
                : "hover:bg-primary-100 opacity-90 hover:opacity-100"
            }`}
            onClick={() => setSelectedContact(contact.id)}
          >
            <div className="">
              {contact.avatar ? (
                <img
                  src={contact.avatar}
                  alt={contact.name}
                  className="w-12 h-12 rounded-full"
                />
              ) : (
                <UserIcon className="w-12 h-12 rounded-full" />
              )}
            </div>

            <div className="flex-1 border-y-[0.2px] border-primary-100 py-3 pr-3">
              <div className="flex justify-between">
                <span className="font-semibold">{contact.name}</span>
                <span className="text-sm text-gray-500">{contact.date}</span>
              </div>

              <div className="flex items-center gap-1">
                <span className="">
                  {React.createElement(
                    messageStatus[contact.status as keyof typeof messageStatus]
                  )}
                </span>
                <p className="text-gray-600 text-sm truncate">
                  {contact.message}
                </p>
              </div>
            </div>
          </div>
        ))
      ) : (
        <div className="pl-[50px] py-[72px]">
          <p className="text-sm text-neutral-50">
            Nenhuma conversa, contato ou mensagem encontrada
          </p>
        </div>
      )}
    </div>
  );
}

export default ScrollArea;
