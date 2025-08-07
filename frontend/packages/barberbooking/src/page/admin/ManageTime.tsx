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
import { WorkingHour } from "../../api/workingHour";
import { getWorkingHours } from "../../api/workingHour";
import { createWorkingDayOverride, WorkingDayOverrideInput } from "../../api/workingDayOverride";

dayjs.extend(utc);
dayjs.extend(timezone);

export interface OverrideDay {
    date: string;
    start_time: string;
    end_time: string;
    IsClosed: boolean;
    reason:string
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

    useEffect(() => {
        const fetchWorkingHours = async () => {
            try {
                if (tenantId !== undefined && branchId !== undefined) {
                    const data = await getWorkingHours({ tenantId, branchId });
                    setWorkingHours(data);
                }

            } catch (err) {
                console.error(err);
            }
        };
        fetchWorkingHours();
    }, [tenantId, branchId]);


    const [selectedDay, setSelectedDay] = useState<WorkingHour | null>(null);
    const [isModalOpen, setIsModalOpen] = useState(false);

    const handleEdit = (updated: WorkingHour) => {
        setWorkingHours((prev) =>
            prev.map((item) =>
                item.week_day === updated.week_day ? updated : item
            )
        );
    };



    const [newOverride, setNewOverride] = useState({
        date: "",
        start_time: "",
        end_time: "",
        reason: "",
    });

    const handleAddOverride = async () => {

        try {
            const input: WorkingDayOverrideInput = {
                branch_id: 1,
                work_date: newOverride.date,
                start_time: isClosed ? "00:00" : newOverride.start_time,
                end_time: isClosed ? "00:00" : newOverride.end_time,
                is_closed: isClosed,
                reason: newOverride.reason.trim(),
            };

            const response = await createWorkingDayOverride(input);

            alert("เพิ่มข้อมูลวันทำการสำเร็จแล้ว!");
            console.log("response:", response);

            setNewOverride({
                date: "",
                start_time: "",
                end_time: "",
                reason: "",
            });
            setIsClosed(false);
        } catch (error) {
            console.error("เกิดข้อผิดพลาดในการบันทึก:", error);
            alert("ไม่สามารถเพิ่มวันทำการได้ โปรดลองใหม่");
        }
    };


    const [isClosed, setIsClosed] = useState(false)

    return (
        <div className="max-w-full mx-auto p-3">
            <h1 className="text-2xl font-bold mb-4">จัดการเวลาทำการของสาขา</h1>
            <section className="mb-8">
                <h2 className="text-xl font-semibold mb-2">เวลาทำการประจำสัปดาห์ </h2>
                <div className="flex flex-wrap gap-3">
                    {workingHours
                        .sort((a, b) => a.week_day - b.week_day)
                        .map((day) => {
                            const isClosed = day.is_closed;
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
                                        <div className="font-semibold">{weekdays[day.week_day]}</div>
                                        <div className="text-sm">{timeDisplay}</div>
                                    </div>
                                </div>
                            );
                        })}
                </div>
            </section>

            {/* Add Override Form */}
            <section>
                <h2 className="text-xl font-semibold mb-2">เพิ่มเวลาเปิด - ปิด กรณีพิเศษ</h2>

                <div className="space-y-4">
                    <input
                        type="date"
                        className="input input-bordered w-full max-w-xs"
                        value={newOverride.date}
                        onChange={(e) => setNewOverride({ ...newOverride, date: e.target.value })}
                    />

                    <div className="flex flex-col">
                        <div className="flex gap-6">
                            <label className="inline-flex items-center">
                                <input
                                    type="checkbox"
                                    checked={!isClosed}
                                    onChange={() => setIsClosed(false)}
                                    className="w-4 h-4 text-green-600 border-gray-300"
                                />
                                <span className="ml-2">เปิดร้าน</span>
                            </label>
                            <label className="inline-flex items-center">
                                <input
                                    type="checkbox"
                                    checked={isClosed}
                                    onChange={() => setIsClosed(true)}
                                    className="w-4 h-4 text-red-600 border-gray-300"
                                />
                                <span className="ml-2">ปิดร้าน</span>
                            </label>
                        </div>


                        <div className="flex space-x-4">

                            {!isClosed && (
                                <>
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
                                </>
                            )}
                            <div>
                                <label className="block text-sm font-medium">
                                    หมายเหตุ
                                </label>
                                <input
                                    type="text"
                                    placeholder="หมายเหตุ"
                                    className={`w-full input input-bordered`}
                                    value={newOverride.reason}
                                    onChange={(e) => setNewOverride({ ...newOverride, reason: e.target.value })}
                                />
                            </div>
                        </div>
                    </div>

                    <button className="bg-blue-500 text-white p-2 rounded-md" onClick={handleAddOverride}>เพิ่มวันทำการ</button>
                </div>
            </section>
            <BranchCalendar workingHours={workingHours} />


            <WorkingHourModal
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

    const [isClosed, setIsClosed] = useState<boolean>(workingHour?.is_closed ?? false);
    useEffect(() => {
        setIsClosed(workingHour?.is_closed ?? false);
    }, [workingHour]);


    const onSubmit = async (data: EditWorkingHourFormData) => {
        if (!tenantId || !branchId || !workingHour) return;

        try {
            const payload = isClosed
                ? [{
                    weekday: workingHour.week_day,
                    start_time: null,
                    end_time: null,
                    is_closed: true,
                }]
                : [{
                    weekday: workingHour.week_day,
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
                is_closed: isClosed,
            });
            onClose();
        } catch (err) {
            console.error("Failed to update working hour", err);
        }
    };


    if (!isOpen || !workingHour) return null;

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


