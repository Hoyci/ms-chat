import Header from "./Header";
import chatBg from "@assets/chat-bg.png";
import SendIcon from "@assets/send.svg?react";
import PlusIcon from "@assets/plus.svg?react";
import PendingIcon from "@assets/pending.svg?react";
import SentIcon from "@assets/sent.svg?react";
import DeliveredIcon from "@assets/delivered.svg?react";

const messages = [
  {
    id: 1,
    text: "Olá, tudo bem?",
    isMine: false,
    timestamp: "15:05",
    status: "read",
  },
  {
    id: 2,
    text: "Tudo ótimo, e você?",
    isMine: true,
    timestamp: "15:05",
    status: "read",
  },
  {
    id: 3,
    text: "Estou bem também!",
    isMine: false,
    timestamp: "15:07",
    status: "delivered",
  },
  {
    id: 4,
    text: "Que bom! O que vai fazer hoje?",
    isMine: false,
    timestamp: "15:09",
    status: "delivered",
  },
  {
    id: 5,
    text: "Vou trabalhar e estudar React.",
    isMine: true,
    timestamp: "15:10",
    status: "sent",
  },
  {
    id: 6,
    text: "Legal! Depois me conta como foi.",
    isMine: false,
    timestamp: "15:10",
    status: "sent",
  },
  {
    id: 7,
    text: "Pode deixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadeixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadeixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadeixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadeixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadeixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadeixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadeixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadeixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadeixaraaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa!",
    isMine: true,
    timestamp: "15:10",
    status: "pending",
  },
];
function Message({
  message,
}: {
  message: {
    id: number;
    text: string;
    isMine: boolean;
    timestamp: string;
    status: "pending" | "sent" | "delivered" | "read";
  };
}) {
  return (
    <div
      className={`max-w-xs p-2 rounded-lg flex flex-col ${
        message.isMine
          ? "bg-secondary-100 text-white self-end"
          : "bg-primary-100 text-white self-start"
      }`}
    >
      <p className="break-words">{message.text}</p>

      <div className="flex justify-end items-center gap-1 mt-1">
        <span className="text-xs opacity-80">{message.timestamp}</span>

        {message.status === "pending" ? (
          <PendingIcon className="w-3 h-3 fill-current" />
        ) : message.status === "sent" ? (
          <SentIcon className="w-3 h-3 fill-current" />
        ) : message.status === "delivered" ? (
          <DeliveredIcon className="w-3 h-3 fill-current" />
        ) : (
          // Status read usa o mesmo ícone de delivered com cor azul
          <DeliveredIcon className="w-3 h-3 fill-current text-blue-400" />
        )}
      </div>
    </div>
  );
}

function Chat() {
  return (
    <div className="relative w-full h-full bg-primary-300 flex flex-col">
      <Header className="relative z-10" />

      <div className="relative flex-1 min-h-0 overflow-hidden">
        <div
          className="absolute inset-0 bg-repeat bg-center opacity-5"
          style={{ backgroundImage: `url(${chatBg})` }}
        ></div>
        <div className="scroll-area relative z-10 flex flex-col p-4 space-y-2 h-full overflow-y-auto">
          {messages.map((msg) => (
            <Message key={msg.id} message={msg} />
          ))}
        </div>
      </div>

      {/* Restante do código permanece igual */}
      <div className="bg-primary-100 py-3 flex items-center gap-4 text-white px-4">
        <PlusIcon className="text-neutral-200" />
        <input
          type="text"
          placeholder="Digite sua mensagem..."
          className="w-full p-2 rounded-lg focus:outline-none bg-primary-200 placeholder-gray-500 focus:placeholder-transparent"
        />
        <SendIcon className="text-neutral-200" />
      </div>
    </div>
  );
}

export default Chat;
