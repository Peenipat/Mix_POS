// src/components/Navbar.tsx
import React, { useState, useRef, useEffect, RefObject } from "react";
import { useNavigate } from "react-router-dom";

interface MenuItem {
  name: string;
  href: string;
}

interface UserMenuItem {
  name: string;
  href: string;
}

export default function Navbar() {
  const navigate = useNavigate();
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [isUserDropdownOpen, setIsUserDropdownOpen] = useState(false);
  const dropdownRef: RefObject<HTMLDivElement> = useRef(null);

  const mainMenu: MenuItem[] = [
    { name: "หน้าหลัก", href: "/" },
    { name: "บริการ", href: "/service" },
    { name: "ข้อมูลช่างตัดผม", href: "/barbers" },
    { name: "จองคิว", href: "/booking" },
    { name: "ประวัติการจอง", href: "/history" },
    { name: "เข้าสู่ระบบ", href: "/login" },
  ];

  const userDropdownMenu: UserMenuItem[] = [
    { name: "Dashboard", href: "/dashboard" },
    { name: "Settings", href: "/settings" },
    { name: "Earnings", href: "/earnings" },
    { name: "Sign out", href: "/login" },
  ];

  // ปิด dropdown เมื่อคลิกนอก
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsUserDropdownOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <nav className=" bg-slate-50 shadow-md border-gray-700">
      <div className="max-w-screen-xl flex flex-wrap items-center justify-between mx-auto p-1">
        {/* Logo */}
        <button
          onClick={() => navigate("/")}
          className="flex items-center space-x-3 text-gray-900"
        >

          <span className="self-center text-xl font-semibold">Barber Shop</span>
        </button>

        {/* User + Mobile toggle */}
        <div className="flex items-center md:order-2">
          {/* Avatar */}

          {/* User dropdown */}
     

          {/* Mobile menu button */}
          <button
            onClick={() => setIsNavOpen((o) => !o)}
            className="inline-flex items-center p-1 ml-1 text-gray-900 rounded-lg md:hidden hover:bg-gray-700 focus:ring-2 focus:ring-gray-600"
          >
            <span className="sr-only">Open main menu</span>
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>

        {/* Main menu */}
        <div
          className={`w-full md:flex md:w-auto md:order-1 ${isNavOpen ? "" : "hidden"
            }`}
        >
          <ul className="flex flex-col mt-4 space-y-2 font-medium md:mt-0 md:space-y-0 md:flex-row md:space-x-6">
            {mainMenu.map((item) => (
              <li key={item.href}>
                <button
                  id={item.href === "/history" ? "nav-history" : undefined}
                  onClick={() => {
                    navigate(item.href);
                    setIsNavOpen(false);
                  }}
                  className="block py-1 px-2 rounded-sm text-gray-900 hover:text-gray-900 hover:bg-gray-200"
                >
                  {item.name}
                </button>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </nav>
  );
}
