import React, { useEffect, useState } from 'react';
import axios from "../../../lib/axios";
import { toast } from 'react-toastify';
import { UserResponseSchema, CreateUserForm, User as UserResponse } from '../../../schemas/userSchema';
import EditUserModal from "../components/EditUserModal";
import CreateUserModal from '../components/CreateUserModal';
import { DataTable, Column } from '../components/DataTable';
import { z } from "zod";

export default function ManageUsers() {
  // State
  const [users, setUsers] = useState<UserResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedUser, setSelectedUser] = useState<UserResponse | null>(null);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  // Fetch users on mount
  useEffect(() => {
    (async () => {
      try {
        const res = await axios.get('/admin/users');
        const parsed = z.array(UserResponseSchema).safeParse(res.data);
        if (!parsed.success) throw parsed.error;
        setUsers(parsed.data);
      } catch (err: any) {
        console.error(err);
        toast.error('Failed to load users');
      } finally {
        setLoading(false);
      }
    })();
  }, []);

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

  const handleCreate = (newUser: CreateUserForm) => {
    // Optimistic add, backend should return real ID
    const created: UserResponse = {
      id: Date.now(),
      username: newUser.username,
      email: newUser.email,
      role: newUser.role,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      deletedAt: null,
    };
    setUsers((prev) => [...prev, created]);
    setIsCreateOpen(false);
  };

  // Columns definition
  const columns: Column<UserResponse>[] = [
    { header: '#', accessor: 'id' },
    { header: 'Username', accessor: 'username' },
    { header: 'Email', accessor: 'email' },
    { header: 'Role', accessor: 'role' },
    { header: 'Created At', accessor: 'createdAt' },
    { header: 'Updated At', accessor: 'updatedAt' },
    { header: 'Deleted At', accessor: 'deletedAt' },
  ];

  if (loading) return <div className="text-center p-4">Loading...</div>;

  return (
    <div className="p-4 space-y-4">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Manage Users</h1>
        <button
          className="btn btn-success"
          onClick={() => setIsCreateOpen(true)}
        >
          Create User
        </button>
      </div>

      <DataTable<UserResponse>
        data={users}
        columns={columns}
        showEdit
        showDelete
        onEdit={handleEdit}
        onDelete={handleDelete}
        actions={[
          {
            label: 'View',
            onClick: (u) => console.log('View', u),
            className: 'text-blue-600',
          },
        ]}
      />

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
    </div>
  );
}