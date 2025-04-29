// import { User } from "@/page/admin/super_admin/ManageUsers";
import api from "./axios"
import { loginResponseSchema } from "@/schemas/userSchema"
import type { loginResponse } from "@/schemas/userSchema"

export async function loginApi(
    credentials: { email: string; password: string }
  ): Promise<loginResponse> {
    const resp = await api.post("/auth/login", credentials);
    const parsed = loginResponseSchema.safeParse(resp.data);
  
    if (!parsed.success) {
      console.log("ZodError:", parsed.error.format());
      throw new Error("Invalid response format");
    }
  
    return parsed.data; 
  }
  
//   export async function fetchMe(): Promise<User> {
//     const res = await api.get("/auth/me"); // token จะถูกแนบจาก cookie
//     return res.data.user;
//   }
  