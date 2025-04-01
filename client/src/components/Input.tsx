import IconButton from "@components/IconButton";
import SearchIcon from "@assets/search.svg?react";

function Input({
  ref,
  onClick,
}: {
  ref: React.RefObject<HTMLInputElement | null>;
  onClick: () => void;
}) {
  return (
    <div className="flex items-center gap-6 w-full bg-primary-100 h-9 rounded-md mt-2 px-2">
      <IconButton
        key={0}
        index={0}
        Icon={SearchIcon}
        disabled={false}
        isSelected={false}
        width={"24"}
        height={"24"}
        onClick={onClick}
      />
      <input
        ref={ref}
        type="text"
        placeholder="Pesquisar"
        aria-label="Campo de pesquisa"
        className="w-full bg-transparent border-none outline-none placeholder-gray-500 focus:placeholder-transparent focus:ring-0 focus:outline-none"
      />
    </div>
  );
}

export default Input;
