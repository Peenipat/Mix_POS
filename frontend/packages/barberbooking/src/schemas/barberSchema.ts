import { z } from "zod";
export const editBarberSchema = z.object({
    username: z
      .string()
      .min(1, "กรุณากรอกชื่อผู้ใช้")
      .max(100, "ชื่อผู้ใช้ต้องไม่เกิน 100 ตัวอักษร"),
    email: z
      .string()
      .email("รูปแบบอีเมล์ไม่ถูกต้อง"),
    phone_number: z
      .string()
      .min(1, "กรุณากรอกเบอร์โทร")
      .regex(/^\d+$/, "เบอร์โทรต้องเป็นตัวเลขเท่านั้น")
      .min(10, "เบอร์โทร 10 หลัก")
      .max(10, "เบอร์โทร 10 หลัก"),
    img_path:z.string().optional(),
    img_name:z.string().optional(),
    branch_id:z.number(),
    description: z.string().optional(),
    roleUser: z.string().optional(),
    profilePicture: z.any().optional(),
  });
  
  export type EditBarberFormData = z.infer<typeof editBarberSchema>;