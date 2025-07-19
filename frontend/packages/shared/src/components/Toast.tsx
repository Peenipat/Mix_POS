import { useEffect, useState } from "react";
import { FiCheckCircle, FiAlertCircle, FiXCircle } from "react-icons/fi";
import { IoWarningOutline } from "react-icons/io5";
import { twMerge } from "tailwind-merge";

interface ToastProps {
    message: string;
    variant?: "success" | "error" | "warning";
    showIcon?: boolean;
    duration?: number;
    position?: "top-left" | "top-right" | "bottom-left" | "bottom-right";
    onClose?: () => void;
    useFixed?: boolean;
    closable?: boolean;
    disableClose?: boolean;
    containerMode?: boolean
}

const ICONS = {
    success: <FiCheckCircle className="w-5 h-5" />,
    error: <FiXCircle className="w-5 h-5" />,
    warning: <IoWarningOutline className="w-5 h-5" />,
};

const POSITION_CLASSES = {
    "top-left": "top-4 left-4",
    "top-right": "top-4 right-4",
    "bottom-left": "bottom-4 left-4",
    "bottom-right": "bottom-4 right-4",
};

const POSITION_CLASSES_CONTAINER = {
    "top-left": "top-4 left-4",
    "top-right": "top-4 right-4",
    "bottom-left": "bottom-4 left-4",
    "bottom-right": "bottom-4 right-4",
};

const COLORS = {
    success: {
        icon: "text-green-500 bg-green-100 dark:bg-green-800 dark:text-green-200",
        ring: "focus:ring-green-300",
    },
    error: {
        icon: "text-red-500 bg-red-100 dark:bg-red-800 dark:text-red-200",
        ring: "focus:ring-red-300",
    },
    warning: {
        icon: "text-orange-500 bg-orange-100 dark:bg-orange-700 dark:text-orange-200",
        ring: "focus:ring-orange-300",
    },
};

export const Toast = ({
    message,
    variant = "success",
    showIcon = true,
    duration = 4000,
    position = "bottom-right", // default
    onClose,
    useFixed,
    closable,
    disableClose,
    containerMode

}: ToastProps) => {
    const [visible, setVisible] = useState(true);

    useEffect(() => {
        if (duration) {
            const timeout = setTimeout(() => {
                setVisible(false);
                onClose?.();
            }, duration);

            return () => clearTimeout(timeout);
        }
    }, [duration, onClose]);

    if (!visible) return null;

    return (
        <div
            className={twMerge(
                containerMode
                    ? "absolute z-10"
                    : useFixed
                        ? "fixed z-50"
                        : "",
                "flex items-center w-full max-w-xs p-4 mb-4 text-gray-600 bg-white rounded-lg dark:text-gray-300 dark:bg-gray-800 border shadow-xl/20",
                containerMode
                    ? POSITION_CLASSES_CONTAINER[position]
                    : useFixed !== false
                        ? POSITION_CLASSES[position]
                        : ""
            )}
            role="alert"
        >

            {showIcon && (
                <div
                    className={twMerge(
                        "inline-flex items-center justify-center shrink-0 w-8 h-8 rounded-lg mr-3",
                        COLORS[variant].icon
                    )}
                >
                    {ICONS[variant]}
                </div>
            )}
            <div className="flex-1 flex items-center text-sm font-medium text-gray-700 leading-5 px-1 py-0.5">
                {message}
            </div>
            {
                closable !== false && (
                    <button
                        type="button"
                        onClick={() => {
                            if (disableClose) return; // <<< ถ้าห้ามปิด → ไม่ทำอะไรเลย
                            setVisible(false);
                            onClose?.();
                        }}
                        disabled={disableClose} // <--- เพื่อ UX ให้ดูเหมือน disable ด้วย
                        className={twMerge(
                            "ml-auto -mx-1.5 -my-1.5 rounded-lg p-1.5 inline-flex items-center justify-center h-8 w-8",
                            "bg-white text-gray-400 hover:text-gray-900 hover:bg-gray-100",
                            "dark:text-gray-500 dark:hover:text-white dark:bg-gray-800 dark:hover:bg-gray-700",
                            COLORS[variant].ring,
                            disableClose && "opacity-50 cursor-not-allowed" // <--- UI disabled
                        )}
                        aria-label="Close"
                    >
                        <svg className="w-3 h-3" viewBox="0 0 14 14" fill="none">
                            <path
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M1 1l6 6m0 0l6 6M7 7l6-6M7 7L1 13"
                            />
                        </svg>
                    </button>
                )
            }
        </div >
    );
};

