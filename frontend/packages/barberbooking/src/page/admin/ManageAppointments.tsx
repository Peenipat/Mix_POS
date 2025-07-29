// src/pages/admin/ManageAppointments.tsx

import React, { useEffect, useState, useRef } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";
import { useAppSelector } from "../../store/hook";
import { getAppointmentsByBranch } from "../../api/appointment"; // <-- import api ใหม่
import type { AppointmentBrief } from "../../api/appointment";   // <-- import type ใหม่

export function ManageAppointments() {
  const me = useAppSelector((state) => state.auth.me);
  const branchId = me?.branch_id;

  const [appointments, setAppointments] = useState<AppointmentBrief[]>([]);
  const [loadingAppts, setLoadingAppts] = useState<boolean>(false);
  const [errorAppts, setErrorAppts] = useState<string | null>(null);
  const didFetchAppts = useRef(false);

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

  const columns: Column<AppointmentBrief>[] = [
    {
      header: "#",
      accessor: (_row, rowIndex) => rowIndex + 1,
    },
    {
      header: "Customer",
      accessor: (row) => row.customer.name,
    },
    {
      header: "Phone",
      accessor: (row) => row.customer.phone,
    },
    {
      header: "Barber",
      accessor: (row) => row.barber.username,
    },
    {
      header: "Service",
      accessor: (row) => row.service.name,
    },
    {
      header: "Date",
      accessor: (row) => row.date,
    },
    {
      header: "Time",
      accessor: (row) => `${row.start} - ${row.end}`,
    },
    {
      header: "Status",
      accessor: (row) => row.status,
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

  return (
    <div>
      <h2 className="text-xl mb-4">
        Appointments for Branch {branchId}
      </h2>
      <DataTable<AppointmentBrief>
        data={appointments}
        columns={columns}
        onRowClick={(r) => console.log("row clicked", r)}
        actions={[viewAction, cancelAction]}
        showEdit={false}
        showDelete={false}
      />
    </div>
  );
}
