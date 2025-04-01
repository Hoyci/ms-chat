import React, { useState } from "react";
import UserIcon from "@assets/user.svg?react";
import { formatMessageDate } from "@utils/date";
import { IRoom, messageStatus } from "@api/rooms/types";
import { useAuthStore } from "@store/authStore";

interface ScrollAreaProps {
  rooms: IRoom[];
  onSelectContact: (contact: IRoom) => void;
}
function ScrollArea({ rooms, onSelectContact }: ScrollAreaProps) {
  const [localSelected, setLocalSelected] = useState<number | null>(null);
  const { user } = useAuthStore();

  console.log(rooms[0].messages[rooms[0].messages.length - 1].status);

  const handleSelect = (contact: IRoom) => {
    setLocalSelected(contact.id);
    onSelectContact(contact);
  };

  return (
    <div className="scroll-area flex flex-col gap-[0.75px] w-full h-[calc(100vh-160px)] overflow-y-auto mt-2">
      {rooms.length ? (
        rooms.map((room: IRoom) => {
          const otherParticipant = room.participants.find(
            (participant) => participant.id !== user!.id
          );

          return (
            <div
              key={room.id}
              className={`flex items-center gap-3 pl-3 cursor-pointer 
              transition-all duration-100 ease-in
              ${
                localSelected === room.id
                  ? "bg-primary-50 opacity-100"
                  : "hover:bg-primary-100 opacity-90 hover:opacity-100"
              }`}
              onClick={() => handleSelect(room)}
            >
              <div className="">
                {otherParticipant ? (
                  <img
                    src={otherParticipant.avatar || ""}
                    alt={otherParticipant.name}
                    className="w-12 h-12 rounded-full"
                  />
                ) : (
                  <UserIcon className="w-12 h-12 rounded-full" />
                )}
              </div>

              <div className="flex-1 border-b-[0.2px] border-primary-100 py-3 pr-3">
                <div className="flex justify-between">
                  <span className="font-semibold">
                    {otherParticipant!.name}
                  </span>
                  {room.messages.length > 0 && (
                    <span className="text-sm text-gray-500">
                      {formatMessageDate(
                        room.messages[room.messages.length - 1].timestamp
                      )}
                    </span>
                  )}
                </div>

                {room.messages.length > 0 ? (
                  <div className="flex items-center gap-1">
                    <span>
                      {React.createElement(
                        messageStatus[
                          room.messages[room.messages.length - 1]
                            .status as keyof typeof messageStatus
                        ]
                      )}
                    </span>
                    <p className="text-gray-600 text-sm truncate">
                      {room.messages[room.messages.length - 1].text}
                    </p>
                  </div>
                ) : (
                  <span className="text-sm">There isn't any message yet</span>
                )}
              </div>
            </div>
          );
        })
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
