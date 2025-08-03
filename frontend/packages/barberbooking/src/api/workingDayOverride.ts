import api from "../lib/axios";
export type WorkingDayOverrideInput = {
  branch_id: number;
  work_date: string; // Format: "YYYY-MM-DD"
  start_time: string; // Format: "HH:mm"
  end_time: string;   // Format: "HH:mm"
  is_closed: boolean;
  reason: string
};

export type WorkingDayOverrideResponse = {
  message: string;
  data: {
    id: number;
    branch_id: number;
    work_date: string;
    start_time: string;
    end_time: string;
    is_closed: boolean;
    reason: string
    created_at: string;
    updated_at: string;
  };
};

export async function createWorkingDayOverride(
  input: WorkingDayOverrideInput
): Promise<WorkingDayOverrideResponse> {
  const resp = await api.post<WorkingDayOverrideResponse>(
    `/barberbooking/working-day-overrides`,
    input
  );
  return resp.data;
}

export type WorkingDayOverride = {
  id: number;
  branch_id: number;
  work_date: string;    // format: "YYYY-MM-DD"
  start_time: string;   // format: "HH:mm"
  end_time: string;     // format: "HH:mm"
  is_closed: boolean;
  reason: string;
  created_at: string;
  updated_at: string;
};

export async function getWorkingDayOverridesByDateRange(params: {
  tenantId: number;
  branchId: number;
  start: string; // format: "YYYY-MM-DD"
  end: string;   // format: "YYYY-MM-DD"
}): Promise<WorkingDayOverride[]> {
  const { tenantId, branchId, start, end } = params;

  const resp = await api.get<WorkingDayOverride[]>(
    `/barberbooking/tenants/${tenantId}/branches/${branchId}/working-day-overrides/date`,
    {
      params: {
        start,
        end,
      },
    }
  );

  return resp.data;
}
