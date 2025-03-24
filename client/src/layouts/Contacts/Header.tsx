import { useEffect, useRef, useState } from "react";
import NewChatIcon from "@assets/new-chat.svg?react";
import MenuIcon from "@assets/menu.svg?react";
import IconButton from "@components/IconButton";

const allIcons = [
  {
    icon: NewChatIcon,
    disabled: false,
  },
  {
    icon: MenuIcon,
    disabled: false,
  },
];

function Header() {
  const [selected, setSelected] = useState<number | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const handleSelect = (index: number) => setSelected(index);

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

        <div className="flex items-center gap-5">
          {allIcons.map(({ icon, disabled }, index) => (
            <IconButton
              key={index}
              index={index}
              Icon={icon}
              disabled={disabled}
              isSelected={selected === index}
              onClick={handleSelect}
              ref={containerRef}
            />
          ))}
        </div>
      </div>
    </header>
  );
}

export default Header;
