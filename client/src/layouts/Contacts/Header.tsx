import { useEffect, useRef, useState } from "react";
import NewChatIcon from "@assets/new-chat.svg?react";
import MenuIcon from "@assets/menu.svg?react";
import IconButton from "@components/IconButton";
import { useAuthStore } from "@store/authStore";
import { Link } from "react-router-dom";

const allIcons = [
  {
    icon: NewChatIcon,
    disabled: false,
    tooltip: "Nova conversa",
  },
  {
    icon: MenuIcon,
    disabled: false,
    tooltip: "Menu",
  },
];

const menuItems = [
  { label: "Nova conversa", action: "new-chat" },
  { label: "Desconectar", action: "logout" },
];

function Header() {
  const [selected, setSelected] = useState<number | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const { logout } = useAuthStore();

  const handleSelect = (index: number) => {
    setSelected(selected === index ? null : index);
  };

  const handleMenuAction = (action: string) => {
    if (action === "logout") {
      logout();
    }
    setSelected(null);
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setSelected(null);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <header className="flex items-center h-16">
      <div className="flex items-center justify-between w-full">
        <h1 className="text-[22px] font-bold text-white">Conversas</h1>

        <div className="flex items-center gap-5" ref={containerRef}>
          {allIcons.map(({ icon, disabled, tooltip }, index) => (
            <div key={index} className="relative">
              <IconButton
                index={index}
                Icon={icon}
                disabled={disabled}
                isSelected={selected === index}
                onClick={handleSelect}
                tooltipText={tooltip}
              />

              {index === 1 && selected === 1 && (
                <div className="absolute right-0 top-12 z-50 w-64 bg-primary-100 rounded-sm shadow-xl py-2">
                  {menuItems.map((item) => (
                    <Link
                      key={item.action}
                      to={item.action === "new-chat" ? "/new-chat" : "#"}
                      className="block px-4 py-3 text-white hover:bg-primary-400 transition-colors"
                      onClick={() => handleMenuAction(item.action)}
                    >
                      {item.label}
                    </Link>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </header>
  );
}

export default Header;
