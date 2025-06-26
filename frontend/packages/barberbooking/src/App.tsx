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
import ServicePage from "./page/ServicePage";
import MainLayout from "./layouts/MainLayout";
import BarberPage from "./page/BarberPage";
import ManageTime from "./page/admin/ManageTime";
import NotReady from "./page/NotReady";
import ContractDev from "./page/admin/ContractDev";

export default function App() {
  const dispatch = useAppDispatch();
  const location = useLocation();
  // initialized เพื่อควบคุม Loading screen
  const [initialized, setInitialized] = useState(false);
  // ref เพื่อบอกว่าโหลด /me ไปแล้วหรือยัง
  const didFetchMe = useRef(false);

  useEffect(() => {
    if (!didFetchMe.current && location.pathname.startsWith("/admin")) {
      didFetchMe.current = true;
      dispatch(loadCurrentUser())
        .catch(() => {
          dispatch(logout());
        })
        .finally(() => {
          setInitialized(true);
        });
    }

    setInitialized(true);
  }, [dispatch,]);


  return (
    <Routes>
      {/* Public */}
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />

      <Route element={<MainLayout />}>
      <Route path="service" element={<ServicePage />} />
      <Route path="barbers" element={<BarberPage />} />
      <Route path="booking" element={<NotReady  message="ขออภัยในความไม่สะดวก"/>} />
      <Route path="reviews" element={<NotReady  message="ขออภัยในความไม่สะดวก"/>} />
      <Route path="history" element={<NotReady  message="ขออภัยในความไม่สะดวก"/>} />
      </Route>

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
        <Route path="service" element={<ManageService />} />
        <Route path="customer" element={<ManageCustomer />} />
        <Route path="appointments" element={<ManageAppointments />} />
        <Route path="working" element={<ManageTime />} />
        
        <Route path="report" element={<NotReady message="ขออภัยในความไม่สะดวก"/>} />
        <Route path="branch" element={<NotReady message="ขออภัยในความไม่สะดวก"/>} />
        <Route path="billing" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="help" element={<ContractDev/>} />

        <Route path="*" element={<Navigate to="/admin/dashboard" replace />} />
      </Route>

      {/* ไม่ match route ใด ๆ → redirect ไป Home */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
