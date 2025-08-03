import { useState, useMemo, useEffect } from "react";
import { Calendar, dateFnsLocalizer, Event, SlotInfo } from "react-big-calendar";
import { format, parse, startOfWeek, getDay, addDays, set } from "date-fns";
import { format as formatDate } from "date-fns";
import { th } from "date-fns/locale/th";
import "react-big-calendar/lib/css/react-big-calendar.css";
import Modal from "@object/shared/components/Modal";
import { getWorkingDayOverridesByDateRange } from "../../../api/workingDayOverride";
import { startOfMonth, endOfMonth } from "date-fns";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { OverrideDay } from "../ManageTime";
import { useAppSelector } from "../../../store/hook";
import { WorkingHour } from "../../../api/workingHour";

const locales = { th };
const localizer = dateFnsLocalizer({ format, parse, startOfWeek, getDay, locales, locale: th });

type OpenStatus = "open" | "closed" | "weekly_closed";

export interface OpenDayEvent extends Event {
  title: string;
  start: Date;
  end: Date;
  status: OpenStatus;
}



export default function BranchCalendar({ workingHours }: { workingHours: WorkingHour[] }) {
  const me = useAppSelector((state) => state.auth.me);
  const statusMe = useAppSelector((state) => state.auth.statusMe);

  const tenantId = me?.tenant_ids?.[0];
  const branchId = Number(me?.branch_id);
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [eventList, setEventList] = useState<OpenDayEvent[]>([]);
  const closedDays = workingHours.filter((w) => w.is_closed === true);
  const overrideDates = eventList.map((e) => formatDate(e.start, "yyyy-MM-dd"));

 const weeklyClosedEvents: OpenDayEvent[] = closedDays.map((item) => {
  const today = new Date();
  const currentMonthStart = startOfMonth(today);
  const currentMonthEnd = endOfMonth(today);

  const allDates: OpenDayEvent[] = [];

  for (
    let d = new Date(currentMonthStart);
    d <= currentMonthEnd;
    d.setDate(d.getDate() + 1)
  ) {
    const day = d.getDay();
    const dateStr = formatDate(d, "yyyy-MM-dd");

    const isOverrideExists = overrideDates.includes(dateStr);
    if (day === item.week_day && !isOverrideExists) {
      const date = new Date(d);
      const start = new Date(date.setHours(0, 0, 0, 0));
      const end = new Date(date.setHours(23, 59, 59, 999));

      allDates.push({
        title: "หยุดประจำสัปดาห์",
        start,
        end,
        status: "weekly_closed",
      });
    }
  }

  return allDates;
}).flat();

  const allEvents = [...eventList, ...weeklyClosedEvents];



  const events = useMemo(() => {
    return allEvents.map((e) => {
      if (e.status === "open") {
        const startTime = formatDate(e.start, "HH:mm");
        const endTime = formatDate(e.end, "HH:mm");
        return {
          ...e,
          title: `เปิดกรณีพิเศษ\n(${startTime} - ${endTime})`,
        };
      } else if (e.status === "weekly_closed") {
        return {
          ...e,
          title: "หยุดประจำสัปดาห์",
        };
      } else {
        return {
          ...e,
          title: "หยุดกรณีพิเศษ",
        };
      }
    });
  }, [eventList, closedDays]);


  const handleSelectSlot = (slotInfo: SlotInfo) => {
    setSelectedDate(slotInfo.start);
    setModalOpen(true);
  };

  const handleAdd = (event: OpenDayEvent) => {
    setEventList((prev) => [...prev, event]);
  };

  function handleEditSuccess(updated: WorkingHour): void {
    throw new Error("Function not implemented.");
  }

  useEffect(() => {
    if (!tenantId || !branchId) return;

    const today = selectedDate || new Date(); // เดือนที่กำลังดูอยู่
    const start = formatDate(startOfMonth(today), "yyyy-MM-dd");
    const end = formatDate(endOfMonth(today), "yyyy-MM-dd");

    getWorkingDayOverridesByDateRange({
      tenantId,
      branchId,
      start,
      end,
    }).then((data) => {
      const events: OpenDayEvent[] = data.map((item) => {
        const date = new Date(item.work_date);
        const [startHour, startMin] = item.start_time.split(":").map(Number);
        const [endHour, endMin] = item.end_time.split(":").map(Number);

        return {
          title: item.is_closed ? "หยุดกรณีพิเศษ" : "เปิดกรณีพิเศษ",
          start: new Date(date.setHours(startHour, startMin)),
          end: new Date(date.setHours(endHour, endMin)),
          status: item.is_closed ? "closed" : "open",
        };
      });

      setEventList(events);
    });
  }, [tenantId, branchId, selectedDate]);



  return (
    <div className="p-6 max-w-full mx-auto">
      <h1 className="text-2xl font-bold mb-4">ปฏิทินเปิด-ปิดร้าน</h1>

      <Calendar
        selectable
        localizer={localizer}
        events={events}
        onNavigate={(date) => setSelectedDate(date)}
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
        eventPropGetter={(event) => {
          let backgroundColor = "#D1FAE5"; // open
          if (event.status === "closed") backgroundColor = "#FECACA";
          if (event.status === "weekly_closed") backgroundColor = "#E5E7EB"; 

          return {
            style: {
              backgroundColor,
              borderRadius: "6px",
              padding: "4px",
              border: "1px solid #ccc",
            },
          };
        }}
        components={{
          event: CustomEvent,
        }}
      />
      {selectedDate && (
        <WorkingHourModal
          isOpen={modalOpen}
          onClose={() => setModalOpen(false)}
          onEdit={handleEditSuccess}
          workingHour={undefined}
          branchId={branchId}
          tenantId={null} />
      )}
    </div>
  );
}

interface WorkingHourModalProps {
  isOpen: boolean;
  onClose: () => void;
  onEdit: (updated: WorkingHour) => void;
  workingHour: WorkingHour | undefined;
  tenantId: number | null;
  branchId: number | null;
}

export function WorkingHourModal({
  isOpen,
  onClose,
  onEdit,
  workingHour,
  tenantId,
  branchId,
}: WorkingHourModalProps) {

  if (!isOpen || !workingHour) return null;
  const [overrideDays, setOverrideDays] = useState<OverrideDay[]>([]);

  const [newOverride, setNewOverride] = useState<OverrideDay>({
    date: "", start_time: "08:00", end_time: "17:00", IsClosed: true,
  });

  const handleAddOverride = () => {
    setOverrideDays([...overrideDays, newOverride]);
    setNewOverride({ date: "", start_time: "08:00", end_time: "17:00", IsClosed: true, });
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="แก้ไขเวลาเปิด-ปิด">
      <section>
        <h2 className="text-xl font-semibold mb-2">เพิ่มเวลาเปิด - ปิด กรณีพิเศษ</h2>
        <div className="space-y-4">
          <input
            type="date"
            className="input input-bordered w-full max-w-xs"
            value={newOverride.date}
            onChange={(e) => setNewOverride({ ...newOverride, date: e.target.value })}
          />
          <div className="flex space-x-4">
            <div>
              <label className="block text-sm font-medium">เวลาเปิด</label>
              <input
                type="time"
                className="input input-bordered"
                value={newOverride.start_time}
                onChange={(e) => setNewOverride({ ...newOverride, start_time: e.target.value })}
              />
            </div>
            <div>
              <label className="block text-sm font-medium">เวลาปิด</label>
              <input
                type="time"
                className="input input-bordered"
                value={newOverride.end_time}
                onChange={(e) => setNewOverride({ ...newOverride, end_time: e.target.value })}
              />
            </div>
            <div>
              <label className="block text-sm font-medium">
                หมายเหตุ
              </label>
              <input
                type="text"
                placeholder="หมายเหตุ"
                className={`w-full input input-bordered`}
              />
            </div>
          </div>
          <button className="bg-blue-500 text-white p-2 rounded-md" onClick={handleAddOverride}>เพิ่มวันทำการ</button>
        </div>
      </section>
    </Modal>
  );
}
type Props = {
  event: OpenDayEvent;
};
export function CustomEvent({ event }: Props) {
  if (event.status === "weekly_closed") {
    return <span className="text-gray-600 font-medium">หยุดประจำสัปดาห์</span>;
  }

  const isOpen = event.status === "open";
  const start = formatDate(event.start, "HH:mm");
  const end = formatDate(event.end, "HH:mm");

  return (
    <div className="flex flex-col">
      <span className={isOpen ? "text-green-900 font-semibold" : "text-red-800 font-semibold"}>
        {isOpen ? "เปิดกรณีพิเศษ" : "หยุดกรณีพิเศษ"}
      </span>
      {isOpen && (
        <span className="text-xs text-gray-700">{`(${start} - ${end})`}</span>
      )}
    </div>
  );
}
