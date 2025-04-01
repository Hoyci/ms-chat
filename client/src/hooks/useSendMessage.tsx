import { IRoom } from "@api/rooms/types";
import { IMessage } from "@store/message";
import { useCallback, useEffect } from "react";

function useSendMessage(
  inputRef: React.RefObject<HTMLInputElement | null>,
  selectedRoom: IRoom | null,
  updateRoom: (
    id: number,
    updates: Partial<IRoom> | ((contact: IRoom) => Partial<IRoom>)
  ) => void
) {
  const sendMessage = useCallback(() => {
    if (!selectedRoom || !inputRef.current) return;

    const message = inputRef.current.value.trim();
    if (message === "") return;

    const messages = selectedRoom.messages;
    const lastMessage = messages[messages.length - 1];
    const newId = lastMessage ? lastMessage.id + 1 : 1;

    const newMessage: IMessage = {
      id: newId,
      room_id: selectedRoom.id,
      sendId: 1,
      text: message,
      status: "pending",
      timestamp: new Date().toISOString(),
    };

    updateRoom(selectedRoom.id, (prevContact) => ({
      messages: [...prevContact.messages, newMessage],
    }));

    inputRef.current.value = "";
  }, [inputRef, selectedRoom, updateRoom]);

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
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [inputRef, sendMessage]);

  return { sendMessage };
}

export default useSendMessage;
