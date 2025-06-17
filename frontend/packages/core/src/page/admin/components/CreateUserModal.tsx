import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import axios from '../../../lib/axios';
import { toast } from 'react-toastify';
import { CreateUserForm } from '../../../schemas/userSchema';
import { CreateUserSchema } from '../../../schemas/userSchema';
import React, { useState, useEffect } from "react";

interface CreateUserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreate: () => void;
}

export function CreateUserModal({ isOpen, onClose, onCreate }: CreateUserModalProps) {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("USER");
  const [branchId, setBranchId] = useState<number | "">("");
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);

  // สมมติดึง branch list มาใส่ dropdown
  const [branches, setBranches] = useState<{ id: number; name: string }[]>([]);
  useEffect(() => {
    const fetchBranches = async () => {
      try {
        const res = await axios.get<{
          status: string;
          data: { id: number; name: string }[];
        }>("/core/branches/all");
        if (res.data.status !== "success") {
          throw new Error(res.data.status);
        }
        // เอาแค่ array ของ branches ลง state
        setBranches(res.data.data);
      } catch (e) {
        console.error(e);
        toast.error("Failed to load branches");
      }
    };
    fetchBranches();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const formData = new FormData();
      formData.append("username", username);
      formData.append("email", email);
      formData.append("password", password);
      formData.append("role", role);
      if (branchId) formData.append("branch_id", branchId.toString());
      if (avatarFile) formData.append("file", avatarFile);
      formData.append("keyprefix", "user_profile")

      await axios.post("/admin/create_users", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      toast.success("User created!");
      onCreate();
      onClose();
    } catch (err: any) {
      console.error(err);
      toast.error(err.response?.data?.error || "Failed to create user");
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
      <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg space-y-4 w-96">
        <h2 className="text-xl font-bold">Create New User</h2>

        <label className="block">
          <span>Username</span>
          <input
            type="text"
            value={username}
            onChange={e => setUsername(e.target.value)}
            required
            className="mt-1 block w-full border rounded p-2"
          />
        </label>

        <label className="block">
          <span>Email</span>
          <input
            type="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
            className="mt-1 block w-full border rounded p-2"
          />
        </label>

        <label className="block">
          <span>Password</span>
          <input
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
            className="mt-1 block w-full border rounded p-2"
          />
        </label>

        <label className="block">
          <span>Role</span>
          <select
            value={role}
            onChange={e => setRole(e.target.value)}
            className="mt-1 block w-full border rounded p-2"
          >
            {["TENANT_ADMIN", "BRANCH_ADMIN", "ASSISTANT_MANAGER", "STAFF", "USER"].map(r => (
              <option key={r} value={r}>{r}</option>
            ))}
          </select>
        </label>

        <label className="block">
          <span>Branch (optional)</span>
          <select
            value={branchId}
            onChange={e => setBranchId(e.target.value ? Number(e.target.value) : "")}
            className="mt-1 block w-full border rounded p-2"
          >
            <option value="">— none —</option>
            {branches.map(b => (
              <option key={b.id} value={b.id}>{b.name} #{b.id}</option>
            ))}
          </select>
        </label>

        <label className="block">
          <span>Avatar (optional)</span>
          <input
            type="file"
            accept="image/*"
            onChange={e => setAvatarFile(e.target.files?.[0] || null)}
            className="mt-1 block w-full"
          />
        </label>

        <div className="flex justify-end space-x-2">
          <button type="button" onClick={onClose} disabled={loading}
            className="px-4 py-2 border rounded">Cancel</button>
          <button type="submit" disabled={loading}
            className="px-4 py-2 bg-blue-600 text-white rounded">
            {loading ? "Saving..." : "Create"}
          </button>
        </div>
      </form>
    </div>
  );
}

