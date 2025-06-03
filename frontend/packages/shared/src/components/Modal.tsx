import { ReactNode } from "react";

interface ModalProps {
  isOpen: boolean;
  title?: string;
  onClose: () => void;
  children: ReactNode;
}

export default function Modal({
  isOpen,
  title,
  onClose,
  children,
}: ModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-lg mx-4 p-6">
        {/* หัวเรื่อง + ปุ่มปิด */}
        <div className="flex items-center justify-between mb-4">
          {title && (
            <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
              {title}
            </h2>
          )}
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
          >
            ✕
          </button>
        </div>

        {/* เนื้อหาเฉพาะหน้า (Create, Edit, Delete ฯลฯ) */}
        <div>{children}</div>
      </div>
    </div>
  );
}
