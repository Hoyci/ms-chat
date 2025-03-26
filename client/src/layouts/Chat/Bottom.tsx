import SendIcon from "@assets/send.svg?react";
import PlusIcon from "@assets/plus.svg?react";
import { useCallback, useEffect } from "react";
import { useContactStore } from "@store/contactStore";
import { IMessage } from "@store/message";

interface BottomProps {
  inputRef: React.RefObject<HTMLInputElement | null>;
}

function Bottom({ inputRef }: BottomProps) {
  const { selectedContact, updateContact } = useContactStore();

  const sendMessage = useCallback(() => {
    if (!selectedContact) return;

    if (!inputRef.current) return;

    const message = inputRef.current.value.trim();
    if (message === "") return;

    const messages = selectedContact.messages;
    const lastMessage = messages[messages.length - 1];
    const newId = lastMessage ? lastMessage.id + 1 : 1;

    const newMessage: IMessage = {
      id: newId,
      sendId: 1,
      text: inputRef.current.value,
      status: "pending",
      timestamp: new Date().toLocaleTimeString("pt-BR"),
    };

    updateContact(selectedContact.id, (prevContact) => ({
      messages: [...prevContact.messages, newMessage],
    }));

    inputRef.current.value = "";
  }, [inputRef, selectedContact, updateContact]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (!inputRef.current) return;

      if (
        (event.key === "Enter" || event.code === "NumpadEnter") &&
        inputRef.current.value.trim() !== ""
      ) {
        event.preventDefault();
        sendMessage();
      }
    };

    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [inputRef, sendMessage]);

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
