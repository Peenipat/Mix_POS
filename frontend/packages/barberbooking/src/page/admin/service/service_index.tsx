// src/page/admin/ManageService.tsx
import React, { useEffect, useState, useRef, FormEvent } from "react";
import { DataTable } from "../../../components/DataTable";
import { useForm } from "react-hook-form";
import type { Action, Column } from "../../../components/DataTable";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAppSelector } from "../../../store/hook";
import axios from "../../../lib/axios";
import Modal from "@object/shared/components/Modal";
import { TableCellsIcon } from "../../../components/icons/TableCellsIcon";
import { GridViewIcon } from "../../../components/icons/GridViewIcon"
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
import { CardViewIcon } from "../../../components/icons/CardViewIcon";
export const editServiceSchema = z.object({
  name: z.string().min(1, "กรุณากรอกชื่อบริการ").max(100, "ชื่อผู้ใช้ต้องไม่เกิน 100 ตัวอักษร"),
  price: z.number().nullable(),
  duration: z.number().nullable()
});

export type EditServiceFormData = z.infer<typeof editServiceSchema>;

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


      {/* <CreateServiceModal
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
      /> */}
    </div>
  );
}
