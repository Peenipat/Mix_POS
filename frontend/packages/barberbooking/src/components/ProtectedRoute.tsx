// src/components/ProtectedRoute.tsx
import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAppSelector } from '../store/hook';
import type { ReactNode } from 'react';

export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const me = useAppSelector((state) => state.auth.me);

  if (!me) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}
