// src/pages/admin/ManageAppointments.tsx

import React, { useEffect, useState, useRef } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";
import { useAppSelector } from "../../store/hook";
import { getAppointmentsByBranch } from "../../api/appointment"; // <-- import api ใหม่
import type { AppointmentBrief } from "../../api/appointment";   // <-- import type ใหม่
import dayjs from "dayjs";

export function ManageAppointments() {
  const me = useAppSelector((state) => state.auth.me);
  const branchId = me?.branch_id;

  const [appointments, setAppointments] = useState<AppointmentBrief[]>([]);
  const [loadingAppts, setLoadingAppts] = useState<boolean>(false);
  const [errorAppts, setErrorAppts] = useState<string | null>(null);
  const didFetchAppts = useRef(false);

  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState<string | null>(null);
  const [currentPage,setCurrentPage] = useState(1)

  useEffect(() => {
    if (!branchId || didFetchAppts.current) return;
    didFetchAppts.current = true;

    const loadAppointments = async () => {
      setLoadingAppts(true);
      setErrorAppts(null);
      try {
        const data = await getAppointmentsByBranch(branchId);
        setAppointments(data);
      } catch (err: any) {
        setErrorAppts(err.message || "Failed to load appointments");
      } finally {
        setLoadingAppts(false);
      }
    };

    loadAppointments();
  }, [branchId]);

  if (!branchId) {
    return <p className="text-red-500">Cannot determine branch information.</p>;
  }
  if (loadingAppts) {
    return <p>Loading appointments…</p>;
  }
  if (errorAppts) {
    return <p className="text-red-500">Error loading appointments: {errorAppts}</p>;
  }
  function mapStatusToThai(status: string): string {
    switch (status) {
      case 'CONFIRMED':
        return 'จองแล้ว';
      case 'CANCELLED':
        return 'ยกเลิกแล้ว';
      case 'COMPLETED':
        return 'เสร็จสิ้น';
      case 'IN_SERVICE':
        return 'กำลังให้บริการ';
      default:
        return 'ไม่ทราบสถานะ';
    }
  }


  const columns: Column<AppointmentBrief>[] = [
    {
      header: "#",
      accessor: (_row, rowIndex) => rowIndex + 1,
    },
    {
      header: "ชื่อลูกค้า",
      accessor: (row) => row.customer.name,
    },
    {
      header: "เบอร์โทร",
      accessor: (row) => row.customer.phone,
    },
    {
      header: "ช่างที่เลือก",
      accessor: (row) => row.barber.username,
    },
    {
      header: "บริการที่เลือก",
      accessor: (row) => row.service.name,
    },
    {
      header: "วันที่จอง",
      accessor: (row) => dayjs(row.date).format('DD/MM/YYYY'),
    },
    {
      header: "ช่วงเวลาที่จอง",
      accessor: (row) => `${row.start} - ${row.end}`,
    },
    {
      header: "สถาะนะการจอง",
      accessor: (row) => mapStatusToThai(row.status),
    },
  ];

  const viewAction: Action<AppointmentBrief> = {
    label: "View",
    onClick: (row) => console.log("view appointment", row),
    className: "text-green-600",
  };

  const cancelAction: Action<AppointmentBrief> = {
    label: "Cancel",
    onClick: (row) => console.log("cancel appointment", row),
    className: "text-red-600",
  };


  const filteredData = appointments.filter((row) => {
    const matchesSearch =
      row.customer.name.includes(searchTerm) ||
      row.customer.phone.includes(searchTerm) ||
      row.barber.username.includes(searchTerm);

    const matchesStatus = !statusFilter || row.status === statusFilter;

    return matchesSearch && matchesStatus;
  });


  return (
    <div>
      <h2 className="text-xl mb-4">
        Appointments for Branch {branchId}
      </h2>
      <div className="flex flex-col sm:flex-row gap-4 mb-4 items-center">
        <input
          type="text"
          className="border border-gray-300 rounded px-3 py-2 w-full sm:w-64"
          placeholder="ค้นหาชื่อลูกค้า / เบอร์โทร / ช่าง"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />

        <select
          value={statusFilter || ""}
          onChange={(e) =>
            setStatusFilter(e.target.value === "" ? null : e.target.value)
          }
          className="border border-gray-300 rounded px-3 py-2"
        >
          <option value="">ทั้งหมด</option>
          <option value="CONFIRMED">จองแล้ว</option>
          <option value="IN_SERVICE">กำลังให้บริการ</option>
          <option value="COMPLETED">เสร็จสิ้น</option>
          <option value="CANCELLED">ยกเลิกแล้ว</option>
        </select>
      </div>

      <DataTable<AppointmentBrief>
        data={appointments}
        columns={columns}
        actions={[viewAction, cancelAction]}
        showEdit={false}
        showDelete={false}
        page={currentPage}
        perPage={10}
        total={appointments.length}
        onPageChange={(p) => setCurrentPage(p)}
      />
    </div>
  );
}
