// src/App.tsx
import { Routes, Route } from "react-router-dom";
import Home from "./page/Home";
import Login from "./page/Login";
import RequireRole from "./components/RequireRole";
import AdminLayout from "./layouts/AdminLayout";
import { RoleName } from "./types/role";
import AdminDashboard from "./page/admin/AdminDashboard";
import { ManageBarber } from "./page/admin/ManageBarber";
export default function App() {
  return (
    <Routes>
      {/* Public Route */}
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />
      {/* <Route path="/register" element={<Register />} />
      <Route path="/unauthorized" element={<Unauthorized />} />  */}

     
      {/* <Route
        path="/dashboard"
        element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        }
      /> */}


      <Route path="/admin" element={
        <RequireRole roles={[RoleName.BranchAdmin]}>
          <AdminLayout />
        </RequireRole>
      }>
        <Route index element={<AdminDashboard />} />
        <Route path="dashboard" element={<AdminDashboard />} />
        <Route path="barber" element={<ManageBarber />} />

      </Route>

      {/* ต้องเป็น BRANCH_ADMIN */}
      {/* <Route
        path="/branch/dashboard"
        element={
          <RequireRole roles={[RoleName.BranchAdmin]}>
            <BranchAdminDashboard />
          </RequireRole>
        }
      /> */}

      {/* <Route
        path="/staff/dashboard"
        element={
          <RequireRole roles={[RoleName.Staff]}>
            <StaffDashboard />
          </RequireRole>
        }
      />*/}
      
    </Routes> 


  );
}
