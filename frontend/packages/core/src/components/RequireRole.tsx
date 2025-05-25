import { Navigate } from "react-router-dom";
import { ReactNode } from "react";
import { useAppSelector } from "../store/hook";
import React from "react";
interface RequireRoleProps {
  roles: string[];
  children: ReactNode;
}

export default function RequireRole({ roles, children }: RequireRoleProps) {
  const user = useAppSelector(state => state.auth.user); // ดึง user จาก redux
  if (!user) {
    return <Navigate to="/" replace />;
  }

  if (!roles.includes(user.role)) {
    return <Navigate to="/unauthorized" replace />;
  }

  return <>{children}</>;
}
