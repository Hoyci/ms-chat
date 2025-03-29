import { useSocket } from "@hooks/useSocket";
import Chat from "@layouts/Chat";
import Contacts from "@layouts/Contacts";
import Header from "@layouts/Header";
import { useAuthStore } from "@store/authStore";

function ChatApp() {
  const { user } = useAuthStore();
  useSocket();

  console.log(user);

  return (
    <div className="h-screen w-full bg-background flex items-center justify-center font-system">
      <div className="flex top-[19px] w-[calc(100%-38px)] max-w-[1600px] h-[calc(100%-38px)]">
        <Header />
        <Contacts />
        <Chat />
      </div>
    </div>
  );
}

export default ChatApp;
