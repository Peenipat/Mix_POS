import * as z from "zod";

export const editWorkingHourSchema = z.object({
  start_time: z.string().regex(/^([01]\d|2[0-3]):([0-5]\d)$/, {
    message: "กรุณากรอกเวลาเริ่มต้นให้ถูกต้อง เช่น 09:00",
  }),
  end_time: z.string().regex(/^([01]\d|2[0-3]):([0-5]\d)$/, {
    message: "กรุณากรอกเวลาสิ้นสุดให้ถูกต้อง เช่น 17:00",
  }),
});

export type EditWorkingHourFormData = z.infer<typeof editWorkingHourSchema>;
