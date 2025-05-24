import { Navigate } from "react-router-dom";
import { ReactNode } from "react";
import { useAppSelector } from "../store/hook";

export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const user = useAppSelector(state => state.auth.user); // ดึง user จาก redux

  if (!user) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}
