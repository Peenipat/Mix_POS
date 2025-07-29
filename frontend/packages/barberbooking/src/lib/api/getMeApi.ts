// src/lib/api/getMeApi.ts
import axios from '../axios'; 

export interface Me {
  id: number;
  username: string;
  email: string;
  role: string; 
  branch_id?: number;
  tenant_ids: number[];
}

export async function getMeApi(): Promise<Me> {
  const resp = await axios.get<{ data: Me }>('/core/user/me');
  return resp.data.data;
}
