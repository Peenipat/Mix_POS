// src/App.tsx
import { Routes, Route } from "react-router-dom";
import Login from "./page/Login";
import SuperAdminDashboard from "./page/admin/SuperAdminDashboard"
import BranchAdminDashboard from "./page/admin/BranchAdminDashboard"
import ManageUsers from "./page/admin/ManageUsers";
import StaffDashboard from "./page/staff/StaffDashboard";
import SuperAdminLayout from "./layouts/SuperAdminLayout";
// import Dashboard from "./pages/Dashboard";
// import AdminDashboard from "./pages/AdminDashboard";
// import BranchDashboard from "./pages/BranchDashboard";
// import Unauthorized from "./pages/Unauthorized";
import Home from "./page/Home";

// import ProtectedRoute from "./components/ProtectedRoute";
import RequireRole from "./components/RequireRole";

export default function App() {
  return (
    <Routes>
      {/* Public Route */}
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />
      {/* <Route path="/unauthorized" element={<Unauthorized />} /> */}

      {/*ต้อง login ก่อนเท่านั้น
      <Route
        path="/dashboard"
        element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        }
      /> */}


      <Route
        path="/admin"
        element={
          <RequireRole roles={["SUPER_ADMIN"]}>
            <SuperAdminLayout />
          </RequireRole>
        }
      >
        <Route path="dashboard" element={<SuperAdminDashboard />} />
        <Route path="users" element={<ManageUsers />} />
      </Route>

      {/* ต้องเป็น BRANCH_ADMIN */}
      <Route
        path="/branch/dashboard"
        element={
          <RequireRole roles={["BRANCH_ADMIN"]}>
            <BranchAdminDashboard />
          </RequireRole>
        }
      />

      <Route
        path="/staff/dashboard"
        element={
          <RequireRole roles={["STAFF"]}>
            <StaffDashboard />
          </RequireRole>
        }
      />
    </Routes>


  );
}
