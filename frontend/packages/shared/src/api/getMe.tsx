import  axios from "../lib/axios";
export interface Me {
    id: number;
    username: string;
    email: string;
    branch_id?: number;
    tenant_ids: number[];
  }
  
  export async function fetchMe(): Promise<Me> {
    const resp = await axios.get<{ data: Me }>("/core/user/me");
    return resp.data.data;
  }