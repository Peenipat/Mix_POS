import api from "../lib/axios"; 

export type AppointmentLockPayload = {
    tenant_id:number;
    branch_id:number;
    barber_id: number;
    customer_id: number;
    start_time: string; // ISO format
    end_time: string;
  };
  
  export type AppointmentLockResponse = {
    id: number;
    barber_id: number;
    customer_id: number;
    start_time: string;
    end_time: string;
    expires_at: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
  };
  
  export async function createAppointmentLock(
    tenantId: number,
    branchId: number,
    payload: AppointmentLockPayload
  ): Promise<AppointmentLockResponse> {
    console.log(payload)
    const resp = await api.post(
      `/barberbooking/tenants/${tenantId}/branches/${branchId}/appointments-lock`,
      payload
    );
    return resp.data;
  }

  export async function releaseAppointmentLock(
    tenantId: number,
    branchId: number,
    lockId: number
  ): Promise<void> {
    await api.delete(
      `/barberbooking/tenants/${tenantId}/branches/${branchId}/appointments-lock/${lockId}`
    );
  }

  
  export type AppointmentLock = {
    id: number;
    barber_id: number;
    customer_id: number;
    start_time: string;
    end_time: string;
    expires_at: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
  };
  
  export async function getAppointmentLocks(
    tenantId: number,
    branchId: number,
    barberId: number,
    date: string // e.g., "2025-07-13"
  ): Promise<AppointmentLock[]> {
    const resp = await api.get(
      `/barberbooking/tenants/${tenantId}/branches/${branchId}/appointments-lock`,
      {
        params: {
          barber_id: barberId,
          branch_id: branchId,
          date: date,
        },
      }
    );
    return resp.data;
  }

  
  