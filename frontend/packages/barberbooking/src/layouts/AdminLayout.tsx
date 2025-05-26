
import { Link, Outlet, useLocation } from "react-router-dom";

const crumbsMap: Record<string, string> = {
  dashboard: "Dashboard",
  users: "Manage Users",
  tenant: "Manage Tenant",
  log: "System Log",
};

interface MenuItem {
  to: string;
  label: string;
}

const menuItems: MenuItem[] = [
  { to: "dashboard",  label: "Dashboard" },
  { to: "barber",     label: "Manage Barber" },
  { to: "service",    label: "Manage Services" },
  { to: "customer",   label: "Manage Customers" },
  { to: "report",     label: "Reports & Analytics"},
  { to: "billing",    label: "Billing & Expenses"},
  { to: "user",       label: "User Management"},
  { to: "help",       label: "Help & Support"}
];

export default function AdminLayout() {
  const { pathname } = useLocation();
  const raw = pathname.split("/").filter(Boolean);
  const segments = raw.filter(seg => seg !== "admin");
  return (
    <>
      {/* Top Navbar */}
      <nav className="fixed top-0 z-50 w-full bg-white border-b border-gray-200 dark:bg-gray-800 dark:border-gray-700">
        
        <div className="px-3 py-3 lg:px-5 lg:pl-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center justify-start">
              <button
                data-drawer-target="logo-sidebar"
                data-drawer-toggle="logo-sidebar"
                aria-controls="logo-sidebar"
                type="button"
                className="inline-flex items-center p-2 text-sm text-gray-500 rounded-lg sm:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
              >
                <span className="sr-only">Open sidebar</span>
                <svg
                  className="w-6 h-6"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    clipRule="evenodd"
                    fillRule="evenodd"
                    d="M2 4.75A.75.75 0 012.75 4h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 4.75zm0 10.5a.75.75 0 01.75-.75h7.5a.75.75 0 010 1.5h-7.5a.75.75 0 01-.75-.75zm0-5a.75.75 0 01.75-.75h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 10z"
                  />
                </svg>
              </button>
              <Link to="/" className="flex ml-2 md:ml-6">
                <span className="font-bai self-center text-xl font-semibold sm:text-2xl whitespace-nowrap dark:text-white">
                  ระบบหลังบ้าน
                </span>
              </Link>
            </div>
            {/* <div className="flex items-center">
              <LogoutButton />
            </div> */}
          </div>
        </div>
      </nav>

      {/* Sidebar */}
      <aside
      id="logo-sidebar"
      className="fixed top-0 left-0 z-40 w-64 h-screen pt-16 transition-transform -translate-x-full bg-[#1f2937] border-r border-gray-200 sm:translate-x-0 dark:bg-gray-800 dark:border-gray-700"
      aria-label="Sidebar"
    >
      <div className="h-full px-3 pb-4 overflow-y-auto">
        <ul className="space-y-2 font-medium">
          {menuItems.map((item) => (
            <li key={item.to}>
              <Link
                to={item.to}
                className="flex items-center p-2 text-white rounded-lg hover:bg-gray-700 dark:hover:bg-gray-700"
              >
                <span className="flex-1 ml-3 whitespace-nowrap">
                  {item.label}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      </div>
    </aside>

      {/* Main Content */}
      <div className="p-4 sm:ml-64 pt-20 bg-gray-100 dark:bg-gray-900 min-h-screen">
        <div className="p-6 bg-white dark:bg-gray-800 rounded-lg shadow-sm">

          {/* Breadcrumb */}
          <nav className="flex mb-4" aria-label="Breadcrumb">
            <ol className="inline-flex items-center space-x-1 md:space-x-2">
              {segments.length === 0 || segments[0] !== "dashboard" ? (
                <li className="inline-flex items-center">
                  <Link
                    to="/admin/dashboard"
                    className="inline-flex items-center text-sm font-medium text-gray-700 hover:text-blue-600"
                  >
                    <svg
                      className="w-3 h-3 me-2.5"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path d="m19.707 9.293-2-2-7-7a1…Z" />
                    </svg>
                    Dashboard
                  </Link>
                </li>
              ) : (
                // กรณี path เริ่มต้นเลยเป็น /dashboard
                <li>
                  <span className="text-sm font-medium text-gray-500">
                    Dashboard
                  </span>
                </li>
              )}

              {segments.map((seg, idx) => {
                // ข้าม segment "dashboard" ไปแล้วเพราะเราใช้มันเป็น home
                if (idx === 0 && seg === "dashboard") return null;
                const to = "/" + segments.slice(0, idx + 1).join("/");
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
                        <span className="ms-1 text-sm font-medium text-gray-500">
                          {label}
                        </span>
                      ) : (
                        <Link
                          to={to}
                          className="ms-1 text-sm font-medium text-gray-700 hover:text-blue-600 md:ms-2"
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
