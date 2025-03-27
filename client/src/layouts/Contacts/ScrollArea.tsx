import React, { useState } from "react";
import UserIcon from "@assets/user.svg?react";
import type { IRoom } from "@store/roomStore";
import { messageStatus } from "@store/message";
import { formatMessageDate } from "@utils/date";

interface ScrollAreaProps {
  contacts: IRoom[];
  onSelectContact: (contact: IRoom) => void;
}

function ScrollArea({ contacts, onSelectContact }: ScrollAreaProps) {
  const [localSelected, setLocalSelected] = useState<number | null>(null);

  const handleSelect = (contact: IRoom) => {
    setLocalSelected(contact.id);
    onSelectContact(contact);
  };

  return (
    <div className="scroll-area flex flex-col gap-[0.75px] w-full h-[calc(100vh-160px)] overflow-y-auto mt-2">
      {contacts.length ? (
        contacts.map((contact) => (
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
                {contact.messages.length > 0 && (
                  <span className="text-sm text-gray-500">
                    {formatMessageDate(contact.messages[contact.messages.length - 1].timestamp)}
                  </span>
                )}
              </div>

              {contact.messages.length > 0 ? (
                <div className="flex items-center gap-1">
                  <span className="">
                    {React.createElement(
                      messageStatus[
                        contact.messages[contact.messages.length - 1]
                          .status as keyof typeof messageStatus
                      ]
                    )}
                  </span>
                  <p className="text-gray-600 text-sm truncate">
                    {contact.messages[contact.messages.length - 1].text}
                  </p>
                </div>
              ) : (
                <span className="text-sm">There isn't any message yet</span>
              )}
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
