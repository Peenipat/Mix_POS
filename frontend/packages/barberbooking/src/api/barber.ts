import api from "../lib/axios";

export interface BarberDetail {
  id: number;
  tenant_id: number;
  branch_id: number;
  user_id: number;
  description: string;
  role_user: string;
  user: {
    id: number;
    username: string;
    email: string;
    phone_number: string;
    role_id: number;
    Img_path: string;
    Img_name: string;
  };

}

export async function getBarberById(barberId: number): Promise<BarberDetail> {
  if (!barberId) throw new Error("barberId is required");

  const resp = await api.get(`/barberbooking/barbers/${barberId}`);

  return resp.data.data as BarberDetail;
}

export async function updateBarber(tenantId: number, barberId: number, formData: FormData): Promise<any> {
  if (!tenantId || !barberId) throw new Error("ต้องระบุ tenantId และ barberId");

  const resp = await api.put(`/barberbooking/tenants/${tenantId}/barbers/${barberId}/update-barber`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });

  return resp.data;
}