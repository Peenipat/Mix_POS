// src/page/admin/ManageCustomer.tsx
import React, { useEffect, useState, useRef } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";
import { useAppSelector } from "../../store/hook";
import axios from "../../lib/axios";
import { useNavigate } from "react-router-dom";
import { loadCustomers } from "../../api/customer";

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

  const handleCustomerDetail = (customerId: number) => {
    navigate(`/admin/customer/${customerId}`)
  }


  const [customers, setCustomers] = useState<Customer[]>([]);
  const [pagination, setPagination] = useState({ total: 0, page: 1, limit: 10 });
  const [loadingCustomers, setLoadingCustomers] = useState(false);
  const [errorCustomers, setErrorCustomers] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [query, setQuery] = useState("");


  useEffect(() => {
    if (tenantId !== undefined && branchId !== undefined) {
      loadCustomers(
        tenantId,
        branchId,
        setCustomers,
        setPagination,
        setLoadingCustomers,
        setErrorCustomers,
        {
          page: pagination.page,
          limit: pagination.limit,
          name: query,
          phone: query,
          sortBy: "updated_at",
          sortOrder: "desc",
        }
      );
    }
  }, [pagination.page, query]);



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
      accessor: (_row, rowIndex) =>
        (pagination.page - 1) * pagination.limit + rowIndex + 1,
    },
    { header: "Name", accessor: "name" },
    { header: "Email", accessor: "email" },
    { header: "Phone", accessor: "phone" },
  ];

  const viewAction: Action<Customer> = {
    label: "ดูประวัติการจอง",
    onClick: (row) => handleCustomerDetail(row.id),
    className: "text-blue-600",
  };





  return (
    <div>
      <h2 className="text-xl mb-4">Customers for Tenant {tenantId}</h2>
      <div className="flex flex-col sm:flex-row gap-4 mb-4 items-center">
        <input
          type="text"
          className="border border-gray-300 rounded px-3 py-2 w-full sm:w-64"
          placeholder="ค้นหาชื่อ / เบอร์โทร"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <button
          onClick={() => {
            setPagination((prev) => ({ ...prev, page: 1 })); 
            setQuery(searchTerm); 
          }}
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          ค้นหา
        </button>
      </div>

      <DataTable<Customer>
        data={customers}
        columns={columns}
        onRowClick={(r) => console.log("row clicked", r)}
        actions={[viewAction]}
        showEdit={false}
        showDelete={false}
        page={pagination.page}
        perPage={pagination.limit}
        total={pagination.total}
        onPageChange={(p) => {
          setPagination((prev) => ({ ...prev, page: p }));
        }}
      />
    </div>
  );
}
