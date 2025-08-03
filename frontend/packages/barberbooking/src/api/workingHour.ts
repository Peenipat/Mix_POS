import api from "../lib/axios";
export type WorkingHour = {
    id: number;
    branch_id: number;
    tenantid: number;
    week_day: number;
    start_time: string;
    end_time: string;
    is_closed: boolean;
};

  export async function getWorkingHours(params: {
    tenantId: number;
    branchId: number;
  }): Promise<WorkingHour[]> {
    const { tenantId, branchId } = params;
  
    const resp = await api.get<{
      status: string;
      message: string;
      data: WorkingHour[];
    }>(`/barberbooking/tenants/${tenantId}/workinghour/branches/${branchId}`);
  
    return resp.data.data;
  }