import Header from "./Header";
import ScrollArea from "./ScrollArea";
import { useContactsStore } from "@store/contactStore";

function ContactList() {
  const { contacts } = useContactsStore();

  return (
    <div className="flex-shrink-0 flex-grow-0 basis-[30%] bg-neutral-700 text-neutral-100 px-4 flex flex-col h-full">
      <Header />
      <ScrollArea contacts={contacts} />
    </div>
  );
}

export default ContactList;
