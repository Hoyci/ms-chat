import { IContact } from "@api/contacts/types";
import UserIcon from "@assets/user.svg?react";
import { useAuthStore } from "@store/authStore";
import { useRoomStore } from "@store/roomStore";

interface ScrollAreaProps {
  contacts: IContact[];
}

function ScrollArea({ contacts }: ScrollAreaProps) {
  const { rooms, setSelectedRoom, addRoom } = useRoomStore();
  const { user } = useAuthStore();

  const handleSelectContact = (contact: IContact) => {
    console.log("contact", contact);
    console.log("rooms", rooms);
    const room = rooms.find((room) =>
      room.participants.some((participant) => participant.id === contact.id)
    );
    console.log(room);
    if (room) {
      setSelectedRoom(room);
      return;
    }
    addRoom({ participants: [contact.id, user!.id] });
  };

  return (
    <div className="scroll-area flex flex-col gap-[0.75px] w-full h-[calc(100vh-160px)] overflow-y-auto mt-2">
      {contacts.length ? (
        contacts.map((contact) => (
          <div
            key={contact.id}
            onClick={() => handleSelectContact(contact)}
            className="h-[64px] flex items-center gap-3 pl-3 cursor-pointer
              hover:bg-primary-100 opacity-90 hover:opacity-100 
              transition-all duration-100 ease-in"
          >
            <div>
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

            <div className="flex-1 border-b-[0.2px] border-primary-100 py-3 pr-3 h-full">
              <div className="flex">
                <span className="font-semibold text-white truncate">
                  {contact.name}
                </span>
              </div>

              {contact.statusMessage && (
                <div>
                  <span className="text-sm text-neutral-200 truncate">
                    {contact.statusMessage}
                  </span>
                </div>
              )}
            </div>
          </div>
        ))
      ) : (
        <div>
          Nenhum contato. Clique no botão acima para criar um novo contato.
        </div>
      )}
    </div>
  );
}

export default ScrollArea;
