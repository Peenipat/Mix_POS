// src/page/admin/ManageBarber.tsx
import React, { useEffect, useState, useRef, useCallback, FormEvent } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";
import type { Barber } from "../../types/barber";
import { useAppSelector } from "../../store/hook";
import axios from "../../lib/axios";


export function ManageBarber() {
  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);

  const tenantId = me?.tenant_ids?.[0];
  const branchId = me?.branch_id;

  const [barbers, setBarbers] = useState<Barber[]>([]);
  const [loadingBarbers, setLoadingBarbers] = useState<boolean>(false);
  const [errorBarbers, setErrorBarbers] = useState<string | null>(null);
  const didFetchBarbers = useRef(false);

  const [isModalOpen, setIsModalOpen] = useState(false);

  // ฟังก์ชันดึงข้อมูล barbers
  const loadBarbers = useCallback(async () => {
    if (!tenantId || !branchId) return;
    setLoadingBarbers(true);
    setErrorBarbers(null);
    try {
      const res = await axios.get<{ status: string; data: Barber[] }>(
        `/barberbooking/tenants/${tenantId}/barbers/branches/${branchId}/barbers`
      );
      if (res.data.status !== "success") {
        throw new Error(res.data.status);
      }
      setBarbers(res.data.data);
    } catch (err: any) {
      setErrorBarbers(err.response?.data?.message || err.message || "Failed to load barbers");
    } finally {
      setLoadingBarbers(false);
    }
  }, [tenantId, branchId]);

  // เรียก loadBarbers ครั้งแรก เมื่อ me พร้อม
  useEffect(() => {
    if (
      statusMe === "succeeded" &&
      me &&
      tenantId &&
      branchId &&
      !didFetchBarbers.current
    ) {
      didFetchBarbers.current = true;
      loadBarbers();
    }
  }, [statusMe, me, tenantId, branchId, loadBarbers]);

  // เมื่อสร้างเสร็จ ให้รีเซ็ต flag แล้ว load ใหม่
  const handleCreated = () => {
    didFetchBarbers.current = false;
    loadBarbers();
  };

  if (statusMe === "loading") return <p>Loading user info…</p>;
  if (statusMe === "succeeded" && !me) return <p className="text-red-500">Not authenticated.</p>;
  if (statusMe === "succeeded" && me && (!tenantId || !branchId))
    return <p className="text-red-500">Cannot determine branch information.</p>;

  return (
    <div className="space-y-6">
      {/* ปุ่มเปิด Modal */}
      <div>
        <button onClick={() => setIsModalOpen(true)} className="btn btn-primary">
          + Create Barber
        </button>
      </div>

      {/* Modal สร้าง Barber */}
      <CreateBarberModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onCreate={handleCreated}
      />

      {/* แสดงสถานะ Loading/Error ของ barbers */}
      {loadingBarbers && <p>Loading barbers…</p>}
      {errorBarbers && <p className="text-red-500">Error loading barbers: {errorBarbers}</p>}

      {/* ตารางบาร์เบอร์ */}
      {!loadingBarbers && !errorBarbers && (
        <>
          <h2 className="text-xl font-semibold">Barbers in Branch {branchId}</h2>
          <DataTable<Barber>
            data={barbers}
            columns={[
              { header: "ID", accessor: "id" },
              { header: "Username", accessor: "username" },
              { header: "Email", accessor: "email" },
              { header: "Phone", accessor: "phone_number" },
            ]}
            onRowClick={(r) => console.log("row clicked", r)}
            actions={[]}
            showEdit={false}
            showDelete={false}
          />
        </>
      )}
    </div>
  );
}


interface User {
  id: number;
  username: string;
  email: string;
}

interface CreateBarberModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreate: () => void; // callback ให้ parent รีเฟรชลิสต์
}

function CreateBarberModal({
  isOpen,
  onClose,
  onCreate,
}: CreateBarberModalProps) {
  // ดึง me จาก Redux (ใช้ tenant_ids[0] กับ branch_id)
  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);

  const tenantId = me?.tenant_ids?.[0];
  const branchId = me?.branch_id;

  // 1) state สำหรับ dropdown รายชื่อ users
  const [users, setUsers] = useState<User[]>([]);
  const [loadingUsers, setLoadingUsers] = useState<boolean>(false);
  const [errorUsers, setErrorUsers] = useState<string | null>(null);

  // 2) ฟิลด์ฟอร์ม: selected user_id และ phoneNumber
  const [selectedUserId, setSelectedUserId] = useState<number | "">("");
  const [phoneNumber, setPhoneNumber] = useState<string>("");

  // 3) สถานะ loading / error ของการสร้าง barber
  const [loadingCreate, setLoadingCreate] = useState<boolean>(false);
  const [errorCreate, setErrorCreate] = useState<string | null>(null);

  // 4) โหลด list ของ Users เมื่อ modal เปิด
  useEffect(() => {
    if (!isOpen) return;

    setUsers([]);
    setErrorUsers(null);
    setLoadingUsers(true);

    // ดึง users ที่มีสิทธิ์เป็น “barber” หรือ “staff” ได้ (ปรับตาม API จริง)
    axios
      .get<{ status: string; data: User[] }>(`/core/tenant-user/tenants/${tenantId}`)
      .then((res) => {
        if (res.data.status !== "success") {
          throw new Error(res.data.status);
        }
        setUsers(res.data.data);
      })
      .catch((err) => {
        setErrorUsers(err.response?.data?.message || err.message || "Failed to load users");
      })
      .finally(() => {
        setLoadingUsers(false);
      });
  }, [isOpen]);

  // 5) รีเซ็ตฟอร์มเมื่อเปิด modal ใหม่
  useEffect(() => {
    if (isOpen) {
      setSelectedUserId("");
      setPhoneNumber("");
      setErrorCreate(null);
      setLoadingCreate(false);
    }
  }, [isOpen]);

  // 6) ฟังก์ชัน submit
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setErrorCreate(null);

    if (!tenantId || !branchId) {
      setErrorCreate("Cannot determine tenant or branch. Please try again later.");
      return;
    }
    if (!selectedUserId) {
      setErrorCreate("Please select a user.");
      return;
    }
    if (!phoneNumber.trim()) {
      setErrorCreate("Please enter a valid phone number.");
      return;
    }

    setLoadingCreate(true);
    try {
      const res = await axios.post<{ status: string }>(
        `/barberbooking/tenants/${tenantId}/barbers/branches/${branchId}`,
        {
          user_id: selectedUserId,
          phone_number: phoneNumber.trim(),
        }
      );
      if (res.data.status !== "success") {
        throw new Error(res.data.status);
      }
      onCreate(); // ให้ parent รีเฟรชลิสต์ใหม่
      onClose();  // ปิด modal
    } catch (err: any) {
      setErrorCreate(err.response?.data?.message || err.message || "Failed to create barber.");
    } finally {
      setLoadingCreate(false);
    }
  };

  // 7) ระหว่างรอโหลด me
  if (statusMe === "loading") return null;
  if (statusMe === "succeeded" && !me) return null;
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <form
        onSubmit={handleSubmit}
        className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-md p-6"
      >
        <h2 className="text-xl font-semibold mb-4 text-gray-900 dark:text-gray-100">
          Create Barber
        </h2>

        {/* 8) Dropdown เลือก User */}
        <div className="mb-4">
          <label htmlFor="user" className="block text-gray-700 dark:text-gray-200 mb-1">
            Select User
          </label>
          {loadingUsers ? (
            <p>Loading users…</p>
          ) : errorUsers ? (
            <p className="text-red-600">{errorUsers}</p>
          ) : (
            <select
              id="user"
              value={selectedUserId}
              onChange={(e) => setSelectedUserId(Number(e.target.value))}
              className="w-full select select-bordered"
            >
              <option value="">-- Choose a user --</option>
              {users.map((u) => (
                <option key={u.id} value={u.id}>
                  {u.username} ({u.email})
                </option>
              ))}
            </select>
          )}
        </div>

        {/* 9) Input เบอร์โทร */}
        <div className="mb-4">
          <label htmlFor="phone" className="block text-gray-700 dark:text-gray-200 mb-1">
            Phone Number
          </label>
          <input
            id="phone"
            type="text"
            value={phoneNumber}
            onChange={(e) => setPhoneNumber(e.target.value)}
            placeholder="0812345678"
            className="w-full input input-bordered"
          />
        </div>

        {errorCreate && <p className="text-sm text-red-600 mb-4">{errorCreate}</p>}

        <div className="flex justify-end space-x-2 mt-6">
          <button
            type="button"
            onClick={onClose}
            className="btn btn-ghost"
            disabled={loadingCreate}
          >
            Cancel
          </button>
          <button
            type="submit"
            className={`btn btn-primary ${loadingCreate ? "opacity-50 cursor-not-allowed" : ""}`}
            disabled={loadingCreate}
          >
            {loadingCreate ? "Creating…" : "Create"}
          </button>
        </div>
      </form>
    </div>
  );
}
