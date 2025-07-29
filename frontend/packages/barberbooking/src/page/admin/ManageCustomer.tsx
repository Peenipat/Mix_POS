  // src/page/admin/ManageCustomer.tsx
import React, { useEffect, useState, useRef } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";
import { useAppSelector } from "../../store/hook";
import axios from "../../lib/axios";
import { useNavigate } from "react-router-dom";

interface Customer {
    id: number;
    name: string;
    email: string;
    phone: string;
  }

export function ManageCustomer() {
  const navigate = useNavigate();
  const me = useAppSelector((state) => state.auth.me);
  const tenantId = me?.tenant_ids[0];
  const branchId = me?.branch_id;

  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loadingCustomers, setLoadingCustomers] = useState<boolean>(false);
  const [errorCustomers, setErrorCustomers] = useState<string | null>(null);
  const didFetchCustomers = useRef(false);

  const handleCustomerDetail = (customerId: number) => {
    navigate(`/admin/customer/${customerId}`)
  }

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
    { header: "Name",        accessor: "name" },
    { header: "Email",       accessor: "email" },
    { header: "Phone",       accessor: "phone" },
  ];

  const viewAction: Action<Customer> = {
    label: "ดูประวัติการจอง",
    onClick: (row) => handleCustomerDetail(row.id),
    className: "text-blue-600",
  };


  return (
    <div>
      <h2 className="text-xl mb-4">Customers for Tenant {tenantId}</h2>
      <DataTable<Customer>
        data={customers}
        columns={columns}
        onRowClick={(r) => console.log("row clicked", r)}
        actions={[viewAction]}
        showEdit={false}
        showDelete={false}
      />
    </div>
  );
}
