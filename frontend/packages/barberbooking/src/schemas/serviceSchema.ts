import { z } from "zod";

export const serviceFormSchema = z.object({
    file: z
      .any()
      .refine(
        (val) => val instanceof FileList && val.length > 0 && val[0] instanceof File,
        { message: "กรุณาอัปโหลดรูปภาพ" }
      ),
    name: z.string().min(1, "กรุณากรอกชื่อบริการ"),
    price: z
      .string()
      .min(1, "กรุณากรอกราคา")
      .refine((val) => !isNaN(Number(val)), {
        message: "ราคาต้องเป็นตัวเลข",
      }),
    duration: z
      .string()
      .min(1, "กรุณากรอกระยะเวลา")
      .refine((val) => !isNaN(Number(val)), {
        message: "ระยะเวลาต้องเป็นตัวเลข",
      }),
  });
  

export type ServiceFormData = z.infer<typeof serviceFormSchema>;
