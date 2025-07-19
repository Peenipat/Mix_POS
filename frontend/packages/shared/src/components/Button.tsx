import * as React from "react";

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary";
  className?: string;
}

export function Button({
  variant = "primary",
  className = "",
  ...props
}: ButtonProps) {
  const base = "px-4 py-2 rounded text-white";
  const colors =
    variant === "primary"
      ? "bg-blue-600 hover:bg-blue-700"
      : "bg-gray-600 hover:bg-gray-700";

  // แก้สัญลักษณ์ template literal ไม่ต้อง escape
  return (
    <button className={`${base} ${colors} ${className}`} {...props} />
  );
}
