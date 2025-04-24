import { z } from "zod"

export const loginSchema = z.object({
    email: z.string().email({message:"Invaild email"}),
    password: z.string().min(6, { message:"Passwrod too short"}),
})

export const registerSchema = loginSchema.extend({})