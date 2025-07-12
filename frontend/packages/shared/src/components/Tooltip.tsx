import { useEffect, useRef } from "react";
import { Tooltip } from "flowbite";
import type { TooltipOptions, TooltipInterface } from "flowbite";

interface CustomTooltipProps {
  id: string;
  children: React.ReactNode;
  content: string;
  trigger?: "hover" | "click";
  placement?: "top" | "bottom" | "left" | "right";
  bgColor?: string;
  textColor?: string;
  textSize?: string;
  className?: string;
}

export default function CustomTooltip({
  id,
  children,
  content,
  trigger = "hover",
  placement = "top",
  bgColor = "bg-gray-900",
  textColor = "text-white",
  textSize = "text-sm",
  className = "",
}: CustomTooltipProps) {
  const tooltipRef = useRef<HTMLDivElement>(null);
  const targetRef = useRef<HTMLSpanElement>(null);
  const tooltipInstance = useRef<TooltipInterface>();

  useEffect(() => {
    if (targetRef.current && tooltipRef.current) {
      tooltipInstance.current = new Tooltip(tooltipRef.current, targetRef.current, {
        triggerType: trigger,
        placement,
      });
    }

    return () => {
      tooltipInstance.current?.destroy();
    };
  }, [trigger, placement]);

  return (
    <>
      <span
        ref={targetRef}
        data-tooltip-target={id}
        className={`inline-flex items-center ${className}`}
        role="tooltip-target"
      >
        {children}
      </span>

      <div
        id={id}
        ref={tooltipRef}
        role="tooltip"
        className={`absolute z-10 invisible inline-block px-3 py-2 font-medium rounded-lg shadow-sm opacity-0 tooltip border border-gray-200 ${bgColor} ${textColor} ${textSize}`}
      >
        {content}
        <div className="tooltip-arrow" data-popper-arrow></div>
      </div>
    </>
  );
}
