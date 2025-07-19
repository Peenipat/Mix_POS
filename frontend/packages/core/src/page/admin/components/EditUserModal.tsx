import React from 'react'
import { useState } from "react";
import { User } from '../../../schemas/userSchema';
interface EditUserModalProps {
  isOpen: boolean;
  user: User;
  onClose: () => void;
  onSave: (updatedUser: User) => void;
}

export default function EditUserModal({ isOpen, user, onClose, onSave }: EditUserModalProps) {
  if (!isOpen) return null;
  const [formData, setFormData] = useState({
    username: user.username,
    email: user.email,
    role: user.role,
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave({ ...user, ...formData });
    onClose();
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/30 z-50">
      <div className="bg-white rounded-lg shadow-lg p-8 w-full max-w-md">
        <h2 className="text-2xl font-bold mb-6 text-center">Edit User</h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block mb-1 font-semibold">Username</label>
            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              className="input input-bordered w-full"
              required
            />
          </div>

          <div>
            <label className="block mb-1 font-semibold">Email</label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              className="input input-bordered w-full"
              required
            />
          </div>

          <div>
            <label className="block mb-1 font-semibold">Role</label>
            <select
              name="role"
              value={formData.role}
              onChange={handleChange}
              className="select select-bordered w-full"
              required
            >
              <option value="BRANCH_ADMIN">BRANCH_ADMIN</option>
              <option value="STAFF">STAFF</option>
              <option value="USER">USER</option>
            </select>
          </div>

          <div className="flex justify-end gap-4 pt-4">
            <button type="button" onClick={onClose} className="btn btn-ghost">
              Cancel
            </button>
            <button type="submit" className="btn btn-primary">
              Save
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
