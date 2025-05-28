import * as React from "react";
export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: "primary" | "secondary";
    className?: string;
}
export declare function Button({ variant, className, ...props }: ButtonProps): import("react/jsx-runtime").JSX.Element;
