import api from "../lib/axios";

export interface CustomerDetail {
    id: number;
    tenant_id: number;
    branch_id: number;
    name: string;
    phone: string;
    email: string;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
  }


  export async function getCustomerById(
    tenantId: number,
    branchId: number,
    customerId: number
  ): Promise<CustomerDetail> {
    if (!tenantId || !branchId || !customerId) {
      throw new Error("ต้องระบุ tenantId, branchId และ customerId");
    }
  
    const resp = await api.get(`/barberbooking/tenants/${tenantId}/branch/${branchId}/customers/${customerId}`);
    
    return resp.data.data as CustomerDetail;
  }