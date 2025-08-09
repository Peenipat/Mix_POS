import { createRoot } from "react-dom/client";
import { Toast } from "@object/shared/components/Toast";

type ToastVariant = "success" | "error" | "warning";

type MakeToastParams = {
  message: string;
  variant?: ToastVariant;
  duration?: number;
  position?: "top-left" | "top-right" | "bottom-left" | "bottom-right";
};

export function makeToast({
  message,
  variant = "success",
  duration = 4000,
  position = "top-right",
}: MakeToastParams) {
  // สร้าง container ชั่วคราว
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);

  const handleClose = () => {
    root.unmount();
    container.remove();
  };

  // Render Toast เข้า document.body
  root.render(
    <Toast
      message={message}
      variant={variant}
      duration={duration}
      position={position}
      onClose={handleClose}
      useFixed
    />
  );
}
