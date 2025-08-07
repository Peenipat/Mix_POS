import { ReactNode, useEffect } from "react";
declare global {
  interface Window {
    adsbygoogle: unknown[];
  }
}
interface ModalProps {
  isOpen: boolean;
  title?: string;
  onClose: () => void;
  modalName?: string;
  children: ReactNode;
  blurBackground?: boolean;
  closeButtonPosition?: "left" | "right";
  showAds?: {
    left?: boolean;
    right?: boolean;
    bottom?: boolean;
  };
  size?: "xs" | "sm" | "md" | "lg" | "xl";
}

export default function Modal({
  isOpen,
  title,
  onClose,
  modalName,
  children,
  blurBackground = false,
  closeButtonPosition = "right",
  showAds = {},
  size = "md",
}: ModalProps) {
  if (!isOpen) return null;

  useEffect(() => {
    const timeout = setTimeout(() => {
      try {
        (window as any).adsbygoogle = (window as any).adsbygoogle || [];
        (window as any).adsbygoogle.push({});
      } catch (e) {
        console.error("Adsense error:", e);
      }
    }, 300);

    return () => clearTimeout(timeout);
  }, [isOpen]);
  const sizeClassMap: Record<NonNullable<ModalProps["size"]>, string> = {
    xs: "max-w-xs",
    sm: "max-w-sm",
    md: "max-w-md",
    lg: "max-w-3xl",
    xl: "max-w-5xl",
  };


  return (
    <div
      className="fixed inset-0 z-50 flex items-start justify-center bg-black/50 overflow-y-auto"
      style={blurBackground ? { backdropFilter: "blur(2px)" } : undefined}
      onClick={onClose}
    >
      {showAds?.left && (
        <div className="fixed left-0 top-1/2 -translate-y-1/2 w-[200px] h-[600px] text-white text-sm font-semibold shadow-lg hidden lg:flex items-center justify-center z-[55] pointer-events-none">
          <ins
            className="adsbygoogle"
            style={{ display: "inline-block", width: "200px", height: "600px" }}
            data-ad-client="ca-pub-8579461004845602"
            data-ad-slot="6884165986"
          ></ins>
        </div>
      )}

      {showAds?.right && (
        <div className="fixed right-0 top-1/2 -translate-y-1/2 w-[200px] h-[600px] text-white text-sm font-semibold shadow-lg hidden lg:flex items-center justify-center z-[55] pointer-events-none">
          <ins
            className="adsbygoogle"
            style={{ display: "inline-block", width: "200px", height: "600px" }}
            data-ad-client="ca-pub-8579461004845602"
            data-ad-slot="6884165986"
          ></ins>
        </div>
      )}

      {showAds?.bottom && (
        <div className="fixed bottom-0 left-1/2 -translate-x-1/2 w-[728px] h-[90px] hidden lg:flex items-center justify-center z-[55] bg-transparent">
          <ins
            className="adsbygoogle"
            style={{ display: "block" }}
            data-ad-client="ca-pub-8579461004845602"
            data-ad-slot="5762656002"
            data-ad-format="auto"
            data-full-width-responsive="true"
          ></ins>
        </div>
      )}

      <div
        className={`bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full mx-1 mt-10 relative max-h-[90vh] flex flex-col z-60 ${sizeClassMap[size]}`}
        onClick={(e) => e.stopPropagation()}
        id={`modal-${modalName || "debug-undefined"}`}
      >
        {/* Header */}
        <div className="mt-2 dark:border-gray-700 relative">
          <button
            onClick={onClose}
            className={`absolute top-0 ${closeButtonPosition === "left" ? "left-4" : "right-4"
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
        <div className="p-3 overflow-y-auto max-h-[70vh] relative">
          {children}
        </div>
      </div>
    </div>
  );
}

