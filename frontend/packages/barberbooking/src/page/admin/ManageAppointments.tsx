// src/pages/admin/ManageAppointments.tsx

import React, { useEffect, useState, useRef, useCallback } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";
import { useAppSelector } from "../../store/hook";
import { getAppointments } from "../../api/appointment"; // <-- import api ใหม่
import type { AppointmentBrief } from "../../api/appointment";   // <-- import type ใหม่
import dayjs from "dayjs";
import { getBarbers } from "../../api/barber";
import { Barber } from "../../types/barber";

export function ManageAppointments() {
  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);
  const branchId = me?.branch_id;
  const tenantId = me?.tenant_ids[0]
  const didFetchBarbers = useRef(false);



  const [appointments, setAppointments] = useState<AppointmentBrief[]>([]);
  const [pagination, setPagination] = useState({ page: 1, limit: 10, total: 0 });
  const [loadingAppts, setLoadingAppts] = useState(false);
  const [errorAppts, setErrorAppts] = useState<string | null>(null);

  // filters

  const [selectedStatuses, setSelectedStatuses] = useState<string[]>([]);

  const [selectedBarberId, setSelectedBarberId] = useState<number | null>(null);
  const [selectedServiceId, setSelectedServiceId] = useState<number | null>(null);
  const [createdStart, setCreatedStart] = useState<string | null>(null);
  const [createdEnd, setCreatedEnd] = useState<string | null>(null);

  const [searchInput, setSearchInput] = useState("");
  const [searchTerm, setSearchTerm] = useState("");

  const [loadingBarbers, setLoadingBarbers] = useState<boolean>(false);
  const [errorBarbers, setErrorBarbers] = useState<string | null>(null);
  const [barbers, setBarbers] = useState<Barber[]>([]);
  const loadBarbers = useCallback(async () => {
    if (!tenantId || !branchId) return;

    setLoadingBarbers(true);
    setErrorBarbers(null);

    try {
      const barbers = await getBarbers(branchId);
      setBarbers(barbers);
    } catch (err: any) {
      setErrorBarbers(err.response?.data?.message || err.message || "Failed to load barbers");
    } finally {
      setLoadingBarbers(false);
    }
  }, [tenantId, branchId]);

  useEffect(() => {
    if (!branchId || !tenantId) return;

    const loadAppointments = async () => {
      setLoadingAppts(true);
      setErrorAppts(null);
      try {
        const { data, total, page, limit } = await getAppointments(branchId!, tenantId!, {
          page: pagination.page,
          limit: pagination.limit,
          search: searchTerm,
          status: selectedStatuses,
          barberId: selectedBarberId || undefined,
          serviceId: selectedServiceId || undefined,
          createdStart: createdStart || undefined,
          createdEnd: createdEnd || undefined,
        });

        setAppointments(data ?? []);
        setPagination({ total, page, limit });
      } catch (err: any) {
        setErrorAppts(err.message || "Failed to load appointments");
      } finally {
        setLoadingAppts(false);
      }
    };

    loadAppointments();
  }, [
    branchId,
    tenantId,
    pagination.page,
    pagination.limit,
    searchTerm,
    selectedStatuses,
    selectedBarberId,
    selectedServiceId,
    createdStart,
    createdEnd,
  ]);

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
      accessor: (_row, rowIndex) =>
        (pagination.page - 1) * pagination.limit + rowIndex + 1,
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


  const handleSearch = () => {
    setSearchTerm(searchInput);
    setPagination((prev) => ({ ...prev, page: 1 }));
  };


  return (
    <div>
      <h2 className="text-xl mb-4">
        Appointments for Branch {branchId}
      </h2>
      <div className="flex flex-col sm:flex-row gap-4 mb-4 items-center">
        <input
          type="text"
          className="border border-gray-300 rounded px-3 py-2 w-full sm:w-64"
          placeholder="ค้นหาชื่อลูกค้า / เบอร์โทร "
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              handleSearch();
            }
          }}
        />

        <button
          onClick={handleSearch}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          ค้นหา
        </button>

        <select
          value={selectedStatuses[0] || ""}
          onChange={(e) =>
            setSelectedStatuses(e.target.value === "" ? [] : [e.target.value])
          }
          className="border border-gray-300 rounded px-3 py-2"
        >
          <option value="">-- สถานะทั้งหมด --</option>
          <option value="CONFIRMED">จองแล้ว</option>
          <option value="IN_SERVICE">กำลังให้บริการ</option>
          <option value="COMPLETED">เสร็จสิ้น</option>
          <option value="CANCELLED">ยกเลิกแล้ว</option>
        </select>


        <select
          value={selectedBarberId ?? ""}
          onChange={(e) => {
            const val = e.target.value;
            setSelectedBarberId(val === "" ? null : Number(val));
          }}
          className="border border-gray-300 rounded px-3 py-2"
        >
          <option value="">-- ช่างทั้งหมด --</option>
          {barbers.map((barber) => (
            <option key={barber.id} value={barber.id}>
              {barber.username}
            </option>
          ))}
        </select>
      </div>

      <DataTable<AppointmentBrief>
        data={appointments}
        columns={columns}
        actions={[viewAction, cancelAction]}
        showEdit={false}
        showDelete={false}
        page={pagination.page}
        perPage={pagination.limit}
        total={pagination.total}
        onPageChange={(p) => {
          setPagination((prev) => ({ ...prev, page: p }));
        }}
      />
    </div>
  );
}
