import api from "../lib/axios";
import qs from "qs";
import { format } from "date-fns";
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

export type GetAppointmentsOptions = {
  start?: string;
  end?: string;
  filter?: "" | "week" | "month";
  excludeStatus?: string[];
  search?: string;
  status?: string[];
  barberId?: number;
  serviceId?: number;
  createdStart?: string;
  createdEnd?: string;
  page?: number;
  limit?: number;
};

export async function getAppointments(
  branchId: number,
  tenantId: number,
  options: GetAppointmentsOptions
): Promise<{ data: AppointmentBrief[]; total: number; page: number; limit: number }> {
  const {
    start,
    end,
    filter,
    excludeStatus,
    search,
    status,
    barberId,
    serviceId,
    createdStart,
    createdEnd,
    page = 1,
    limit = 10,
  } = options;

  const params: Record<string, string> = {
    tenant_id: String(tenantId),
    page: String(page),
    limit: String(limit),
  };

  if (start) params.start = start;
  if (end) params.end = end;
  if (filter) params.filter = filter;
  if (search) params.search = search;
  if (barberId) params.barber_id = String(barberId);
  if (serviceId) params.service_id = String(serviceId);
  if (createdStart) params.created_start = createdStart;
  if (createdEnd) params.created_end = createdEnd;
  if (excludeStatus && excludeStatus.length > 0) {
    params.exclude_status = excludeStatus.join(",");
  }
  if (status && status.length > 0) {
    params.status = status.join(",");
  }

  const resp = await api.get(`/barberbooking/tenants/${tenantId}/branches/${branchId}/appointments`, {
    params,
  });

  const rawData = resp.data.data;
  const pagination = resp.data.meta?.pagination || { total: 0, page, limit };

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

  return {
    data: transformed,
    total: pagination.total,
    page: pagination.page,
    limit: pagination.limit,
  };
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