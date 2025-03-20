import { useState } from "react";
import MessageIcon from "../assets/message.svg?react";
import StatusIcon from "../assets/status.svg?react";
import ChannelsIcon from "../assets/channels.svg?react";
import CommunitiesIcon from "../assets/communities.svg?react";
import ConfigIcon from "../assets/config.svg?react";
import UserIcon from "../assets/user.svg?react";

const allIcons = [
  {
    icon: MessageIcon,
    disabled: false,
    tooltip: "Conversas",
  },
  {
    icon: StatusIcon,
    disabled: true,
    tooltip: "Status",
  },
  {
    icon: ChannelsIcon,
    disabled: true,
    tooltip: "Canais",
  },
  {
    icon: CommunitiesIcon,
    disabled: true,
    tooltip: "Comunidades",
  },
  {
    icon: ConfigIcon,
    disabled: false,
    tooltip: "Configurações",
  },
  {
    icon: UserIcon,
    disabled: false,
    width: "w-8",
    height: "h-8",
    tooltip: "Perfil",
  },
];

function Header() {
  const [selected, setSelected] = useState<number>(0);

  const topIcons = allIcons.slice(0, -2);
  const bottomIcons = allIcons.slice(-2);

  const handleSelect = (index: number) => setSelected(index);

  return (
    <header className="w-16 h-full px-3 bg-primary-100 text-neutral-100">
      <div className="flex flex-col items-center h-full justify-between py-2.5">
        <div className="flex flex-col items-center gap-2.5">
          {topIcons.map(({ icon, disabled, width, height, tooltip }, index) => (
            <IconButton
              key={index}
              index={index}
              Icon={icon}
              disabled={disabled}
              width={width}
              height={height}
              tooltipText={tooltip}
              isSelected={selected === index}
              onClick={handleSelect}
            />
          ))}
        </div>

        <div className="flex flex-col items-center gap-2.5 h-">
          {bottomIcons.map(
            ({ icon, disabled, width, height, tooltip }, index) => {
              const originalIndex = allIcons.length - 2 + index;
              return (
                <IconButton
                  key={originalIndex}
                  index={originalIndex}
                  Icon={icon}
                  disabled={disabled}
                  width={width}
                  height={height}
                  tooltipText={tooltip}
                  isSelected={selected === originalIndex}
                  onClick={handleSelect}
                />
              );
            }
          )}
        </div>
      </div>
    </header>
  );
}

const IconButton = ({
  index,
  Icon,
  isSelected,
  disabled,
  width,
  height,
  onClick,
  tooltipText = "",
}: {
  index: number;
  Icon: React.FC<React.SVGProps<SVGSVGElement>>;
  isSelected: boolean;
  disabled: boolean;
  width?: string;
  height?: string;
  onClick: (index: number) => void;
  tooltipText?: string;
}) => (
  <div className="group relative">
    <button
      type="button"
      disabled={disabled}
      onClick={() => !disabled && onClick(index)}
      className={`w-10 h-10 flex items-center justify-center
          transition-all duration-200 ${
            disabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"
          } ${isSelected ? "bg-neutral-50 rounded-full" : ""}`}
    >
      <Icon
        className={`${width || "w-6"} ${height || "h-6"} ${
          disabled ? "text-neutral-400" : ""
        }`}
      />
    </button>

    {!disabled && tooltipText && (
      <div
        className="absolute left-full top-1/2 -translate-y-1/2 ml-2
            opacity-0 group-hover:opacity-100 transition-opacity duration-200
            px-3 py-1.5 bg-white text-neutral-600 text-xs font-medium
            rounded-2xl shadow-lg whitespace-nowrap pointer-events-none"
      >
        {tooltipText}
      </div>
    )}
  </div>
);

export default Header;
