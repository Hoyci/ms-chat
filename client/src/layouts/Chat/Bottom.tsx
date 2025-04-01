import SendIcon from "@assets/send.svg?react";
import PlusIcon from "@assets/plus.svg?react";
import { useRoomStore } from "@store/roomStore";
import useSendMessage from "@hooks/useSendMessage";

interface BottomProps {
  inputRef: React.RefObject<HTMLInputElement | null>;
}

function Bottom({ inputRef }: BottomProps) {
  const { selectedRoom, updateRoom } = useRoomStore();
  const { sendMessage } = useSendMessage(inputRef, selectedRoom, updateRoom);

  return (
    <div className="bg-primary-100 py-3 flex items-center gap-4 text-white px-4">
      <PlusIcon className="text-neutral-200" />
      <input
        ref={inputRef}
        type="text"
        placeholder="Digite sua mensagem..."
        className="w-full p-2 rounded-lg focus:outline-none bg-primary-200 placeholder-gray-500 focus:placeholder-transparent"
      />
      <SendIcon className="text-neutral-200" onClick={sendMessage} />
    </div>
  );
}

export default Bottom;
