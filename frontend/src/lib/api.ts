import api from "./axios"
import { loginResponseSchema } from "@/schemas/userSchema"
import type { loginResponse } from "@/schemas/userSchema"

export async function loginApi(
    credentials: { email: string; password: string } // ค่าที่ส่งเข้า backend
  ): Promise<loginResponse> {
    const resp = await api.post("/auth/login", credentials);
    const parsed = loginResponseSchema.safeParse(resp.data); // ตรวจสอบ type ให้ตรงตามที่กำหนด
  
    // ถ้าไม่ตรง แสดง error
    if (!parsed.success) {
      console.log("ZodError:", parsed.error.format());
      throw new Error("Invalid response format");
    }
  // return ข้อมูล ซึ่งมีรูปแบบตรงกับ loginResponse
    return parsed.data; 
  }
  

  