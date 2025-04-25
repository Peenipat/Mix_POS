import { Navigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";
import { ReactNode } from "react";

interface RequireRoleProps {
  roles: string[];
  children: ReactNode;
}

interface DecodedToken {
  user_id: number;
  role: string;
  exp: number;
}

export default function RequireRole({ roles, children }: RequireRoleProps) {
  const token = localStorage.getItem("token");

  if (!token) {
    return <Navigate to="/" replace />;
  }

  try {
    const decoded = jwtDecode<DecodedToken>(token);

    if (!roles.includes(decoded.role)) {
      return <Navigate to="/unauthorized" replace />;
    }

    return <>{children}</>;
  } catch (err) {
    console.error("Invalid token", err);
    return <Navigate to="/" replace />;
  }
}
