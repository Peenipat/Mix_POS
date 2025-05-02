import { z } from "zod";

export const UserResponseSchema = z.object({
  id: z.number(), 
  username: z.string().min(1, { message: "Username is required" }),
  email: z.string().email({ message: "Invalid email address" }),
  role_id: z.number(),
  role: z.enum(["SAAS_SUPER_ADMIN", "BRANCH_ADMIN", "STAFF", "USER"]),
  createdAt: z.string().datetime({ message: "Invalid createdAt format" }).optional(),
  updatedAt: z.string().datetime({ message: "Invalid updatedAt format" }).optional(),
  deletedAt: z.string().datetime({ message: "Invalid deletedAt format "}).nullable().optional(),
});

export const loginResponseSchema = z.object({
  user: UserResponseSchema
})

export type loginResponse = z.infer<typeof loginResponseSchema>

export const EditUserFromAdmin = z.object({
    username: z.string().min(1, {message: "Username is required"}),
    email: z.string().email({ message: "Invalid email address"}),
    role: z.enum(["BRANCH_ADMIN", "STAFF", "USER"])
})
