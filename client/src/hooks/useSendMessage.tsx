import { IRoom } from "@store/roomStore";
import { IMessage } from "@store/message";
import { useCallback, useEffect } from "react";

function useSendMessage(
    inputRef: React.RefObject<HTMLInputElement | null>,
    selectedContact: IRoom | null,
    updateContact:  (id: number, updates: Partial<IRoom> | ((contact: IRoom) => Partial<IRoom>)) => void
  ) {
    const sendMessage = useCallback(() => {
      if (!selectedContact || !inputRef.current) return;
  
      const message = inputRef.current.value.trim();
      if (message === "") return;
  
      const messages = selectedContact.messages;
      const lastMessage = messages[messages.length - 1];
      const newId = lastMessage ? lastMessage.id + 1 : 1;
  
      const newMessage: IMessage = {
        id: newId,
        room_id: selectedContact.id,
        sendId: 1,
        text: message,
        status: "pending",
        timestamp: new Date().toISOString(),
      };
  
      console.log(newMessage);
  
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
      return () => document.removeEventListener("keydown", handleKeyDown);
    }, [inputRef, sendMessage]);
  
    return { sendMessage };
  }
  
  export default useSendMessage