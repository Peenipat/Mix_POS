import { useEffect, useState } from "react";
import axios from "@/lib/axios";
import { UserResponseSchema } from "@/schemas/userSchema";
import EditUserModal from "../admin/components/EditUserModel";
import { z } from "zod"

export interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string | null;
}

export default function ManageUsers() {
  const usersResponseSchema = z.array(UserResponseSchema);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const token = localStorage.getItem("token");
        const res = await axios.get("/admin/users", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const parsed = usersResponseSchema.safeParse(res.data);

        if (!parsed.success) {
          console.error("Invalid API Response:", parsed.error);
          setUsers([]);
          return
        }
        setUsers(parsed.data);
      } catch (err) {
        setUsers([]);
        console.error("Fetch users error:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
  }, []);

  const handleEdit = (user: User) => {
    setSelectedUser(user);
    setIsModalOpen(true);
  };
  
  const handleSave = (updatedUser: User) => {
    setUsers((prev) => prev.map(u => (u.id === updatedUser.id ? updatedUser : u)));
    // TODO: ยิง axios.put("/admin/users/:id") ไป save ฝั่ง backend จริง ๆ
  };

  if (loading) return <div className="text-center">Loading...</div>

  return (
    <div className="overflow-x-auto p-4">
      <h1 className="text-3xl font-bold mb-6">Manage Users</h1>

      <table className="table table-zebra w-full">
        {/* head */}
        <thead className="bg-gray-200 text-gray-700">
          <tr>
            <th className="text-center">#</th>
            <th className="text-center">Username</th>
            <th className="text-center">Email</th>
            <th className="text-center">Role</th>
            <th className="text-center">Created At</th>
            <th className="text-center">Updated At</th>
            <th className="text-center">Deleted At</th>
            <th className="text-center">Actions</th>
          </tr>
        </thead>

        <tbody>
          {users.map((user,index) => (
            <tr key={user.id}>
              <td className="text-center" >{index + 1}</td>
              <td className="text-center" >{user.username}</td>
              <td >{user.email}</td>
              <td className="text-center">{user.role}</td>
              <td className="text-center">
                {user.createdAt
                  ? new Date(user.createdAt).toLocaleDateString('th-TH')
                  : "-"}
              </td>
              <td className="text-center">
                {user.updatedAt
                  ? new Date(user.updatedAt).toLocaleDateString('th-TH')
                  : "-"}
              </td>
              <td className="text-center">
                {user.deletedAt
                  ? new Date(user.deletedAt).toLocaleDateString('th-TH')
                  : "-"}
              </td>
              <td className="flex justify-between text-center">
                <button className="btn btn-warning"  onClick={() => handleEdit(user)}>Edit</button>
                <button className="btn btn-error">Delete</button>
                </td>
            </tr>
          ))}
        </tbody>
      </table>
      {isModalOpen && selectedUser && (
        <EditUserModal
          user={selectedUser}
          onClose={() => setIsModalOpen(false)}
          onSave={handleSave}
        />
      )}
    </div>
  );
}
