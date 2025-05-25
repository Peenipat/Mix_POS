import { z } from "zod";

export const UserResponseSchema = z.object({
  id: z.number(), 
  username: z.string().min(1, { message: "Username is required" }),
  email: z.string().email({ message: "Invalid email address" }),
  role: z.enum(["SAAS_SUPER_ADMIN", "BRANCH_ADMIN", "STAFF", "USER","TENANT","ASSISTANT_MANAGER","TENANT_ADMIN"]),
  createdAt: z.string().datetime({ message: "Invalid createdAt format" }).optional(),
  updatedAt: z.string().datetime({ message: "Invalid updatedAt format" }).optional(),
  deletedAt: z.string().datetime({ message: "Invalid deletedAt format "}).nullable().optional(),
});
export const UsersSchema = z.array(UserResponseSchema);
export type User = z.infer<typeof UserResponseSchema>;

export const CreateUserSchema = z.object({
  username: z.string().min(3, 'Username ต้องมีอย่างน้อย 3 ตัวอักษร'),
  email: z.string().email('รูปแบบอีเมลไม่ถูกต้อง'),
  role: z.enum(['SAAS_SUPER_ADMIN', 'BRANCH_ADMIN', 'STAFF', 'USER', 'TENANT'], {
    errorMap: () => ({ message: 'กรุณาเลือก Role ให้ถูกต้อง' }),
  }),
});

export type CreateUserForm = z.infer<typeof CreateUserSchema>;


export const loginResponseSchema = z.object({
  user: UserResponseSchema
})

export type loginResponse = z.infer<typeof loginResponseSchema>

export const EditUserFromAdmin = z.object({
    username: z.string().min(1, {message: "Username is required"}),
    email: z.string().email({ message: "Invalid email address"}),
    role: z.enum(["BRANCH_ADMIN", "STAFF", "USER"])
})
