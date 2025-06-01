// src/lib/api/getMeApi.ts
import axios from '../axios'; // ควรเป็น instance เดิมของคุณ ที่ตั้ง withCredentials: true แล้ว

export interface Me {
  id: number;
  username: string;
  email: string;
  role: string; 
  branch_id?: number;
  tenant_ids: number[];
}

export async function getMeApi(): Promise<Me> {
  // GET /core/user/me → backend จะตอบ { status: "success", data: { ... } }
  // หรือในบางการ implement อาจตอบแค่ data โดยตรง เช่น { data: { id, username, ... } }
  const resp = await axios.get<{ data: Me }>('/core/user/me');
  return resp.data.data;
}
