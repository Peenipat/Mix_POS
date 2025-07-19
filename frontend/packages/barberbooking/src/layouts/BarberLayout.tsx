import React, { useState } from "react";
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useAppDispatch } from "../store/hook";
import { logout } from "../store/authSlice";

interface MenuItem {
  to: string;
  label: string;
  status?: "new" | "maintenance" | "comingsoon";
}

interface MenuGroup {
  label: string;
  key: string;
  items: MenuItem[];
}

const crumbsMap: Record<string, string> = {
  dashboard: "หน้าหลัก",
  appointments_history: "ประวัติการนัดหมาย",
  income: "รายได้",
  calendar: "ปฏิทินการนัดหมาย",
  tax: "คำนวณภาษี",
  feedback: "รีวิวจากลูกค้า",
  help: "ความช่วยเหลือ & ติดต่อผู้พัฒนา",
  contact: "ติดต่อผู้พัฒนา",
  profile: "ข้อมูลส่วนตัว",
  setting:"ตั้งค่า"
};

const groupedMenu: MenuGroup[] = [
  {
    key: "customer",
    label: "ข้อมูลลูกค้า",
    items: [
      { to: "feedback", label: "รีวิวจากลูกค้า",},
      { to: "appointments_history", label: "ประวัตินัดหมาย" ,status: "comingsoon" },
    ],
  },
  {
    key: "finance",
    label: "การเงิน / บัญชี",
    items: [
      { to: "income", label: "รายได้", },
      { to: "tax", label: "คำนวณภาษี", status: "comingsoon" },
    ],
  },
  {
    key: "etc",
    label: "อื่น ๆ",
    items: [
      { to: "setting", label: "ตั้งค่า" ,status: "comingsoon"},
      { to: "contact", label: "ติดต่อผู้พัฒนา" },
    ],
  },
];

const renderStatusBadge = (status?: string) => {
  if (!status) return null;

  const badgeClass = {
    new: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
    maintenance: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
    comingsoon: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
  }[status];

  const label = {
    new: "ใหม่",
    maintenance: "ปิดปรับปรุง",
    comingsoon: "ยังไม่พร้อมใช้งาน",
  }[status];

  return (
    <span className={`text-[11px] font-medium px-1 py-0.5 rounded ${badgeClass}`}>
      {label}
    </span>
  );
};



export default function BarberLayout() {
  const { pathname } = useLocation();
  const raw = pathname.split("/").filter(Boolean);
  const segments = raw.filter((seg) => seg !== "barber");

  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const [openGroups, setOpenGroups] = useState<Record<string, boolean>>({});

  const toggleGroup = (key: string) => {
    setOpenGroups((prev) => ({
      ...prev,
      [key]: !prev[key],
    }));
  };

  const handleLogout = () => {
    dispatch(logout());
    navigate("/login");
  };

  return (
    <>
      {/* Top Navbar */}
      <nav className="fixed top-0 z-50 w-full bg-white border-b border-gray-200 dark:bg-gray-800 dark:border-gray-700">
        <div className="px-3 py-3 lg:px-5 lg:pl-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center justify-start">
              <Link to="/" className="flex ml-2 md:ml-6">
                <span className="font-bai self-center text-xl font-semibold sm:text-2xl whitespace-nowrap dark:text-white">
                  ระบบหลังบ้าน
                </span>
              </Link>
            </div>
            <div className="flex items-center">
              <button
                onClick={handleLogout}
                className="text-gray-700 hover:text-red-600 dark:text-gray-300 dark:hover:text-red-400 text-sm font-medium"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* Sidebar */}
      <aside
        id="logo-sidebar"
        className="fixed top-0 left-0 z-40 w-64 h-screen pt-16 transition-transform -translate-x-full bg-[#1f2937] border-r border-gray-200 sm:translate-x-0 dark:bg-gray-800 dark:border-gray-700"
        aria-label="Sidebar"
      >
        <div className="h-full px-1 pb-4 overflow-y-auto">
          <ul className="space-y-4 font-medium text-white">
            <li>
              <Link to="dashboard" className="block p-2 rounded-lg hover:bg-gray-700 font-bold">
                หน้าหลัก
              </Link>
            </li>
            {groupedMenu.map((group) => (
              <li key={group.key}>
                <button
                  onClick={() => toggleGroup(group.key)}
                  className="flex justify-between items-center w-full px-3 py-2 text-sm font-semibold text-gray-300 uppercase hover:text-white"
                >
                  <span>{group.label}</span>
                  <svg
                    className={`w-4 h-4 transition-transform ${openGroups[group.key] ? "rotate-90" : "rotate-0"}`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
                {openGroups[group.key] && (
                  <ul className="pl-5 mt-1 space-y-1">
                    {group.items.map((item) => (
                      <li key={item.to}>
                        <Link to={item.to} className="block p-2 text-sm rounded-lg hover:bg-gray-700">
                          <span className="mr-2.5">{item.label}</span>
                          {renderStatusBadge(item.status)}
                        </Link>
                      </li>
                    ))}
                  </ul>
                )}
              </li>
            ))}
            <li>
              <Link
                to="/barber/profile"
                className="block p-2 text-sm rounded-lg hover:bg-gray-700"
              >
              ข้อมูลส่วนตัว
              </Link>
            </li>
          </ul>
        </div>
      </aside>

      {/* Main Content */}
      <div className="p-4 sm:ml-64 pt-20 bg-gray-100 dark:bg-gray-900 min-h-screen">
        <div className="p-6 bg-white dark:bg-gray-800 rounded-lg shadow-sm">
          <nav className="flex mb-4" aria-label="Breadcrumb">
            <ol className="inline-flex items-center space-x-1 md:space-x-2">
              {segments.length === 0 || segments[0] !== "dashboard" ? (
                <li className="inline-flex items-center">
                  <Link
                    to="/barber/dashboard"
                    className="inline-flex items-center text-sm font-medium text-gray-700 hover:text-blue-600 dark:text-gray-300 dark:hover:text-blue-500"
                  >
                    หน้าหลัก
                  </Link>
                </li>
              ) : (
                <li>
                  <span className="text-sm font-medium text-gray-500 dark:text-gray-400">หน้าหลัก</span>
                </li>
              )}
              {segments.map((seg, idx) => {
                if (idx === 0 && seg === "dashboard") return null;
                const basePath = "/barber";
                const to = `${basePath}/${segments.slice(0, idx + 1).join("/")}`;
                const isLast = idx === segments.length - 1;
                const label = crumbsMap[seg] || seg;
                return (
                  <li key={to} aria-current={isLast ? "page" : undefined}>
                    <div className="flex items-center">
                      <svg
                        className="w-3 h-3 text-gray-400 mx-1"
                        fill="none"
                        viewBox="0 0 6 10"
                      >
                        <path
                          stroke="currentColor"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="m1 9 4-4-4-4"
                        />
                      </svg>
                      {isLast ? (
                        <span className="ms-1 text-sm font-medium text-gray-500 dark:text-gray-400">{label}</span>
                      ) : (
                        <Link
                          to={to}
                          className="ms-1 text-sm font-medium text-gray-700 hover:text-blue-600 md:ms-2 dark:text-gray-300 dark:hover:text-blue-500"
                        >
                          {label}
                        </Link>
                      )}
                    </div>
                  </li>
                );
              })}
            </ol>
          </nav>

          <Outlet />
        </div>
      </div>
    </>
  );
}