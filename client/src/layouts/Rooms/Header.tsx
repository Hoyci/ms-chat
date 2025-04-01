import { useEffect, useRef, useState } from "react";
import NewChatIcon from "@assets/new-chat.svg?react";
import MenuIcon from "@assets/menu.svg?react";
import IconButton from "@components/IconButton";
import { useAuthStore } from "@store/authStore";
import useLayoutStore from "@store/layoutStore";

const allIcons = [
  {
    name: "newChat",
    icon: NewChatIcon,
    disabled: false,
    tooltip: "Nova conversa",
  },
  {
    name: "Menu",
    icon: MenuIcon,
    disabled: false,
    tooltip: "Menu",
  },
];

const menuItems = [{ label: "Desconectar", action: "logout" }];

function Header() {
  const [selected, setSelected] = useState<number | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const { logout } = useAuthStore();
  const { setLayout } = useLayoutStore();

  const handleSelect = (index: number) => {
    if (allIcons[index].name === "newChat") {
      setLayout("newChat");
    } else {
      setSelected(selected === index ? null : index);
    }
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
                    <span
                      key={item.action}
                      className="block px-4 py-3 text-white hover:bg-primary-400 hover:cursor-pointer transition-colors"
                      onClick={() => handleMenuAction(item.action)}
                    >
                      {item.label}
                    </span>
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
