const IconButton = ({
  index,
  Icon,
  isSelected,
  disabled,
  width,
  height,
  onClick,
  tooltipText = "",
  ref,
}: {
  index: number;
  Icon: React.FC<React.SVGProps<SVGSVGElement>>;
  isSelected?: boolean;
  disabled: boolean;
  width?: string;
  height?: string;
  onClick: (index: number) => void;
  tooltipText?: string;
  ref?: React.RefObject<HTMLDivElement | null>;
}) => (
  <div className="group relative" ref={ref}>
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
              rounded-2xl shadow-lg whitespace-nowrap pointer-events-none z-100"
      >
        {tooltipText}
      </div>
    )}
  </div>
);

export default IconButton;
