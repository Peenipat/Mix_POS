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
import { ManageBarber } from "./page/admin/barber/Barber_index";
import { ManageService } from "./page/admin/ManageService";
import { ManageCustomer } from "./page/admin/ManageCustomer";
import { ManageAppointments } from "./page/admin/ManageAppointments";
import ServicePage from "./page/ServicePage";
import MainLayout from "./layouts/MainLayout";
import BarberPage from "./page/BarberPage";
import ManageTime from "./page/admin/ManageTime";
import NotReady from "./page/NotReady";
import ContractDev from "./page/admin/ContractDev";
import HelpPage from "./page/admin/HelpPage";
import BarberDetail from "./page/admin/barber/Barber_id";
import BarberLayout from "./layouts/BarberLayout";
import BarberDashboard from "./page/barbers/BarberDashboard";
import BarberAppointment from "./page/barbers/AppointmentsHistory";
import AppointmentsHistory from "./page/barbers/AppointmentsHistory";
import BarberProfile from "./page/barbers/BarberProfile";
import CustomerReview from "./page/barbers/CustomerReview";
import BarberIncome from "./page/barbers/BarberIncome";
import HistoryPage from "./page/members/HistoryPage";
import AppointmentsPage from "./page/members/AppointmentsPage";

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

      <Route path="/login" element={<Login />} />

      <Route element={<MainLayout />}>
        <Route path="/" element={<Home />} />
        <Route path="service" element={<ServicePage />} />
        <Route path="barbers" element={<BarberPage />} />
        <Route path="booking" element={<AppointmentsPage />} />
        <Route path="history" element={<HistoryPage />} />
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
        <Route path="barber/:id" element={<BarberDetail />} />
        <Route path="service" element={<ManageService />} />
        <Route path="customer" element={<ManageCustomer />} />
        <Route path="appointments" element={<ManageAppointments />} />
        <Route path="working" element={<ManageTime />} />
        <Route path="help" element={<HelpPage />} />

        <Route path="inventory" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="tax" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="layout" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="calendar" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="feedback" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="report" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="branch" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="billing" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />

        <Route path="contact" element={<ContractDev />} />

        <Route path="*" element={<Navigate to="/admin/dashboard" replace />} />
      </Route>

      <Route
        path="/barber"
        element={
          <RequireRole roles={[RoleName.Staff]}>
            <BarberLayout />
          </RequireRole>
        }
      >
        <Route path="dashboard" element={<BarberDashboard />} />
        <Route path="appointments_history" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />

        <Route path="calendar" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="setting" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="tax" element={<NotReady message="ขออภัยในความไม่สะดวก" />} />
        <Route path="income" element={<BarberIncome />} />
        <Route path="feedback" element={<CustomerReview />} />

        <Route path="contact" element={<ContractDev />} />
        <Route path="profile" element={<BarberProfile />} />

      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
