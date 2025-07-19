// src/components/RequireRole.tsx
import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAppSelector } from '../store/hook';

interface RequireRoleProps {
  roles: string[];
  children: React.ReactNode;
}

export default function RequireRole({ roles, children }: RequireRoleProps) {
  const me       = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);

  // 1) หากกำลังรอ fetch /me (statusMe==='loading') ให้แสดง loading หรือคงไว้บนหน้าปัจจุบัน
  if (statusMe === 'loading') {
    return <div>Loading…</div>; // หรือ return null ก็ได้ (render ค้างบนหน้าเดิม)
  }

  // 2) ถ้า /me สำเร็จ (statusMe==='succeeded') แต่ me === null (เช่น token หมดอายุ) → redirect ไป login
  if (statusMe === 'succeeded' && !me) {
    return <Navigate to="/login" replace />;
  }

  // 3) ถ้า /me สำเร็จและ me.role ไม่ตรงกับ roles → redirect ไป unauthorized (หรือ login)
  if (statusMe === 'succeeded' && me && !roles.includes(me.role)) {
    return <Navigate to="/unauthorized" replace />;
  }

  // 4) ถ้า statusMe ยังไม่เคยเรียก (เช่นหน้าแรกก่อน load /me) ก็ถือว่า “ยังไม่ล็อกอิน” → redirect ไป login
  if (statusMe === 'idle' && !me) {
    return <Navigate to="/login" replace />;
  }

  // 5) ถ้า /me สำเร็จและ role ตรงเงื่อนไข → render children
  if (statusMe === 'succeeded' && me && roles.includes(me.role)) {
    return <>{children}</>;
  }

  // ส่วนอื่น ๆ (เช่น statusMe==='failed') → redirect login
  return <Navigate to="/login" replace />;
}
