import { useState, useMemo } from "react";
import { Calendar, dateFnsLocalizer, Event, SlotInfo } from "react-big-calendar";
import { format, parse, startOfWeek, getDay } from "date-fns";
import { format as formatDate } from "date-fns";
import { th } from "date-fns/locale/th";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import Modal from "@object/shared/components/Modal";
import Toggle from "@object/shared/components/Toggle";
import { z } from "zod";

import "react-big-calendar/lib/css/react-big-calendar.css";

const locales = { th };
const localizer = dateFnsLocalizer({ format, parse, startOfWeek, getDay, locales, locale: th });

type OpenStatus = "open" | "closed";

export interface OpenDayEvent extends Event {
  title: string;
  start: Date;
  end: Date;
  status: OpenStatus;
}

const mockOpenDays: OpenDayEvent[] = [
  {
    title: "เปิดทำการ",
    start: new Date(2025, 5, 24, 0, 0),
    end: new Date(2025, 5, 25, 0, 0), 
    status: "open",
  },
  {
    title: "หยุด",
    start: new Date(2025, 5, 29, 0, 0),
    end: new Date(2025, 5, 30, 0, 0), 
    status: "closed",
  },
  {
    title: "เปิดชั่วคราว 10:00–14:00",
    start: new Date(2025, 5, 26, 10, 0),
    end: new Date(2025, 5, 26, 14, 0),
    status: "open",
  },
];



export default function BranchCalendar() {
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [eventList, setEventList] = useState<OpenDayEvent[]>(mockOpenDays);

  const events = useMemo(() => {
    return eventList.map((e) => ({
      ...e,
      title: e.status === "open" ? e.title : "🚫 หยุดทำการ",
    }));
  }, [eventList]);

  const handleSelectSlot = (slotInfo: SlotInfo) => {
    setSelectedDate(slotInfo.start);
    setModalOpen(true);
  };

  const handleAdd = (event: OpenDayEvent) => {
    setEventList((prev) => [...prev, event]);
  };

  return (
    <div className="p-6 max-w-full mx-auto">
      <h1 className="text-2xl font-bold mb-4">ปฏิทินเปิด-ปิดร้าน</h1>

      <Calendar
        selectable
        localizer={localizer}
        events={events}
        startAccessor="start"
        endAccessor="end"
        style={{ height: 600 }}
        min={new Date(1970, 1, 1, 0, 0)}
        max={new Date(1970, 1, 1, 23, 59)}
        step={30}
        timeslots={1}
        scrollToTime={new Date(1970, 1, 1, 6, 0)}
        messages={{
          month: "เดือน", week: "สัปดาห์", day: "วัน", agenda: "รายการ",
          today: "วันนี้", next: "ถัดไป", previous: "ย้อนกลับ",
          showMore: (total) => `+ เพิ่มอีก ${total} รายการ`,
        }}
        formats={{
          monthHeaderFormat: (date) => `${formatDate(date, "MMMM", { locale: th })} ${date.getFullYear() + 543}`,
          dayHeaderFormat: (date) => `${formatDate(date, "EEEE d MMMM", { locale: th })} ${date.getFullYear() + 543}`,
          dayRangeHeaderFormat: ({ start, end }) =>
            `${formatDate(start, "d MMM", { locale: th })} – ${formatDate(end, "d MMM", { locale: th })}`,
          timeGutterFormat: (date) => formatDate(date, "HH:mm", { locale: th }), 
        }}
        onSelectSlot={handleSelectSlot}
        eventPropGetter={(event: OpenDayEvent) => {
          const backgroundColor = event.status === "open" ? "#D1FAE5" : "#FECACA";
          const color = event.status === "open" ? "#065F46" : "#B91C1C";
          return {
            style: {
              backgroundColor,
              color,
              border: "1px solid #ccc",
              borderRadius: "4px",
              padding: "4px",
            },
          };
        }}
      />


      <AddWorkingDayModal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        onAdd={handleAdd}
        date={selectedDate}
      />
    </div>
  );
}



const schema = z.object({
  start_time: z.string().min(1, "กรุณาระบุเวลาเปิด"),
  end_time: z.string().min(1, "กรุณาระบุเวลาปิด"),
});

type FormData = z.infer<typeof schema>;

interface AddWorkingDayModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAdd: (event: OpenDayEvent) => void;
  date: Date | null;
}

export function AddWorkingDayModal({
  isOpen,
  onClose,
  onAdd,
  date,
}: AddWorkingDayModalProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { isSubmitting, errors },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      start_time: "",
      end_time: "",
    },
  });

  const [isClosed, setIsClosed] = useState(false);

  useEffect(() => {
    if (isOpen) {
      reset({ start_time: "", end_time: "" });
      setIsClosed(false);
    }
  }, [isOpen, reset]);

  const onSubmit = (data: FormData) => {
    if (!date) return;
    const startDate = isClosed
      ? new Date(date.setHours(0, 0, 0, 0))
      : new Date(`${format(date, "yyyy-MM-dd")}T${data.start_time}:00`);
    const endDate = isClosed
      ? new Date(date.setHours(23, 59, 59, 999))
      : new Date(`${format(date, "yyyy-MM-dd")}T${data.end_time}:00`);

    const event: OpenDayEvent = {
      title: isClosed ? "หยุดทำการ" : `เปิดทำการ ${data.start_time}–${data.end_time}`,
      start: startDate,
      end: endDate,
      status: isClosed ? "closed" : "open",
      allDay: isClosed,
    };

    onAdd(event);
    onClose();
  };

  if (!isOpen || !date) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="เพิ่มเวลาเปิด - ปิด กรณีพิเศษ">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {!isClosed && (
          <>
            <div>
              <label className="block mb-1">เวลาเปิด</label>
              <input
                type="time"
                {...register("start_time")}
                className="input input-bordered w-full"
              />
              {errors.start_time && (
                <p className="text-red-500 text-sm">{errors.start_time.message}</p>
              )}
            </div>
            <div>
              <label className="block mb-1">เวลาปิด</label>
              <input
                type="time"
                {...register("end_time")}
                className="input input-bordered w-full"
              />
              {errors.end_time && (
                <p className="text-red-500 text-sm">{errors.end_time.message}</p>
              )}
            </div>
          </>
        )}

        <div>
          <Toggle checked={isClosed} onChange={setIsClosed} label=" วันหยุดทั้งวัน" />
        </div>

        <div className="flex justify-end gap-2">
          <button type="button" onClick={onClose} className="btn btn-ghost" disabled={isSubmitting}>
            ยกเลิก
          </button>
          <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
            {isSubmitting ? "กำลังบันทึก..." : "ยืนยัน"}
          </button>
        </div>
      </form>
    </Modal>
  );
}

