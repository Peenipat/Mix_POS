// src/App.tsx
import { Routes, Route } from "react-router-dom";
import Login from "./page/Login"
import SuperAdminDashboard from "./page/admin/SuperAdminDashboard"
import BranchAdminDashboard from "./page/admin/BranchAdminDashboard"
import ManageUsers from "./page/admin/super_admin/ManageUsers";
import StaffDashboard from "./page/staff/StaffDashboard";
import SuperAdminLayout from "./layouts/SuperAdminLayout";
import Register from "./page/Register";
import Dashboard from "./page/user/Dashboard";
import LogTablePage from "./page/admin/super_admin/Log";
import { RoleName } from "./types/role";
import Home from "./page/Home";

import ProtectedRoute from "./components/ProtectedRoute";
import RequireRole from "./components/RequireRole";
import Unauthorized from "./page/unauthorized";
import React from "react";
import ManageTenant from "./page/admin/super_admin/ManageTenant";
import TenantDetail from "./page/admin/super_admin/Tenant$id";

export default function App() {
  return (
    <Routes>
      {/* Public Route */}
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/unauthorized" element={<Unauthorized />} />

      ต้อง login ก่อนเท่านั้น
      <Route
        path="/dashboard"
        element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        }
      />


      <Route path="/admin" element={
        <RequireRole roles={[RoleName.SaaSSuperAdmin]}>
          <SuperAdminLayout />
        </RequireRole>
      }>
        <Route index element={<SuperAdminDashboard />} />
        <Route path="dashboard" element={<SuperAdminDashboard />} />
        <Route path="users" element={<ManageUsers />} />
        <Route path="log" element={<LogTablePage />} />
        <Route path="tenant" element={<ManageTenant/>}/>
        <Route path="tenant/:id" element={<TenantDetail/>}/>
      </Route>

      {/* ต้องเป็น BRANCH_ADMIN */}
      <Route
        path="/branch/dashboard"
        element={
          <RequireRole roles={[RoleName.BranchAdmin]}>
            <BranchAdminDashboard />
          </RequireRole>
        }
      />

      <Route
        path="/staff/dashboard"
        element={
          <RequireRole roles={[RoleName.Staff]}>
            <StaffDashboard />
          </RequireRole>
        }
      />
    </Routes>


  );
}
