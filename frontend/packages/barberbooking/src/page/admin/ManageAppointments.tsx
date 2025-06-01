// src/pages/admin/ManageAppointments.tsx

import React, { useEffect, useState, useRef } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";
import { useAppSelector } from "../../store/hook";
import axios from "../../lib/axios";

// src/types/appointment.ts

export interface Role {
  id: number;
  name: string;
  description: string;
  created_at: string;  // ISO
  updated_at: string;  // ISO
  deleted_at: string | null;
}

export interface Barber {
  id: number;
  username: string;
  email: string;
  role_id: number;
  role: Role;
  created_at: string;   // ISO
  updated_at: string;   // ISO
  deleted_at: string | null;
}

export interface Service {
  ID: number;
  name: string;
  tenant_id: number;
  duration: number;
  price: number;
  created_at: string;  // ISO
  updated_at: string;  // ISO
  deleted_at: string | null;
}

export interface Customer {
  ID: number;
  TenantID: number;
  Name: string;
  Phone: string;
  email: string;
  CreatedAt: string;   // ISO
  updated_at: string;  // ISO
  deleted_at: string | null;
}

export interface Appointment {
  id: number;
  branch_id: number;
  service_id: number;
  /**
   * อ็อบเจกต์ service ที่ฝังมาใน response
   * จะมี field ชื่อ name, price, duration เป็นต้น
   */
  service: Service;
  barber_id: number;
  /**
   * อ็อบเจกต์ barber ที่ฝังมาใน response
   * จะมี field ชื่อ username, email, role_id, role (object)
   */
  barber: Barber;
  customer_id: number;
  /**
   * อ็อบเจกต์ customer ที่ฝังมาใน response
   * จะมี field ชื่อ Name, Phone, email เป็นต้น
   */
  customer: Customer;
  tenant_id: number;
  start_time: string; // ISO
  end_time: string;   // ISO
  status: string;
  notes: string;
  created_at: string; // ISO
  updated_at: string; // ISO
  deleted_at: string | null;
}


export function ManageAppointments() {
  // อ่าน me จาก Redux store
  const me = useAppSelector((state) => state.auth.me);

  // ดึง tenantId และ branchId
  const tenantId = me?.tenant_ids[0];
  const branchId = me?.branch_id;

  // state สำหรับเก็บรายการ Appointment
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [loadingAppts, setLoadingAppts] = useState<boolean>(false);
  const [errorAppts, setErrorAppts] = useState<string | null>(null);

  // ref เพื่อบล็อกไม่ให้ fetch ซ้ำ
  const didFetchAppts = useRef(false);

  useEffect(() => {
    if (!tenantId || !branchId || didFetchAppts.current) return;
    didFetchAppts.current = true;

    const loadAppointments = async () => {
      setLoadingAppts(true);
      setErrorAppts(null);
      try {
        const res = await axios.get<{ status: string; data: Appointment[] }>(
          `/barberbooking/tenants/${tenantId}/appointments`,
          {
            params: { branch_id: branchId },
          }
        );
        if (res.data.status !== "success") {
          throw new Error(res.data.status);
        }
        setAppointments(res.data.data);
      } catch (err: any) {
        setErrorAppts(
          err.response?.data?.message || err.message || "Failed to load appointments"
        );
      } finally {
        setLoadingAppts(false);
      }
    };

    loadAppointments();
  }, [tenantId, branchId]);

  // แสดงสถานะก่อนเข้า DataTable
  if (!tenantId || !branchId) {
    return <p className="text-red-500">Cannot determine tenant or branch information.</p>;
  }
  if (loadingAppts) {
    return <p>Loading appointments…</p>;
  }
  if (errorAppts) {
    return <p className="text-red-500">Error loading appointments: {errorAppts}</p>;
  }

  // กำหนด columns ตามโครงสร้างใหม่
  const columns: Column<Appointment>[] = [
    {
      header: "#",
      accessor: (_row, rowIndex) => rowIndex + 1,
    },
    {
      header: "Customer",
      // ดึงชื่อจากลูกค้า: customer.Name
      accessor: (row: Appointment) => row.customer.Name,
    },
    {
      header: "Barber",
      // ดึงชื่อช่างจาก obj barber.username
      accessor: (row: Appointment) => row.barber.username,
    },
    {
      header: "Service",
      // ดึงชื่อบริการจาก service.name
      accessor: (row: Appointment) => row.service.name,
    },
    {
      header: "Time",
      // แปลง ISO เป็น locale string
      accessor: (row: Appointment) => {
        const dt = new Date(row.start_time);
        return dt.toLocaleString();
      },
    },
    {
      header: "Status",
      accessor: (row: Appointment) => row.status,
    },
  ];

  const viewAction: Action<Appointment> = {
    label: "View",
    onClick: (row) => console.log("view appointment", row),
    className: "text-green-600",
  };
  const cancelAction: Action<Appointment> = {
    label: "Cancel",
    onClick: (row) => console.log("cancel appointment", row),
    className: "text-red-600",
  };

  return (
    <div>
      <h2 className="text-xl mb-4">
        Appointments for Branch {branchId} (Tenant {tenantId})
      </h2>
      <DataTable<Appointment>
        data={appointments}
        columns={columns}
        onRowClick={(r) => console.log("row clicked", r)}
        actions={[viewAction, cancelAction]}
        showEdit={false}
        showDelete={false}
      />
    </div>
  );
}
