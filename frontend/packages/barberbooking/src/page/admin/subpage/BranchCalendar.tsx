import { useState, useMemo, useEffect } from "react";
import { Calendar, dateFnsLocalizer, Event, SlotInfo } from "react-big-calendar";
import { format, parse, startOfWeek, getDay, addDays, set } from "date-fns";
import { format as formatDate } from "date-fns";
import { th } from "date-fns/locale/th";
import "react-big-calendar/lib/css/react-big-calendar.css";
import Modal from "@object/shared/components/Modal";
import { getWorkingDayOverridesByDateRange, WorkingDayOverride } from "../../../api/workingDayOverride";
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

  const tenantId = Number(me?.tenant_ids?.[0]);
  const branchId = Number(me?.branch_id);
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [eventList, setEventList] = useState<OpenDayEvent[]>([]);
  const closedDays = workingHours.filter((w) => w.is_closed === true);
  const overrideDates = eventList.map((e) => formatDate(e.start, "yyyy-MM-dd"));
  const [currentViewDate, setCurrentViewDate] = useState(new Date());

  const weeklyClosedEvents: OpenDayEvent[] = closedDays.map((item) => {
    const today = new Date();
    const currentMonthStart = startOfMonth(currentViewDate);
    const currentMonthEnd = endOfMonth(currentViewDate);

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
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    if (slotInfo.start < today) {

      return;
    }

    setSelectedDate(slotInfo.start);
    setModalOpen(true);
  };

  useEffect(() => {
    if (!tenantId || !branchId) return;

    const today = selectedDate || new Date();
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
  }, [tenantId, branchId, currentViewDate]);



  return (
    <div className="p-6 max-w-full mx-auto">
      <h1 className="text-2xl font-bold mb-4">ปฏิทินเปิด-ปิดร้าน</h1>

      <Calendar
        selectable
        localizer={localizer}
        events={events}
        onNavigate={(date) => {
          setSelectedDate(date);
          setCurrentViewDate(date);
        }}
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
        view="month"
        dayPropGetter={(date) => {
          const today = new Date();
          today.setHours(0, 0, 0, 0);
          const isPast = date < today;

          return isPast
            ? {
              style: {
                backgroundColor: "#e5e7eb",
                color: "#9ca3af",           
                cursor: "not-allowed",
                pointerEvents: "none",     
                filter: "grayscale(0.6)",  
              },
            }
            : {};
        }}

        components={{
          event: CustomEvent,
        }}
      />
      {selectedDate && (
        <WorkingHourModal
          isOpen={modalOpen}
          onClose={() => setModalOpen(false)}
          onEdit={() => console.log("ss")}
          branchId={branchId}
          tenantId={tenantId}
          selectedDate={selectedDate} />
      )}
    </div>
  );
}

interface WorkingHourModalProps {
  isOpen: boolean;
  onClose: () => void;
  onEdit: (updated: WorkingHour) => void;
  tenantId: number;
  branchId: number;
  selectedDate: Date;
}

export function WorkingHourModal({
  isOpen,
  onClose,
  tenantId,
  branchId,
  selectedDate
}: WorkingHourModalProps) {

  if (!isOpen) return null;
  const [overrideDays, setOverrideDays] = useState<WorkingDayOverride[]>([]);

  const [newOverride, setNewOverride] = useState<OverrideDay>({
    date: format(selectedDate, "yyyy-MM-dd"),
    start_time: "08:00",
    end_time: "17:00",
    IsClosed: false,
    reason: "",
  });


  useEffect(() => {

    const dateStr = format(selectedDate, "yyyy-MM-dd");

    getWorkingDayOverridesByDateRange({
      tenantId,
      branchId,
      start: dateStr,
      end: dateStr,
    }).then((data) => {
      const matched = data.find(item => item.work_date.slice(0, 10) === dateStr);
      if (matched) {
        setNewOverride({
          date: matched.work_date,
          start_time: matched.start_time,
          end_time: matched.end_time,
          IsClosed: matched.is_closed,
          reason: matched.reason,
        });
      } else {
        setNewOverride({
          date: dateStr,
          start_time: "",
          end_time: "",
          IsClosed: true,
          reason: "",
        });
      }
      setOverrideDays(data);
    });
  }, [tenantId, branchId, selectedDate]);



  return (
    <Modal isOpen={isOpen} onClose={onClose} title={overrideDays.length === 0
      ? "เพิ่มเวลาเปิด - ปิด กรณีพิเศษ"
      : "แก้เวลาเปิด - ปิด กรณีพิเศษ"} size="md">
      <div className="w-full p-3">


        <div className="flex flex-col space-y-4">
          <div>
            <input
              type="date"
              className="input input-bordered w-full"
              value={newOverride.date}
              onChange={(e) => setNewOverride({ ...newOverride, date: e.target.value })}
            />
          </div>

          <div className="flex gap-6">
            <label className="inline-flex items-center">
              <input
                type="checkbox"
                checked={!newOverride.IsClosed}
                onChange={() => setNewOverride({ ...newOverride, IsClosed: false })}
                className="w-4 h-4 text-green-600 border-gray-300"
              />
              <span className="ml-2">เปิดร้าน</span>
            </label>
            <label className="inline-flex items-center">
              <input
                type="checkbox"
                checked={newOverride.IsClosed}
                onChange={() => setNewOverride({ ...newOverride, IsClosed: true })}
                className="w-4 h-4 text-red-600 border-gray-300"
              />
              <span className="ml-2">ปิดร้าน</span>
            </label>
          </div>
          <div className="flex gap-3">
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
                value={newOverride.reason}
                onChange={(e) => setNewOverride({ ...newOverride, reason: e.target.value })}
                className={`w-full input input-bordered`}
              />
            </div>
          </div>
          <div className="w-full flex gap-3">
            <button className="bg-green-500 text-white p-2 rounded-md w-full" onClick={() => console.log("d")}>
              {overrideDays.length === 0
                ? "เพิ่มเวลากรณีพิเศษ"
                : "แก้เวลากรณีพิเศษ"}</button>
                {overrideDays.length !== 0 ? <button className="bg-red-500 text-white p-2 rounded-md  w-full" onClick={() => console.log("d")}>เปลี่ยนเป็นเวลาปกติ</button> :""}
            
          </div>
        </div>
      </div>
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
