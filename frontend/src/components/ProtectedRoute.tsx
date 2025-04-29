import { Navigate } from "react-router-dom";
import { ReactNode } from "react";
import { useAppSelector } from "@/store/hook";

export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const user = useAppSelector(state => state.auth.user);

  console.log("🧩 [ProtectedRoute] user =", user); // 🔥 เพิ่มตรงนี้ดู

  if (!user) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}
