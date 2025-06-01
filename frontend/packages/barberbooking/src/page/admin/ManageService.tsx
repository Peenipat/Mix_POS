// src/page/admin/ManageService.tsx
import React, { useEffect, useState, useRef } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";

import { useAppSelector } from "../../store/hook";
import axios from "../../lib/axios";

interface Service {
    id: number;
    name: string;
    price: number;
    duration: number;
  }

export function ManageService() {
  // 1) อ่านข้อมูลโปรไฟล์ผู้ใช้ (me) จาก Redux store
  const me = useAppSelector((state) => state.auth.me);

  // 2) ดึง tenantId (สมมติ service อยู่ในระดับ tenant)
  const tenantId = me?.tenant_ids[0];

  // 3) state สำหรับเก็บรายการ service
  const [services, setServices] = useState<Service[]>([]);
  const [loadingServices, setLoadingServices] = useState<boolean>(false);
  const [errorServices, setErrorServices] = useState<string | null>(null);

  // 4) ref เพื่อบล็อกไม่ให้ fetch ซ้ำ
  const didFetchServices = useRef(false);

  // === useEffect: โหลดรายการ Service เมื่อมี tenantId ===
  useEffect(() => {
    if (!tenantId || didFetchServices.current) return;
    didFetchServices.current = true;

    const loadServices = async () => {
      setLoadingServices(true);
      setErrorServices(null);
      try {
        // เรียก API สมมติให้เป็น /barberbooking/tenants/{tenantId}/services
        const res = await axios.get<{ status: string; data: Service[] }>(
          `/barberbooking/tenants/${tenantId}/services`
        );
        if (res.data.status !== "success") {
          throw new Error(res.data.status);
        }
        setServices(res.data.data);
      } catch (err: any) {
        setErrorServices(
          err.response?.data?.message || err.message || "Failed to load services"
        );
      } finally {
        setLoadingServices(false);
      }
    };

    loadServices();
  }, [tenantId]);

  // === แสดงสถานะ Loading / Error ของ Service ===
  if (!tenantId) {
    return <p className="text-red-500">Cannot determine tenant information.</p>;
  }
  if (loadingServices) {
    return <p>Loading services…</p>;
  }
  if (errorServices) {
    return <p className="text-red-500">Error loading services: {errorServices}</p>;
  }

  // === กำหนด columns และ actions สำหรับ DataTable ===
  const columns: Column<Service>[] = [
    {
      header: "#",
      accessor: (_row, rowIndex) => rowIndex + 1,
    },
    { header: "Name",           accessor: "name" },
    { header: "Price",          accessor: "price" },
    { header: "Duration (min)", accessor: "duration" },
  ];
  const editAction: Action<Service> = {
    label: "Edit",
    onClick: (row) => console.log("edit service", row),
    className: "text-blue-600",
  };
  const deleteAction: Action<Service> = {
    label: "Delete",
    onClick: (row) => console.log("delete service", row),
    className: "text-red-600",
  };

  // === ถ้าทุกอย่างพร้อมแล้ว ให้แสดง DataTable ===
  return (
    <div>
      <h2 className="text-xl mb-4">Services for Tenant {tenantId}</h2>
      <DataTable<Service>
        data={services}
        columns={columns}
        onRowClick={(r) => console.log("row clicked", r)}
        actions={[editAction, deleteAction]}
        showEdit={false}
        showDelete={false}
      />
    </div>
  );
}
