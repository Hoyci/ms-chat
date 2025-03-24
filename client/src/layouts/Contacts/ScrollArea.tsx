import React, { useState } from "react";
import UserIcon from "@assets/user.svg?react";
import type { IContact } from "@store/contactStore"
import { messageStatus } from "@types/message";

interface ScrollAreaProps {
  contacts: IContact[]
  onSelectContact: (contact: IContact) => void;
}

function ScrollArea({ contacts, onSelectContact }: ScrollAreaProps) {
  const [localSelected, setLocalSelected] = useState<number | null>(null);

  const handleSelect = (contact: IContact) => {
    setLocalSelected(contact.id);
    onSelectContact(contact);
  };

  return (
    <div className="scroll-area flex flex-col gap-[0.75px] w-full h-[calc(100vh-160px)] overflow-y-auto mt-2">
      {contacts.length ? (
        contacts.map((contact: IContact) => (
          <div
            key={contact.id}
            className={`flex items-center gap-3 pl-3 cursor-pointer 
            transition-all duration-100 ease-in
            ${
              localSelected === contact.id
                ? "bg-primary-50 opacity-100"
                : "hover:bg-primary-100 opacity-90 hover:opacity-100"
            }`}
            onClick={() => handleSelect(contact)}
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
                <span className="text-sm text-gray-500">{contact.lastMessage.timestamp}</span>
              </div>

              <div className="flex items-center gap-1">
                <span className="">
                  {React.createElement(
                    messageStatus[contact.lastMessage.status as keyof typeof messageStatus]
                  )}
                </span>
                <p className="text-gray-600 text-sm truncate">
                  {contact.lastMessage.text}
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
