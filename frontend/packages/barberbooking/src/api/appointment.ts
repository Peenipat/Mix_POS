import api from "../lib/axios"; 

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
  end?: string
): Promise<AppointmentBrief[]> {
  const params: Record<string, string> = {};
  if (start) params.start = start;
  if (end) params.end = end;

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

