// src/layouts/MainLayout.tsx
import React from "react";
import { Outlet } from "react-router-dom";
import Navbar from "../components/Navbar";

export default function MainLayout() {
  return (
    <div className="flex flex-col h-screen"> 
      <Navbar />

      <main className="flex-grow">
        <Outlet />
      </main>
    </div>
  );
}

