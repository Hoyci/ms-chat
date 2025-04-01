import UserIcon from "@assets/user.svg?react";
import SearchIcon from "@assets/search.svg?react";
import MenuIcon from "@assets/menu.svg?react";
import { IRoom } from "@api/rooms/types";
import { useAuthStore } from "@store/authStore";

function Header({ room, className }: { room: IRoom; className: string }) {
  const { user } = useAuthStore();

  console.log(room);

  const otherParticipant = room.participants.find(
    (participant) => participant.id !== user!.id
  );

  return (
    otherParticipant && (
      <div
        className={`z-10 w-full h-16 bg-primary-100 px-4 py-2.5 flex items-center justify-between text-neutral-400 ${className}`}
      >
        <div className="flex  items-center gap-4">
          {otherParticipant.avatar ? (
            <img
              src={otherParticipant.avatar}
              width={40}
              height={40}
              className="rounded-full"
            />
          ) : (
            <UserIcon width={40} height={40} />
          )}
          <h1 className="font-semibold flex-col">{otherParticipant.name}</h1>
        </div>
        <div className="flex items-center gap-6 text-neutral-100">
          <SearchIcon width={30} height={30} />
          <MenuIcon width={24} height={24} />
        </div>
      </div>
    )
  );
}

export default Header;
