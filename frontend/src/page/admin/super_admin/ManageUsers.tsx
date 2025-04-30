import { useEffect, useState } from "react";
import axios from "@/lib/axios";
import { toast } from 'react-toastify'
import { UserResponseSchema } from "@/schemas/userSchema";
import EditUserModal from "../components/EditUserModal";
import { z } from "zod"

// interface กำหนดโครงสร้างข้อมูล user
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
  const usersResponseSchema = z.array(UserResponseSchema); // schema สำหรับ validate array ของ users  ที่ตอบกลับมาจาก Database
  const [users, setUsers] = useState<User[]>([]); // เก็บข้อมูล user ทั้งหมดที่ดึงจาก backend
  const [loading, setLoading] = useState(true); // สถานะโหลดหน้า
  const [selectedUser, setSelectedUser] = useState<User | null>(null);  // user ที่ถูกเลือกเพื่อแก้ไข
  const [isModalOpen, setIsModalOpen] = useState(false); // สถานะเปิด modal แก้ไข

  // ดึงข้อมูล users ตอนเปิด
  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const res = await axios.get("/admin/users")
        const parsed = usersResponseSchema.safeParse(res.data); // validate response ของ user ว่าถูกต้องตาม schema ไหม

        if (!parsed.success) {
          console.error("Invalid API Response:", parsed.error);
          setUsers([]); //ไม่ผ่านให้ส่ง array เปล่า
          return
        }
        setUsers(parsed.data); // ผ่าน set user เข้า useState
      } catch (err) {
        setUsers([]);
        console.error("Fetch users error:", err);
      } finally {
        setLoading(false); // หยุดโหลดไม่ว่าจะสำเร็จหรือล้มเหลว
      }
    };

    fetchUsers();
  }, []);

  // ฟังก์ชันเปิด modal พร้อมส่ง user ที่ต้องการแก้ไข
  const handleEdit = (user: User) => {
    setSelectedUser(user);
    setIsModalOpen(true);
  };

  // ฟังก์ชัน save user หลังจากแก้ไขจาก modal
  const handleSave = async (updatedUser: User) => {
    try {
      await axios.put('/admin/change_role', {
        id: updatedUser.id,
        role: updatedUser.role,
      })
      // อัปเดต state เมื่อสำเร็จ
      setUsers(prev =>
        prev.map(u => (u.id === updatedUser.id ? updatedUser : u))
      )
      toast.success('User updated successfully')
      setIsModalOpen(false)  // ปิด modal แทน onClose()
    } catch (err: any) {
      console.error(err)
      toast.error(err.response?.data?.error || 'Failed to update user')
    }
  }

  // กรณีโหลดข้อมูลยังไม่เสร็จ
  if (loading) return <div className="text-center">Loading...</div>

  return (
    <div className="overflow-x-auto p-4">
      <h1 className="text-3xl font-bold mb-6">Manage Users</h1>
      <button className="btn btn-success">Create user</button>
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
          {users.map((user, index) => (
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
                <button className="btn btn-warning" onClick={() => handleEdit(user)}>Edit</button>
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
