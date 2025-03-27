import { useRef } from "react";
import { useContactStore } from "@store/roomStore";
import Bottom from "./Bottom";
import Header from "./Header";
import ScrollArea from "./ScrollArea";

function Chat() {
  const { selectedContact } = useContactStore();
  const inputRef = useRef<HTMLInputElement | null>(null);

  return selectedContact ? (
    <div className="relative w-full h-full bg-primary-300 flex flex-col">
      <Header contact={selectedContact} className="relative z-10" />
      <ScrollArea messages={selectedContact.messages} />
      <Bottom inputRef={inputRef} />
    </div>
  ) : (
    <div className="flex-1 bg-primary-300 flex items-center justify-center">
      <p className="text-neutral-200">Select a contact to start chatting</p>
    </div>
  );
}

export default Chat;
