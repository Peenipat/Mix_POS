
import { useState, useCallback, useEffect } from "react";
import axios from "../lib/axios";
import type { Barber } from "../types/barber";

import { Calendar, dateFnsLocalizer, Event, SlotInfo } from "react-big-calendar";
import { format, parse, startOfWeek, getDay, parseISO } from "date-fns";
import { format as formatDate } from "date-fns";
import { th } from "date-fns/locale/th";
import "react-big-calendar/lib/css/react-big-calendar.css";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { getWorkingHourRangeAxios } from "../components/TimeSelector"
import { MdFilterAlt } from "react-icons/md";
import { AppointmentBrief, getAppointmentsByBranch } from "../api/appointment";
import { AppointmentLock } from "../api/appointmentLock";

import isSameOrBefore from "dayjs/plugin/isSameOrBefore";
import customParseFormat from "dayjs/plugin/customParseFormat";
import dayjs from 'dayjs';
import Modal from "@object/shared/components/Modal";
import React from "react";
dayjs.extend(isSameOrBefore);
dayjs.extend(customParseFormat);

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

    const haddleClose = () => {
        setModalOpen(false)
    }


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

                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                        <BarberScheduleModal isOpen={isModalOpen} onClose={haddleClose} barberName="123" />

                    </div>) : (
                        <div>
                            {/* <TotalBarberSchedule barbers={barbers} /> */}
                        </div>
                    )}

            </div>
        </div>
    );
}
export const TotalBarberSchedule = ({
    barbers,
    onClick,
    appointments = [],
    locks = [],
    selectedDate,
    setSelectedDate

}: {
    barbers: Barber[];
    onClick?: (date: string, barberId: number, time: string) => void;
    appointments?: AppointmentBrief[]; // for mock
    locks?: AppointmentLock[];
    selectedDate: string;
    setSelectedDate: (date: string) => void;
}) => {
    const [slotMap, setSlotMap] = useState<Record<string, string[]>>({});
    const today = format(new Date(), "yyyy-MM-dd")
    const [startTime, setStartTime] = useState<string>("");
    const [endTime, setEndTime] = useState<string>("");
    const [selectedOption, setSelectedOption] = useState<"week" | "month" | "" | null>("");

    const fetchSlot = useCallback(async () => {
        try {
            const result = await getWorkingHourRangeAxios(
                1,
                1,
                new Date(selectedDate),
                selectedOption
                    ? {
                        filter: selectedOption as "week" | "month",
                        fromTime: startTime || undefined,
                        toTime: endTime || undefined,
                    }
                    : undefined
            );

            if (!result) {
                setSlotMap({});
                return;
            }
            console.log("🔍 result from API:", result);

            if (result.type === "range") {
                const generatedMap: Record<string, string[]> = {};

                // 🔍 หาวันแรก–วันสุดท้ายใน range
                const allDates = Object.keys(result.data).sort();
                const startDate = dayjs(allDates[0]);
                const endDate = dayjs(allDates[allDates.length - 1]);

                for (let d = startDate; d.isSameOrBefore(endDate); d = d.add(1, 'day')) {
                    const dateStr = d.format("YYYY-MM-DD");
                    const range = result.data[dateStr];

                    // ✅ ตรวจว่ามี start/end ไหม
                    if (range && typeof range === "object" && "start" in range && "end" in range) {
                        generatedMap[dateStr] = generateTimeSlots(range.start, range.end);
                    } else {
                        generatedMap[dateStr] = []; // ร้านปิด หรือไม่มีข้อมูล
                    }
                }
                Object.entries(result.data).forEach(([date, val]) => {
                    console.log("📅", date, "➡️", val);
                });

                setSlotMap(generatedMap);
            }
            else if (result.type === "single") {
                // 📦 แบบวันเดียว
                const generated = generateTimeSlots(result.start, result.end);
                setSlotMap({ [result.date]: generated });
            }
        } catch (err) {
            console.error("❌ Failed to fetch slots:", err);
            setSlotMap({});
        }
    }, [selectedDate, selectedOption, startTime, endTime]);


    useEffect(() => {
        fetchSlot();
    }, [fetchSlot, selectedDate, selectedOption, startTime, endTime]);


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

    function addMinutes(time: string, mins: number): string {
        const [h, m] = time.split(":").map(Number);
        const date = new Date();
        date.setHours(h, m + mins, 0, 0);
        return date.toTimeString().substring(0, 5);
    }

    function timeToMinutes(time: string): number {
        const [h, m] = time.split(":").map(Number);
        return h * 60 + m;
    }

    function getBookingFraction(
        barberId: number,
        slot: string,
        appointments: AppointmentBrief[] | undefined,
        selectedDate: string
    ): "full" | "top" | "bottom" | null {
        if (!appointments) return null;

        const slotStart = timeToMinutes(slot);
        const slotEnd = timeToMinutes(addMinutes(slot, 30));

        for (const app of appointments) {
            if (app.barber_id !== barberId || app.date !== selectedDate) continue;

            const appStart = timeToMinutes(app.start);
            const appEnd = timeToMinutes(app.end);

            if (appStart <= slotStart && appEnd >= slotEnd) return "full";
            if (appStart <= slotStart && appEnd > slotStart && appEnd <= slotEnd) return "top";
            if (appStart >= slotStart && appStart < slotEnd && appEnd >= slotEnd) return "bottom";
        }

        return null;
    }


    function isSlotLocked(barberId: number, slotStartTime: string): boolean {
        const slotStart = timeToMinutes(slotStartTime);
        const slotEnd = timeToMinutes(addMinutes(slotStartTime, 30));

        for (const appt of appointments) {
            if (appt.barber_id !== barberId) continue;

            const apptStart = timeToMinutes(appt.start);
            const apptEnd = timeToMinutes(appt.end);

            if (slotStart < apptEnd && slotEnd > apptStart) {
                return false;
            }
        }

        for (const lock of locks || []) {
            if (lock.barber_id !== barberId || !lock.is_active) continue;

            const lockStart = timeToMinutes(lock.start_time.slice(11, 16));
            const lockEnd = timeToMinutes(lock.end_time.slice(11, 16));

            if (slotStart < lockEnd && slotEnd > lockStart) {
                return true;
            }
        }

        return false;
    }

    function isPastSlot(dateStr: string, timeStr: string): boolean {
        const slotDateTime = dayjs(`${dateStr} ${timeStr}`, "YYYY-MM-DD HH:mm");
        return slotDateTime.isBefore(dayjs());
    }

    const options = [
        { label: "สัปดาห์นี้", value: "week" },
        { label: "เดือนนี้", value: "month" },
    ] as const;

    const [openFilter, setOpenFilter] = useState<boolean>(false)

    const allTimeSlots: string[] = (() => {
        const slots: string[] = [];
        const current = new Date();
        current.setHours(0, 0, 0, 0);

        const end = new Date();
        end.setHours(23, 59, 0, 0);

        while (current <= end) {
            const hour = current.getHours().toString().padStart(2, "0");
            const minute = current.getMinutes().toString().padStart(2, "0");
            slots.push(`${hour}:${minute}`);
            current.setMinutes(current.getMinutes() + 60);
        }

        return slots;
    })();

    const [appointmentList, setAppointmentList] = useState<AppointmentBrief[]>()
    // useEffect(() => {
    //     async function fetchAppointment() {
    //         const appointments = await getAppointmentsByBranch(1, selectedDate, selectedDate, selectedOption, ["CANCELLED"]);
    //         setAppointmentList(appointments ?? []);
    //     }
    //     fetchAppointment();
    // }, [selectedDate]);

    function handdlesFilter() {
        setSelectedOption("")
        setStartTime("")
        setEndTime("")
        setOpenFilter(!openFilter)
    }

    return (
        <div className="flex flex-col border rounded-md overflow-auto shadow-md bg-white">
            <div className="p-4 border-b border-gray-200 bg-gray-50">
                <div className="flex items-center gap-4 flex-wrap">
                    <label className="text-lg font-semibold">เลือกวันที่คุณว่าง</label>
                    <input
                        id="date-picker"
                        type="date"
                        min={today}
                        value={selectedDate}
                        disabled={selectedOption !== ""}
                        onChange={(e) => setSelectedDate(e.target.value)}
                        className={`input input-bordered p-2 rounded-md border border-gray-300 ${selectedOption !== "" ? "bg-gray-100 text-gray-400 cursor-not-allowed" : ""
                            }`}
                    />
                    <button
                        id="filter-button"
                        onClick={() => handdlesFilter()}
                        className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1.5 rounded-md flex items-center transition"
                    >
                        <MdFilterAlt size={20} className="mr-1" />
                        ตัวกรอง
                    </button>
                </div>
                <div
                    id="all-option"
                    className={`mt-4 space-y-3 transition-all duration-300 ${openFilter ? "block" : "hidden"}`}
                >
                    <div className="mt-4 space-y-3">
                        <h4 className="text-md font-semibold">ตัวกรองช่วงเวลา</h4>
                        <div className="flex gap-5 w-max" id="select-option">
                            {options.map((option) => (
                                <label key={option.value} className="flex items-center gap-2">
                                    <input
                                        type="checkbox"
                                        checked={selectedOption === option.value}
                                        onChange={() =>
                                            setSelectedOption(selectedOption === option.value ? "" : option.value)
                                        }
                                        className="w-4 h-4 text-blue-600 border-gray-300"
                                    />
                                    <span>{option.label}</span>
                                </label>
                            ))}
                        </div>

                        {/* Time Selectors */}
                        <div className="flex gap-3 items-center flex-wrap">
                            <div id="time-option" className="flex gap-3 items-center flex-wrap">
                                <select value={startTime} onChange={(e) => setStartTime(e.target.value)}>
                                    <option value="">เวลาเริ่มต้น</option>
                                    {allTimeSlots.map((time) => (
                                        <option key={`start-${time}`} value={time}>{time}</option>
                                    ))}
                                </select>
                                <span>-</span>
                                <select value={endTime} onChange={(e) => setEndTime(e.target.value)}>
                                    <option value="">เวลาสิ้นสุด</option>
                                    {allTimeSlots.map((time) => (
                                        <option key={`start-${time}`} value={time}>{time}</option>
                                    ))}
                                </select>
                            </div>

                            <button
                                id="search-button"
                                onClick={fetchSlot}
                                className="bg-gray-700 hover:bg-blue-700 text-white px-4 py-2 rounded-md transition"
                            >
                                ค้นหา
                            </button>


                        </div>
                    </div>
                </div>
            </div>
            <div id="data-box">
                <div className="w-full">
                    {Object.entries(slotMap).length === 0 ? (
                        <div className="text-center text-red-500 font-semibold py-6">
                            ไม่พบเวลาทำงานในช่วงนี้
                        </div>
                    ) : (
                        Object.entries(slotMap).map(([date, times]) => (
                            <div key={date} className="mb-6">

                                <div className="bg-gray-200 text-gray-800 font-bold px-4 py-2">
                                    วันที่ {date.split("-").reverse().join("/")}
                                </div>
                                <div className="flex text-center font-semibold bg-gray-100 divide-x divide-gray-300 border-b border-gray-300">
                                    <div className="py-2 w-[120px]">เวลา</div>
                                    {barbers.map((barber) => (
                                        <div key={barber.id} className="py-2 flex-1" id="barber-name">
                                            {barber.username}
                                        </div>
                                    ))}
                                </div>

                                {times.length === 0 ? (
                                    <div className="text-center text-red-500 font-semibold py-4">ร้านปิด</div>
                                ) : (
                                    times.map((time, timeIndex) => (
                                        <div key={timeIndex} className="flex" id="barber-time">
                                            <div className="w-[120px] py-2 bg-gray-100 font-semibold text-center">
                                                {time}
                                            </div>
                                            {barbers.map((barber) => {
                                                const fraction = getBookingFraction(barber.id, time, appointmentList, date);
                                                return (
                                                    <div key={barber.id + "_" + timeIndex} className="relative h-14 flex-1 border-l">
                                                        {/* จองแล้ว (red) */}
                                                        {fraction === "full" && (
                                                            <div className="absolute top-0 h-full w-full bg-red-400 flex items-end justify-center text-white font-semibold pointer-events-none z-20">
                                                                <p>จองแล้ว</p>
                                                            </div>
                                                        )}

                                                        {/* ว่าง (green) */}
                                                        {(fraction === "bottom" || fraction === null) &&
                                                            !isSlotLocked(barber.id, time) &&
                                                            !isPastSlot(date, time) && (
                                                                <div
                                                                    className={`absolute top-0 ${fraction === "bottom" ? "h-1/2" : "h-full"} w-full bg-green-100 text-green-600 font-semibold flex items-center justify-center hover:bg-green-200 cursor-pointer z-10`}
                                                                    onClick={() => onClick?.(date, barber.id, time)}
                                                                >
                                                                    ว่าง
                                                                </div>
                                                            )}

                                                        {fraction === "top" &&
                                                            !isSlotLocked(barber.id, time) &&
                                                            !isPastSlot(date, time) && (
                                                                <div
                                                                    className="absolute bottom-0 h-1/2 w-full bg-green-100 text-green-600 font-semibold flex items-center justify-center hover:bg-green-200 cursor-pointer z-10"
                                                                    onClick={() => onClick?.(date, barber.id, addMinutes(time, 15))}
                                                                >
                                                                    ว่าง
                                                                </div>
                                                            )}

                                                        {/* Slot หมดเวลา */}
                                                        {isPastSlot(date, time) && (
                                                            <div className="absolute inset-0 bg-gray-300/70 flex items-center justify-center text-gray-600 font-semibold pointer-events-none z-0">
                                                                หมดเวลา
                                                            </div>
                                                        )}


                                                        {/* จองครึ่ง (red overlay) */}
                                                        {fraction === "top" && (
                                                            <div className="absolute top-0 h-1/2 w-full bg-red-400 flex justify-center items-end pb-[2px] text-white font-semibold text-sm pointer-events-none z-20" />
                                                        )}
                                                        {fraction === "bottom" && (
                                                            <div className="absolute bottom-0 h-1/2 w-full bg-red-400 flex justify-center items-start pt-[2px] text-white font-semibold text-sm pointer-events-none z-20" />
                                                        )}

                                                        {isSlotLocked(barber.id, time) && (
                                                            <div className="absolute inset-0 bg-yellow-300/80 flex items-end justify-center text-black font-semibold pointer-events-none z-0">
                                                                มีคนกำลังจองอยู่...
                                                            </div>
                                                        )}
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    ))
                                )}
                            </div>
                        ))
                    )}
                </div>
            </div>


        </div >
    );
};

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
    barberName,
}: BarberScheduleModalProps) {
    const timeSlots = ["09:00", "10:00", "11:00", "13:00"];

    const dates = ["01/08/68", "02/08/68", "03/08/68", "04/08/68", "05/08/68", "06/08/68", "07/08/68"];

    type TimeSlotStatus = "available" | "unavailable" | "booked";

    type BarberScheduleMatrix = {
        [time: string]: {
            [date: string]: TimeSlotStatus;
        };
    };

    const scheduleMatrix: BarberScheduleMatrix = {
        "09:00": {
            "01/08/68": "available",
            "02/08/68": "booked",
            "03/08/68": "available",
            "04/08/68": "unavailable",
            "05/08/68": "unavailable",
            "06/08/68": "unavailable",
            "07/08/68": "unavailable",
        },
        "10:00": {
            "01/08/68": "available",
            "02/08/68": "booked",
            "03/08/68": "available",
            "04/08/68": "unavailable",
            "05/08/68": "unavailable",
            "06/08/68": "unavailable",
            "07/08/68": "unavailable",
        },
        "11:00": {
            "01/08/68": "available",
            "02/08/68": "booked",
            "03/08/68": "available",
            "04/08/68": "unavailable",
            "05/08/68": "unavailable",
            "06/08/68": "unavailable",
            "07/08/68": "unavailable",
        },
        "13:00": {
            "01/08/68": "available",
            "02/08/68": "booked",
            "03/08/68": "available",
            "04/08/68": "unavailable",
            "05/08/68": "unavailable",
            "06/08/68": "unavailable",
            "07/08/68": "unavailable",
        },
    };


    if (!isOpen) return null;
    const today = format(new Date(), "yyyy-MM-dd")
    const [selectedOption, setSelectedOption] = useState<"week" | "month" | "" | null>("");
    const [selectedDate, setSelectedDate] = useState<string>(format(new Date(), "yyyy-MM-dd"));



    const handleClickSlot = (time: string, date: string, status: TimeSlotStatus) => {
        console.log(`คลิก: ${date} ${time} | สถานะ: ${status}`);
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} blurBackground showAds={{ left: true, right: true, bottom: true }}>
            <div className="p-4 border-b border-gray-200 bg-gray-50">
                <div className="flex items-center gap-4 flex-wrap">
                    <label className="text-lg font-semibold">เลือกวันที่คุณว่าง</label>
                </div>
            </div>
            <div className="w-full">
                <div className="mb-6">
                    <div className="overflow-x-auto">
                        <div
                            className="grid"
                            style={{
                                gridTemplateColumns: `100px repeat(${dates.length}, minmax(100px, 1fr))`,
                            }}
                        >
                            {/* Header row */}
                            <div className="bg-gray-200 font-bold border border-gray-300 px-4 py-2">เวลา</div>
                            {dates.map((date) => (
                                <div key={date} className="bg-gray-200 font-bold border border-gray-300 px-4 py-2 text-center">
                                    {date}
                                </div>
                            ))}

                            {/* Time rows */}
                            {timeSlots.map((time) => (
                                <React.Fragment key={time}>
                                    {/* Time column */}
                                    <div className="bg-gray-100 font-medium border border-gray-300 px-4 py-2 text-center">{time}</div>

                                    {/* Status columns per day */}
                                    {dates.map((date) => {
                                        const status = scheduleMatrix[time]?.[date];
                                        const baseStyle = "border border-gray-300 px-4 py-2 text-center cursor-pointer select-none transition";

                                        let statusStyle = "";
                                        if (status === "available") statusStyle = "bg-green-100 text-green-700 hover:bg-green-200";
                                        else if (status === "booked") statusStyle = "bg-yellow-100 text-yellow-700 hover:bg-yellow-200";
                                        else statusStyle = "bg-gray-100 text-gray-500 cursor-not-allowed";

                                        return (
                                            <div
                                                key={`${time}-${date}`}
                                                className={`${baseStyle} ${statusStyle}`}
                                                onClick={() =>
                                                    status !== "unavailable" && handleClickSlot(time, date, status)
                                                }
                                            >
                                                {status === "available" && "ว่าง"}
                                                {status === "booked" && "จองแล้ว"}
                                                {status === "unavailable" && "ไม่ว่าง"}
                                            </div>
                                        );
                                    })}
                                </React.Fragment>
                            ))}
                        </div>
                    </div>

                </div>
            </div>
        </Modal>
    );
}