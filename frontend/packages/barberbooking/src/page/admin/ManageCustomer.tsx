  // src/page/admin/ManageCustomer.tsx
import React, { useEffect, useState, useRef } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";
import { useAppSelector } from "../../store/hook";
import axios from "../../lib/axios";



interface Customer {
    id: number;
    Name: string;
    email: string;
    Phone: string;
  }

export function ManageCustomer() {
  // 1) อ่านข้อมูลโปรไฟล์ผู้ใช้ (me) จาก Redux store
  const me = useAppSelector((state) => state.auth.me);

  // 2) ดึง tenantId (สมมติว่า customer ถูกผูกกับ tenant เท่านั้น)
  const tenantId = me?.tenant_ids[0];
  const branchId = me?.branch_id;

  // 3) state สำหรับเก็บรายการ customer
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loadingCustomers, setLoadingCustomers] = useState<boolean>(false);
  const [errorCustomers, setErrorCustomers] = useState<string | null>(null);

  // 4) ref เพื่อบล็อกไม่ให้ fetch ซ้ำ
  const didFetchCustomers = useRef(false);

  // === useEffect: โหลดรายการ Customer เมื่อมี tenantId ===
  useEffect(() => {
    if (!tenantId || didFetchCustomers.current) return;
    didFetchCustomers.current = true;

    const loadCustomers = async () => {
      setLoadingCustomers(true);
      setErrorCustomers(null);
      try {
        // เรียก API สมมติเป็น /barberbooking/tenants/{tenantId}/customers
        const res = await axios.get<{ status: string; data: Customer[] }>(
          `/barberbooking/tenants/${tenantId}/branch/${branchId}/customers`
        );
        if (res.data.status !== "success") {
          throw new Error(res.data.status);
        }
        setCustomers(res.data.data);
      } catch (err: any) {
        setErrorCustomers(
          err.response?.data?.message || err.message || "Failed to load customers"
        );
      } finally {
        setLoadingCustomers(false);
      }
    };

    loadCustomers();
  }, [tenantId]);

  // === แสดงสถานะ Loading / Error ของ Customer ===
  if (!tenantId) {
    return <p className="text-red-500">Cannot determine tenant information.</p>;
  }
  if (loadingCustomers) {
    return <p>Loading customers…</p>;
  }
  if (errorCustomers) {
    return <p className="text-red-500">Error loading customers: {errorCustomers}</p>;
  }

  // === กำหนดคอลัมน์และ actions สำหรับ DataTable ===
  const columns: Column<Customer>[] = [
    {
      header: "#",
      accessor: (_row, rowIndex) => rowIndex + 1, 
    },
    { header: "Name",        accessor: "Name" },
    { header: "Email",       accessor: "email" },
    { header: "Phone",       accessor: "Phone" },
  ];

  const viewAction: Action<Customer> = {
    label: "ดูประวัติการจอง",
    onClick: (row) => console.log("edit customer", row),
    className: "text-blue-600",
  };

  const editAction: Action<Customer> = {
    label: "Edit",
    onClick: (row) => console.log("edit customer", row),
    className: "text-blue-600",
  };
  const deleteAction: Action<Customer> = {
    label: "Delete",
    onClick: (row) => console.log("delete customer", row),
    className: "text-red-600",
  };

  // === ถ้าทุกอย่างพร้อมแล้ว ให้แสดง DataTable ===
  return (
    <div>
      <h2 className="text-xl mb-4">Customers for Tenant {tenantId}</h2>
      <DataTable<Customer>
        data={customers}
        columns={columns}
        onRowClick={(r) => console.log("row clicked", r)}
        actions={[viewAction,editAction, deleteAction]}
        showEdit={false}
        showDelete={false}
      />
    </div>
  );
}
