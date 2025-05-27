import React, { useEffect, useState, FormEvent } from 'react';
import axios from "../../../lib/axios";
import { toast } from 'react-toastify';
import { UserResponseSchema, CreateUserForm, User as UserResponse } from '../../../schemas/userSchema';
import { User } from '../../../schemas/userSchema';
import EditUserModal from "../components/EditUserModal";
import { CreateUserModal } from '../components/CreateUserModal';
import { DataTable, Column, Action } from '../components/DataTable';
import { z } from "zod";
import { Button } from '../components/Button';
import { UserList } from '../components/UserList';

export default function ManageUsers() {
  // State
  const [users, setUsers] = useState<UserResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedUser, setSelectedUser] = useState<UserResponse | null>(null);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);


  const [isAssignOpen, setIsAssignOpen] = useState(false);
  const [assignUser, setAssignUser] = useState<User | null>(null);
  const [tenants, setTenants] = useState<{ id: number; name: string }[]>([]);
  const [selectedTenantId, setSelectedTenantId] = useState<number>();
  const [currentTenant, setCurrentTenant] = useState<{ id: number; name: string } | null>(null);
  // Fetch users on mount
  const fetchUsers = async () => {
    setLoading(true);
    try {
      const res = await axios.get("/admin/users");
      const parsed = z.array(UserResponseSchema).safeParse(res.data);
      if (!parsed.success) throw parsed.error;
      setUsers(parsed.data);
    } catch (e: any) {
      console.error(e);
      toast.error("Failed to load users");
    } finally {
      setLoading(false);
    }
  };

  // 2. ฟังก์ชัน fetchTenants
  const fetchTenants = async () => {
    try {
      const res = await axios.get<{
        status: string;
        data: { id: number; name: string }[];
      }>("/core/tenant-route?active=true");
      if (res.data.status !== "success") throw new Error(res.data.status);
      setTenants(res.data.data);
    } catch (e) {
      console.error(e);
      toast.error("Failed to load tenants");
    }
  };

  // 3. เรียกทั้งสองใน useEffect
  useEffect(() => {
    fetchUsers();
    fetchTenants();
  }, []);

  useEffect(() => {
    if (!assignUser) return

    (async () => {
      try {
        const res = await axios.get<{
          status: string
          data: Array<{ id: number; name: string }>
        }>(`/core/tenant-user/user/${assignUser.id}`)

        if (res.data.status !== 'success') {
          throw new Error(res.data.status)
        }

        const list = res.data.data
        // เอา entry แรก ถ้าไม่มีเลยก็เป็น null
        setCurrentTenant(
          list.length > 0
            ? { id: list[0].id, name: list[0].name }
            : null
        )
      } catch (err) {
        console.error(err)
        setCurrentTenant(null)
      }
    })()
  }, [assignUser])

  const handleAssignClick = (u: User) => {
    setAssignUser(u);
    setSelectedTenantId(tenants[0]?.id);
    setIsAssignOpen(true);
  };

  // on submit assign
  const handleAssign = async (e: FormEvent) => {
    e.preventDefault();
    if (!assignUser || !selectedTenantId) return;
    try {
      await axios.post(
        `/core/tenant-user/tenants/${selectedTenantId}/users/${assignUser.id}`
      );
      toast.success(`User #${assignUser.id} assigned to tenant ${selectedTenantId}`);
      setIsAssignOpen(false);
    } catch (err: any) {
      toast.error(err.response?.data?.message || 'Assignment failed');
    }
  };

  // Handlers
  const handleEdit = (user: UserResponse) => {
    setSelectedUser(user);
    setIsEditOpen(true);
  };

  const handleSave = async (updated: UserResponse) => {
    try {
      await axios.put('/admin/change_role', {
        id: updated.id,
        role: updated.role,
      });
      setUsers((prev) =>
        prev.map((u) => (u.id === updated.id ? updated : u))
      );
      toast.success('User updated');
    } catch (err: any) {
      console.error(err);
      toast.error('Update failed');
    } finally {
      setIsEditOpen(false);
    }
  };

  const handleDelete = async (user: UserResponse) => {
    if (!confirm(`Delete ${user.username}?`)) return;
    try {
      await axios.delete(`/admin/users/${user.id}`);
      setUsers((prev) => prev.filter((u) => u.id !== user.id));
      toast.success('User deleted');
    } catch (err: any) {
      console.error(err);
      toast.error('Delete failed');
    }
  };

  const handleCreate = () => {
    setIsCreateOpen(false);
    fetchUsers();
  };

  // Columns definition
  const columns: Column<UserResponse>[] = [
    { header: '#', accessor: 'id' },
    { header: 'Username', accessor: 'username' },
    { header: 'Email', accessor: 'email' },
    { header: 'Role', accessor: 'role' },
    // { header: 'Created At', accessor: 'createdAt' },
    // { header: 'Updated At', accessor: 'updatedAt' },
    // { header: 'Deleted At', accessor: 'deletedAt' },
  ];

  const assignAction: Action<User> = {
    label: 'Add to tenant',
    onClick: handleAssignClick,
    className: 'text-green-600',
  };

  if (loading) return <div className="text-center p-4">Loading...</div>;

  return (
    <div className="p-4 space-y-4">
      <div className="flex justify-between items-center">

        <h1 className="text-3xl font-bold">Manage Users</h1>
        <Button color="default" onClick={() => setIsCreateOpen(true)}>Create User</Button>
      </div>

      <DataTable<UserResponse>
        data={users}
        columns={columns}
        showEdit
        showDelete
        onEdit={handleEdit}
        onDelete={handleDelete}
        actions={[assignAction]}
      />

      {isAssignOpen && assignUser && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <form
            onSubmit={handleAssign}
            className="bg-white p-6 rounded-lg shadow-lg w-full max-w-sm"
          >
            <h2 className="text-xl font-semibold mb-4">
              Assign User to Tenant
            </h2>
            <p className="mb-2">
              <strong>User ID:</strong> {assignUser.id}
            </p>
            <p className="mb-4">
              <strong>Username:</strong> {assignUser.username}
            </p>
            <p className="mb-4">
              <strong>Current Tenant:</strong>{' '}
              {currentTenant
                ? `${currentTenant.name} (#${currentTenant.id})`
                : 'None'}
            </p>
            <label className="block mb-2 font-medium">Select Tenant</label>
            <select
              value={selectedTenantId}
              onChange={(e) => setSelectedTenantId(+e.target.value)}
              className="w-full input input-bordered mb-4"
            >
              {tenants.map((t) => (
                <option key={t.id} value={t.id}>
                  {t.name} (#{t.id})
                </option>
              ))}
            </select>
            <div className="flex justify-end space-x-2">
              <button
                type="button"
                onClick={() => setIsAssignOpen(false)}
                className="btn btn-ghost"
              >
                Cancel
              </button>
              <button type="submit" className="btn btn-primary">
                Assign
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Edit Modal */}
      {selectedUser && (
        <EditUserModal
          isOpen={isEditOpen}
          user={selectedUser}
          onClose={() => setIsEditOpen(false)}
          onSave={handleSave}
        />
      )}

      {/* Create Modal */}
      <CreateUserModal
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onCreate={handleCreate}
      />

      <UserList users={users} />
    </div>
  );
}