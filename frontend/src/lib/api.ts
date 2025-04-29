import api from "./axios"
import { loginResponseSchema } from "@/schemas/userSchema"
import type { loginResponse } from "@/schemas/userSchema"

export async function loginApi(
    creadentails :{email:string; password:string}
): Promise<loginResponse>{
    const resp = await api.post("/auth/login",creadentails)
    const parsed = loginResponseSchema.safeParse(resp.data)
    if (!parsed.success){
        console.log("ZodError : ", parsed.error.format())
    }
    return loginResponseSchema.parse(resp.data)
}
