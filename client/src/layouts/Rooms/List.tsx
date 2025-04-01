import { useRoomStore } from "@store/roomStore";
import ListHeader from "./ListHeader";
import ScrollArea from "./ScrollArea";

function List() {
  const { rooms, setSelectedRoom } = useRoomStore();

  return (
    <div className="flex-1">
      <ListHeader />
      <ScrollArea rooms={rooms} onSelectContact={setSelectedRoom} />
    </div>
  );
}

export default List;
