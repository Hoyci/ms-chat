import { useContactStore } from "@store/roomStore";
import ListHeader from "./ListHeader";
import ScrollArea from "./ScrollArea";

function List() {
  const { contacts, setSelectedContact } = useContactStore();

  return (
    <div className="flex-1">
      <ListHeader />
      <ScrollArea contacts={contacts} onSelectContact={setSelectedContact} />
    </div>
  );
}

export default List;
