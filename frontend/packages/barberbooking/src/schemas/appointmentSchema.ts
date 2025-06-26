import { z } from "zod";

export const appointmentSchema = z.object({
    cusName:        z.string().nonempty("กรุณากรอกชื่อผู้จอง"),
    phoneNumber:    z.string().regex(/^\d{10}$/, "เบอร์โทรศัพท์ต้องมี 10 หลัก"),
    barberId:       z.number().gt(0, "กรุณาเลือกช่างตัดผม"),
    serviceId:      z.number().gt(0, "กรุณาเลือกบริการ"),
    date:           z.string().nonempty("กรุณาเลือกวันที่ต้องการจอง"),
    time:           z.string().nonempty("กรุณาเลือกเวลาที่ต้องการจอง"),
    note:           z.string().optional()
  });

export type appointmentForm = z.infer<typeof appointmentSchema>;