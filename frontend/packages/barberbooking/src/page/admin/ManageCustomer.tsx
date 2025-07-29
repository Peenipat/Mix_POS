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
  const me = useAppSelector((state) => state.auth.me);
  const tenantId = me?.tenant_ids[0];
  const branchId = me?.branch_id;

  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loadingCustomers, setLoadingCustomers] = useState<boolean>(false);
  const [errorCustomers, setErrorCustomers] = useState<string | null>(null);
  const didFetchCustomers = useRef(false);

  useEffect(() => {
    if (!tenantId || didFetchCustomers.current) return;
    didFetchCustomers.current = true;

    const loadCustomers = async () => {
      setLoadingCustomers(true);
      setErrorCustomers(null);
      try {
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

  if (!tenantId) {
    return <p className="text-red-500">Cannot determine tenant information.</p>;
  }
  if (loadingCustomers) {
    return <p>Loading customers…</p>;
  }
  if (errorCustomers) {
    return <p className="text-red-500">Error loading customers: {errorCustomers}</p>;
  }

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
