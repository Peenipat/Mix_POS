import api from "../lib/axios";
import qs from "qs";
// 🎯 Type สำหรับ Customer
type CustomerInfo = {
  name: string;
  phone: string;
};

export type BookAppointmentPayload = {
  barber_id: number;
  branch_id: number;
  service_id: number;
  start_time: string;
  notes?: string;
  customer_id: number;
  customer?: CustomerInfo;
};


export type BookAppointmentResponse = {
  id: number;
  start_time: string;
  end_time: string;
  status: string;
  customer_id: number;
  service_id: number;
  notes?: string;
  created_at: string;
};

export async function bookAppointment(
  tenantId: number,
  payload: BookAppointmentPayload
): Promise<BookAppointmentResponse> {
  const resp = await api.post(`/barberbooking/tenants/${tenantId}/appointments`, payload);
  return resp.data;
}

export type AppointmentBrief = {
  id: number;
  branch_id: number;
  service_id: number;
  service: {
    name: string;
    description: string;
    duration: number;
    price: number;
  };
  barber_id: number;
  barber: {
    username: string;
  };
  customer_id: number;
  customer: {
    name: string;
    phone: string;
  };
  date: string;
  start: string;
  end: string;
  status: string;
};

import { format } from "date-fns";
export async function getAppointmentsByBranch(
  branchId: number,
  start?: string,
  end?: string,
  filter?: "" | "week" | "month" | null,
  excludeStatus?: string[] 
): Promise<AppointmentBrief[]> {
  const params: Record<string, string> = {};

  if (start) params.start = start;
  if (end) params.end = end;
  if (filter) params.filter = filter;

  if (excludeStatus && excludeStatus.length > 0) {
    params.exclude_status = excludeStatus.join(","); 
  }

  const resp = await api.get(`/barberbooking/branches/${branchId}/appointments`, {
    params,
  });

  const rawData = resp.data.data;

  const transformed: AppointmentBrief[] = rawData.map((a: any) => {
    const startDate = new Date(a.start_time);
    const endDate = new Date(a.end_time);

    return {
      id: a.id,
      branch_id: a.branch_id,
      service_id: a.service_id,
      service: a.service,
      barber_id: a.barber_id,
      barber: a.barber,
      customer_id: a.customer_id,
      customer: a.customer,
      status: a.status,
      date: format(startDate, "yyyy-MM-dd"),
      start: format(startDate, "HH:mm"),
      end: format(endDate, "HH:mm"),
    };
  });

  return transformed;
}

export type GetAppointmentsByBarberQuery = {
  start?: string;                // "yyyy-MM-dd"
  end?: string;                  // "yyyy-MM-dd"
  status?: string[];            // ["COMPLETED", "CONFIRMED"]
  mode?: "today" | "week" | "past";
};

export type GetAppointmentsByBarberResponse = {
  status: string;
  data: AppointmentBrief[];
};

export async function getAppointmentsByBarber(
  barberId: number,
  query?: GetAppointmentsByBarberQuery
): Promise<GetAppointmentsByBarberResponse> {
  const queryString = qs.stringify(
    {
      ...query,
      status: query?.status?.join(","),
    },
    { skipNulls: true }
  );

  const resp = await api.get(`/barberbooking/barbers/${barberId}/appointments?${queryString}`);

  const transformed = (resp.data.data ?? []).map((item: any): AppointmentBrief => ({
    ...item,
    start: item.start_time,
    end: item.end_time,
  }));

  return {
    status: resp.data.status,
    data: transformed,
  };
}

export async function getAppointmentsByPhone(phone: string): Promise<AppointmentBrief[]> {
  if (!phone) throw new Error("Phone number is required");

  const resp = await api.get("/barberbooking/appointments/by-phone", {
    params: { phone },
  });

  const rawData = resp.data.data;

  const transformed: AppointmentBrief[] = rawData.map((a: any) => {
    const startDate = new Date(a.start_time);
    const endDate = new Date(a.end_time);

    return {
      id: a.id,
      branch_id: a.branch_id,
      service_id: a.service_id,
      service: a.service,
      barber_id: a.barber_id,
      barber: a.barber,
      customer_id: a.customer_id,
      customer: a.customer,
      status: a.status,
      date: format(startDate, "yyyy-MM-dd"),
      start: format(startDate, "HH:mm"),
      end: format(endDate, "HH:mm"),
    };
  });

  return transformed;
}


export async function updateAppointmentStatus(
  tenantId: number,
  appointmentId: number,
  status: string,
  userId?: number
) {
  const payload: any = {
    status,
  };

  if (userId) {
    payload.user_id = userId;
  }

  const resp = await api.put(
    `/barberbooking/tenants/${tenantId}/appointments/${appointmentId}`,
    payload
  );

  return resp.data.data; 
}