import { Outlet } from "react-router-dom";

export default function SuperAdminLayout() {
  return (
    <div className="flex min-h-screen">
      {/* Sidebar */}
      <aside className="w-64 bg-gray-800 text-white p-4">
        <h2 className="text-2xl font-bold mb-8">Admin Panel</h2>
        <ul className="space-y-4">
          <li><a href="/admin/dashboard" className="hover:underline">Dashboard</a></li>
          <li><a href="/admin/users" className="hover:underline">Users</a></li>
        </ul>
      </aside>

      {/* Main Content */}
      <div className="flex-1 p-6 bg-gray-100">
        <Outlet /> 
        {/* 🔥 ตรงนี้จะเปลี่ยนไปตาม Route ย่อย เช่น Dashboard หรือ ShowUsers */}
      </div>
    </div>
  );
}
