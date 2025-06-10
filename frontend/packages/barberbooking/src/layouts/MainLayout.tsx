// src/layouts/MainLayout.tsx
import React from "react";
import { Outlet } from "react-router-dom";
import Navbar from "../components/Navbar";

export default function MainLayout() {
  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />

      <main className="">
        <Outlet />
      </main>

    </div>
  );
}
