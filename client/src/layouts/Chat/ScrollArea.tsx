import type { IMessage } from "@store/message";
import Message from "./Message";
import chatBg from "@assets/chat-bg.png";

function ScrollArea({ messages }: { messages: IMessage[] }) {
  return (
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
  );
}

export default ScrollArea;
