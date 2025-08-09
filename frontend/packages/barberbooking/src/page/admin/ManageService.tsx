// src/page/admin/ManageService.tsx
import React, { useEffect, useState, useRef, FormEvent } from "react";
import { DataTable } from "../../components/DataTable";
import { FieldError, useForm } from "react-hook-form";
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
import { ServiceFormData, serviceFormSchema } from "../../schemas/serviceSchema";
import { makeToast } from "../../utils/makeToast";
export const editServiceSchema = z.object({
  name: z.string().min(1).max(100),
  price: z.number({ required_error: "กรุณากรอกราคา" }).min(0),
  duration: z.number({ required_error: "กรุณากรอกระยะเวลา" }).min(1),
  file: z.any().optional().refine(
    (val) => val === undefined || val instanceof FileList || val instanceof File,
    { message: "รูปภาพไม่ถูกต้อง" }
  ),
});

// ✅ ใช้ input แทน infer
export type EditServiceFormData = z.input<typeof editServiceSchema>;

export function ManageService() {

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
        <button onClick={() => setIsCreateOpen(true)} className="bg-blue-600 text-white p-1.5 rounded-md">
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
  const [loadingCreate, setLoadingCreate] = useState<boolean>(false);
  const [errorCreate, setErrorCreate] = useState<string | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<ServiceFormData>({
    resolver: zodResolver(serviceFormSchema),
    defaultValues: {
      name: "",
      price: "",
      duration: "",
      file: undefined,
    },
  });

  const file = watch("file");

  useEffect(() => {
    if (!isOpen) return;
    setErrorUsers(null);
    setLoadingUsers(true);
  }, [isOpen]);

  useEffect(() => {
    if (file && file instanceof FileList && file.length > 0) {
      const objectUrl = URL.createObjectURL(file[0]);
      setPreviewUrl(objectUrl);
      return () => URL.revokeObjectURL(objectUrl);
    }
  }, [file]);

  const onValidSubmit = async (data: ServiceFormData) => {
    setLoadingCreate(true);
    setErrorCreate(null);
    try {
      const file = data.file[0];
      const formData = new FormData();
      formData.append("name", data.name.trim());
      formData.append("price", data.price);
      formData.append("duration", data.duration);
      formData.append("file", file);

      const res = await axios.post<{ status: string }>(
        `/barberbooking/tenants/${tenantId}/branch/${branchId}/services`,
        formData,
        { headers: { "Content-Type": "multipart/form-data" } }
      );

      if (res.data.status !== "success") throw new Error(res.data.status);

      if (res.data.status === "success") {
        makeToast({
          message: "เพิ่มข้อมูลสำเร็จแล้ว!",
          variant: "success",
        });
      } else {
        makeToast({
          message: "เกิดข้อผิดพลาด: " + ("ไม่ทราบสาเหตุ"),
          variant: "error",
        });
      }

      onCreate();
      onClose();
      reset();
    } catch (err: any) {
      setErrorCreate(err.response?.data?.message || err.message || "เกิดข้อผิดพลาดในการเพิ่มข้อมูลบริการ");
      makeToast({
        message: "ไม่สามารถเพิ่มข้อมูลได้ โปรดลองอีกครั้งในภายหลัง",
        variant: "error",
      });
    } finally {
      setLoadingCreate(false);
    }
  };

  if (statusMe === "loading") return null;
  if (statusMe === "succeeded" && !me) return null;
  if (!isOpen) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="เพิ่มบริการใหม่">
      <form onSubmit={handleSubmit(onValidSubmit)}>
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">รูปภาพโปรไฟล์</label>
          <input
            type="file"
            accept="image/*"
            {...register("file")}
            className="file-input file-input-bordered w-full"
          />
          {errors.file && (
            <p className="text-red-600 text-sm mt-1">
              {(errors.file as FieldError).message}
            </p>
          )}

          {previewUrl && (
            <div className="mt-3">
              <img
                src={previewUrl}
                alt="Preview"
                className="h-20 w-24 object-cover rounded-md border object-top"
              />
            </div>
          )}
        </div>



        {/* ✅ ชื่อบริการ + ราคา */}
        <div className="flex gap-3">
          <div className="mb-4 w-full">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">ชื่อบริการ</label>
            <input
              type="text"
              {...register("name")}
              placeholder="กรุณากรอกชื่อบริการ"
              className="w-full input input-bordered"
            />

            {errors.name && <p className="text-red-600 text-sm mt-1">{errors.name.message}</p>}
          </div>

          <div className="mb-4 w-full">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">ราคา</label>
            <input
              type="number"
              {...register("price")}
              placeholder="กรุณากรอกราคา"
              className="w-full input input-bordered"
            />
            {errors.price && <p className="text-red-600 text-sm mt-1">{errors.price.message}</p>}
          </div>
        </div>

        {/* ✅ ระยะเวลา */}
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">ระยะเวลาโดยประมาณ (นาที)</label>
          <input
            type="number"
            {...register("duration")}
            placeholder="กรุณากรอกระยะเวลา"
            className="w-full input input-bordered"
          />
          {errors.duration && <p className="text-red-600 text-sm mt-1">{errors.duration.message}</p>}
        </div>

        {errorCreate && <p className="text-sm text-red-600 mb-4">{errorCreate}</p>}

        <div className="flex justify-end space-x-2 mt-6">
          <button type="button" onClick={onClose} className="btn btn-ghost" disabled={loadingCreate}>
            ยกเลิก
          </button>
          <button type="submit" className="btn btn-primary" disabled={loadingCreate}>
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
    watch,
    formState: { errors, isSubmitting },
  } = useForm<EditServiceFormData>({
    resolver: zodResolver(editServiceSchema),
    defaultValues: {
      name: "",
      price: 0,
      duration: 0,
      file: undefined,
    },
  });



  useEffect(() => {
    if (isOpen && service) {
      reset({
        name: service.name ?? "",
        price: service.price ?? 0,
        duration: service.duration ?? 0,
        file: undefined,
      });
    }
  }, [isOpen, service, reset]);




  const onSubmit = async (data: EditServiceFormData) => {
    const formData = new FormData();
    formData.append("name", data.name);
    formData.append("price", data.price.toString());
    formData.append("duration", data.duration.toString());

    if (data.file instanceof FileList && data.file.length > 0) {
      formData.append("file", data.file[0]);
    }

    const res = await axios.put(
      `/barberbooking/tenants/${tenantId}/branch/${branchId}/services/${service?.id}`,
      formData,
      { headers: { "Content-Type": "multipart/form-data" } }
    );

    if (res.data.status !== "success") throw new Error(res.data.status);

    if (res.data.status === "success") {
      makeToast({
        message: "แก้ไขข้อมูลสำเร็จแล้ว!",
        variant: "success",
      });
    } else {
      makeToast({
        message: "เกิดข้อผิดพลาด: " + (res.data.message || "ไม่ทราบสาเหตุ"),
        variant: "error",
      });
    }

    onEdit(res.data.data);
    onClose();
  };



  const file = watch("file");
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  useEffect(() => {
    if (file && file instanceof FileList && file.length > 0) {
      const objectUrl = URL.createObjectURL(file[0]);
      setPreviewUrl(objectUrl);
      return () => URL.revokeObjectURL(objectUrl);
    } else if (service?.Img_path && service?.Img_name) {
      setPreviewUrl(`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${service.Img_path}/${service.Img_name}`);
    }
  }, [file, service]);
  if (statusMe === "loading") return null;
  if (statusMe === "succeeded" && !me) return null;
  if (!isOpen || !service) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="แก้ไขข้อมูลบริการ">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        {/* ชื่อบริการ */}
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            รูปภาพประกอบบริการ
          </label>
          <input
            type="file"
            accept="image/*"
            {...register("file")}
            className="file-input file-input-bordered w-full"
          />
          {previewUrl && (
            <div className="mt-3">
              <img
                src={previewUrl}
                alt="Preview"
                className="h-20 w-24 object-cover rounded-md border object-top"
              />
            </div>
          )}
        </div>

        <div>
          <label htmlFor="name" className="block text-sm font-medium text-gray-700 dark:text-gray-200">
            ชื่อบริการ
          </label>
          <input
            id="name"
            type="text"
            placeholder="กรุณากรอกชื่อบริการ"
            {...register("name")}
            className={`mt-1 w-full input input-bordered rounded-md shadow-sm transition-all duration-150 focus:outline-none focus:ring-2 ${errors.name ? "border-red-500 ring-red-300" : "focus:ring-primary"}`}
          />
          {errors.name && (
            <p className="text-red-600 text-xs mt-1">{errors.name.message}</p>
          )}
        </div>

        {/* ราคา */}
        <div>
          <label htmlFor="price" className="block text-sm font-medium text-gray-700 dark:text-gray-200">
            ราคา (บาท)
          </label>
          <input
            id="price"
            type="number"
            placeholder="กรุณากรอกราคา"
            {...register("price", { valueAsNumber: true })}
            className={`mt-1 w-full input input-bordered rounded-md shadow-sm transition-all duration-150 focus:outline-none focus:ring-2 ${errors.price ? "border-red-500 ring-red-300" : "focus:ring-primary"}`}
          />
          {errors.price && (
            <p className="text-red-600 text-xs mt-1">{errors.price.message}</p>
          )}
        </div>

        {/* ระยะเวลา */}
        <div>
          <label htmlFor="duration" className="block text-sm font-medium text-gray-700 dark:text-gray-200">
            ระยะเวลาโดยประมาณ (นาที)
          </label>
          <input
            id="duration"
            type="number"
            placeholder="กรุณากรอกระยะเวลา"
            {...register("duration", { valueAsNumber: true })}
            className={`mt-1 w-full input input-bordered rounded-md shadow-sm transition-all duration-150 focus:outline-none focus:ring-2 ${errors.duration ? "border-red-500 ring-red-300" : "focus:ring-primary"}`}
          />
          {errors.duration && (
            <p className="text-red-600 text-xs mt-1">{errors.duration.message}</p>
          )}
        </div>

        {/* ปุ่ม */}
        <div className="flex justify-end gap-3 pt-4">
          <button
            type="button"
            onClick={onClose}
            disabled={isSubmitting}
            className="btn btn-outline"
          >
            ยกเลิก
          </button>
          <button
            type="submit"
            disabled={isSubmitting}
            className={`btn btn-primary ${isSubmitting ? "opacity-50 cursor-not-allowed" : ""}`}
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

        makeToast({
          message: "ลบข้อมูลสำเร็จแล้ว!",
          variant: "success",
        });
        onDelete()
        onClose()

      } else {
        makeToast({
          message: "เกิดข้อผิดพลาด: " + ("ไม่ทราบสาเหตุ"),
          variant: "error",
        });
      }
    } catch (err: any) {
      makeToast({
        message: "ไม่สามารถลบข้อมูลได้ โปรดลองอีกครั้งในภายหลัง",
        variant: "error",
      });
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
