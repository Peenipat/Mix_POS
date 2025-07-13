import api from "../lib/axios"; // เปลี่ยน path ถ้าอยู่คนละที่

// 🎯 Type สำหรับ Customer
type CustomerInfo = {
  name: string;
  phone: string;
};

// 🎯 Payload รวม (ใช้ได้ทั้ง guest และ member)
export type BookAppointmentPayload = {
  barber_id: number;
  branch_id: number;
  service_id: number;
  start_time: string;
  notes?: string;
  customer_id: number;         // 👈 ส่ง 0 ถ้าเป็น guest
  customer?: CustomerInfo;     // 👈 ต้องมีถ้า customer_id = 0
};

// 🎯 Optional: Response DTO (แก้ตามของ backend)
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
