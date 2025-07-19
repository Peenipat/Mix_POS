import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from '../../../lib/axios';
import { DataTable, Column } from '../components/DataTable';

interface Tenant {
  id: number;
  name: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
  // ถ้ามี domain ให้เพิ่มที่นี่
}

interface Branch {
  id: number;
  name: string;
  address: string;
  createdAt: string;
  updatedAt: string;
}

// กำหนด columns สำหรับ branches
const branchColumns: Column<Branch>[] = [
  { header: 'ID', accessor: 'id' },
  { header: 'Name', accessor: 'name' },
  { header: 'Address', accessor: 'address' },
  { header: 'Created At', accessor: 'createdAt' },
  { header: 'Updated At', accessor: 'updatedAt' },
];

export default function TenantDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [branches, setBranches] = useState<Branch[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    (async () => {
      try {
        // fetch tenant
        const resT = await axios.get<{ status: string; data: Tenant }>(
          `/core/tenant-route/${id}`
        );
        if (resT.data.status !== 'success') throw new Error(resT.data.status);
        setTenant(resT.data.data);

        // fetch branches
        const resB = await axios.get<{ status: string; data: Branch[] }>(
          `core/tenants/${id}/branches`
        );
        if (resB.data.status !== 'success') throw new Error(resB.data.status);
        setBranches(resB.data.data);
      } catch (e: any) {
        console.error(e);
        setError(e.message);
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  if (loading) {
    return <div className="p-6 text-center">Loading…</div>;
  }
  if (error) {
    return <div className="p-6 text-center text-red-500">Error: {error}</div>;
  }
  if (!tenant) {
    return <div className="p-6 text-center">Tenant not found</div>;
  }

  return (
    <div className="p-6 space-y-8 max-w-5xl mx-auto">
      <button
        className="btn btn-ghost mb-2"
        onClick={() => navigate(-1)}
      >
        ← Back
      </button>

      {/* Tenant Card */}
      <a
        className="block w-full p-6 bg-white border border-gray-200 rounded-lg shadow-sm hover:bg-gray-100 dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700"
      >
        <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          Tenant #{tenant.id}: {tenant.name}
        </h5>
        <p className="font-normal text-gray-700 dark:text-gray-400">
          <strong>Active:</strong> {tenant.active ? 'Yes' : 'No'}
        </p>
        <p className="font-normal text-gray-700 dark:text-gray-400">
          <strong>Created At:</strong>{' '}
          {new Date(tenant.createdAt).toLocaleString()}
        </p>
        <p className="font-normal text-gray-700 dark:text-gray-400">
          <strong>Updated At:</strong>{' '}
          {new Date(tenant.updatedAt).toLocaleString()}
        </p>
      </a>

      {/* Branches Section */}
      <div className="space-y-4">
        <h2 className="text-xl font-semibold">Branches</h2>
        <div className="overflow-x-auto">
          <DataTable<Branch>
            data={branches}
            columns={branchColumns}
            showEdit={false}
            showDelete={false}
            actions={[
              {
                label: 'View',
                onClick: (b) => navigate(`/admin/tenant/${id}/branch/${b.id}`),
                className: 'text-blue-600 dark:text-blue-400'
              },
            ]}
          />
        </div>
      </div>
    </div>
  );
}
