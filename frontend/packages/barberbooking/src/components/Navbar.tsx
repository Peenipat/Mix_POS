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
    { name: "บริการของทางร้าน", href: "/service" },
    { name: "ข้อมูลช่างตัดผม", href: "/barbers" },
    { name: "จองคิวออนไลน์", href: "/booking" },
    { name: "รีวิว & ให้คะแนน", href: "/reviews" },
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
      <div className="max-w-screen-xl flex flex-wrap items-center justify-between mx-auto p-4">
        {/* Logo */}
        <button
          onClick={() => navigate("/")}
          className="flex items-center space-x-3 text-gray-900"
        >
    
          <span className="self-center text-2xl font-semibold">Barber Shop</span>
        </button>

        {/* User + Mobile toggle */}
        <div className="flex items-center md:order-2">
          {/* Avatar */}
          <button
            onClick={() => setIsUserDropdownOpen((o) => !o)}
            className="flex text-sm bg-gray-700 rounded-full focus:ring-2 focus:ring-gray-600"
          >
            <span className="sr-only">Open user menu</span>
            <img
              className="w-8 h-8 rounded-full"
              src="/docs/images/people/profile-picture-3.jpg"
              alt="user photo"
            />
          </button>

          {/* User dropdown */}
          {isUserDropdownOpen && (
            <div
              ref={dropdownRef}
              className="z-50 absolute top-16 right-4 w-48 bg-gray-700 divide-y divide-gray-600 rounded-lg shadow-lg"
            >
              <div className="px-4 py-3">
                <span className="block text-sm text-gray-900">Bonnie Green</span>
                <span className="block text-xs text-gray-900 truncate">
                  name@flowbite.com
                </span>
              </div>
              <ul className="py-2">
                {userDropdownMenu.map((item) => (
                  <li key={item.name}>
                    <button
                      onClick={() => {
                        setIsUserDropdownOpen(false);
                        navigate(item.href);
                      }}
                      className="w-full text-left px-4 py-2 text-sm text-gray-900 hover:bg-gray-600"
                    >
                      {item.name}
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          )}

          {/* Mobile menu button */}
          <button
            onClick={() => setIsNavOpen((o) => !o)}
            className="inline-flex items-center p-2 ml-2 text-gray-900 rounded-lg md:hidden hover:bg-gray-700 focus:ring-2 focus:ring-gray-600"
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
          className={`w-full md:flex md:w-auto md:order-1 ${
            isNavOpen ? "" : "hidden"
          }`}
        >
          <ul className="flex flex-col mt-4 space-y-2 font-medium md:mt-0 md:space-y-0 md:flex-row md:space-x-6">
            {mainMenu.map((item) => (
              <li key={item.href}>
                <button
                  onClick={() => {
                    navigate(item.href);
                    setIsNavOpen(false);
                  }}
                  className="block py-2 px-3 rounded-sm text-gray-900 hover:text-gray-900 hover:bg-gray-200"
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
