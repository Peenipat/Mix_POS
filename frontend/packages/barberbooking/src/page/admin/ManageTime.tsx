// pages/ManageTime.tsx

import { useCallback, useState, useEffect, useRef } from "react";
import { format, parseISO } from "date-fns";
import { th } from "date-fns/locale";
import axios from "../../lib/axios";
import BranchCalendar from "./subpage/BranchCalendar";
import { useAppSelector } from "../../store/hook";
import { PencilIcon } from "@heroicons/react/24/outline";
import Modal from "@object/shared/components/Modal";
import Toggle from "@object/shared/components/Toggle";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { editWorkingHourSchema, EditWorkingHourFormData } from "../../schemas/WorkingHourSchema";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";

dayjs.extend(utc);
dayjs.extend(timezone);

type WorkingHour = {
    id: number;
    BranchID: number;
    TenantID: number;
    Weekday: number;
    start_time: string;
    end_time: string;
    IsClosed: boolean;
};

interface OverrideDay {
    date: string;
    start_time: string;
    end_time: string;
    IsClosed: boolean;
}

const weekdays = [
    "อาทิตย์", "จันทร์", "อังคาร", "พุธ", "พฤหัสบดี", "ศุกร์", "เสาร์"
];

export default function ManageTime() {
    const me = useAppSelector((state) => state.auth.me);
    const statusMe = useAppSelector((state) => state.auth.statusMe);
    const tenantId = me?.tenant_ids?.[0];
    const branchId = Number(me?.branch_id);
    const didFetchWorkHours = useRef(false);

    const [workingHours, setWorkingHours] = useState<WorkingHour[]>([]);
    const [loadingHours, setLoadingHours] = useState(false);
    const [errorHours, setErrorHours] = useState<string | null>(null);
    const [isClosedToggle, setIsClosedToggle] = useState(false);

    const loadWorkingHours = useCallback(async () => {
        if (!tenantId || !branchId) return;

        setLoadingHours(true);
        setErrorHours(null);

        try {
            const res = await axios.get<{ status: string; data: WorkingHour[] }>(
                `/barberbooking/tenants/${tenantId}/workinghour/branches/${branchId}`
            );

            if (res.data.status !== "success") {
                throw new Error(res.data.status);
            }

            setWorkingHours(res.data.data);
        } catch (err: any) {
            setErrorHours(err.response?.data?.message || err.message || "Failed to load working hours");
        } finally {
            setLoadingHours(false);
        }
    }, [tenantId, branchId]);

    useEffect(() => {
        if (
            statusMe === "succeeded" &&
            me &&
            tenantId &&
            branchId &&
            !didFetchWorkHours.current
        ) {
            didFetchWorkHours.current = true;
            loadWorkingHours();
        }
    }, [statusMe, me, tenantId, branchId, loadWorkingHours]);

    const [selectedDay, setSelectedDay] = useState<WorkingHour | null>(null);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [hours, setHours] = useState(workingHours);

    const handleEdit = (updated: WorkingHour) => {
        setWorkingHours((prev) =>
            prev.map((item) =>
                item.Weekday === updated.Weekday ? updated : item
            )
        );
    };


    const [overrideDays, setOverrideDays] = useState<OverrideDay[]>([
        { date: "2025-06-24", start_time: "10:00", end_time: "15:00", IsClosed: true },
        { date: "2025-06-30", start_time: "09:00", end_time: "14:00", IsClosed: false },
    ]);

    const [newOverride, setNewOverride] = useState<OverrideDay>({
        date: "", start_time: "08:00", end_time: "17:00", IsClosed: true,
    });

    const handleAddOverride = () => {
        setOverrideDays([...overrideDays, newOverride]);
        setNewOverride({ date: "", start_time: "08:00", end_time: "17:00", IsClosed: true, });
    };

    return (
        <div className="max-w-full mx-auto p-3">
            <h1 className="text-2xl font-bold mb-4">จัดการเวลาทำการของสาขา</h1>
            <section className="mb-8">
                <h2 className="text-xl font-semibold mb-2">เวลาทำการประจำสัปดาห์ </h2>
                <div className="flex flex-wrap gap-3">
                    {workingHours
                        .sort((a, b) => a.Weekday - b.Weekday)
                        .map((day) => {
                            const isClosed = day.IsClosed;
                            const boxClass = isClosed
                                ? "bg-red-100 text-red-700"
                                : "bg-white text-gray-800";

                            let timeDisplay = "-";
                            if (!isClosed && day.start_time && day.end_time) {
                                const startTime = format(parseISO(day.start_time), "HH:mm", { locale: th });
                                const endTime = format(parseISO(day.end_time), "HH:mm", { locale: th });
                                timeDisplay = `${startTime} - ${endTime}`;
                            }

                            return (
                                <div className="relative">
                                    <button
                                        onClick={() => {
                                            setSelectedDay(day);
                                            setIsModalOpen(true);
                                        }}
                                        className="absolute top-1 right-1 text-gray-400 hover:text-gray-600"
                                        aria-label="แก้ไขเวลา"
                                    >
                                        <PencilIcon className="w-4 h-4" />
                                    </button>
                                    <div
                                        className={`p-1 border rounded shadow-sm flex flex-col items-center justify-center text-center w-[120px] ${boxClass}`}
                                    >
                                        <div className="font-semibold">{weekdays[day.Weekday]}</div>
                                        <div className="text-sm">{timeDisplay}</div>
                                    </div>
                                </div>
                            );
                        })}
                </div>
            </section>

            {/* Add Override Form */}
            <section>
                <h2 className="text-xl font-semibold mb-2">เพิ่มวันทำการเฉพาะกิจ</h2>
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
                    </div>
                    <button className="bg-blue-500 text-white p-2 rounded-md" onClick={handleAddOverride}>เพิ่มวันทำการ</button>
                </div>
            </section>
            <BranchCalendar />


            <EditWorkingHourModal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onEdit={handleEdit}
                workingHour={selectedDay ?? undefined}
                tenantId={tenantId ?? null}
                branchId={branchId ?? null}
            />
        </div>
    );
}

interface EditWorkingHourModalProps {
    isOpen: boolean;
    onClose: () => void;
    onEdit: (updated: WorkingHour) => void;
    workingHour: WorkingHour | undefined;
    tenantId: number | null;
    branchId: number | null;
}

export function EditWorkingHourModal({
    isOpen,
    onClose,
    onEdit,
    workingHour,
    tenantId,
    branchId,
}: EditWorkingHourModalProps) {
    const {
        register,
        handleSubmit,
        reset,
        formState: { errors, isSubmitting },
    } = useForm<EditWorkingHourFormData>({
        resolver: zodResolver(editWorkingHourSchema),
        defaultValues: {
            start_time: "",
            end_time: "",
        },
    });

    console.log(workingHour)

    useEffect(() => {
        if (isOpen && workingHour) {
            const validStart = workingHour.start_time ? new Date(workingHour.start_time) : null;
            const validEnd = workingHour.end_time ? new Date(workingHour.end_time) : null;

            reset({
                start_time: validStart && !isNaN(validStart.getTime()) ? format(validStart, "HH:mm") : "",
                end_time: validEnd && !isNaN(validEnd.getTime()) ? format(validEnd, "HH:mm") : "",
            });
        }
    }, [isOpen, workingHour, reset]);

    const [isClosed, setIsClosed] = useState<boolean>(workingHour?.IsClosed ?? false);
    useEffect(() => {
        setIsClosed(workingHour?.IsClosed ?? false);
    }, [workingHour]);


    const onSubmit = async (data: EditWorkingHourFormData) => {
        if (!tenantId || !branchId || !workingHour) return;

        try {
            const payload = isClosed
            ? [{
                weekday: workingHour.Weekday,
                start_time: null,
                end_time: null,
                is_closed: true,
              }]
            : [{
                weekday: workingHour.Weekday,
                start_time: dayjs.tz(`2025-01-01T${data.start_time}:00`, "Asia/Bangkok").toISOString(),
                end_time: dayjs.tz(`2025-01-01T${data.end_time}:00`, "Asia/Bangkok").toISOString(),
                is_closed: false,
              }];

            console.log(payload)

            const res = await axios.put(
                `/barberbooking/tenants/${tenantId}/workinghour/branches/${branchId}`,
                payload
            );

            if (res.data.status !== "success") throw new Error("Update failed");

            onEdit({
                ...workingHour,
                start_time: payload[0].start_time ?? "",
                end_time: payload[0].end_time ?? "",
                IsClosed: isClosed,
            });
            onClose();
        } catch (err) {
            console.error("Failed to update working hour", err);
        }
    };


    if (!isOpen || !workingHour) return null;
    console.log("📆 Working day selected:", workingHour?.Weekday);

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="แก้ไขเวลาเปิด-ปิด">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                {!isClosed && (
                    <>
                        <div className="mb-4">
                            <label className="block mb-1">เวลาเปิด</label>
                            <input type="time" {...register("start_time")} defaultValue={workingHour?.start_time?.slice(11, 16)} className="input input-bordered w-full" />
                        </div>
                        <div className="mb-4">
                            <label className="block mb-1">เวลาปิด</label>
                            <input type="time" {...register("end_time")} defaultValue={workingHour?.end_time?.slice(11, 16)} className="input input-bordered w-full" />

                        </div>
                    </>)}

                <div className="mb-4">
                    <label className="inline-flex items-center cursor-pointer">
                        <Toggle checked={isClosed} onChange={setIsClosed} label=" วันหยุดประจำสัปดาห์" />
                    </label>
                </div>
                <div className="flex justify-end space-x-2">
                    <button type="button" onClick={onClose} className="btn btn-ghost" disabled={isSubmitting}>
                        ยกเลิก
                    </button>
                    <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
                        {isSubmitting ? "กำลังบันทึก…" : "ยืนยัน"}
                    </button>
                </div>
            </form>
        </Modal>
    );
}


