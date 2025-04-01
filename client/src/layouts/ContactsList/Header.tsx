import ArrowIcon from "@assets/arrow.svg?react";
import Input from "@components/Input";
import useLayoutStore from "@store/layoutStore";
import { useRef } from "react";

function Header() {
  // const {} = useR
  const { setLayout } = useLayoutStore();
  const inputRef = useRef<HTMLInputElement>(null);

  const handleArrowClick = () => {
    setLayout("rooms");
  };

  const handleInputIconClick = () => {
    inputRef.current?.focus();
  };

  return (
    <header className="flex flex-col pt-6">
      <div className="flex items-center gap-6 pl-2">
        <ArrowIcon
          className="hover:cursor-pointer"
          onClick={handleArrowClick}
        />
        <h1 className="text-white">Nova conversa</h1>
      </div>
      <div className="mt-4">
        <Input ref={inputRef} onClick={handleInputIconClick} />
      </div>
      <div></div>
    </header>
  );
}

export default Header;
