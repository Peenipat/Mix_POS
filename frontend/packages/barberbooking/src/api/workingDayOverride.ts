import api from "../lib/axios";
export type WorkingDayOverrideInput = {
    branch_id: number;
    work_date: string; // Format: "YYYY-MM-DD"
    start_time: string; // Format: "HH:mm"
    end_time: string;   // Format: "HH:mm"
    is_closed: boolean;
    reason:string
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
      reason:string
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
  