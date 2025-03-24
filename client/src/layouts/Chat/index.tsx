import { useRef } from "react";
import { useContactStore } from "@store/contactStore";
import Bottom from "./Bottom";
import Header from "./Header";
import ScrollArea from "./ScrollArea";
import { IMessage } from "@types/message";

const messages: IMessage[] = [
  {
    id: 1,
    text: "Olá, tudo bem?",
    sendId: 2,
    timestamp: "15:05",
    status: "delivered",
  },
  {
    id: 2,
    text: "Tudo ótimo, e você?",
    sendId: 1,
    timestamp: "15:05",
    status: "delivered",
  },
  {
    id: 3,
    text: "Estou bem também!",
    sendId: 2,
    timestamp: "15:07",
    status: "sent",
  },
  {
    id: 4,
    text: "Que bom! O que vai fazer hoje?",
    sendId: 2,
    timestamp: "15:09",
    status: "sent",
  },
  {
    id: 5,
    text: "Vou trabalhar e estudar React.",
    sendId: 1,
    timestamp: "15:10",
    status: "sent",
  },
  {
    id: 6,
    text: "Legal! Depois me conta como foi.",
    sendId: 2,
    timestamp: "15:10",
    status: "sent",
  },
  {
    id: 7,
    text: "Pode deixar",
    sendId: 1,
    timestamp: "15:10",
    status: "pending",
  },
];

function Chat() {
  const { selectedContact } = useContactStore();
  const inputRef = useRef<HTMLInputElement | null>(null);

  const sendMessage = () => {
    if (!inputRef.current) return;

    const message = inputRef.current.value.trim();
    if (message === "") return;

    console.log(message)

    inputRef.current.value = "";
  };

  return (
    selectedContact ? (
      <div className="relative w-full h-full bg-primary-300 flex flex-col">
        <Header contact={selectedContact} className="relative z-10" />
        <ScrollArea messages={messages}/>
        <Bottom inputRef={inputRef} sendMessage={sendMessage} />
    </div>
    ) : (
      <div className="flex-1 bg-primary-300 flex items-center justify-center">
          <p className="text-neutral-200">Select a contact to start chatting</p>
      </div>
    ) 
  );
}

export default Chat;
