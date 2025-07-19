// src/page/admin/ManageBarber.tsx
import React, { useEffect, useState, useRef, useCallback, FormEvent } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { DataTable } from "../../../components/DataTable";
import type { Action, Column } from "../../../components/DataTable";
import type { Barber } from "../../../types/barber";
import { useAppSelector } from "../../../store/hook";
import axios from "../../../lib/axios";
import Modal from "@object/shared/components/Modal";
import type { EditBarberFormData } from "../../../schemas/barberSchema";
import { editBarberSchema } from "../../../schemas/barberSchema";
import type { ChangeEvent } from "react";
import { Card } from "@object/shared/components/Card";
import { TableCellsIcon } from "../../../components/icons/TableCellsIcon";
import { GridViewIcon } from "../../../components/icons/GridViewIcon"
import { useNavigate } from "react-router-dom";

export function ManageBarber() {
  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);
  const navigate = useNavigate();

  const tenantId = me?.tenant_ids?.[0];
  const branchId = Number(me?.branch_id);

  const [barbers, setBarbers] = useState<Barber[]>([]);
  const [editBarber, setEditBarber] = useState<Barber>()
  const [deleteBarber, setDeleteBarber] = useState<Barber>()
  const [loadingBarbers, setLoadingBarbers] = useState<boolean>(false);
  const [errorBarbers, setErrorBarbers] = useState<string | null>(null);
  const didFetchBarbers = useRef(false);

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = useState(false)

  const [viewMode, setViewMode] = useState<"table" | "card">("card");

  const handleDeleteSuccess = () => {
    if (deleteBarber) {
      setBarbers(prev => prev.filter(b => b.id !== deleteBarber.id));
    }
  };
  const handleView = (barberId:number)=>{
    console.log(barberId)
    navigate(`/admin/barber/${barberId}`)
  }

  const loadBarbers = useCallback(async () => {
    if (!tenantId || !branchId) return;
    setLoadingBarbers(true);
    setErrorBarbers(null);
    try {
      const res = await axios.get<{ status: string; data: Barber[] }>(
        `/barberbooking/branches/${branchId}/barbers`
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
  }, [branchId]);

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

  const editAction: Action<Barber> = {
    label: "Edit",
    onClick: (row) => console.log("edit Barber", row),
    className: "text-blue-600",
  };
  const deleteAction: Action<Barber> = {
    label: "Delete",
    onClick: (row) => console.log("delete Barber", row),
    className: "text-red-600",
  };

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
      {/* ปุ่มควบคุม */}
      <div className="flex flex-wrap items-center gap-4 justify-between">
        <button
          onClick={() => setIsCreateOpen(true)}
          className="btn btn-primary"
        >
          + เพิ่มช่างคนใหม่
        </button>

        <div className="flex gap-2">
          <button
            onClick={() => setViewMode("card")}
            className={`px-4 py-2 rounded-md border ${viewMode === "card"
              ? "bg-blue-600 text-white"
              : "bg-white text-gray-800 border-gray-300"
              }`}
          >
            <GridViewIcon className="w-6 h-6" />
          </button>
          <button
            onClick={() => setViewMode("table")}
            className={`px-4 py-2 rounded-md border ${viewMode === "table"
              ? "bg-blue-600 text-white"
              : "bg-white text-gray-800 border-gray-300"
              }`}
          >
            <TableCellsIcon className="w-6 h-6" />
          </button>

        </div>
      </div>

      {/* Loading/Error */}
      {loadingBarbers && <p>Loading barbers…</p>}
      {errorBarbers && (
        <p className="text-red-500">Error loading barbers: {errorBarbers}</p>
      )}

      {/* แสดงข้อมูล */}
      {!loadingBarbers && !errorBarbers && (
        <>
          <h2 className="text-xl font-semibold">
            ช่างในสาขาที่ {branchId}
          </h2>

          {viewMode === "table" ? (
            <DataTable<Barber>
              data={barbers}
              columns={[
                { header: "#", accessor: (_row, i) => i + 1 },
                { header: "ชื่อผู้ใช้", accessor: "username" },
                { header: "อีเมล์", accessor: "email" },
                { header: "เบอร์โทร", accessor: "phone_number" },
              ]}
              onEdit={(b) => {
                setEditBarber(b);
                setIsEditOpen(true);
              }}
              onDelete={(b) => {
                setDeleteBarber(b);
                setIsDeleteOpen(true);
              }}
              showEdit
              showDelete
            />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              {/* {barbers.map((barber) => (
                <Card
                  key={barber.id}
                  onView={() => handleView(barber.id)}
                  onEdit={() => {
                    setEditBarber(barber);
                    setIsEditOpen(true);
                  }}
                  onDelete={() => {
                    setDeleteBarber(barber);
                    setIsDeleteOpen(true);
                  }}
                >
                  <img
                    src={`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${barber.img_path}/${barber.img_name}`}
                    alt={barber.avatarAlt || "Barber"}
                    className="w-full h-64 object-cover object-top"
                  />
                  <div className="mt-2">
                    <h3 className="text-lg font-semibold">{barber.username}</h3>
                    {barber.jobTitle && (
                      <p className="text-sm text-gray-500">{barber.jobTitle}</p>
                    )}
                    {barber.email && (
                      <p className="text-sm text-gray-400">📧 {barber.email}</p>
                    )}
                    {barber.phone_number && (
                      <p className="text-sm text-gray-400">📞 {barber.phone_number}</p>
                    )}
                    {barber.sex && (
                      <p className="text-sm text-gray-400">👤 {barber.sex}</p>
                    )}
                  </div>
                </Card>
              ))} */}
            </div>

          )}

          {/* โมดัล */}
          <DeleteBarberModal
            isOpen={isDeleteOpen}
            barber={deleteBarber}
            onDelete={handleDeleteSuccess}
            onClose={() => setIsDeleteOpen(false)}
            onCreate={handleCreated}
          />
          <EditBarberModal
            isOpen={isEditOpen}
            barber={editBarber}
            onClose={() => setIsEditOpen(false)}
            onCreate={handleCreated}
          />
          <CreateBarberModal
            isOpen={isCreateOpen}
            onClose={() => setIsCreateOpen(false)}
            onCreate={handleCreated}
          />
        </>
      )}
    </div>
  );
}


interface User {
  user_id: number;
  username: string;
  email: string;
}
interface CreateBarberModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreate: () => void;
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
  const [selectedUserId, setSelectedUserId] = useState<number | "">();
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
      .get<{ status: string; data: User[] }>(`/barberbooking/tenants/${tenantId}/branches/${branchId}/users`)
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
  // useEffect(() => {
  //   if (isOpen) {
  //     setSelectedUserId("");
  //     setPhoneNumber("");
  //     setErrorCreate(null);
  //     setLoadingCreate(false);
  //   }
  // }, [isOpen]);

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
        `/barberbooking/tenants/${tenantId}/branches/${branchId}/create-barber`,
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
              value={selectedUserId || undefined}
              onChange={(e) => {
                const v = e.target.value;
                setSelectedUserId(v === "" ? "" : Number(v));
              }}
              className="w-full select select-bordered"
            >
              {users.length != 0 ? (
                <option value="">-- เลือกผู้ใช้ --</option>
              ) : <option value="">-- ไม่พบข้อมูล --</option>}


              {users.map((u) => (
                <option key={u.user_id} value={u.user_id}>
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

interface EditBarberModalProps {
  isOpen: boolean;
  barber: Barber | undefined;
  onClose: () => void;
  onCreate: () => void;
}
function EditBarberModal({
  isOpen,
  barber,
  onClose,
  onCreate,
}: EditBarberModalProps) {
  if (!isOpen || !barber) {
    return null;
  }

  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);

  const tenantId = me?.tenant_ids?.[0];
  const branchId = me?.branch_id;

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<EditBarberFormData>({
    resolver: zodResolver(editBarberSchema),
    defaultValues: {
      username: "",
      email: "",
      phone_number: "",
      img_path: "",
      img_name: "",
      profilePicture: "",
    },
  });

  useEffect(() => {
    if (isOpen && barber) {
      reset({
        username: barber.username,
        email: barber.email,
        phone_number: barber.phone_number,
        img_path: barber.img_path,
        img_name: barber.img_name,
      });
    }
  }, [isOpen, barber, reset]);

  if (statusMe === "loading") return null;
  if (statusMe === "succeeded" && !me) return null;
  const [previewUrl, setPreviewUrl] = useState<string | null>(`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${barber.img_path}/${barber.img_name}`);

  const onSubmit = async (data: EditBarberFormData) => {
    if (!tenantId || !branchId) {
      return;
    }

    try {
      const res = await axios.put<{ status: string }>(
        `/barberbooking/tenants/${tenantId}/barbers/${barber.id}`,
        {
          branch_id: branchId,
          user_id: barber.user_id,
          phone_number: data.phone_number,
          username: data.username,
          email: data.email,
        }
      );
      if (res.data.status !== "success") {
        throw new Error(res.data.status);
      }
      onCreate();
      onClose();
    } catch (err: any) {
      console.error(err);
    }
  };


  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Edit Barber">
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            รูปภาพโปรไฟล์
          </label>
          <input
            type="file"
            accept="image/*"
            {...register("profilePicture")}
            className="file-input file-input-bordered w-full"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) {
                const url = URL.createObjectURL(file);
                setPreviewUrl(url);
              }
            }}
          />
          {/* {errors.profilePicture && (
            <p className="text-red-600 text-sm mt-1">
              {errors.profilePicture.message}
            </p>
          )} */}
          {previewUrl && previewUrl !== "" && (
            <div className="mt-3">
              <img
                src={previewUrl}
                alt="Preview"
                className="h-20 w-24 object-cover rounded-md border object-top"
              />
            </div>
          )}
        </div>
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            ชื่อผู้ใช้
          </label>
          <input
            type="text"
            {...register("username")}
            className={`w-full input input-bordered ${errors.username ? "border-red-500" : ""
              }`}
          />
          {errors.username && (
            <p className="text-red-600 text-sm mt-1">
              {errors.username.message}
            </p>
          )}
        </div>

        {/* Email */}
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            อีเมล์
          </label>
          <input
            type="email"
            {...register("email")}
            readOnly
            className={`w-full input input-bordered bg-gray-200 ${errors.email ? "border-red-500" : ""}`}
          />
          {errors.email && (
            <p className="text-red-600 text-sm mt-1">
              {errors.email.message}
            </p>
          )}
        </div>

        {/* Phone Number */}
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            เบอร์โทร
          </label>
          <input
            type="text"
            {...register("phone_number")}
            placeholder="0812345678"
            className={`w-full input input-bordered ${errors.phone_number ? "border-red-500" : ""
              }`}
          />
          {errors.phone_number && (
            <p className="text-red-600 text-sm mt-1">
              {errors.phone_number.message}
            </p>
          )}
        </div>

        <div className="flex justify-end space-x-2 mt-6">
          <button
            type="button"
            onClick={onClose}
            className="btn btn-ghost"
            disabled={isSubmitting}
          >
            ยกเลิก
          </button>
          <button
            type="submit"
            className={`btn btn-primary ${isSubmitting ? "opacity-50 cursor-not-allowed" : ""
              }`}
            disabled={isSubmitting}
          >
            {isSubmitting ? "กำลังบันทึกข้อมูล…" : "ยืนยัน"}
          </button>
        </div>
      </form>
    </Modal>
  );
}

interface DeleteBarberModalProps {
  isOpen: boolean;
  barber: Barber | undefined;
  onDelete: () => void;
  onClose: () => void;
  onCreate: () => void;
}

function DeleteBarberModal({
  isOpen,
  barber,
  onDelete,
  onClose,
  onCreate,
}: DeleteBarberModalProps) {

  if (!barber) {
    return null;
  }
  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);

  const tenantId = me?.tenant_ids?.[0];
  const branchId = me?.branch_id;

  const [loadingCreate, setLoadingCreate] = useState<boolean>(false);
  const [errorCreate, setErrorCreate] = useState<string | null>(null);

  const handleDelete = async () => {
    try {
      const res = await axios.delete<{ status: string }>(
        `/barberbooking/tenants/${tenantId}/barbers/${barber.id}`
      );
      if (res.data.status === "success") {
        onDelete()
        onClose()
      }
    } catch (err: any) {
      console.error(err);
    }
  };

  if (statusMe === "loading") return null;
  if (statusMe === "succeeded" && !me) return null;
  if (!isOpen) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="ลบข้อมูลช่าง">

      <p>
        คุณต้องการลบข้อมูลช่าง <span className="text-red-500">{barber.username}</span> ใช่มั้ย
      </p>



      {errorCreate && <p className="text-sm text-red-600 mb-4">{errorCreate}</p>}

      <div className="flex justify-end space-x-2 mt-6">
        <button
          type="button"
          onClick={onClose}
          className="btn btn-ghost"
          disabled={loadingCreate}
        >
          ยกเลิก
        </button>
        <button
          type="button"
          className={`btn btn-primary ${loadingCreate ? "opacity-50 cursor-not-allowed" : ""}`}
          disabled={loadingCreate}
          onClick={handleDelete}
        >
          {loadingCreate ? "กำลังบันทึกข้อมูล…" : "ยืนยัน"}
        </button>
      </div>
    </Modal>
  );
}
