import React from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import axios from '../../../lib/axios';
import { toast } from 'react-toastify';
import { CreateUserForm } from '../../../schemas/userSchema';
import { CreateUserSchema } from '../../../schemas/userSchema';



interface CreateUserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreate: (user: CreateUserForm) => void;
}

export default function CreateUserModal({
  isOpen,
  onClose,
  onCreate,
}: CreateUserModalProps) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<CreateUserForm>({
    resolver: zodResolver(CreateUserSchema),
    defaultValues: { username: '', email: '', role: 'USER' },
  });

  const submitHandler = async (data: CreateUserForm) => {
    try {
      const resp = await axios.post('/admin/create_users', data);
      toast.success('User created successfully');
      onCreate(data);
      reset();
      onClose();
    } catch (err: any) {
      console.error(err);
      toast.error(err.response?.data?.error || 'Failed to create user');
    }
  };

  if (!isOpen) return null;
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-md p-6">
        <h2 className="text-xl font-semibold mb-4 text-gray-900 dark:text-gray-100">
          Create New User
        </h2>
        <form onSubmit={handleSubmit(submitHandler)}>
          <div className="mb-4">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">
              Username
            </label>
            <input
              {...register('username')}
              className="w-full input input-bordered"
              placeholder="Enter username"
            />
            {errors.username && (
              <p className="mt-1 text-sm text-red-500">
                {errors.username.message}
              </p>
            )}
          </div>

          <div className="mb-4">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">
              Email
            </label>
            <input
              {...register('email')}
              className="w-full input input-bordered"
              placeholder="Enter email"
            />
            {errors.email && (
              <p className="mt-1 text-sm text-red-500">
                {errors.email.message}
              </p>
            )}
          </div>

          <div className="mb-4">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">
              Role
            </label>
            <select
              {...register('role')}
              className="w-full select select-bordered"
            >
              <option value="SAAS_SUPER_ADMIN">SaaS Super Admin</option>
              <option value="BRANCH_ADMIN">Branch Admin</option>
              <option value="STAFF">Staff</option>
              <option value="USER">User</option>
              <option value="TENANT">Tenant</option>
            </select>
            {errors.role && (
              <p className="mt-1 text-sm text-red-500">
                {errors.role.message}
              </p>
            )}
          </div>

          <div className="flex justify-end space-x-2 mt-6">
            <button
              type="button"
              onClick={() => {
                reset();
                onClose();
              }}
              className="btn btn-ghost"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="btn btn-primary"
            >
              {isSubmitting ? 'Creating...' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
