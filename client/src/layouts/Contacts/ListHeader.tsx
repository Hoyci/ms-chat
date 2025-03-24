import { useRef, useEffect } from "react";
import IconButton from "@components/IconButton";
import SearchIcon from "@assets/search.svg?react";

function ListHeader() {
  const containerRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleIconClick = () => {
    inputRef.current?.focus();
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        inputRef.current?.blur();
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div ref={containerRef}>
      <div className="flex items-center gap-6 w-full bg-primary-100 h-9 rounded-md mt-2 px-2">
        <IconButton
          key={0}
          index={0}
          Icon={SearchIcon}
          disabled={false}
          isSelected={false}
          width={"24"}
          height={"24"}
          onClick={handleIconClick}
        />
        <input
          ref={inputRef}
          type="text"
          placeholder="Pesquisar"
          aria-label="Campo de pesquisa"
          className="w-full bg-transparent border-none outline-none placeholder-gray-500 focus:placeholder-transparent focus:ring-0 focus:outline-none"
        />
      </div>
    </div>
  );
}

export default ListHeader;
