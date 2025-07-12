import { ReactNode } from "react";

interface ModalProps {
  isOpen: boolean;
  title?: string;
  onClose: () => void;
  modalName?: string;
  children: ReactNode;
  blurBackground?: boolean;
  closeButtonPosition?: "left" | "right";
}

export default function Modal({
  isOpen,
  title,
  onClose,
  modalName,
  children,
  blurBackground = false,
  closeButtonPosition = "right",
}: ModalProps) {
  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-start justify-center bg-black/50 overflow-y-auto"
      style={blurBackground ? { backdropFilter: "blur(2px)" } : undefined}
      onClick={onClose}
    >
      <div
        className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-4xl mx-1 mt-10 relative max-h-[90vh] flex flex-col"
        onClick={(e) => e.stopPropagation()}
        id={`modal-${modalName || "debug-undefined"}`}
      >
        {/* Header */}
        <div className="mt-2 dark:border-gray-700 relative">
          <button
            onClick={onClose}
            className={`absolute top-4 ${
              closeButtonPosition === "left" ? "left-4" : "right-4"
            } text-gray-500 hover:text-gray-800 dark:hover:text-white text-xl font-bold focus:outline-none`}
            aria-label="Close modal"
            id={`close-${modalName || "debug-undefined"}`}
          >
            &times;
          </button>

          {title && (
            <h2 className="text-2xl font-semibold text-gray-900 dark:text-gray-100 text-center">
              {title}
            </h2>
          )}
        </div>

        {/* Scrollable Body */}
        <div className="p-3 overflow-y-auto max-h-[70vh]">
          {children}
        </div>
      </div>
    </div>
  );
}
