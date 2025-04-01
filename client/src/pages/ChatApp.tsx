// import { useSocket } from "@hooks/useSocket";
import Chat from "@layouts/Chat";
import Rooms from "@layouts/Rooms";
import Header from "@layouts/Header";
// import { useAuthStore } from "@store/authStore";
import ContactList from "@layouts/ContactsList";
import useLayoutStore from "@store/layoutStore";

const layouts = {
  newChat: ContactList,
  rooms: Rooms,
};

function ChatApp() {
  // useSocket();
  // const { user } = useAuthStore();
  const { currentLayout } = useLayoutStore();

  const LayoutComponent = layouts[currentLayout] || Rooms;

  return (
    <div className="h-screen w-full bg-background flex items-center justify-center font-system">
      <div className="flex top-[19px] w-[calc(100%-38px)] max-w-[1600px] h-[calc(100%-38px)]">
        <Header />
        <LayoutComponent />
        <Chat />
      </div>
    </div>
  );
}

export default ChatApp;
