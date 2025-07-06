
import { useState, useCallback, useEffect } from "react";
import axios from "../lib/axios";
import type { Barber } from "../types/barber";

import { Calendar, dateFnsLocalizer, Event, SlotInfo } from "react-big-calendar";
import { format, parse, startOfWeek, getDay } from "date-fns";
import { format as formatDate } from "date-fns";
import { th } from "date-fns/locale/th";
import "react-big-calendar/lib/css/react-big-calendar.css";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { getWorkingHourRangeAxios } from "../components/TimeSelector"
import { MdFilterAlt } from "react-icons/md";

export default function BarberPage() {
    const [barbers, setBarbers] = useState<Barber[]>([]);
    const [loadingBarbers, setLoadingBarbers] = useState<boolean>(false);
    const [errorBarbers, setErrorBarbers] = useState<string | null>(null);
    const [displayBarbers, setDisplayBarbers] = useState<boolean>(true);
    const loadBarbers = useCallback(async () => {
        setLoadingBarbers(true);
        setErrorBarbers(null);
        try {
            const res = await axios.get<{ status: string; data: Barber[] }>(
                `/barberbooking/branches/1/barbers`
            );
            if (res.data.status !== "success") {
                throw new Error(res.data.status);
            }
            setBarbers(res.data.data);
        } catch (err: any) {
            setErrorBarbers(err.response?.data?.message || err.message || "Failed to load barbers");
        } finally {
            setLoadingBarbers(false);
        }
    }, []);

    useEffect(() => {
        loadBarbers();
    }, []);
    const [isModalOpen, setModalOpen] = useState(false);


    return (
        <div className="min-h-screen bg-gradient-to-b from-white via-slate-100 to-slate-200 text-gray-900">
            <div className="container mx-auto py-6 px-6">
                <h1 className="text-4xl font-extrabold mb-4">ช่างของเรา</h1>
                <div className="flex gap-3 mb-3">
                    <button className="px-3 py-1.5 bg-gray-400 rounded text-white hover:bg-blue-700 transition"
                        onClick={() => setDisplayBarbers(true)}
                    >
                        ดูคิวของช่างรายคน
                    </button>

                    <button className="px-3 py-1.5 bg-gray-400 rounded text-white hover:bg-blue-700 transition"
                        onClick={() => setDisplayBarbers(false)}
                    >
                        ดูคิวของช่างทั้งร้าน
                    </button>
                </div>


                {displayBarbers ?
                    (<div>
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
                            {loadingBarbers && <p>Loading barbers…</p>}
                            {errorBarbers && <p className="text-red-500">Error loading barbers: {errorBarbers}</p>}
                            {barbers.map((barber) => (
                                <div
                                    key={barber.id}
                                    className="bg-gray-200 rounded-lg shadow-lg hover:shadow-xl transition flex flex-col h-full"
                                >
                                    <img
                                        src={`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${barber.img_path}/${barber.img_name}`}
                                        alt={barber.username}
                                        className="w-full h-64 object-cover object-top rounded-t-lg"
                                    />

                                    <div className="p-4 flex flex-col justify-between flex-1">
                                        <div>
                                            <h2 className="text-2xl font-semibold mb-2">
                                                {barber.username}
                                            </h2>
                                            <p className="text-gray-400 text-sm mb-2">
                                                {barber.description}
                                            </p>
                                        </div>
                                        <div className="flex items-center justify-between text-gray-900 mt-4">
                                            <span className="text-sm bg-gray-300 p-1.5 rounded">
                                                ⭐ 4.5
                                            </span>

                                            <button className="px-3 py-1.5 bg-gray-400 rounded text-white hover:bg-blue-700 transition"
                                                onClick={() => setModalOpen(true)}>
                                                ดูคิวของช่าง
                                            </button>
                                            {/* <BarberScheduleModal
                                        isOpen={isModalOpen}
                                        onClose={() => setModalOpen(false)}
                                        slots={dummySlots}
                                        barberName="สมชาย ศรีสุข"
                                    /> */}

                                            <BarberScheduleContainer
                                                tenantId={1}
                                                branchId={1}
                                                barberName="123"
                                                isOpen={isModalOpen}
                                                onClose={() => setModalOpen(false)}
                                            />

                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>

                    </div>) : (
                        <div>
                            <TotalBarberSchedule barbers={barbers} />
                        </div>
                    )}

            </div>
        </div>
    );
}
const TotalBarberSchedule = ({ barbers }: { barbers: Barber[] }) => {
    const [slot, setSlot] = useState<string[]>([]);
    const [startTime, setStartTime] = useState<string>("");
    const [endTime, setEndTime] = useState<string>("");
    const today = format(new Date(), "yyyy-MM-dd");
    const [selectedDate, setSelectedDate] = useState(today);
    interface WorkingHourResult {
        start: string;
        end: string;
    }
    useEffect(() => {
        async function fetchSlot() {

            const result = await getWorkingHourRangeAxios(1, 1, new Date(selectedDate));
            if (result?.start && result?.end) {
                const generated = generateTimeSlots(result.start, result.end);
                setSlot(generated);
            } else {
                setSlot([]);
            }
        }

        fetchSlot();
    }, [selectedDate]);

    function generateTimeSlots(start: string, end: string): string[] {
        const slots: string[] = [];

        const [startHour, startMinute] = start.split(":").map(Number);
        const [endHour, endMinute] = end.split(":").map(Number);

        let current = new Date();
        current.setHours(startHour, startMinute, 0, 0);

        const endTime = new Date();
        endTime.setHours(endHour, endMinute, 0, 0);

        while (current <= endTime) {
            const hour = current.getHours().toString().padStart(2, "0");
            const minute = current.getMinutes().toString().padStart(2, "0");
            slots.push(`${hour}:${minute}`);

            current.setMinutes(current.getMinutes() + 30);
        }

        return slots;
    }

    function filterQuarterSlots(slots: string[]) {
        return slots.filter((time) => {
            const minute = time.split(":")[1];
            return ["00", "30"].includes(minute);
        });
    }
    const options = [
        { label: "สัปดาห์นี้", value: 1 },
        { label: "เดือนนี้", value: 2 },
    ];
    const [selectedOption, setSelectedOption] = useState<number | null>(null);
    const [openFilter, setOpenFilter] = useState<boolean>(false)

    return (
        <>
            <div className="flex flex-col border rounded-md overflow-hidden">

                <div className="p-2 border-b border-gray-200 flex items-center">
                    <div className="flex flex-col gap-3 ">
                        <div className="flex gap-2 items-center p-1 ">
                            <label className="text-lg font-bold">เลือกวันที่คุณว่าง </label>
                            <input
                                type="date"
                                min={today}
                                value={selectedDate}
                                onChange={(e) => setSelectedDate(e.target.value)}
                                className="input input-bordered rounded p-1"
                            />
                            <button className="bg-gray-300 p-1 rounded-md" onClick={() => setOpenFilter(prev => !prev)}><MdFilterAlt size={24} color="white" /></button>
                        </div>


                        <div className="flex flex-col gap-3 ">
                            {openFilter ? (
                                <div className="flex flex-col gap-3">
                                    <h4 className="text-lg font-bold">ตัวกรองช่วงเวลา</h4>
                                    <div className="flex gap-4 px-8 items-center">
                                        {options.map((option) => (
                                            <label
                                                key={option.value}
                                                className="flex items-center space-x-2 cursor-pointer"
                                            >
                                                <input
                                                    type="checkbox"
                                                    checked={selectedOption === option.value}
                                                    onChange={() =>
                                                        setSelectedOption(
                                                            selectedOption === option.value ? null : option.value
                                                        )
                                                    }
                                                    className="w-4 h-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
                                                />
                                                <span className="text-gray-800">{option.label}</span>
                                            </label>
                                        ))}
                                    </div>
                                    <div className="flex gap-2 items-center px-8">
                                        <select
                                            className="border py-1 rounded"
                                            value={startTime}
                                            onChange={(e) => setStartTime(e.target.value)}
                                        >
                                            <option value="">เวลาเริ่มต้น</option>
                                            {filterQuarterSlots(slot).map((time) => (
                                                <option key={`start-${time}`} value={time}>
                                                    {time}
                                                </option>
                                            ))}
                                        </select>

                                        <span>-</span>

                                        <select
                                            className="border py-1 rounded"
                                            value={endTime}
                                            onChange={(e) => setEndTime(e.target.value)}
                                        >
                                            <option value="">เวลาสิ้นสุด</option>
                                            {filterQuarterSlots(slot).map((time) => (
                                                <option key={`end-${time}`} value={time}>
                                                    {time}
                                                </option>
                                            ))}
                                        </select>
                                    </div>
                                </div>

                            ) : (
                                <div>
                                    {/* ตัวกรองปิดอยู่ */}
                                </div>
                            )}

                        </div>
                    </div>
                </div>
                <div className="flex bg-gray-100 text-center font-semibold divide-x divide-gray-200 border-b border-gray-200">
                    <div className="py-2 w-[120px]">เวลา</div>
                    {barbers.map((barber) => (
                        <div key={barber.id} className="py-2 flex-1">
                            {barber.username}
                        </div>
                    ))}
                </div>


                <div className="w-full">
                    {slot.length === 0 ? (
                        <div className="text-center text-red-500 py-4">วันนี้ร้านปิด</div>
                    ) : (
                        slot.map((time, timeIndex) => (
                            <div
                                key={timeIndex}
                                className={`flex text-center divide-x divide-gray-200 border-b border-gray-200 ${timeIndex === slot.length - 1 ? "border-b-0" : ""
                                    }`}
                            >
                                <div className="py-2 w-[120px] bg-gray-50">{time}</div>
                                {barbers.map((barber) => (
                                    <div key={barber.id + "_" + timeIndex} className="py-2 flex-1">
                                        ว่าง
                                    </div>
                                ))}
                            </div>
                        ))
                    )}
                </div>
            </div>
        </>
    );

}

const locales = { th };
const localizer = dateFnsLocalizer({ format, parse, startOfWeek, getDay, locales, locale: th });
interface TimeSlot {
    id: string;
    date: string;
    time: string; // e.g. '09:00'
    available: boolean;
}

interface OpenDayEvent {
    title: string;
    start: Date;
    end: Date;
    status: "open" | "closed";
}

interface BarberScheduleModalProps {
    isOpen: boolean;
    onClose: () => void;
    slots: TimeSlot[];
    barberName: string;
}

interface TimeSlot {
    id: string;
    date: string;
    time: string;
    available: boolean;
}

interface RawSlot {
    start: string;
    end: string;
    status: "open" | "closed";
}

export async function fetchAvailableSlots(
    tenantId: number,
    branchId: number,
    startDate: string,
    endDate: string
): Promise<TimeSlot[]> {
    const res = await axios.get<{ status: string; data: RawSlot[] }>(
        `/barberbooking/tenants/${tenantId}/branches/${branchId}/available-slots?start=${startDate}&end=${endDate}`
    );

    if (res.data.status !== "success") {
        throw new Error(res.data.status);
    }

    const timeSlots: TimeSlot[] = res.data.data.map((item) => {
        const start = new Date(item.start);
        return {
            id: `${item.start}-${item.end}`,
            date: format(start, "yyyy-MM-dd"),
            time: format(start, "HH:mm"),
            available: item.status === "open",
        };
    });

    return timeSlots;
}



export function BarberScheduleModal({
    isOpen,
    onClose,
    slots,
    barberName,
}: BarberScheduleModalProps) {
    if (!isOpen) return null;

    const events: OpenDayEvent[] = slots.map((slot) => {
        const dateTime = new Date(`${slot.date}T${slot.time}`);
        return {
            title: slot.available ? "ว่าง" : "",
            start: dateTime,
            end: new Date(dateTime.getTime() + 30 * 60000),
            status: slot.available ? "open" : "closed",
        };
    });

    console.log(events)

    const minTime = events.length > 0
        ? new Date(Math.min(...events.map((e) => e.start.getTime())))
        : new Date(1970, 1, 1, 9, 0);

    const maxTime = events.length > 0
        ? new Date(Math.max(...events.map((e) => e.end.getTime())))
        : new Date(1970, 1, 1, 17, 0);

    const CustomAgendaDate = ({ date }: { date: Date }) => {
        return (
            <span>
                {formatDate(date, "EEEEที่ d MMMM yyyy", { locale: th }).replace(
                    `${date.getFullYear()}`,
                    `${date.getFullYear() + 543}`
                )}
            </span>
        );
    };
    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
            onClick={onClose}
        >
            <div
                className="bg-white rounded-lg w-full max-w-5xl mx-4 shadow-lg overflow-hidden"
                onClick={(e) => e.stopPropagation()}
            >
                {/* Header */}
                <div className="flex justify-between items-center px-6 py-4 border-b text-black">
                    <h3 className="text-xl font-semibold">
                        ตารางงานของ {barberName}
                    </h3>
                    <button
                        className="text-gray-500 hover:text-gray-700"
                        onClick={onClose}
                    >
                        ✕
                    </button>
                </div>

                {/* Calendar */}
                <div className="p-6 max-h-[70vh] overflow-auto text-black">
                    <Calendar
                        localizer={localizer}
                        events={events}
                        startAccessor="start"
                        endAccessor="end"
                        style={{ height: 600 }}
                        step={30}
                        timeslots={1}
                        scrollToTime={minTime}
                        min={minTime}
                        max={maxTime}
                        defaultView="week"
                        views={["week", "day"]}
                        messages={{
                            date: "วันที่",
                            time: "เวลา",
                            event: "เหตุการณ์",
                            week: "สัปดาห์",
                            day: "วัน",
                            today: "วันนี้",
                            previous: "ย้อนกลับ",
                            next: "ถัดไป",
                            showMore: (total) => `+ เพิ่มอีก ${total} รายการ`,
                        }}
                        formats={{
                            monthHeaderFormat: (date) =>
                                `${formatDate(date, "MMMM", { locale: th })} ${date.getFullYear() + 543}`,
                            dayHeaderFormat: (date) =>
                                `${formatDate(date, "EEEE d MMMM", { locale: th })} ${date.getFullYear() + 543}`,

                            dayRangeHeaderFormat: ({ start, end }) =>
                                `${formatDate(start, "d MMM", { locale: th })} – ${formatDate(end, "d MMM", { locale: th })}`,
                            timeGutterFormat: (date) => formatDate(date, "HH:mm", { locale: th }),
                            dayFormat: (date) =>
                                `${formatDate(date, "dd", { locale: th })} ${formatDate(date, "EEE", { locale: th })}`,
                            eventTimeRangeFormat: ({ start, end }, culture, localizer) => {
                                const s = formatDate(start, "HH:mm");
                                const e = formatDate(end, "HH:mm");
                                return `${s} - ${e}`;
                            },
                            agendaDateFormat: (date) =>
                                `${formatDate(date, "EEEEที่ d MMMM", { locale: th })} ${date.getFullYear() + 543}`,
                            agendaHeaderFormat: ({ start, end }) =>
                                `${formatDate(start, "d MMM", { locale: th })} – ${formatDate(end, "d MMM", { locale: th })}`,

                        }}
                        eventPropGetter={(event) => {
                            const backgroundColor = event.status === "open" ? "#D1FAE5" : "#FECACA";
                            const color = event.status === "open" ? "#065F46" : "#B91C1C";
                            return {
                                style: {
                                    backgroundColor,
                                    color,
                                    border: "1px solid #ccc",
                                    borderRadius: "4px",
                                    height: "100%",
                                    padding: "0",
                                    display: "flex",
                                    alignItems: "center",
                                    justifyContent: "center",
                                },
                            };

                        }}
                    />
                </div>

                {/* Footer */}
                <div className="flex justify-end px-6 py-4 border-t">
                    <button
                        className="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 transition"
                        onClick={onClose}
                    >
                        ปิด
                    </button>
                </div>
            </div>
        </div>
    );
}


export function BarberScheduleContainer({
    tenantId,
    branchId,
    barberName,
    isOpen,
    onClose,
}: {
    tenantId: number;
    branchId: number;
    barberName: string;
    isOpen: boolean;
    onClose: () => void;
}) {
    const [slots, setSlots] = useState<TimeSlot[]>([]);
    const [loading, setLoading] = useState(false);

    // useEffect(() => {
    //     if (!isOpen) return;

    //     const today = new Date();
    //     const next7 = new Date();
    //     next7.setDate(today.getDate() + 6);

    //     const startDate = format(today, "yyyy-MM-dd");
    //     const endDate = format(next7, "yyyy-MM-dd");

    //     setLoading(true);
    //     fetchAvailableSlots(tenantId, branchId, startDate, endDate)
    //         .then(setSlots)
    //         .catch(console.error)
    //         .finally(() => setLoading(false));
    // }, [isOpen, tenantId, branchId]);
    useEffect(() => {
        if (isOpen) {
            const mockSlots: TimeSlot[] = [
                { id: "1", date: "2025-07-01", time: "09:00", available: false },
                { id: "2", date: "2025-07-01", time: "09:30", available: false },
                { id: "3", date: "2025-07-01", time: "10:00", available: false },
                { id: "4", date: "2025-07-02", time: "13:00", available: false },
                { id: "5", date: "2025-07-03", time: "15:30", available: false },
            ];
            setSlots(mockSlots);
        }
    }, [isOpen]);

    return (
        <BarberScheduleModal
            isOpen={isOpen}
            onClose={onClose}
            slots={slots}
            barberName={barberName}
        />
    );
}

