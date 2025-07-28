// src/page/admin/ManageService.tsx
import React, { useEffect, useState, useRef, FormEvent } from "react";
import { DataTable } from "../../components/DataTable";
import { useForm } from "react-hook-form";
import type { Action, Column } from "../../components/DataTable";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAppSelector } from "../../store/hook";
import axios from "../../lib/axios";
import Modal from "@object/shared/components/Modal";
import { TableCellsIcon } from "../../components/icons/TableCellsIcon";
import { GridViewIcon } from "../../components/icons/GridViewIcon"
import { Card } from "@object/shared/components/Card";

interface Service {
  id: number;
  name: string;
  price: number | null;
  duration: number | null;
  description: string;
  Img_path: string;
  Img_name: string;
  branch_id: number;
  tenant_id: number;
}

import { z } from "zod";
import { CardViewIcon } from "../../components/icons/CardViewIcon";
export const editServiceSchema = z.object({
  name: z.string().min(1, "กรุณากรอกชื่อบริการ").max(100, "ชื่อผู้ใช้ต้องไม่เกิน 100 ตัวอักษร"),
  price: z.number().nullable(),
  duration: z.number().nullable()
});

export type EditServiceFormData = z.infer<typeof editServiceSchema>;

export function ManageServiceOld() {

  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);
  const tenantId = me?.tenant_ids[0]
  const branchId = me?.branch_id

  // 3) state สำหรับเก็บรายการ service
  const [services, setServices] = useState<Service[]>([]);
  const [editService, setEditService] = useState<Service>()
  const [deleteService, setDeleteService] = useState<Service>()

  const [loadingServices, setLoadingServices] = useState<boolean>(false);
  const [errorServices, setErrorServices] = useState<string | null>(null);

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = useState(false)

  const [viewMode, setViewMode] = useState<"table" | "card">("card");

  const handleDeleteSuccess = () => {
    if (deleteService) {
      setServices(prev => prev.filter(b => b.id !== deleteService.id));
    }
  };

  const handleEditSuccess = () => {
    if (!editService) return;
    setServices(prevServices =>
      prevServices.map(s => {
        const currentId = String(s.id);
        const editedId = String(editService.id);

        return currentId === editedId
          ? editService
          : s;
      })
    );
  };

  // 4) ref เพื่อบล็อกไม่ให้ fetch ซ้ำ
  const didFetchServices = useRef(false);

  const loadServices = async () => {
    setLoadingServices(true);
    setErrorServices(null);
    try {
      const res = await axios.get<{ status: string; data: Service[] }>(
        `/barberbooking/tenants/${tenantId}/branch/${branchId}/services`
      );
      if (res.data.status !== "success") {
        throw new Error(res.data.status);
      }
      setServices(res.data.data);
    } catch (err: any) {
      setErrorServices(
        err.response?.data?.message || err.message || "Failed to load services"
      );
    } finally {
      setLoadingServices(false);
    }
  };


  useEffect(() => {
    if (
      statusMe === "succeeded" &&
      me &&
      tenantId &&
      branchId &&
      !didFetchServices.current
    ) {
      didFetchServices.current = true;
      loadServices();
    }
  }, [statusMe, me, tenantId, branchId, didFetchServices]);

  if (!tenantId) {
    return <p className="text-red-500">Cannot determine tenant information.</p>;
  }
  if (loadingServices) {
    return <p>Loading services…</p>;
  }
  if (errorServices) {
    return <p className="text-red-500">Error loading services: {errorServices}</p>;
  }


  // === กำหนด columns และ actions สำหรับ DataTable ===
  const columns: Column<Service>[] = [
    {
      header: "#",
      accessor: (_row, rowIndex) => rowIndex + 1,
    },
    { header: "ชื่อบริการ", accessor: "name" },
    { header: "คำอธิบาย", accessor: "description" },
    { header: "ราคา", accessor: "price" },
    { header: "ระยะเวลา (นาที)", accessor: "duration" },
  ];

  const handleCreated = () => {
    didFetchServices.current = false;
    loadServices();
  };

  return (
    <div>
      <div className="flex flex-wrap items-center gap-4 justify-between">
        <button onClick={() => setIsCreateOpen(true)} className="btn btn-primary">
          + เพิ่มบริการใหม่
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
            <TableCellsIcon className="w-6 h-6 " />
          </button>
        </div>
      </div>

      <h2 className="text-xl mt-4">Services for Tenant {tenantId}</h2>

      {viewMode === "table" && (
        <DataTable<Service>
          data={services}
          columns={columns}
          onRowClick={(r) => console.log("row clicked", r)}
          actions={[]}
          showEdit
          onEdit={(service) => {
            setEditService(service);
            setIsEditOpen(true);
          }}
          showDelete
          onDelete={(service) => {
            setDeleteService(service);
            setIsDeleteOpen(true);
          }}
        />
      )}

      {/* --- มุมมอง Card --- */}
      {viewMode === "card" && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {services.map((sv) => (
            <Card
              key={sv.id}
              onEdit={() => {
                setEditService(sv);
                setIsEditOpen(true);
              }}
              onDelete={() => {
                setDeleteService(sv);
                setIsDeleteOpen(true);
              }}
              showActions={true} // หรือไม่ใส่ก็ได้ ถ้า default = true
            >
              {sv.Img_path && sv.Img_name && (
                <img
                  src={`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${sv.Img_path}/${sv.Img_name}`}
                  alt={sv.name}
                  className="w-full h-48 object-cover rounded"
                />
              )}

              <div className="mt-2  flex justify-between flex-1">
                <div>
                  <h3 className="text-lg font-semibold mb-1">{sv.name}</h3>
                  <p className="text-gray-400 text-sm line-clamp-3">
                    {sv.description}
                  </p>
                  <p className="text-gray-400 text-sm">⏱ {sv.duration} นาที</p>
                </div>

                <div className="mt-4 flex items-center justify-between">
                  <span className="font-bold text-lg text-green-400">
                    ฿{sv.price}
                  </span>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}


      <CreateServiceModal
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onCreate={handleCreated}
      />

      <EditServiceModal
        isOpen={isEditOpen}
        service={editService}
        onEdit={(updatedService) => {
          setServices(prev =>
            prev.map(s =>
              s.id === updatedService.id ? updatedService : s
            )
          );
        }}
        onClose={() => setIsEditOpen(false)}
        onCreate={handleCreated}
      />

      <DeleteServiceModal
        isOpen={isDeleteOpen}
        service={deleteService}
        onDelete={handleDeleteSuccess}
        onClose={() => setIsDeleteOpen(false)}
        onCreate={handleCreated}
      />
    </div>
  );
}

interface CreateServiceModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreate: () => void;
}
function CreateServiceModal({
  isOpen,
  onClose,
  onCreate,
}: CreateServiceModalProps) {
  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);

  const tenantId = me?.tenant_ids?.[0];
  const branchId = me?.branch_id;

  const [loadingUsers, setLoadingUsers] = useState<boolean>(false);
  const [errorUsers, setErrorUsers] = useState<string | null>(null);

  const [selectedUserId, setSelectedUserId] = useState<number | "">();
  const [serviceName, setServiceName] = useState<string>("");
  const [servicePrice, setServicePrice] = useState<Number | undefined>();
  const [serviceDuration, setServiceDuration] = useState<Number | undefined>();

  // 3) สถานะ loading / error ของการสร้าง barber
  const [loadingCreate, setLoadingCreate] = useState<boolean>(false);
  const [errorCreate, setErrorCreate] = useState<string | null>(null);

  // 4) โหลด list ของ Users เมื่อ modal เปิด
  useEffect(() => {
    if (!isOpen) return;

    setErrorUsers(null);
    setLoadingUsers(true);

    // ดึง users ที่มีสิทธิ์เป็น “barber” หรือ “staff” ได้ (ปรับตาม API จริง)

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
    if (!serviceName.trim()) {
      setErrorCreate("กรุณากรอกชื่อบริการ.");
      return;
    }
    if (!servicePrice) {
      setErrorCreate("กรุณากรอกราคา")
    }
    if (!serviceDuration) {
      setErrorCreate("กรุณากรอกระยะเวลา")
    }

    setLoadingCreate(true);
    try {
      const res = await axios.post<{ status: string }>(
        `/barberbooking/tenants/${tenantId}/branch/${branchId}/services`,
        {
          name: serviceName.trim(),
          price: servicePrice,
          duration: serviceDuration,
        }
      );
      if (res.data.status !== "success") {
        throw new Error(res.data.status);
      }
      onCreate(); // ให้ parent รีเฟรชลิสต์ใหม่
      onClose();  // ปิด modal
    } catch (err: any) {
      setErrorCreate(err.response?.data?.message || err.message || "Failed to create service.");
    } finally {
      setLoadingCreate(false);
    }
  };

  // 7) ระหว่างรอโหลด me
  if (statusMe === "loading") return null;
  if (statusMe === "succeeded" && !me) return null;
  if (!isOpen) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="เพิ่มบริการใหม่">
      <form onSubmit={handleSubmit}>

        {/* 9) Input เบอร์โทร */}
        <div className="mb-4">
          <label htmlFor="phone" className="block text-gray-700 dark:text-gray-200 mb-1">
            ชื่อบริการ
          </label>
          <input
            id="name"
            type="text"
            value={serviceName}
            onChange={(e) => setServiceName(e.target.value)}
            placeholder="กรุณากรอกชื่อบริการ"
            className="w-full input input-bordered"
          />
        </div>

        <div className="mb-4">
          <label htmlFor="phone" className="block text-gray-700 dark:text-gray-200 mb-1">
            ราคา
          </label>
          <input
            id="price"
            type="number"
            value={Number(servicePrice)}
            onChange={(e) => setServicePrice(Number(e.target.value))}
            placeholder="กรุณากรอกราคา"
            className="w-full input input-bordered"
          />
        </div>

        <div className="mb-4">
          <label htmlFor="phone" className="block text-gray-700 dark:text-gray-200 mb-1">
            ระยะเวลาโดยประมาณ (นาที)
          </label>
          <input
            id="duration"
            type="number"
            value={Number(serviceDuration)}
            onChange={(e) => setServiceDuration(Number(e.target.value))}
            placeholder="กรุณากรอกระยะเวลา"
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
            ยกเลิก
          </button>
          <button
            type="submit"
            className={`btn btn-primary ${loadingCreate ? "opacity-50 cursor-not-allowed" : ""}`}
            disabled={loadingCreate}
          >
            {loadingCreate ? "กำลังบันทึก..." : "ยืนยัน"}
          </button>
        </div>
      </form>
    </Modal>
  );
}

interface EditServiceModalProps {
  isOpen: boolean;
  service: Service | undefined;
  onEdit: (updatedService: Service) => void;
  onClose: () => void;
  onCreate: () => void;
}
function EditServiceModal({
  isOpen,
  service,
  onClose,
  onEdit,
}: EditServiceModalProps) {
  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);

  const tenantId = me?.tenant_ids?.[0];
  const branchId = me?.branch_id;
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<EditServiceFormData>({
    resolver: zodResolver(editServiceSchema),
    defaultValues: {
      name: "",
      price: null,
      duration: null,
    },
  });

  useEffect(() => {
    if (isOpen && service) {
      reset({
        name: service.name,
        price: service.price,
        duration: service.duration,
      });
    }
  }, [isOpen, service, reset]);

  if (statusMe === "loading") return null;
  if (statusMe === "succeeded" && !me) return null;
  if (!isOpen || !service) return null;


  const onSubmit = async (data: EditServiceFormData) => {
    if (!tenantId || !branchId) {
      return;
    }

    try {
      const res = await axios.put<{ status: string }>(
        `/barberbooking/tenants/${tenantId}/branch/${branchId}/services/${Number(service.id)}`,
        {
          name: data.name,
          price: data.price,
          duration: data.duration
        }
      );
      if (res.data.status !== "success") {
        throw new Error(res.data.status);
      }
      const updatedService: Service = {
        ...service,
        name: data.name,
        price: data.price,
        duration: data.duration,
      };
      onEdit(updatedService)
      onClose();
    } catch (err: any) {
      console.error(err);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="แก้ไขข้อมูลบริการ">
      <form onSubmit={handleSubmit(onSubmit)}>
        {/* Username */}
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            ชื่อบริการ
          </label>
          <input
            type="text"
            {...register("name")}
            className={`w-full input input-bordered ${errors.name ? "border-red-500" : ""
              }`}
          />
          {errors.name && (
            <p className="text-red-600 text-sm mt-1">
              {errors.name.message}
            </p>
          )}
        </div>

        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            ราคา
          </label>
          <input
            type="number"
            {...register("price")}
            className={`w-full input input-bordered ${errors.price ? "border-red-500" : ""
              }`}
          />
          {errors.price && (
            <p className="text-red-600 text-sm mt-1">
              {errors.price.message}
            </p>
          )}
        </div>

        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            ระยะเวลาโดยประมาณ (นาที)
          </label>
          <input
            type="number"
            {...register("duration")}
            placeholder="กรุณากรอกระยะเวลา"
            className={`w-full input input-bordered ${errors.duration ? "border-red-500" : ""
              }`}
          />
          {errors.duration && (
            <p className="text-red-600 text-sm mt-1">
              {errors.duration.message}
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

interface DeleteServiceModalProps {
  isOpen: boolean;
  service: Service | undefined;
  onDelete: () => void;
  onClose: () => void;
  onCreate: () => void;
}

function DeleteServiceModal({
  isOpen,
  service,
  onDelete,
  onClose,
  onCreate

}: DeleteServiceModalProps) {

  if (!service) {
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
        `/barberbooking/services/${Number(service.id)}`
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
    <Modal isOpen={isOpen} onClose={onClose} title="ลบบริการ">

      <p>
        คุณต้องการลบบริการ <span className="text-red-500">{service.name}</span> ใช่มั้ย
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
