// src/App.tsx
import React, { useEffect, useState, useRef } from "react";
import { Routes, Route, Navigate, useLocation } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "./store/hook";
import { loadCurrentUser, logout } from "./store/authSlice";

import Home from "./page/Home";
import Login from "./page/Login";
import RequireRole from "./components/RequireRole";
import AdminLayout from "./layouts/AdminLayout";
import { RoleName } from "./types/role";
import AdminDashboard from "./page/admin/AdminDashboard";
import { ManageBarber } from "./page/admin/ManageBarber";
import { ManageService } from "./page/admin/ManageService";
import { ManageCustomer } from "./page/admin/ManageCustomer";
import { ManageAppointments } from "./page/admin/ManageAppointments";

export default function App() {
  const dispatch = useAppDispatch();
  const location = useLocation();
  // initialized เพื่อควบคุม Loading screen
  const [initialized, setInitialized] = useState(false);
  // ref เพื่อบอกว่าโหลด /me ไปแล้วหรือยัง
  const didFetchMe = useRef(false);

  useEffect(() => {
    // ถ้ายังไม่เคย fetch /me และ path เริ่มต้นด้วย "/admin" → ให้ fetch
    if (!didFetchMe.current && location.pathname.startsWith("/admin")) {
      didFetchMe.current = true; // ตั้ง flag ว่าเรียก /me ไปแล้ว
      dispatch(loadCurrentUser())
        .catch(() => {
          // ถ้า fetch /me พลาด (เช่น token ไม่ valid) ให้ logout
          dispatch(logout());
        })
        .finally(() => {
          // ไม่ว่า /me จะสำเร็จหรือ fail ก็ถือว่า initialization เสร็จ
          setInitialized(true);
        });
    }
    // แต่ถ้า path ไม่ใช่ "/admin" เลย (เช่น /login, /) → ข้ามการ fetch /me ไป
    
    setInitialized(true);
    // ✅ สังเกตว่า dependency array ไม่ได้ใส่ location.pathname หรือ dispatch ซ้ำ
    // ให้ run โค้ดนี้แค่ครั้งแรกที่ App mount เท่านั้น
    console.log("call")
  }, [dispatch,]);


  return (
    <Routes>
      {/* Public */}
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />

      {/* Protected /admin/* */}
      <Route
        path="/admin"
        element={
          <RequireRole roles={[RoleName.BranchAdmin]}>
            <AdminLayout />
          </RequireRole>
        }
      >
        <Route index element={<AdminDashboard />} />
        <Route path="dashboard" element={<AdminDashboard />} />
        <Route path="barber" element={<ManageBarber />} />
        <Route path="service" element={<ManageService/>}/>
        <Route path="customer" element={<ManageCustomer />} />
        <Route path="appointments" element={<ManageAppointments />} />
      </Route>

      {/* ไม่ match route ใด ๆ → redirect ไป Home */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
