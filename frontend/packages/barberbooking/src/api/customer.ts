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

export interface CustomerList {
  id: number;
  name: string;
  phone: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface CustomerListResponse {
  status: string;
  message: string;
  data: CustomerList[];
  meta: {
    pagination: {
      total: number;
      page: number;
      limit: number;
    };
  };
}

export interface LoadCustomersParams {
  page?: number;
  limit?: number;
  name?: string;
  phone?: string;
  sortBy?: "created_at" | "updated_at";
  sortOrder?: "asc" | "desc";
}

export const loadCustomers = async (
  tenantId: number,
  branchId: number,
  setCustomers: (data: CustomerList[]) => void,
  setPagination: (p: { total: number; page: number; limit: number }) => void,
  setLoadingCustomers: (loading: boolean) => void,
  setErrorCustomers: (err: string | null) => void,
  params: LoadCustomersParams = {}
) => {
  setLoadingCustomers(true);
  setErrorCustomers(null);

  try {
    const query = new URLSearchParams();
    if (params.page) query.append("page", String(params.page));
    if (params.limit) query.append("limit", String(params.limit));
    if (params.name) query.append("name", params.name);
    if (params.phone) query.append("phone", params.phone);
    if (params.sortBy) query.append("sortBy", params.sortBy);
    if (params.sortOrder) query.append("sortOrder", params.sortOrder);

    const res = await api.get<CustomerListResponse>(
      `/barberbooking/tenants/${tenantId}/branch/${branchId}/customers?${query.toString()}`
    );

    if (res.data.status !== "success") {
      throw new Error(res.data.message);
    }

    setCustomers(res.data.data);
    const { total, page, limit } = res.data.meta.pagination;
    setPagination({ total, page, limit });
  } catch (err: any) {
    setErrorCustomers(
      err.response?.data?.message || err.message || "Failed to load customers"
    );
  } finally {
    setLoadingCustomers(false);
  }
};


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