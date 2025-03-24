import { useState } from "react";
import MessageIcon from "../assets/message.svg?react";
import StatusIcon from "../assets/status.svg?react";
import ChannelsIcon from "../assets/channels.svg?react";
import CommunitiesIcon from "../assets/communities.svg?react";
import ConfigIcon from "../assets/config.svg?react";
import UserIcon from "../assets/user.svg?react";
import IconButton from "../components/IconButton";

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
    <header className="w-16 h-full px-3 bg-primary-100 text-neutral-100 border-r-2 border-primary-50">
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

export default Header;
