import { useEffect, useState } from "react";
import { DataTable } from "../../components/DataTable";
import type { Action, Column } from "../../components/DataTable";
import type { Barber } from "../../types/barber";
import { fetchMe } from "../../types/getMe";
import type { Me } from "../../types/getMe";
import axios from '../../lib/axios';

export function ManageBarber() {
  // 1) state
  const [me, setMe] = useState<Me | null>(null);
  const [barbers, setBarbers] = useState<Barber[]>([]);
  const [loadingMe, setLoadingMe] = useState<boolean>(true);
  const [loadingBarbers, setLoadingBarbers] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // 2) load the user
  useEffect(() => {
    fetchMe()
      .then(user => setMe(user))
      .catch(e => setError(e.response?.data?.message || e.message))
      .finally(() => setLoadingMe(false));
  }, []);

  const tenantId = me?.tenant_ids[0];
  const branchId = me?.branch_id;

  // 3) once we have tenant/branch, fetch barbers
  useEffect(() => {
    if (!tenantId || !branchId) return;

    const loadBarbers = async () => {
      setLoadingBarbers(true);
      try {
        const res = await axios.get<{status: string;data: Barber[];}>(`/barberbooking/tenants/${tenantId}/barbers/branches/${branchId}/barbers`,)
        if (res.data.status !== 'success') throw new Error(res.data.status);
        setBarbers(res.data.data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoadingBarbers(false);
      }
    };

    loadBarbers();
  }, [tenantId, branchId]);

  // 4) table set-up
  const columns: Column<Barber>[] = [
    { header: "ID",   accessor: "id" },
    { header: "Username", accessor: "username" },
    { header: "Email",accessor: "email" },
    { header: "Phone",accessor: "phone_number" },
  ];
  const editAction: Action<Barber> = {
    label: "Edit",
    onClick: row => console.log("edit", row),
    className: "text-blue-600",
  };
  const deleteAction: Action<Barber> = {
    label: "Delete",
    onClick: row => console.log("delete", row),
    className: "text-red-600",
  };

  // 5) render loading / error states
  if (loadingMe) return <p>Loading user…</p>;
  if (error)     return <p className="text-red-500">Error: {error}</p>;
  if (loadingBarbers) return <p>Loading barbers…</p>;

  return (
    <div>
      <h2 className="text-xl mb-4">
        Barbers in Branch {branchId}
      </h2>
      <DataTable<Barber>
        data={barbers}
        columns={columns}
        onRowClick={r => console.log("row clicked", r)}
        actions={[editAction, deleteAction]}
        showEdit={false}
        showDelete={false}
      />
    </div>
  );
}
