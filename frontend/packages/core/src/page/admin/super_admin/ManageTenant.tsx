import React, { useState, useEffect, FormEvent } from 'react';
import axios from '../../../lib/axios';
import { DataTable, Column } from '../components/DataTable';
import { useNavigate } from 'react-router-dom';

interface Tenant {
  id: number;
  name: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

const columns: Column<Tenant>[] = [
  { header: 'ID', accessor: 'id' },
  { header: 'Name', accessor: 'name' },
  { header: 'Active', accessor: 'active' },
  { header: 'Created', accessor: 'createdAt' },
  { header: 'Updated', accessor: 'updatedAt' },
];

export default function ManageTenant() {
  const navigate = useNavigate();

  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // state for edit modal
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [editingTenant, setEditingTenant] = useState<Tenant | null>(null);
  const [isCreateOpen, setIsCreateOpen] = useState(false);


  useEffect(() => {
    (async () => {
      try {
        const res = await axios.get<{ status: string; data: Tenant[] }>('/core/tenant-route?active=true');
        if (res.data.status !== 'success') throw new Error(res.data.status);
        setTenants(res.data.data);
      } catch (e: any) {
        console.error(e);
        setError(e.message);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const handleSave = (updated: Tenant) => {
    setTenants((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
    setIsEditOpen(false);
  };

  const handleCreate = (newTenant: Tenant) => {
    setTenants(prev => [...prev, newTenant]);
  };

  if (loading) return <div className="p-4 text-center">Loading tenants…</div>;
  if (error) return <div className="p-4 text-center text-red-500">Error: {error}</div>;

  return (
    <div className="p-4">
      <h1 className="text-3xl font-bold mb-4">Manage Tenants</h1>
      <button
        className="btn btn-success"
        onClick={() => setIsCreateOpen(true)}
      >
        Create Tenant
      </button>
      <DataTable<Tenant>
        data={tenants}
        columns={columns}
        showEdit={true}
        showDelete={false}
        onEdit={(tenant) => {
          setEditingTenant(tenant);
          setIsEditOpen(true);
        }}
        actions={[
          {
            label: 'View',
            onClick: (row) => navigate(`/admin/tenant/${row.id}`),
            className: 'text-green-600'
          }
        ]}
      />

      {/* Inline Edit Modal */}
      {editingTenant && (
        <EditTenantModal
          isOpen={isEditOpen}
          tenant={editingTenant}
          onClose={() => setIsEditOpen(false)}
          onSave={handleSave}
        />
      )}

      <CreateTenantModal
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onCreate={handleCreate}
      />
    </div>
  );
}

// Edit modal defined in same file
interface EditTenantModalProps {
  isOpen: boolean;
  tenant: Tenant;
  onClose: () => void;
  onSave: (t: Tenant) => void;
}
function EditTenantModal({ isOpen, tenant, onClose, onSave }: EditTenantModalProps) {
  const [name, setName] = useState(tenant.name);
  const [active, setActive] = useState(tenant.active);

  useEffect(() => {
    setName(tenant.name);
    setActive(tenant.active);
  }, [tenant]);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await axios.put(
        `/core/tenant-route/${tenant.id}`,
        {
          name,
          isActive: active,
          // domain: xxx     // ถ้ามี field Domain ใน modal ก็ส่งมาได้เลย
        }
      );
      onSave({
        ...tenant,
        name,
        active,
        updatedAt: new Date().toISOString(),
      });
    } catch (err: any) {
      console.error(err);
      alert(err.response?.data?.message || 'Failed to update');
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <form
        onSubmit={submit}
        className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-md p-6"
      >
        <h2 className="text-xl font-semibold mb-4 text-gray-900 dark:text-gray-100">
          Edit Tenant
        </h2>
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">Name</label>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full input input-bordered"
          />
        </div>
        <div className="mb-4 flex items-center">
          <input
            id="active"
            type="checkbox"
            checked={active}
            onChange={(e) => setActive(e.target.checked)}
            className="checkbox checkbox-primary mr-2"
          />
          <label htmlFor="active" className="text-gray-700 dark:text-gray-200">
            Active
          </label>
        </div>
        <div className="flex justify-end space-x-2 mt-6">
          <button
            type="button"
            onClick={onClose}
            className="btn btn-ghost"
          >
            Cancel
          </button>
          <button type="submit" className="btn btn-primary">
            Save
          </button>
        </div>
      </form>
    </div>
  );
}



interface CreateTenantModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreate: (newTenant: Tenant) => void;
}
function CreateTenantModal({ isOpen, onClose, onCreate }: CreateTenantModalProps) {
  const [name, setName] = useState('');
  const [domain, setDomain] = useState('');
  const [active, setActive] = useState(true);
  const [saving, setSaving] = useState(false);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      // เรียก API สร้าง tenant ใหม่
      const res = await axios.post<{
        status: string;
        data: { id: number };
      }>('/core/tenant-route/create', {
        name,
        domain,
        isActive: active,
      });
      if (res.data.status !== 'success') {
        throw new Error(`API status: ${res.data.status}`);
      }
      // สร้าง object ใหม่ขึ้นมาใช้ในตาราง
      const newTenant: Tenant = {
        id: res.data.data.id,
        name,
        active,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
      onCreate(newTenant);
      // เคลียร์ฟอร์มและปิด modal
      setName('');
      setDomain('');
      setActive(true);
      onClose();
    } catch (err: any) {
      console.error(err);
      alert(err.response?.data?.message || 'Failed to create tenant');
    } finally {
      setSaving(false);
    }
  };

  if (!isOpen) return null;
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <form
        onSubmit={submit}
        className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-md p-6"
      >
        <h2 className="text-xl font-semibold mb-4 text-gray-900 dark:text-gray-100">
          Create Tenant
        </h2>

        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            Name
          </label>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full input input-bordered"
            placeholder="Tenant name"
            required
          />
        </div>

        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            Domain
          </label>
          <input
            value={domain}
            onChange={(e) => setDomain(e.target.value)}
            className="w-full input input-bordered"
            placeholder="yourdomain.com"
            required
          />
        </div>

        <div className="mb-4 flex items-center">
          <input
            id="create-active"
            type="checkbox"
            checked={active}
            onChange={(e) => setActive(e.target.checked)}
            className="checkbox checkbox-primary mr-2"
          />
          <label htmlFor="create-active" className="text-gray-700 dark:text-gray-200">
            Active
          </label>
        </div>

        <div className="flex justify-end space-x-2 mt-6">
          <button
            type="button"
            onClick={onClose}
            className="btn btn-ghost"
            disabled={saving}
          >
            Cancel
          </button>
          <button
            type="submit"
            className="btn btn-primary"
            disabled={saving}
          >
            {saving ? 'Creating...' : 'Create'}
          </button>
        </div>
      </form>
    </div>
  );
}
